package a04closure

// This file contains closure method call.

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

type App struct{}

type Message struct {
	text string
}

// ReturnOkClosureError defines closure and returns its result.
func (app *App) ReturnOkClosureError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnOkClosureError:"okFunc"
	okClosure := func() error {
		return connect.NewError(connect.CodeInternal, errors.New("ReturnOkClosureError"))
	}
	if err := okClosure(); err != nil {
		return nil, err
	}
	return connect.NewResponse(&Message{"ReturnOkClosureError"}), nil
}

// ReturnBadClosureError defines closure and returns its result.
func (app *App) ReturnBadClosureError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnBadClosureError:"badFunc" ".*RPC method ReturnBadClosureError returns error.*"
	badClosure := func() error {
		return errors.New("ReturnBadClosureError")
	}
	if err := badClosure(); err != nil {
		return nil, err // want ".*RPC method ReturnBadClosureError returns error.*"
	}
	return connect.NewResponse(&Message{"ReturnBadClosureError"}), nil
}

// ReturnConstClosureError defines const closure (no body) and returns its result.
func (app *App) ReturnConstClosureError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnConstClosureError:"badFunc" ".*RPC method ReturnConstClosureError returns error.*"
	var constFunc func() error
	if err := constFunc(); err != nil {
		return nil, err // want ".*RPC method ReturnConstClosureError returns error.*"
	}
	return connect.NewResponse(&Message{"ReturnConstClosureError"}), nil
}
