package signature

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/gostaticanalysis/analysisutil"
)

func ErrIshIndices(fn *ssa.Function) []int {
	sig := fn.Signature
	retLen := sig.Results().Len()
	indices := make([]int, 0, retLen)
	for i := range retLen {
		v := sig.Results().At(i)
		if isErrorIsh(v) {
			indices = append(indices, i)
		}
	}
	return indices
}

func isErrorIsh(v *types.Var) bool {
	if analysisutil.ImplementsError(v.Type()) {
		return true
	}

	// higher order function
	sig, ok := v.Type().(*types.Signature)
	if !ok {
		return false
	}
	retLen := sig.Results().Len()
	for i := range retLen {
		w := sig.Results().At(i)
		if isErrorIsh(w) {
			return true
		}
	}
	return false
}
