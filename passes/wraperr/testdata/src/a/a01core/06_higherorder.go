package a01core

// This file contains higher order methods.
// Note: Fact is record even for higher order func if that calling such func will return error in the end.
// Note: The wraperr does not report higher order methods.

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

// ReturnOkClosure returns closure that returns wrap error.
// Return anonymous func.
func (app *App) ReturnOkClosure(_ context.Context, _ *connect.Request[Message]) func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnOkClosure:"okFunc"
	return func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
		return nil, connect.NewError(connect.CodeInternal, errors.New("ReturnOkClosure"))
	}
}

// CallReturnOkClosure calls and returns error directly.
func (app *App) CallReturnOkClosure(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnOkClosure:"okFunc"
	return app.ReturnOkClosure(ctx, req)(ctx, req)
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnBadClosure returns closure that returns unwrap error.
// Return anonymous func.
func (app *App) ReturnBadClosure(_ context.Context, _ *connect.Request[Message]) func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnBadClosure:"badFunc"
	return func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
		return nil, errors.New("ReturnBadClosure")
	}
}

// CallReturnBadClosure calls and returns error directly.
func (app *App) CallReturnBadClosure(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnBadClosure:"badFunc" ".*RPC method CallReturnBadClosure returns error.*"
	return app.ReturnBadClosure(ctx, req)(ctx, req) // want ".*RPC method CallReturnBadClosure returns error.*"
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnOKMethod returns member method that returns wrap error
func (app *App) ReturnOKMethod(_ context.Context, _ *connect.Request[Message]) func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnOKMethod:"okFunc"
	return app.ReturnWrapError
}

// CallReturnOKMethod calls and returns error directly.
func (app *App) CallReturnOKMethod(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnOKMethod:"okFunc"
	return app.ReturnOKMethod(ctx, req)(ctx, req)
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnBadMethod returns member method that returns unwrap error
func (app *App) ReturnBadMethod(_ context.Context, _ *connect.Request[Message]) func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnBadMethod:"badFunc"
	return app.ReturnUnwrapError
}

// CallReturnBadMethod calls and returns error directly.
func (app *App) CallReturnBadMethod(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnBadMethod:"badFunc" ".*RPC method CallReturnBadMethod returns error.*"
	return app.ReturnBadMethod(ctx, req)(ctx, req) // want ".*RPC method CallReturnBadMethod returns error.*"
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnOKFunc returns function that returns wrap error
func (app *App) ReturnOKFunc(_ context.Context, _ *connect.Request[Message]) func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnOKFunc:"okFunc"
	return returnWrapErrorFunc
}

// CallReturnOKFunc calls and returns error directly.
func (app *App) CallReturnOKFunc(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnOKFunc:"okFunc"
	return app.ReturnOKFunc(ctx, req)(ctx, req)
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnBadFunc returns function that returns unwrap error
func (app *App) ReturnBadFunc(_ context.Context, _ *connect.Request[Message]) func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnBadFunc:"badFunc"
	return returnUnwrapErrorFunc
}

// CallReturnBadFunc calls and returns error directly.
func (app *App) CallReturnBadFunc(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnBadFunc:"badFunc" ".*RPC method CallReturnBadFunc returns error.*"
	return app.ReturnBadFunc(ctx, req)(ctx, req) // want ".*RPC method CallReturnBadFunc returns error.*"
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnNilClosure returns nil closure.
// nil closure is considered as bad closure, then the wraperr records badFunc Fact.
func (app *App) ReturnNilClosure(_ context.Context, _ *connect.Request[Message]) func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnNilClosure:"badFunc"
	return nil
}
