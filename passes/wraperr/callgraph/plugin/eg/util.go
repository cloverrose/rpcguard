package eg

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/tools/go/ssa"
)

// isGo returns true if fn is (*golang.org/x/sync/errgroup.Group).Go()
func isGo(fn *ssa.Function) bool {
	if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg == nil {
		return false
	}
	path := fn.Pkg.Pkg.Path()
	const errGroupPath = "golang.org/x/sync/errgroup"
	return (path == errGroupPath || strings.HasSuffix(path, "vendor/"+errGroupPath)) && fn.Name() == "Go"
}

// isWait returns true if fn is (*golang.org/x/sync/errgroup.Group).Wait()
func isWait(fn *ssa.Function) bool {
	if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg == nil {
		return false
	}
	path := fn.Pkg.Pkg.Path()
	const errGroupPath = "golang.org/x/sync/errgroup"
	return (path == errGroupPath || strings.HasSuffix(path, "vendor/"+errGroupPath)) && fn.Name() == "Wait"
}

// getGoArg returns eg.Go args receiver and func-ish value
//
//nolint:ireturn // interface ssa.Value is ok to return.
func getGoArg(call ssa.CallCommon) (receiver, arg ssa.Value, err error) {
	if len(call.Args) != 2 {
		// eg.Go(f) <- Args should be [receiver, f]
		return nil, nil, fmt.Errorf("expected 2 arguments, got %d", len(call.Args))
	}
	return call.Args[0], call.Args[1], nil
}

// getWaitReceiver returns eg.Wait receiver
//
//nolint:ireturn // interface ssa.Value is ok to return.
func getWaitReceiver(call ssa.CallCommon) (ssa.Value, error) {
	if len(call.Args) != 1 {
		// eg.Wait() <- Args should be [receiver]
		return nil, errors.New("unexpected: len(call.Args) != 1")
	}
	return call.Args[0], nil
}
