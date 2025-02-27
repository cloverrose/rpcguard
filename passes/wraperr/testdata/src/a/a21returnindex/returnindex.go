package a21returnindex

import (
	"errors"

	"connectrpc.com/connect"
)

func ReturnErrorAt3() (string, string, error) { // want ReturnErrorAt3:"okFunc"
	err := ReturnErrorAt1("hello")
	if err != nil {
		return "", "", err
	}
	v, err2 := ReturnErrorAt2("hello")
	if err2 != nil {
		return "", "", err2
	}
	return "ok", v, nil
}

func ReturnErrorAt1[T any](_ T) error { // want ReturnErrorAt1:"okFunc"
	return connect.NewError(connect.CodeInternal, errors.New("error at 1"))
}

func ReturnErrorAt2[T any](_ T) (string, error) { // want ReturnErrorAt2:"okFunc"
	return "", connect.NewError(connect.CodeInternal, errors.New("error at 2"))
}

func BadReturnErrorAt3() (string, string, error) { // want BadReturnErrorAt3:"badFunc"
	err := BadReturnErrorAt1("hello")
	if err != nil {
		return "", "", err
	}
	v, err2 := BadReturnErrorAt2("hello")
	if err2 != nil {
		return "", "", err2
	}
	return "ok", v, nil
}

func BadReturnErrorAt1[T any](_ T) error { // want BadReturnErrorAt1:"badFunc"
	return errors.New("error at 1")
}

func BadReturnErrorAt2[T any](_ T) (string, error) { // want BadReturnErrorAt2:"badFunc"
	return "", errors.New("error at 2")
}
