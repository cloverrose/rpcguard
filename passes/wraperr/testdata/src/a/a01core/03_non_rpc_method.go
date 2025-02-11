package a01core

// This file contains methods those signatures don't match RPC method.
// Note: The wraperr specifically reports only on RPC methods.

import (
	"context"
	"errors"
)

// returnOnlyUnwrapError returns unwrap error.
// This method signature does not match RPC method, so no report.
func (app *App) returnOnlyUnwrapError(_ context.Context) error { // want returnOnlyUnwrapError:"badFunc"
	return errors.New("returnOnlyUnwrapError")
}
