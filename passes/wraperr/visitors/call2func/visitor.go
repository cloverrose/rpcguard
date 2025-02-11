package call2func

import (
	"errors"

	"golang.org/x/tools/go/ssa"

	"github.com/cloverrose/rpcguard/pkg/errtrace/ssawalk"
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
	}
}

func GetFuncFromCall(value ssa.Value) (*ssa.Function, error) {
	plugin := &visitorPlugin{}
	visitor := ssawalk.NewDefaultVisitorWith(plugin.createOptions()...)
	if err := ssawalk.Walk(visitor, value); err != nil {
		return nil, err
	}
	if plugin.ret == nil {
		// srcFnc calls func that is not able to analyze like func in parameter.
		return nil, nil
	}
	return plugin.ret, nil
}
