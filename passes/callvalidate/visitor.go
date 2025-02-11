package callvalidate

import (
	"errors"

	"golang.org/x/tools/go/ssa"

	"github.com/cloverrose/rpcguard/pkg/errtrace/ssawalk"
)

// visitorPlugin visit ssa.Value and collects ssa.Function that are source of returned value.
type visitorPlugin struct {
	fn *ssa.Function
}

func (p *visitorPlugin) VisitFunction(val *ssa.Function) error {
	if p.fn != nil {
		return errors.New("fn is not nil")
	}
	p.fn = val
	return nil
}

func (p *visitorPlugin) createOptions() []ssawalk.Option {
	return []ssawalk.Option{
		ssawalk.WithVisitFunction(p.VisitFunction),
	}
}
