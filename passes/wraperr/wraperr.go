package wraperr

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"

	"github.com/cloverrose/rpcguard/pkg/factutil"
	"github.com/cloverrose/rpcguard/pkg/filter"
	"github.com/cloverrose/rpcguard/pkg/graph"
	"github.com/cloverrose/rpcguard/pkg/logger"
	"github.com/cloverrose/rpcguard/pkg/rpcmethod"
	"github.com/cloverrose/rpcguard/pkg/signature"

	"github.com/cloverrose/rpcguard/passes/wraperr/callgraph"
	"github.com/cloverrose/rpcguard/passes/wraperr/callgraph/plugin/eg"
)

const (
	doc               = "rpc_wraperr checks if errors returned in RPC method is wrapped by connect.NewError."
	reportMsg         = "RPC method %s returns error that is not wrapped with connect.NewError"
	packageKey        = "package"
	unexpectedUnknown = "unexpected KindUnknown"
)

// Analyzer checks if RPC method returns error properly.
var Analyzer = &analysis.Analyzer{
	Name: "rpc_wraperr",
	Doc:  doc,
	Run:  setupAndRun,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
	Flags: *flag.NewFlagSet("rpc_wraperr", flag.ExitOnError),
	FactTypes: []analysis.Fact{
		&isErrorHandler{},
	},
}

func init() {
	Analyzer.Flags.StringVar(&LogConfig.Level, "log.level", LogConfig.Level, "logging level. debug, info, warn, error")
	Analyzer.Flags.StringVar(&LogConfig.File, "log.file", LogConfig.File, "log file path.")
	Analyzer.Flags.StringVar(&LogConfig.Format, "log.format", LogConfig.Format, "logging format. json or text")
	Analyzer.Flags.StringVar(&ReportMode, "ReportMode", ReportMode, "reporting mode (RETURN, FUNCTION, BOTH)")
	Analyzer.Flags.StringVar(&IncludePackages, "IncludePackages", IncludePackages, "include packages")
	Analyzer.Flags.StringVar(&ExcludePackages, "ExcludePackages", ExcludePackages, "exclude packages")
	Analyzer.Flags.StringVar(&ExcludeFiles, "ExcludeFiles", ExcludeFiles, "exclude files")
	Analyzer.Flags.BoolVar(&EnableErrGroupAnalyzer, "EnableErrGroupAnalyzer", EnableErrGroupAnalyzer, "enable ErrGroupAnalyzer (default true)")
}

func setupAndRun(pass *analysis.Pass) (any, error) {
	closer, err := logger.SetDefault(LogConfig, pass)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := closer(); err != nil {
			fmt.Println(err)
		}
	}()

	packageFilter, err = filter.New(IncludePackages, ExcludePackages)
	if err != nil {
		return nil, err
	}

	// any files that are not excluded are target.
	fileFilter, err = filter.New(`.*`, ExcludeFiles)
	if err != nil {
		return nil, err
	}

	return run(pass)
}

//nolint:gocognit,gocyclo,cyclop // main routine
func run(pass *analysis.Pass) (interface{}, error) {
	currentPackage := pass.Pkg.Path()
	slog.Debug("analyzing package", slog.String(packageKey, currentPackage))

	// Phase1: Package is target?
	if !packageFilter.IsTarget(currentPackage) {
		slog.Debug("skip package (not target package)", slog.String(packageKey, currentPackage))
		return nil, nil
	}

	// Phase 2: Get SSA
	ssaData, ok := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	if !ok {
		panic("failed to get SSA")
	}

	// Phase 3: Func is target?
	targetSrcFuncs := make([]*ssa.Function, 0, len(ssaData.SrcFuncs))
	for _, srcFunc := range ssaData.SrcFuncs {
		if isTargetFunc(pass, srcFunc) {
			targetSrcFuncs = append(targetSrcFuncs, srcFunc)
		}
	}

	// Phase 4: Build Call Graph
	cg := callgraph.New(signature.ErrIshIndices)
	for _, srcFunc := range targetSrcFuncs {
		if err := cg.Scan(srcFunc); err != nil {
			return nil, err
		}
	}
	slog.Debug("build callgraph", slog.Any("callgraph", cg))

	// Phase 5: Analyze go.Wait and go.Go
	if EnableErrGroupAnalyzer {
		for _, srcFunc := range targetSrcFuncs {
			if err := cg.ScanWithPlugin(eg.Scan, srcFunc); err != nil {
				return nil, err
			}
		}
		slog.Debug("build callgraph with errgroup", slog.Any("callgraph", cg))
	}

	// closure func (ends with $1) Object() == nil, then we can't export facts.
	// Thus, use this localFacts during package check.
	factWrapper := factutil.NewFactWrapper[*isErrorHandler](pass)

	// Phase 6: Export obvious facts.
	for _, srcFunc := range targetSrcFuncs {
		info := cg.GetReturnInfo(srcFunc)
		if info == nil {
			panic(fmt.Sprintf("unexpected info not found for srcFunc: %s", srcFunc.Name()))
		}
		if info.IsObviouslyBad() {
			factWrapper.Export(srcFunc, &isErrorHandler{Kind: KindBad})
		}
		if info.IsObviouslyOK() {
			factWrapper.Export(srcFunc, &isErrorHandler{Kind: KindOK})
		}
	}

	// Phase 7: Create SCCs (this sccs are topologically sorted)
	g := cg.Convert()
	slog.Debug("build graph", slog.Any("graph", g))
	sccs := graph.Decomposition(g)
	slog.Debug("Strongly Connected Components", slog.Any("sccs", sccs))

	if err := markSCCs(pass, sccs, factWrapper, cg); err != nil {
		return nil, err
	}

	// Phase 8: Check RPC method is marked with bad or not.
	rpcChecker := rpcmethod.BuildChecker(pass)
	if rpcChecker == nil {
		slog.Debug("skip package (no rpc method types)", slog.String(packageKey, currentPackage))
		return nil, nil
	}
	for _, fn := range targetSrcFuncs {
		if !rpcChecker.IsRPCMethod(fn) {
			continue
		}
		slog.Info("found RPC method", logger.Attr(fn))

		fact, ok := factWrapper.Import(fn)
		if ok {
			switch fact.Kind {
			case KindUnknown:
				panic(unexpectedUnknown)
			case KindBad:
				if ReportMode == reportModeFunction || ReportMode == reportModeBoth {
					pass.Reportf(fn.Pos(), reportMsg, fn.Name())
				}
				if ReportMode == reportModeReturn || ReportMode == reportModeBoth {
					info := cg.GetReturnInfo(fn)
					for _, rtn := range info.GetReturns() {
						reportReturn(pass, factWrapper, info, fn, rtn)
					}
				}
			case KindOK:
			}
		}
	}

	return nil, nil
}

func reportReturn(
	pass *analysis.Pass,
	factWrapper *factutil.FactWrapper[*isErrorHandler],
	info *callgraph.FuncInfo,
	fn *ssa.Function,
	rtn *ssa.Return,
) {
	for _, toFunc := range info.GetToFuncs(rtn) {
		if fact, ok := factWrapper.Import(toFunc); ok && fact.Kind == KindBad {
			pass.Reportf(rtn.Pos(), reportMsg, fn.Name())
			return
		}
	}
	if info.IsObviouslyBadReturn(rtn) {
		pass.Reportf(rtn.Pos(), reportMsg, fn.Name())
	}
}

func isTargetFunc(pass *analysis.Pass, srcFunc *ssa.Function) bool {
	if srcFunc == nil {
		panic("srcFunc is nil")
	}
	fileName := pass.Fset.Position(srcFunc.Pos()).Filename
	if srcFunc.Pkg == nil {
		panic("srcFunc.Pkg is nil")
	}
	if srcFunc.Pkg.Pkg == nil {
		panic("srcFunc.Pkg.Pkg is nil")
	}
	if !packageFilter.IsTarget(srcFunc.Pkg.Pkg.Path()) {
		panic("!packageFilter.IsTarget(srcFunc.Pkg.Pkg.Path())")
	}
	if !fileFilter.IsTarget(fileName) {
		slog.Debug("skip Function (non target file)", logger.Attr(srcFunc))
		return false
	}

	// srcFunc signature is not target?
	indices := signature.ErrIshIndices(srcFunc)
	return len(indices) != 0
}

type SCCs [][]*ssa.Function

func (sccs SCCs) LogValue() slog.Value {
	coreData := make([][]string, len(sccs))
	for i, scc := range sccs {
		coreData2 := make([]string, len(scc))
		for j, fn := range scc {
			coreData2[j] = fn.String()
		}
		coreData[i] = coreData2
	}
	jsonBytes, err := json.Marshal(coreData)
	if err != nil {
		return slog.Value{}
	}
	return slog.StringValue(string(jsonBytes))
}
