package ssawalk

import (
	"fmt"
	"regexp"

	"golang.org/x/tools/go/ssa"
)

type options struct {
	visitFunction   func(val *ssa.Function) error
	visitConst      func(val *ssa.Const) error
	visitAlloc      func(val *ssa.Alloc) error
	visitComplex    func(val ssa.Value) error
	visitCall       func(val *ssa.Call) error
	visitCallInvoke func(val *ssa.Call) error
}

type Option interface {
	apply(opts *options)
}

type visitFunctionOption func(val *ssa.Function) error

func (f visitFunctionOption) apply(opts *options) {
	opts.visitFunction = f
}

//nolint:ireturn // for Uber option pattern.
func WithVisitFunction(f func(val *ssa.Function) error) Option {
	return visitFunctionOption(f)
}

type visitConstOption func(val *ssa.Const) error

func (f visitConstOption) apply(opts *options) {
	opts.visitConst = f
}

//nolint:ireturn // for Uber option pattern.
func WithVisitConst(f func(val *ssa.Const) error) Option {
	return visitConstOption(f)
}

type visitAllocOption func(val *ssa.Alloc) error

func (f visitAllocOption) apply(opts *options) {
	opts.visitAlloc = f
}

//nolint:ireturn // for Uber option pattern.
func WithVisitAlloc(f func(val *ssa.Alloc) error) Option {
	return visitAllocOption(f)
}

type visitComplexOption func(val ssa.Value) error

func (f visitComplexOption) apply(opts *options) {
	opts.visitComplex = f
}

//nolint:ireturn // for Uber option pattern.
func WithVisitComplex(f func(val ssa.Value) error) Option {
	return visitComplexOption(f)
}

type visitCallOption func(val *ssa.Call) error

func (f visitCallOption) apply(opts *options) {
	opts.visitCall = f
}

//nolint:ireturn // for Uber option pattern.
func WithVisitCall(f func(val *ssa.Call) error) Option {
	return visitCallOption(f)
}

type visitCallInvokeOption func(val *ssa.Call) error

func (f visitCallInvokeOption) apply(opts *options) {
	opts.visitCallInvoke = f
}

//nolint:ireturn // for Uber option pattern.
func WithVisitCallInvoke(f func(val *ssa.Call) error) Option {
	return visitCallInvokeOption(f)
}

type DefaultVisitor struct {
	opts *options
}

func NewDefaultVisitorWith(opts ...Option) *DefaultVisitor {
	op := &options{}
	for _, o := range opts {
		o.apply(op)
	}
	return &DefaultVisitor{
		opts: op,
	}
}

//nolint:ireturn,gocognit,cyclop // for Visitor pattern.
func (v DefaultVisitor) Visit(value ssa.Value) (Visitor, error) {
	switch value := value.(type) {
	case *ssa.Function:
		if v.opts.visitFunction != nil {
			if err := v.opts.visitFunction(value); err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ssa.Const:
		if v.opts.visitConst != nil {
			if err := v.opts.visitConst(value); err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ssa.Alloc:
		if isOKDefer(value) {
			return nil, nil
		}
		if v.opts.visitAlloc != nil {
			if err := v.opts.visitAlloc(value); err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ssa.FreeVar, *ssa.Parameter, *ssa.Global, *ssa.MakeMap, *ssa.MakeChan,
		*ssa.MakeSlice, *ssa.Slice, *ssa.FieldAddr, *ssa.Field, *ssa.IndexAddr,
		*ssa.Index, *ssa.Lookup, *ssa.Select, *ssa.Range, *ssa.Next:
		if v.opts.visitComplex != nil {
			if err := v.opts.visitComplex(value); err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ssa.Builtin, *ssa.BinOp, *ssa.Convert, *ssa.MultiConvert:
		return nil, fmt.Errorf("unexpected value: %s %T", value, value)
	case *ssa.Call:
		if value.Call.IsInvoke() {
			if v.opts.visitCallInvoke != nil {
				if err := v.opts.visitCallInvoke(value); err != nil {
					return nil, err
				}
			}
			return nil, nil
		}
		if v.opts.visitCall != nil {
			if err := v.opts.visitCall(value); err != nil {
				return nil, err
			}
			return nil, nil
		}
		return v, nil // without opts, call should be walked further.
	}
	return v, nil
}

// okDeferPattern are `local error ()` or `local error (err1)`
// err1 in paren is variable name.
// okDefer means func has a defer block and the defer block does not assign error.
/*
	func Foo() (x string, err1 error) {
		defer func () {
			x = "set in defer"
		}
		return "ok", nil
	}
*/
var okDeferPattern = regexp.MustCompile(`local error \(.*\)`)

func isOKDefer(value *ssa.Alloc) bool {
	str := value.String()
	return okDeferPattern.MatchString(str)
}
