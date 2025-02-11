package rtn

import (
	"golang.org/x/tools/go/ssa"
)

type ReturnAndValue struct {
	Return *ssa.Return
	Value  ssa.Value
}

// GetReturnsAt returns returns at index for the given fn.
// func Hello() (int, error) if indices = []{1}, returns error related returns.
func GetReturnsAt(fn *ssa.Function, indices []int) []ReturnAndValue {
	rtns := getReturns(fn)
	ats := make([]ReturnAndValue, 0, len(rtns))
	for _, rtn := range rtns {
		for _, index := range indices {
			ats = append(ats, ReturnAndValue{
				Return: rtn,
				Value:  rtn.Results[index],
			})
		}
	}
	return ats
}

// getReturns returns all returns for the given fn.
func getReturns(fn *ssa.Function) []*ssa.Return {
	var rtns []*ssa.Return
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if rtn, ok := instr.(*ssa.Return); ok {
				rtns = append(rtns, rtn)
			}
		}
	}
	return rtns
}
