package wraperr

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"

	"github.com/cloverrose/rpcguard/pkg/factutil"
	"github.com/cloverrose/rpcguard/pkg/logger"

	"github.com/cloverrose/rpcguard/passes/wraperr/callgraph"
)

type factImporter interface {
	Import(fn *ssa.Function) (*isErrorHandler, bool)
}

// markSCCs checks and records facts.
func markSCCs(pass *analysis.Pass, sccs [][]*ssa.Function, factWrapper *factutil.FactWrapper[*isErrorHandler],
	cg *callgraph.CallGraph,
) error {
	for _, scc := range sccs {
		bad, err := checkSCC(pass, scc, factWrapper, cg)
		if err != nil {
			return err
		}
		propagateMarkToSCC(scc, bad, factWrapper)
	}
	return nil
}

// checkSCC checks if given scc has bad error sources or not.
func checkSCC(
	pass *analysis.Pass,
	scc []*ssa.Function,
	factWrapper factImporter,
	cg *callgraph.CallGraph,
) (bool, error) {
	for _, fromFunc := range scc {
		bad, err := checkFunc(pass, factWrapper, fromFunc, scc, cg)
		if err != nil {
			return false, err
		}
		if bad {
			return true, nil
		}
	}
	return false, nil
}

func checkFunc(
	pass *analysis.Pass,
	factWrapper factImporter,
	srcFunc *ssa.Function,
	scc []*ssa.Function,
	cg *callgraph.CallGraph,
) (bool, error) {
	slog.Debug("check srcFunc", logger.Attr(srcFunc))

	fact, ok := factWrapper.Import(srcFunc)
	if ok {
		switch fact.Kind {
		case KindUnknown:
			panic(unexpectedUnknown)
		case KindBad:
			return true, nil
		case KindOK:
			return false, nil
		}
	}
	if isConnectNewError(srcFunc) {
		return false, nil
	}
	if srcFunc == nil {
		panic("unexpected srcFunc is nil")
	}
	if srcFunc.Pkg == nil {
		// This happens when interface method is assigned to local variable.
		// E.g. fn := app.handler.Handle
		slog.Debug("found bad func (srcFunc.Pkg is nil)", logger.Attr(srcFunc))
		return true, nil
	}
	if srcFunc.Pkg.Pkg != pass.Pkg {
		// srcFunc is defined in different package.
		// and fact is unknown, so it is unknown bad func.
		slog.Debug("found bad func (srcFunc.Pkg.Pkg != pass.Pkg)", logger.Attr(srcFunc))
		return true, nil
	}

	var bad bool
	info := cg.GetReturnInfo(srcFunc)
	if info == nil {
		panic(fmt.Sprintf("unexpected info not found for srcFunc: %s", srcFunc.Name()))
	}
	for _, toFunc := range info.GetAllToFuncs() {
		if slices.Contains(scc, toFunc) {
			slog.Debug("toFunc is in the same SCC", logger.Attr(toFunc))
			continue
		}
		if isConnectNewError(toFunc) {
			continue
		}
		bad = checkBad(toFunc, factWrapper)
		if bad {
			break
		}
	}
	if bad {
		slog.Debug("func is bad func", logger.Attr(srcFunc))
	} else {
		slog.Debug("func is not bad (still suspicious)", logger.Attr(srcFunc))
	}
	return bad, nil
}

func checkBad(toFunc *ssa.Function, factWrapper factImporter) bool {
	fact, ok := factWrapper.Import(toFunc)
	if ok {
		switch fact.Kind {
		case KindUnknown:
			panic(unexpectedUnknown)
		case KindBad:
			return true
		case KindOK:
			return false
		}
	}
	// toFunc is not in facts => badFunc
	return true
}

// isConnectNewError returns true if fn is connect.NewError.
func isConnectNewError(fn *ssa.Function) bool {
	if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg == nil {
		return false
	}
	path := fn.Pkg.Pkg.Path()
	const connectPath = "connectrpc.com/connect"
	return (path == connectPath || strings.HasSuffix(path, "vendor/"+connectPath)) && fn.Name() == "NewError"
}

func propagateMarkToSCC(scc []*ssa.Function, bad bool, factWrapper *factutil.FactWrapper[*isErrorHandler]) {
	for _, fn := range scc {
		slog.Debug("propagate mark", logger.Attr(fn), slog.Bool("bad", bad))
		// Export kind
		if bad {
			factWrapper.Export(fn, &isErrorHandler{Kind: KindBad})
		} else {
			factWrapper.Export(fn, &isErrorHandler{Kind: KindOK})
		}
	}
}
