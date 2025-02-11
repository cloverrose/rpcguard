package ssawalk

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

type Visitor interface {
	Visit(value ssa.Value) (Visitor, error)
}

func Walk(visitor Visitor, val ssa.Value) error {
	visited := make(valueSet)
	return walk(visitor, val, visited)
}

//nolint:gocognit,gocyclo,cyclop // walk visit pattern
func walk(visitor Visitor, val ssa.Value, visited valueSet) error {
	if visited.includes(val) {
		return nil
	}
	visited.add(val)

	visitor, err := visitor.Visit(val)
	if err != nil {
		return err
	}
	if visitor == nil {
		return nil
	}

	switch val := val.(type) {
	case *ssa.Function, *ssa.Const, *ssa.Alloc:
		panic(fmt.Sprintf("unexpected: visitor should return nil visitor when reach %T (terminal)", val))
	case *ssa.FreeVar, *ssa.Parameter, *ssa.Global, *ssa.MakeMap, *ssa.MakeChan, *ssa.MakeSlice,
		*ssa.Slice, *ssa.FieldAddr, *ssa.Field, *ssa.IndexAddr, *ssa.Index, *ssa.Lookup, *ssa.Select, *ssa.Range, *ssa.Next:
		panic(fmt.Sprintf("unexpected: visitor should return nil visitor when reach %T (too complex)", val))
	case *ssa.Builtin, *ssa.BinOp, *ssa.Convert, *ssa.MultiConvert:
		panic(fmt.Sprintf("unexpected: must not reach %T", val))
	case *ssa.Phi:
		for _, edge := range val.Edges {
			if err := walk(visitor, edge, visited); err != nil {
				return err
			}
		}
	case *ssa.Call:
		if val.Call.IsInvoke() {
			panic(fmt.Sprintf("unexpected: visitor should return nil visitor when reach invoke mode Call: %s", val.String()))
		}
		if err := walk(visitor, val.Call.Value, visited); err != nil {
			return err
		}
	case *ssa.Extract:
		if err := walk(visitor, val.Tuple, visited); err != nil {
			return err
		}
	case *ssa.MakeClosure:
		if err := walk(visitor, val.Fn, visited); err != nil {
			return err
		}
	case *ssa.UnOp:
		if err := walk(visitor, val.X, visited); err != nil {
			return err
		}
	case *ssa.ChangeType:
		if err := walk(visitor, val.X, visited); err != nil {
			return err
		}
	case *ssa.ChangeInterface:
		if err := walk(visitor, val.X, visited); err != nil {
			return err
		}
	case *ssa.SliceToArrayPointer:
		if err := walk(visitor, val.X, visited); err != nil {
			return err
		}
	case *ssa.MakeInterface:
		if err := walk(visitor, val.X, visited); err != nil {
			return err
		}
	case *ssa.TypeAssert:
		if err := walk(visitor, val.X, visited); err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("unknown type %T", val))
	}
	return nil
}
