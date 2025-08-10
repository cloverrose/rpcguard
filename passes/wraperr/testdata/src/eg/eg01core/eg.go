package eg01core

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"golang.org/x/sync/errgroup"
)

type handler interface {
	Handle() error
}
type App struct {
	handler handler
}

type Message struct {
	text string
}

// GoOKFuncs calls OK funcs in eg.Go and returns eg.Wait error
func (app *App) GoOKFuncs(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoOKFuncs:"okFunc"
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

func (app *App) ReturnWrapError() error { // want ReturnWrapError:"okFunc"
	return connect.NewError(connect.CodeInternal, errors.New("ReturnWrapError"))
}

func (app *App) ReturnReturnWrapError() func() error { // want ReturnReturnWrapError:"okFunc"
	return app.ReturnWrapError
}

func (app *App) ReturnReturnReturnWrapError() func() func() error { // want ReturnReturnReturnWrapError:"okFunc"
	return app.ReturnReturnWrapError
}

//----------------------------------------------------------------------------------------------------------------------

// GoBadClosure calls bad closure in eg.Go and returns eg.Wait error
func (app *App) GoBadClosure(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoBadClosure:"badFunc" ".*RPC method GoBadClosure returns error.*"
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
func (app *App) GoBadMethod(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoBadMethod:"badFunc" ".*RPC method GoBadMethod returns error.*"
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(app.ReturnUnwrapError)

	if err := eg.Wait(); err != nil {
		return nil, err // want ".*RPC method GoBadMethod returns error.*"
	}
	return connect.NewResponse(&Message{"GoBadMethod"}), nil
}

func (app *App) ReturnUnwrapError() error { // want ReturnUnwrapError:"badFunc"
	return errors.New("ReturnUnwrapError")
}

//----------------------------------------------------------------------------------------------------------------------

// GoNestedBadMethod calls bad member method in eg.Go and returns eg.Wait error
func (app *App) GoNestedBadMethod(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoNestedBadMethod:"badFunc" ".*RPC method GoNestedBadMethod returns error.*"
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(app.ReturnReturnUnwrapError())

	if err := eg.Wait(); err != nil {
		return nil, err // want ".*RPC method GoNestedBadMethod returns error.*"
	}
	return connect.NewResponse(&Message{"GoNestedBadMethod"}), nil
}

func (app *App) ReturnReturnUnwrapError() func() error { // want ReturnReturnUnwrapError:"badFunc"
	return func() error {
		return errors.New("ReturnReturnUnwrapError")
	}
}

//----------------------------------------------------------------------------------------------------------------------

// GoNestedBadMethod2 calls bad member method in eg.Go and returns eg.Wait error
func (app *App) GoNestedBadMethod2(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want GoNestedBadMethod2:"badFunc" ".*RPC method GoNestedBadMethod2 returns error.*"
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(app.ReturnReturnReturnUnwrapError()())

	if err := eg.Wait(); err != nil {
		return nil, err // want ".*RPC method GoNestedBadMethod2 returns error.*"
	}
	return connect.NewResponse(&Message{"GoNestedBadMethod2"}), nil
}

func (app *App) ReturnReturnReturnUnwrapError() func() func() error { // want ReturnReturnReturnUnwrapError:"badFunc"
	return app.ReturnReturnUnwrapError
}

//----------------------------------------------------------------------------------------------------------------------

// NoGoCall returns eg.Wait error but there is no eg.Go call.
// TODO: this should be treated as okFunc.
func (app *App) NoGoCall(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want NoGoCall:"badFunc" ".*RPC method NoGoCall returns error.*"
	eg, _ := errgroup.WithContext(context.Background())
	if err := eg.Wait(); err != nil {
		return nil, err // want ".*RPC method NoGoCall returns error.*"
	}
	return connect.NewResponse(&Message{"NoGoCall"}), nil
}

//----------------------------------------------------------------------------------------------------------------------

// Infinite1 and Infinite2 call each other. It makes infinite loop.
// Since both don't return unwrap error, the wraperr consider they are okFunc.
func (app *App) Infinite1() error { // want Infinite1:"okFunc"
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(app.Infinite2)
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (app *App) Infinite2() error { // want Infinite2:"okFunc"
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(app.Infinite1)
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------

// MultiOKOK calls eg.Wait two times with different eg1 and eg2.
// Both eg1 and eg2 Wait's errors are wrap error.
func (app *App) MultiOKOK(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want MultiOKOK:"okFunc"
	eg1, _ := errgroup.WithContext(context.Background())
	eg1.Go(func() error {
		return connect.NewError(connect.CodeInternal, errors.New("wrap err"))
	})
	if err := eg1.Wait(); err != nil {
		return nil, err
	}

	eg2, _ := errgroup.WithContext(context.Background())
	eg2.Go(func() error {
		return connect.NewError(connect.CodeInternal, errors.New("wrap err"))
	})
	if err := eg2.Wait(); err != nil {
		return nil, err
	}

	return connect.NewResponse(&Message{"MultiOKOK"}), nil
}

// MultiOKBad calls eg.Wait two times with different eg1 and eg2.
// Only eg2.Wait's error is unwrap error.
func (app *App) MultiOKBad(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want MultiOKBad:"badFunc" ".*RPC method MultiOKBad returns error.*"
	eg1, _ := errgroup.WithContext(context.Background())
	eg1.Go(func() error {
		return connect.NewError(connect.CodeInternal, errors.New("wrap err"))
	})
	if err := eg1.Wait(); err != nil {
		return nil, err
	}

	eg2, _ := errgroup.WithContext(context.Background())
	eg2.Go(func() error {
		return errors.New("unwrap err")
	})
	if err := eg2.Wait(); err != nil {
		return nil, err // want ".*RPC method MultiOKBad returns error.*"
	}

	return connect.NewResponse(&Message{"MultiOKOK"}), nil
}

// MultiBadOK calls eg.Wait two times with different eg1 and eg2.
// Only eg1.Wait's error is unwrap error.
func (app *App) MultiBadOK(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want MultiBadOK:"badFunc" ".*RPC method MultiBadOK returns error.*"
	eg1, _ := errgroup.WithContext(context.Background())
	eg1.Go(func() error {
		return errors.New("unwrap err")
	})
	if err := eg1.Wait(); err != nil {
		return nil, err // want ".*RPC method MultiBadOK returns error.*"
	}

	eg2, _ := errgroup.WithContext(context.Background())
	eg2.Go(func() error {
		return connect.NewError(connect.CodeInternal, errors.New("wrap err"))
	})
	if err := eg2.Wait(); err != nil {
		return nil, err
	}

	return connect.NewResponse(&Message{"MultiBadOK"}), nil
}

// MultiBadBad calls eg.Wait two times with different eg1 and eg2.
// Both eg1 and eg2 Wait's errors are unwrap error.
func (app *App) MultiBadBad(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want MultiBadBad:"badFunc" ".*RPC method MultiBadBad returns error.*"
	eg1, _ := errgroup.WithContext(context.Background())
	eg1.Go(func() error {
		return errors.New("unwrap err")
	})
	if err := eg1.Wait(); err != nil {
		return nil, err // want ".*RPC method MultiBadBad returns error.*"
	}

	eg2, _ := errgroup.WithContext(context.Background())
	eg2.Go(func() error {
		return errors.New("unwrap err")
	})
	if err := eg2.Wait(); err != nil {
		return nil, err // want ".*RPC method MultiBadBad returns error.*"
	}

	return connect.NewResponse(&Message{"MultiBadBad"}), nil
}

//----------------------------------------------------------------------------------------------------------------------

// PhiEgErrors return err, this error has multiple eg.Wait sources.
func (app *App) PhiEgErrors(ctx context.Context, req *connect.Request[string]) (*connect.Response[Message], error) { // want PhiEgErrors:"badFunc" ".*RPC method PhiEgErrors returns error.*"
	eg1, _ := errgroup.WithContext(context.Background())
	eg1.Go(func() error {
		return connect.NewError(connect.CodeInternal, errors.New("wrap err"))
	})
	err1 := eg1.Wait()

	eg2, _ := errgroup.WithContext(context.Background())
	eg2.Go(func() error {
		return errors.New("unwrap error")
	})
	err2 := eg2.Wait()

	var err error
	if req.Msg == nil {
		err = err1
	} else {
		err = err2
	}
	return nil, err // want ".*RPC method PhiEgErrors returns error.*"
}
