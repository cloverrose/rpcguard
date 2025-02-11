package a01core

// This file contains methods which calls other methods or functions.

import (
	"context"

	"connectrpc.com/connect"
)

// CallOKRPCMethod calls member method that is okFunc and return error directly.
func (app *App) CallOKRPCMethod(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallOKRPCMethod:"okFunc"
	return app.ReturnWrapError(ctx, req)
}

// CallBadRPCMethod calls member method that is badFunc and return error directly.
func (app *App) CallBadRPCMethod(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallBadRPCMethod:"badFunc" ".*RPC method CallBadRPCMethod returns error.*"
	return app.ReturnUnwrapError(ctx, req) // want ".*RPC method CallBadRPCMethod returns error.*"
}

// CallBadNonRPCMethod calls member method that is badFunc
func (app *App) CallBadNonRPCMethod(ctx context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want CallBadNonRPCMethod:"badFunc" ".*RPC method CallBadNonRPCMethod returns error.*"
	if err := app.returnOnlyUnwrapError(ctx); err != nil {
		return nil, err // want ".*RPC method CallBadNonRPCMethod returns error.*"
	}
	return connect.NewResponse(&Message{"CallBadNonRPCMethod"}), nil
}

// CallBadFunc calls function that is badFunc and return error directly.
func (app *App) CallBadFunc(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallBadFunc:"badFunc" ".*RPC method CallBadFunc returns error.*"
	return returnUnwrapErrorFunc(ctx, req) // want ".*RPC method CallBadFunc returns error.*"
}

// CallBadMethodButWrap calls member method that is badFunc but wrap error
func (app *App) CallBadMethodButWrap(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallBadMethodButWrap:"okFunc"
	res, err := app.ReturnUnwrapError(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return res, nil
}

// CallNestedBadMethod calls member method that is badFunc
func (app *App) CallNestedBadMethod(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallNestedBadMethod:"badFunc" ".*RPC method CallNestedBadMethod returns error.*"
	return app.CallBadRPCMethod(ctx, req) // want ".*RPC method CallNestedBadMethod returns error.*"
}
