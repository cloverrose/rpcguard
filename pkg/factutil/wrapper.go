package factutil

import (
	"fmt"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

type FactWrapper[T analysis.Fact] struct {
	pass       *analysis.Pass
	localFacts map[*ssa.Function]T
}

func NewFactWrapper[T analysis.Fact](pass *analysis.Pass) *FactWrapper[T] {
	return &FactWrapper[T]{
		pass:       pass,
		localFacts: make(map[*ssa.Function]T),
	}
}

func (f *FactWrapper[T]) Export(fn *ssa.Function, fact T) {
	if fn == nil {
		panic("fn == nil")
	}
	if fn.Pkg == nil || fn.Pkg.Pkg == nil || fn.Pkg.Pkg != f.pass.Pkg || fn.Object() == nil {
		f.localFacts[fn] = fact
	} else {
		f.pass.ExportObjectFact(fn.Object(), fact)
	}
}

//nolint:ireturn // Need to return interface for generics.
func (f *FactWrapper[T]) Import(fn *ssa.Function) (T, bool) {
	if fn == nil {
		panic("fn == nil")
	}
	if lf, ok := f.localFacts[fn]; ok {
		return lf, true
	}
	if fn.Object() == nil {
		var t T
		return t, false
	}
	var fact T
	factType := reflect.TypeOf(fact).Elem()
	factValue, ok := reflect.New(factType).Interface().(T)
	if !ok {
		panic(fmt.Sprintf("fact type %s is not T", factType))
	}
	if f.pass.ImportObjectFact(fn.Object(), factValue) {
		return factValue, true
	}
	var t T
	return t, false
}
