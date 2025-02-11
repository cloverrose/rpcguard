package a

import (
	"context"

	"connectrpc.com/connect"
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

// CallInterfaceMethod calls interface's method and returns error directly.
func (app *App) CallInterfaceMethod(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want CallInterfaceMethod:"badFunc" ".*RPC method CallInterfaceMethod returns error.*"
	if err := app.handler.Handle(); err != nil {
		return nil, err // want ".*RPC method CallInterfaceMethod returns error.*"
	}
	return connect.NewResponse(&Message{"CallInterfaceMethod"}), nil
}

// CallInterfaceMethod2 calls interface's method and returns error directly.
func (app *App) CallInterfaceMethod2(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want CallInterfaceMethod2:"badFunc" ".*RPC method CallInterfaceMethod2 returns error.*"
	fn := app.handler.Handle

	if err := fn(); err != nil {
		return nil, err // want ".*RPC method CallInterfaceMethod2 returns error.*"
	}
	return connect.NewResponse(&Message{"CallInterfaceMethod2"}), nil
}
