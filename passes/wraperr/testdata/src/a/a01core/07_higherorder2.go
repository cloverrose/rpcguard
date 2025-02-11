package a01core

// This file contains nested higher order methods.
// Note: Fact is record even for higher order func if that calling such func will return error in the end.
// Note: The wraperr does not report higher order methods.

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

// ReturnReturnOkClosure returns nested closure that returns wrap error.
// Return anonymous func.
func (app *App) ReturnReturnOkClosure(_ context.Context, _ *connect.Request[Message]) func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnReturnOkClosure:"okFunc"
	return func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
		return func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
			return nil, connect.NewError(connect.CodeInternal, errors.New("ReturnReturnOkClosure"))
		}
	}
}

// CallReturnReturnOkClosure calls and returns error directly.
func (app *App) CallReturnReturnOkClosure(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnReturnOkClosure:"okFunc"
	return app.ReturnReturnOkClosure(ctx, req)()(ctx, req)
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnReturnBadClosure returns nested closure that returns unwrap error.
// Return anonymous func.
func (app *App) ReturnReturnBadClosure(_ context.Context, _ *connect.Request[Message]) func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnReturnBadClosure:"badFunc"
	return func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
		return func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
			return nil, errors.New("ReturnReturnBadClosure")
		}
	}
}

// CallReturnReturnBadClosure calls and returns error directly.
func (app *App) CallReturnReturnBadClosure(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnReturnBadClosure:"badFunc" ".*RPC method CallReturnReturnBadClosure returns error.*"
	return app.ReturnReturnBadClosure(ctx, req)()(ctx, req) // want ".*RPC method CallReturnReturnBadClosure returns error.*"
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnReturnOKMethod returns closure that returns member method that returns wrap error
func (app *App) ReturnReturnOKMethod(_ context.Context, _ *connect.Request[Message]) func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnReturnOKMethod:"okFunc"
	return func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
		return app.ReturnWrapError
	}
}

// CallReturnReturnOKMethod calls and returns error directly.
func (app *App) CallReturnReturnOKMethod(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnReturnOKMethod:"okFunc"
	return app.ReturnReturnOKMethod(ctx, req)()(ctx, req)
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnReturnBadMethod returns closure that returns member method that returns unwrap error
func (app *App) ReturnReturnBadMethod(_ context.Context, _ *connect.Request[Message]) func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnReturnBadMethod:"badFunc"
	return func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
		return app.ReturnUnwrapError
	}
}

// CallReturnReturnBadMethod calls and returns error directly.
func (app *App) CallReturnReturnBadMethod(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnReturnBadMethod:"badFunc" ".*RPC method CallReturnReturnBadMethod returns error.*"
	return app.ReturnReturnBadMethod(ctx, req)()(ctx, req) // want ".*RPC method CallReturnReturnBadMethod returns error.*"
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnReturnOKFunc returns closure that returns function that returns wrap error
func (app *App) ReturnReturnOKFunc(_ context.Context, _ *connect.Request[Message]) func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnReturnOKFunc:"okFunc"
	return func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
		return returnWrapErrorFunc
	}
}

// CallReturnReturnOKFunc calls and returns error directly.
func (app *App) CallReturnReturnOKFunc(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnReturnOKFunc:"okFunc"
	return app.ReturnReturnOKFunc(ctx, req)()(ctx, req)
}

//----------------------------------------------------------------------------------------------------------------------

// ReturnReturnBadFunc returns closure that returns function that returns unwrap error
func (app *App) ReturnReturnBadFunc(_ context.Context, _ *connect.Request[Message]) func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnReturnBadFunc:"badFunc"
	return func() func(context.Context, *connect.Request[Message]) (*connect.Response[Message], error) {
		return returnUnwrapErrorFunc
	}
}

// CallReturnReturnBadFunc calls and returns error directly.
func (app *App) CallReturnReturnBadFunc(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnReturnBadFunc:"badFunc" ".*RPC method CallReturnReturnBadFunc returns error.*"
	return app.ReturnReturnBadFunc(ctx, req)()(ctx, req) // want ".*RPC method CallReturnReturnBadFunc returns error.*"
}
