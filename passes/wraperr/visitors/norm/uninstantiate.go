package norm

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/tools/go/ssa"

	"github.com/cloverrose/rpcguard/pkg/errtrace/ssawalk"

	"github.com/cloverrose/rpcguard/passes/wraperr/rtn"
)

type visitorPlugin struct {
	ret *ssa.Function
}

func (p *visitorPlugin) VisitFunction(val *ssa.Function) error {
	if p.ret != nil {
		return errors.New("unexpected ret is not nil")
	}
	p.ret = val
	return nil
}

func (p *visitorPlugin) createOptions() []ssawalk.Option {
	return []ssawalk.Option{
		ssawalk.WithVisitFunction(p.VisitFunction),
		ssawalk.WithVisitConst(func(_ *ssa.Const) error {
			panic("unexpected const")
		}),
		ssawalk.WithVisitAlloc(func(_ *ssa.Alloc) error {
			panic("unexpected alloc")
		}),
		ssawalk.WithVisitComplex(func(_ ssa.Value) error {
			panic("unexpected complex")
		}),
		ssawalk.WithVisitCallInvoke(func(_ *ssa.Call) error {
			panic("unexpected call invoke")
		}),
	}
}

func uninstantiate(fn *ssa.Function, indices []int) (*ssa.Function, error) {
	if !strings.HasPrefix(fn.Synthetic, "instantiation wrapper of ") {
		// return given fn.
		return fn, nil
	}
	plugin := &visitorPlugin{}
	visitor := ssawalk.NewDefaultVisitorWith(plugin.createOptions()...)
	for _, val := range rtn.GetReturnsAt(fn, indices) {
		if err := ssawalk.Walk(visitor, val.Value); err != nil {
			return nil, err
		}
	}
	if plugin.ret == nil {
		return nil, fmt.Errorf("unexpected: walk fail to unwrap function %v", fn.Name())
	}
	return plugin.ret, nil
}
