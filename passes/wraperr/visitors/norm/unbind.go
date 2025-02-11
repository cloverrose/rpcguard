package norm

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/ssa"
)

// This file contains bounded method wrapper support.
// Below sample code,
// Hello() method returns `a.hi`.
// Since it has receiver "a", a.hi here is bound method wrapper for `a.hi` by synthetic function.
// Unbind() unbinds and return original `a.hi`.
/**
type A struct{}

func (a *A) Hello() func() error {
	return a.hi // (1) this a.hi is `(*A).hi$bound`, we want to get `(*A).hi` so that we can match with (2)
}

func (a *A) hi() error { // (2) this hi is `(*A).hi`
	return nil
}
**/

// unbind unbinds wrapper func and return original func.
// If fn is not wrapper func returns fn itself.
func unbind(fn *ssa.Function) (*ssa.Function, error) {
	if !isWrapper(fn) {
		return fn, nil
	}
	return unwrap(fn)
}

// unwrap returns a function that is wrapped by wrapperFunc.
func unwrap(wrapperFunc *ssa.Function) (*ssa.Function, error) {
	if len(wrapperFunc.Blocks) != 1 {
		panic(fmt.Sprintf("unexpected wrapper func (len(Blocks) != 1) %v", wrapperFunc.Name()))
	}
	block := wrapperFunc.Blocks[0]
	if len(block.Instrs) < 2 {
		// Instrs should be Call, (Extract,.., Extract), Return.
		panic(fmt.Sprintf("unexpected wrapper func (len(Instrs) != 2) %v", wrapperFunc.Name()))
	}
	call, ok := block.Instrs[0].(*ssa.Call)
	if !ok {
		panic(fmt.Sprintf("unexpected wrapper func (Instrs[0] != *ssa.Call) %v", wrapperFunc.Name()))
	}
	if call.Call.IsInvoke() {
		// eg.Go calls interface method
		return wrapperFunc, nil
	}
	fn, ok := call.Call.Value.(*ssa.Function)
	if !ok {
		panic(fmt.Sprintf("unexpected wrapper func (Call.Value != *ssa.Function) %v", wrapperFunc.Name()))
	}
	return fn, nil
}

// isWrapper returns true if fn is a bound method wrapper for some other func.
func isWrapper(fn *ssa.Function) bool {
	if fn.Synthetic == "" {
		return false
	}
	if strings.HasPrefix(fn.Synthetic, "bound method wrapper for ") {
		return true
	}
	return false
}
