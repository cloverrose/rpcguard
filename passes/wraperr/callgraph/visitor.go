package callgraph

import (
	"fmt"

	"golang.org/x/tools/go/ssa"

	"github.com/gostaticanalysis/analysisutil"

	"github.com/cloverrose/rpcguard/pkg/errtrace/ssawalk"
)

// visitorPlugin visit ssa.Value and collects ssa.Function that are source of returned value.
// If it visits ssa.Value that is not func, ignore.
type visitorPlugin struct {
	normalizeFunc func(fn *ssa.Function) (*ssa.Function, error)

	toFuncs []*ssa.Function
	isBad   bool
}

func (p *visitorPlugin) VisitFunction(val *ssa.Function) error {
	normed, err := p.normalizeFunc(val)
	if err != nil {
		return err
	}
	p.toFuncs = append(p.toFuncs, normed)
	return nil
}

func (p *visitorPlugin) VisitConst(val *ssa.Const) error {
	if !analysisutil.ImplementsError(val.Type()) {
		// srcFunc calls const func. e.g.
		// func Hello() error {
		//   var constFunc func() error
		//   return constFunc()
		// }
		p.isBad = true
	}
	if !val.IsNil() {
		panic(fmt.Sprintf("unexpected const %s\n", val.Name()))
	}
	return nil
}

func (p *visitorPlugin) VisitAlloc(val *ssa.Alloc) error {
	// usually alloc err is badFunc
	p.isBad = true
	return nil
}

func (p *visitorPlugin) VisitComplex(val ssa.Value) error {
	// can't analyze further due to complexity, treat badFunc.
	p.isBad = true
	return nil
}

func (p *visitorPlugin) VisitCallInvoke(val *ssa.Call) error {
	// can't analyze further due to interface method.
	// However, interface method is usually defined in different package.
	// Then it can be considered badFunc.
	p.isBad = true
	return nil
}

func (p *visitorPlugin) createOptions() []ssawalk.Option {
	return []ssawalk.Option{
		ssawalk.WithVisitFunction(p.VisitFunction),
		ssawalk.WithVisitConst(p.VisitConst),
		ssawalk.WithVisitAlloc(p.VisitAlloc),
		ssawalk.WithVisitComplex(p.VisitComplex),
		ssawalk.WithVisitCallInvoke(p.VisitCallInvoke),
	}
}
