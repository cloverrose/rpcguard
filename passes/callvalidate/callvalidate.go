package callvalidate

import (
	"flag"
	"fmt"
	"log/slog"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"

	"github.com/cloverrose/rpcguard/pkg/filter"
	"github.com/cloverrose/rpcguard/pkg/logger"
	"github.com/cloverrose/rpcguard/pkg/rpcmethod"
)

const (
	doc       = "rpc_callvalidate checks if RPC method uses Validate method properly."
	reportMsg = "RPC method %s does not use protovalidate.Validate properly"
)

// Analyzer checks if RPC method uses Validate method properly.
var Analyzer = &analysis.Analyzer{
	Name: "rpc_callvalidate",
	Doc:  doc,
	Run:  setupAndRun,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
	Flags: *flag.NewFlagSet("rpc_callvalidate", flag.ExitOnError),
}

func init() {
	Analyzer.Flags.StringVar(&LogConfig.Level, "log.level", LogConfig.Level, "logging level. debug, info, warn, error")
	Analyzer.Flags.StringVar(&LogConfig.File, "log.file", LogConfig.File, "log file path.")
	Analyzer.Flags.StringVar(&LogConfig.Format, "log.format", LogConfig.Format, "logging format. json or text")
	Analyzer.Flags.StringVar(&ExcludeFiles, "ExcludeFiles", ExcludeFiles, "exclude files")
	Analyzer.Flags.StringVar(&ValidateMethods, "ValidateMethods", ValidateMethods, "Validate methods")
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

	// any files that are not excluded are target.
	fileFilter, err = filter.New(`.*`, ExcludeFiles)
	if err != nil {
		return nil, err
	}

	return run(pass)
}

func run(pass *analysis.Pass) (interface{}, error) {
	currentPackage := pass.Pkg.Path()
	slog.Debug("analyzing package", slog.String("package", currentPackage))

	// Phase 1: Get SSA
	ssaData, ok := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	if !ok {
		panic("failed to get SSA")
	}

	// Phase 2: Func is target?
	targetSrcFuncs := make([]*ssa.Function, 0, len(ssaData.SrcFuncs))
	for _, srcFunc := range ssaData.SrcFuncs {
		if isTargetFunc(pass, srcFunc) {
			targetSrcFuncs = append(targetSrcFuncs, srcFunc)
		}
	}

	// Phase 3: Func is RPC method?.
	rpcAnalyzer := rpcmethod.BuildChecker(pass)
	if rpcAnalyzer == nil {
		slog.Debug("skip package (no rpc method types)", slog.String("package", currentPackage))
		return nil, nil
	}
	rpcMethods := make([]*ssa.Function, 0, len(targetSrcFuncs))
	for _, fn := range targetSrcFuncs {
		if rpcAnalyzer.IsRPCMethod(fn) {
			rpcMethods = append(rpcMethods, fn)
		}
	}

	// Phase 4: Check validate call
	validateMethods, err := parseMethods(ValidateMethods)
	if err != nil {
		return nil, err
	}
	for _, srcFunc := range rpcMethods {
		ok, err := checkCallValidate(srcFunc, validateMethods)
		if err != nil {
			return nil, err
		}
		if !ok {
			pass.Reportf(srcFunc.Pos(), reportMsg, srcFunc.Name())
		}
	}

	return nil, nil
}

func isTargetFunc(pass *analysis.Pass, srcFunc *ssa.Function) bool {
	if srcFunc == nil {
		panic("srcFunc is nil")
	}
	fileName := pass.Fset.Position(srcFunc.Pos()).Filename
	if !fileFilter.IsTarget(fileName) {
		slog.Debug("skip Function (non target file)", logger.Attr(srcFunc))
		return false
	}
	return true
}
