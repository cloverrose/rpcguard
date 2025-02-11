package generics

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"golang.org/x/sync/errgroup"
)

type handler[T any] interface {
	Handle() error
}

type App[T any] struct {
	handler handler[T]
}

type Message struct {
	text string
}

// GoOKFuncs calls OK funcs in eg.Go and returns eg.Wait error
func (app *App[T]) GoOKFuncs(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoOKFuncs:"okFunc"
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return nil
	})
	eg.Go(func() error {
		return connect.NewError(connect.CodeInternal, errors.New("wrap err"))
	})
	eg.Go(app.ReturnWrapError)
	eg.Go(app.ReturnReturnWrapError())
	eg.Go(app.ReturnReturnReturnWrapError()())

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return connect.NewResponse(&Message{"GoOKFuncs"}), nil
}

func (app *App[T]) ReturnWrapError() error { // want ReturnWrapError:"okFunc"
	return connect.NewError(connect.CodeInternal, errors.New("ReturnWrapError"))
}

func (app *App[T]) ReturnReturnWrapError() func() error { // want ReturnReturnWrapError:"okFunc"
	return app.ReturnWrapError
}

func (app *App[T]) ReturnReturnReturnWrapError() func() func() error { // want ReturnReturnReturnWrapError:"okFunc"
	return app.ReturnReturnWrapError
}

//----------------------------------------------------------------------------------------------------------------------

// GoBadClosure calls bad closure in eg.Go and returns eg.Wait error
func (app *App[T]) GoBadClosure(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoBadClosure:"badFunc" ".*RPC method GoBadClosure returns error.*"
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return errors.New("unwrap err")
	})

	if err := eg.Wait(); err != nil {
		return nil, err // want ".*RPC method GoBadClosure returns error.*"
	}
	return connect.NewResponse(&Message{"GoBadClosure"}), nil
}

//----------------------------------------------------------------------------------------------------------------------

// GoBadMethod calls bad member method in eg.Go and returns eg.Wait error
func (app *App[T]) GoBadMethod(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoBadMethod:"badFunc" ".*RPC method GoBadMethod returns error.*"
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(app.ReturnUnwrapError)

	if err := eg.Wait(); err != nil {
		return nil, err // want ".*RPC method GoBadMethod returns error.*"
	}
	return connect.NewResponse(&Message{"GoBadMethod"}), nil
}

func (app *App[T]) ReturnUnwrapError() error { // want ReturnUnwrapError:"badFunc"
	return errors.New("ReturnUnwrapError")
}

//----------------------------------------------------------------------------------------------------------------------

// GoNestedBadMethod calls bad member method in eg.Go and returns eg.Wait error
func (app *App[T]) GoNestedBadMethod(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoNestedBadMethod:"badFunc" ".*RPC method GoNestedBadMethod returns error.*"
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(app.ReturnReturnUnwrapError())

	if err := eg.Wait(); err != nil {
		return nil, err // want ".*RPC method GoNestedBadMethod returns error.*"
	}
	return connect.NewResponse(&Message{"GoNestedBadMethod"}), nil
}

func (app *App[T]) ReturnReturnUnwrapError() func() error { // want ReturnReturnUnwrapError:"badFunc"
	return func() error {
		return errors.New("ReturnReturnUnwrapError")
	}
}

//----------------------------------------------------------------------------------------------------------------------

// GoNestedBadMethod2 calls bad member method in eg.Go and returns eg.Wait error
func (app *App[T]) GoNestedBadMethod2(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoNestedBadMethod2:"badFunc" ".*RPC method GoNestedBadMethod2 returns error.*"
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(app.ReturnReturnReturnUnwrapError()())

	if err := eg.Wait(); err != nil {
		return nil, err // want ".*RPC method GoNestedBadMethod2 returns error.*"
	}
	return connect.NewResponse(&Message{"GoNestedBadMethod2"}), nil
}

func (app *App[T]) ReturnReturnReturnUnwrapError() func() func() error { // want ReturnReturnReturnUnwrapError:"badFunc"
	return app.ReturnReturnUnwrapError
}

//----------------------------------------------------------------------------------------------------------------------

//func (app *App[T]) CallInterfaceMethod() error { //  CallInterfaceMethod:"badFunc" ".*RPC method CallInterfaceMethod returns error.*"
//	eg, _ := errgroup.WithContext(context.Background())
//	eg.Go(app.handler.Handle)
//	if err := eg.Wait(); err != nil {
//		return err //  ".*RPC method CallInterfaceMethod returns error.*"
//	}
//	return nil
//}
