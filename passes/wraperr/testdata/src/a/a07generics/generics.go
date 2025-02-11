package a07generics

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
)

type handler[T any] interface {
	Handle(in *T) error
}

type App[T any] struct {
	handler handler[T]
}

type Message struct {
	text string
}

// ReturnNil returns nil
func (app *App[T]) ReturnNil(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnNil:"okFunc"
	return connect.NewResponse(&Message{"ReturnNil"}), nil
}

// ReturnWrapError returns connect.NewError
func (app *App[T]) ReturnWrapError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnWrapError:"okFunc"
	return nil, connect.NewError(connect.CodeInternal, errors.New("ReturnWrapError"))
}

// ReturnUnwrapError returns unwrap error
func (app *App[T]) ReturnUnwrapError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want ReturnUnwrapError:"badFunc"  ".*RPC method ReturnUnwrapError returns error.*"
	return nil, errors.New("ReturnUnwrapError") // want ".*RPC method ReturnUnwrapError returns error.*"
}

//----------------------------------------------------------------------------------------------------------------------

// CallReturnUnwrapError returns unwrap error
func (app *App[T]) CallReturnUnwrapError(ctx context.Context, req *connect.Request[Message]) (*connect.Response[Message], error) { // want CallReturnUnwrapError:"badFunc"  ".*RPC method CallReturnUnwrapError returns error.*"
	return app.ReturnUnwrapError(ctx, req) // want ".*RPC method CallReturnUnwrapError returns error.*"
}

// CallGenericsInterfaceMethod returns unwrap error
func (app *App[T]) CallGenericsInterfaceMethod(ctx context.Context, req *connect.Request[T]) (*connect.Response[Message], error) { // want CallGenericsInterfaceMethod:"badFunc"  ".*RPC method CallGenericsInterfaceMethod returns error.*"
	if err := app.handler.Handle(req.Msg); err != nil {
		return nil, err // want ".*RPC method CallGenericsInterfaceMethod returns error.*"
	}
	return connect.NewResponse(&Message{"CallGenericsInterfaceMethod"}), nil
}

//----------------------------------------------------------------------------------------------------------------------

// CallComplexMethodError returns unwrap error
func (app *App[T]) CallComplexMethodError(_ context.Context, _ *connect.Request[Message]) (*connect.Response[Message], error) { // want CallComplexMethodError:"badFunc"  ".*RPC method CallComplexMethodError returns error.*"
	var m map[string]int
	_, err := Reversed(m)
	if err != nil {
		return nil, err // want ".*RPC method CallComplexMethodError returns error.*"
	}
	return connect.NewResponse(&Message{"CallComplexMethodError"}), nil
}

func Reversed[K, V comparable](original map[K]V) (map[V]K, error) { // want Reversed:"badFunc"
	reversed := make(map[V]K)
	for key, value := range original {
		if _, ok := reversed[value]; ok {
			return nil, fmt.Errorf("duplicate value found: %v", value)
		}
		reversed[value] = key
	}
	return reversed, nil
}
