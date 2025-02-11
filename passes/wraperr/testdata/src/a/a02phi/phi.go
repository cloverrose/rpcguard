package a02phi

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

type App struct{}

type Message struct {
	text string
}

// PhiAllWrapErrorIf returns error its sources are all wrap error
func (app *App) PhiAllWrapErrorIf(_ context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want PhiAllWrapErrorIf:"okFunc"
	var err error
	if req.Msg.text == "hello" {
		err = connect.NewError(connect.CodeInternal, errors.New("wrap error 1"))
	} else {
		err = connect.NewError(connect.CodeNotFound, errors.New("wrap error 2"))
	}
	return connect.NewResponse(&Message{"PhiAllWrapErrorIf"}), err
}

// PhiUnwrapErrorIf returns error its sources contain unwrap error
func (app *App) PhiUnwrapErrorIf(_ context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want PhiUnwrapErrorIf:"badFunc" ".*RPC method PhiUnwrapErrorIf returns error.*"
	var err error
	if req.Msg.text == "hello" {
		err = connect.NewError(connect.CodeInternal, errors.New("wrap error"))
	} else {
		err = errors.New("unwrap error")
	}
	return connect.NewResponse(&Message{"PhiUnwrapErrorIf"}), err // want ".*RPC method PhiUnwrapErrorIf returns error.*"
}

//----------------------------------------------------------------------------------------------------------------------

// PhiAllWrapErrorSwitch returns error its sources are all wrap error
func (app *App) PhiAllWrapErrorSwitch(_ context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want PhiAllWrapErrorSwitch:"okFunc"
	var err error
	switch req.Msg.text {
	case "hello":
		err = connect.NewError(connect.CodeInternal, errors.New("wrap error 1"))
	default:
		err = connect.NewError(connect.CodeNotFound, errors.New("wrap error 2"))
	}
	return connect.NewResponse(&Message{"PhiAllWrapErrorSwitch"}), err
}

// PhiUnwrapErrorSwitch returns error its sources contain unwrap error
func (app *App) PhiUnwrapErrorSwitch(_ context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want PhiUnwrapErrorSwitch:"badFunc" ".*RPC method PhiUnwrapErrorSwitch returns error.*"
	var err error
	switch req.Msg.text {
	case "hello":
		err = connect.NewError(connect.CodeInternal, errors.New("wrap error"))
	default:
		err = errors.New("unwrap error")
	}
	return connect.NewResponse(&Message{"PhiUnwrapErrorSwitch"}), err // want ".*RPC method PhiUnwrapErrorSwitch returns error.*"
}

//----------------------------------------------------------------------------------------------------------------------

// NestedPhiUnwrapError returns error its sources contain unwrap error
func (app *App) NestedPhiUnwrapError(_ context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want NestedPhiUnwrapError:"badFunc" ".*RPC method NestedPhiUnwrapError returns error.*"
	var err error
	if req.Msg.text == "hello" {
		err = connect.NewError(connect.CodeInternal, errors.New("wrap error"))
	} else {
		err = errors.New("unwrap error")
	}

	switch req.Msg.text {
	case "world":
		err = connect.NewError(connect.CodeInternal, errors.New("wrap error"))
	}
	return connect.NewResponse(&Message{"PhiUnwrapErrorIf"}), err // want ".*RPC method NestedPhiUnwrapError returns error.*"
}
