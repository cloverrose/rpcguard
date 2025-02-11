package eg

import (
	"fmt"

	"golang.org/x/tools/go/ssa"

	"github.com/cloverrose/rpcguard/pkg/errtrace/ssawalk"

	"github.com/cloverrose/rpcguard/passes/wraperr/rtn"
)

// This file contains errgroup support.
// Below sample code,
// Foo: eg.Wait() returns connect.NewError
// Bar: eg.Wait() returns normal error
// However, ssa can't connect Foo to Foo$1, instead ssa connects Foo to eg.Go.
// Thus without special cares, we need to mark Foo as bad func because we can't analyze further.
// This file provides that special cares.

/**
import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"golang.org/x/sync/errgroup"
)

func Foo(ctx context.Context, x int) (int, error) {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { // <- in ssa, this anonymous func name is Foo$1
		return okFunc1(ctx, x)
	})
	if err := eg.Wait(); err != nil {
		return 0, err
	}
	return x, nil
}

func Bar(ctx context.Context, x int) (int, error) {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { // <- in ssa, this anonymous func name is Bar$1
		return badFunc1(ctx, x)
	})
	if err := eg.Wait(); err != nil {
		return 0, err
	}
	return x, nil
}

// okFunc1 use connect.NewError
func okFunc1(ctx context.Context, x int) error {
	if x == 0 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("err"))
	}
	return nil
}

// badFunc1 doesn't use connect.NewError
func badFunc1(ctx context.Context, x int) error {
	if x == 0 {
		return errors.New("err")
	}
	return nil
}
**/

func Scan(fn *ssa.Function, indicesFunc func(fn *ssa.Function) []int) (map[*ssa.Return]map[*ssa.Function][]*ssa.Function, error) {
	indices := indicesFunc(fn)
	if len(indices) == 0 {
		// srcFunc does not return error-ish values.
		return nil, fmt.Errorf("no indices for fn: %s", fn)
	}

	rtnToReplacements := make(map[*ssa.Return]map[*ssa.Function][]*ssa.Function)
	for _, val := range rtn.GetReturnsAt(fn, indices) {
		plugin := newVisitorPlugin()
		visitor := ssawalk.NewDefaultVisitorWith(plugin.createOptions()...)
		if err := ssawalk.Walk(visitor, val.Value); err != nil {
			return nil, err
		}
		rtnToReplacements[val.Return] = plugin.replacements
	}
	return rtnToReplacements, nil
}
