package includedpkg

import (
	"errors"
)

func OKFunc() error { // want OKFunc:"okFunc"
	return nil
}

func BadFunc() error { // want BadFunc:"badFunc"
	return errors.New("BadFunc")
}
