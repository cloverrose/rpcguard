package a01core

// This file contains RPC methods.
// Note: The wraperr specifically reports only on RPC methods.

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

// ReturnNil returns nil
func (app *App) ReturnNil(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnNil:"okFunc"
	return connect.NewResponse(&Message{"ReturnNil"}), nil
}

// ReturnWrapError returns connect.NewError
func (app *App) ReturnWrapError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnWrapError:"okFunc"
	return nil, connect.NewError(connect.CodeInternal, errors.New("ReturnWrapError"))
}

// ReturnUnwrapError returns unwrap error
func (app *App) ReturnUnwrapError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnUnwrapError:"badFunc"  ".*RPC method ReturnUnwrapError returns error.*"
	return nil, errors.New("ReturnUnwrapError") // want ".*RPC method ReturnUnwrapError returns error.*"
}

//----------------------------------------------------------------------------------------------------------------------

type myError struct{}

func (e *myError) Error() string {
	return "myError"
}

// ReturnAllocError returns Alloc error without wrap
func (app *App) ReturnAllocError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnAllocError:"badFunc"  ".*RPC method ReturnAllocError returns error.*"
	return nil, &myError{} // want ".*RPC method ReturnAllocError returns error.*"
}
