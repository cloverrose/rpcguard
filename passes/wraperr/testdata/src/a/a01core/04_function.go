package a01core

// This file contains standalone functions without receivers.
// Note: The wraperr specifically reports only on methods.

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

// returnNilFunc returns nil
func returnNilFunc(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want returnNilFunc:"okFunc"
	return connect.NewResponse(&Message{"returnNilFunc"}), nil
}

// returnWrapErrorFunc returns connect.NewError
func returnWrapErrorFunc(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want returnWrapErrorFunc:"okFunc"
	return nil, connect.NewError(connect.CodeInternal, errors.New("returnWrapErrorFunc"))
}

// returnUnwrapErrorFunc returns unwrap error
// This method signature does not match RPC method, so no report.
func returnUnwrapErrorFunc(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want returnUnwrapErrorFunc:"badFunc"
	return nil, errors.New("returnUnwrapErrorFunc")
}
