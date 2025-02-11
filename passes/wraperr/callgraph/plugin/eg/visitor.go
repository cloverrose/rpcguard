package eg

import (
	"golang.org/x/tools/go/ssa"

	"github.com/cloverrose/rpcguard/pkg/errtrace/ssawalk"

	"github.com/cloverrose/rpcguard/passes/wraperr/visitors/call2func"
	"github.com/cloverrose/rpcguard/passes/wraperr/visitors/norm"
)

type visitorPlugin struct {
	replacements map[*ssa.Function][]*ssa.Function
}

func newVisitorPlugin() *visitorPlugin {
	return &visitorPlugin{
		replacements: make(map[*ssa.Function][]*ssa.Function),
	}
}

func (p *visitorPlugin) createOptions() []ssawalk.Option {
	return []ssawalk.Option{
		ssawalk.WithVisitCall(p.VisitCall),
	}
}

func (p *visitorPlugin) VisitCall(call *ssa.Call) error {
	// Check if call is eg.Wait()
	waitFunc, err := call2func.GetFuncFromCall(call)
	if err != nil {
		return err
	}
	if waitFunc == nil {
		return nil
	}
	if !isWait(waitFunc) {
		return nil
	}

	// Find eg.Wait() receiver eg
	receiver, err := getWaitReceiver(call.Call)
	if err != nil {
		return err
	}

	// Find eg.Go(func) funcs
	refs := receiver.Referrers()
	if refs == nil {
		// eg.Wait is called but there is no eg.Go.
		return nil
	}
	for _, instr := range *refs {
		goCall, ok := instr.(*ssa.Call)
		if !ok {
			continue
		}
		goArgFunc, err := getGoArgFunc(goCall, receiver)
		if err != nil {
			return err
		}
		if goArgFunc != nil {
			p.replacements[waitFunc] = append(p.replacements[waitFunc], goArgFunc)
		}
	}
	return nil
}

func getGoArgFunc(egCall *ssa.Call, waitReceiver ssa.Value) (*ssa.Function, error) {
	// Check if call is eg.Go()
	goFunc, err := call2func.GetFuncFromCall(egCall)
	if err != nil {
		return nil, err
	}
	if goFunc == nil {
		return nil, nil
	}
	if !isGo(goFunc) {
		return nil, nil
	}

	// Get eg.Go(fun) receiver and argument fun.
	goReceiver, goArg, err := getGoArg(egCall.Call)
	if err != nil {
		return nil, err
	}

	// Check if eg.Go() receiver is the same with eg.Wait() receiver.
	if goReceiver != waitReceiver {
		panic("unexpected: goReceiver != waitReceiver")
	}

	// Find eg.Go(fun) argument func from goArg.
	// goArg can be not *ssa.Function
	//   goArg is *ssa.Function: eg.Go(fun): return fun.
	//   goArg is *ssa.Call: eg.Go(higherOrderFun()): return higherOrderFun.
	goArgFunc, err := call2func.GetFuncFromCall(goArg)
	if err != nil {
		return nil, err
	}
	if goArgFunc == nil {
		return nil, nil
	}

	// normalize
	normed, err := norm.NewNormalizeFunc([]int{0})(goArgFunc)
	if err != nil {
		return nil, err
	}
	return normed, nil
}
