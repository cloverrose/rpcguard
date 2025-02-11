package norm

import (
	"golang.org/x/tools/go/ssa"
)

func NewNormalizeFunc(indices []int) func(fn *ssa.Function) (*ssa.Function, error) {
	return func(fn *ssa.Function) (*ssa.Function, error) {
		return normalize(fn, indices)
	}
}

func normalize(fn *ssa.Function, indices []int) (*ssa.Function, error) {
	normed, err := unbind(fn)
	if err != nil {
		return nil, err
	}
	if normed == nil {
		return nil, nil
	}
	normed, err = uninstantiate(normed, indices)
	if err != nil {
		return nil, err
	}
	return normed, nil
}
