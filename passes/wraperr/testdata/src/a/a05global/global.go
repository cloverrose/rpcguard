package a05global

import (
	"context"

	"connectrpc.com/connect"
)

var (
	globalErr  error
	globalFunc func() error
)

type App struct{}

type Message struct {
	text string
}

// ReturnGlobalError returns globalErr
func (app *App) ReturnGlobalError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnGlobalError:"badFunc" ".*RPC method ReturnGlobalError returns error.*"
	return nil, globalErr // want ".*RPC method ReturnGlobalError returns error.*"
}

// CallGlobalFunc calls globalFunc and returns its error.
func (app *App) CallGlobalFunc(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want CallGlobalFunc:"badFunc" ".*RPC method CallGlobalFunc returns error.*"
	if err := globalFunc(); err != nil {
		return nil, err // want ".*RPC method CallGlobalFunc returns error.*"
	}
	return connect.NewResponse(&Message{"CallGlobalFunc"}), nil
}
