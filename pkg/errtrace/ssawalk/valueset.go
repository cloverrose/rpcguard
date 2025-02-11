package ssawalk

import (
	"golang.org/x/tools/go/ssa"
)

// valueSet is a set of ssa.Values that can be used to track
// the values that have been visited during a traversal. This
// is used to prevent infinite recursion, and to prevent
// visiting the same value multiple times.
type valueSet map[ssa.Value]struct{}

// includes returns true if the value is in the set.
func (v valueSet) includes(sv ssa.Value) bool {
	if v == nil {
		return false
	}
	_, ok := v[sv]
	return ok
}

// add adds the value to the set.
func (v valueSet) add(value ssa.Value) {
	v[value] = struct{}{}
}
