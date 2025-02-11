package a09cyclic

import (
	"connectrpc.com/connect"
	"errors"
)

// ReachWrapError1 and ReachWrapError2 call each other. It eventually returns wrap error.
func ReachWrapError1(x int) error { // want ReachWrapError1:"okFunc"
	if x == 0 {
		return connect.NewError(connect.CodeInternal, errors.New("x is 0"))
	}
	return ReachWrapError2(x - 1)
}

func ReachWrapError2(x int) error { // want ReachWrapError2:"okFunc"
	if x == 0 {
		return connect.NewError(connect.CodeInternal, errors.New("x is 0"))
	}
	return ReachWrapError1(x - 1)
}

//----------------------------------------------------------------------------------------------------------------------

// ReachUnwrapError1 and ReachUnwrapError2 call each other. It eventually returns unwrap error.
func ReachUnwrapError1(x int) error { // want ReachUnwrapError1:"badFunc"
	if x == 0 {
		return errors.New("unwrap err")
	}
	return ReachUnwrapError2(x - 1)
}

func ReachUnwrapError2(x int) error { // want ReachUnwrapError2:"badFunc"
	if x == 0 {
		return connect.NewError(connect.CodeInternal, errors.New("wrap err"))
	}
	return ReachUnwrapError1(x - 1)
}

//----------------------------------------------------------------------------------------------------------------------

// Infinite1 and Infinite2 call each other. It makes infinite loop.
// Since both don't return unwrap error, the wraperr consider they are okFunc.
func Infinite1() error { // want Infinite1:"okFunc"
	return Infinite2()
}

func Infinite2() error { // want Infinite2:"okFunc"
	return Infinite1()
}

//----------------------------------------------------------------------------------------------------------------------

// SelfInfinite makes self infinite loop.
// Since SelfInfinite doesn't return unwrap error, the wraperr consider it is okFunc.
func SelfInfinite() error { // want SelfInfinite:"okFunc"
	return SelfInfinite()
}

//----------------------------------------------------------------------------------------------------------------------

// NoCyclicErrorRelation1 and NoCyclicErrorRelation2 call each other.
// But both wrap error with connect.NewError, then the wraperr doesn't put then inside the same SCC.
func NoCyclicErrorRelation1(x int) error { // want NoCyclicErrorRelation1:"okFunc"
	if err := NoCyclicErrorRelation2(x - 1); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	return nil
}

func NoCyclicErrorRelation2(x int) error { // want NoCyclicErrorRelation2:"okFunc"
	if err := NoCyclicErrorRelation1(x - 1); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	return nil
}
