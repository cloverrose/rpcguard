package a

import (
	"context"
	"errors"
	"fmt"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type App struct{}

type Message struct {
	text string
}

func (m Message) ProtoReflect() protoreflect.Message {
	panic("implement me")
}

func (app *App) CallValidateAndReturnError(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // OK
	if err := protovalidate.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&Message{"hello"}), nil
}

func (app *App) CallValidateButReturnOtherError(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // OK
	if err := protovalidate.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("validation error"))
	}
	return connect.NewResponse(&Message{"hello"}), nil
}

func (app *App) NoValidate(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want `RPC method NoValidate does not use validate method properly, accepted validate methods are buf.build/go/protovalidate:Validate,a:customValidate`
	return connect.NewResponse(&Message{"hello"}), nil
}

func (app *App) CallValidateButUseDifferently(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want `RPC method CallValidateButUseDifferently does not use validate method properly, accepted validate methods are buf.build/go/protovalidate:Validate,a:customValidate`
	err := protovalidate.Validate(req.Msg)
	if err == nil {
		return connect.NewResponse(&Message{"hello"}), nil
	}
	return nil, err
}

func (app *App) CallValidateButIgnore(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want `RPC method CallValidateButIgnore does not use validate method properly, accepted validate methods are buf.build/go/protovalidate:Validate,a:customValidate`
	validateErr := protovalidate.Validate(req.Msg)
	fmt.Println(validateErr)
	return connect.NewResponse(&Message{"hello"}), nil
}

func (app *App) CallValidateInNestedFuncAndReturnError(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // OK
	if err := app.validate(req); err != nil {
		return nil, err
	}
	return connect.NewResponse(&Message{"hello"}), nil
}

func (app *App) validate(req *connect.Request[Message]) error {
	if err := protovalidate.Validate(req.Msg); err != nil {
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("validation error"))
	}
	return nil
}

func (app *App) CallValidateInVerboseWrapper(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want `RPC method CallValidateInVerboseWrapper does not use validate method properly, accepted validate methods are buf.build/go/protovalidate:Validate,a:customValidate`
	if err := app.verboseValidate(req); err != nil {
		return nil, err
	}
	return connect.NewResponse(&Message{"hello"}), nil
}

func (app *App) verboseValidate(req *connect.Request[Message]) error {
	// verboseValidate doesn't have if err != nil { return ..., err } structure.
	// verboseValidate uses protovalidate.Validate properly, but verboseValidate is just a verbose.
	// so linter consider this is issue.
	return protovalidate.Validate(req.Msg)
}

func (app *App) CallValidateInNestedFuncWrongly(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want `RPC method CallValidateInNestedFuncWrongly does not use validate method properly, accepted validate methods are buf.build/go/protovalidate:Validate,a:customValidate`
	if err := app.badValidate(req); err != nil {
		return nil, err
	}
	return connect.NewResponse(&Message{"hello"}), nil
}

func (app *App) badValidate(req *connect.Request[Message]) error {
	if err := protovalidate.Validate(req.Msg); err != nil {
		return nil
	}
	return nil
}

func (app *App) CallValidateInMultiNestedFuncAndReturnError(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // OK
	if err := app.validateL1(req); err != nil {
		return nil, err
	}
	return connect.NewResponse(&Message{"hello"}), nil
}

func (app *App) validateL1(req *connect.Request[Message]) error {
	if err := app.validateL2(req); err != nil {
		return err
	}
	return nil
}

func (app *App) validateL2(req *connect.Request[Message]) error {
	if err := app.validateL3(req); err != nil {
		return err
	}
	return nil
}

func (app *App) validateL3(req *connect.Request[Message]) error {
	if err := protovalidate.Validate(req.Msg); err != nil {
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("validation error"))
	}
	return nil
}

func (app *App) IfCompareString(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want `RPC method IfCompareString does not use validate method properly, accepted validate methods are buf.build/go/protovalidate:Validate,a:customValidate`
	var s *string
	if s != nil {
		return nil, errors.New("err")
	}
	return connect.NewResponse(&Message{"hello"}), nil
}

func (app *App) CallCustomValidateAndReturnError(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // OK
	if err := customValidate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&Message{"hello"}), nil
}

// This method can be used for valid validate.
func customValidate(msg *Message) error {
	return errors.New("err")
}
