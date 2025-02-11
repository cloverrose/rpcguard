package a10defer

import (
	"errors"

	"connectrpc.com/connect"
)

func DeferFuncOKNoName() (string, error) { // want DeferFuncOKNoName:"okFunc"
	defer func() {}()
	return "ok", nil
}

func DeferFuncOKNamed() (x string, err1 error) { // want DeferFuncOKNamed:"okFunc"
	defer func() {
		x = "set in defer"
	}()
	return "ok", nil
}

func DeferFuncBad() (x string, err error) { // want DeferFuncBad:"badFunc"
	defer func() {
		x = "set in defer"
		err = errors.New("set in defer")
	}()
	return "ok", nil
}

// DeferFuncBadConnectError has a defer block.
// err is assigned with connect error but it is treated as badFunc.
func DeferFuncBadConnectError() (x string, err error) { // want DeferFuncBadConnectError:"badFunc"
	defer func() {
		x = "set in defer"
		err = connect.NewError(connect.CodeInternal, errors.New("set in defer"))
	}()
	return "ok", nil
}
