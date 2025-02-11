package callvalidate

import (
	"go/token"
	"strings"

	"golang.org/x/tools/go/ssa"

	"github.com/gostaticanalysis/analysisutil"

	"github.com/cloverrose/rpcguard/pkg/errtrace/ssawalk"
)

// checkCallValidate checks if func f calls Validate method and return error when Validate method returns error.
func checkCallValidate(f *ssa.Function, validateMethods []Method) (bool, error) {
	for _, block := range f.Blocks {
		if len(block.Instrs) == 0 {
			continue
		}
		// Because If instruction is the last instruction of its containing BasicBlock.
		instr := block.Instrs[len(block.Instrs)-1]

		// if ...
		ifInstr, ok := instr.(*ssa.If)
		if !ok {
			continue
		}
		// if validateErr != nil { ...
		validateErr, ok := isNilCheck(ifInstr.Cond)
		if !ok {
			continue
		}
		// if ... { return nil, err }
		if !isReturnErr(ifInstr.Block().Succs[0]) {
			continue
		}

		// validateErr := Validate()
		validateFn, err := scanVal(validateErr)
		if err != nil {
			return false, err
		}
		if validateFn == nil {
			continue
		}
		if isValidate(validateFn, validateMethods) {
			return true, nil
		}
		// nested case
		ok, err = checkCallValidate(validateFn, validateMethods)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// if err != nil { ... }
//
//nolint:ireturn // interface ssa.Value is ok to return.
func isNilCheck(cond ssa.Value) (ssa.Value, bool) {
	binop, ok := cond.(*ssa.BinOp)
	if !ok {
		return nil, false
	}
	if binop.Op != token.NEQ {
		// "if X Op Y {" exists but Op is not !=. We want to use "if X != Y {"
		return nil, false
	}
	c, ok := binop.Y.(*ssa.Const)
	if !ok {
		// "if X != Y {" exists but Y is not const. We want to use "if X != nil {"
		return nil, false
	}
	if !c.IsNil() {
		// "if X != Y {" exists but Y is not nil. We want to use "if X != nil {"
		return nil, false
	}
	if !analysisutil.ImplementsError(c.Type()) {
		// "if X != Y {" exists but X and Y's type are not error.  We want to use "if X:error != nil {"
		return nil, false
	}
	return binop.X, true
}

// if ... { return nil, err }
func isReturnErr(block *ssa.BasicBlock) bool {
	for _, instr := range block.Instrs {
		rt, ok := instr.(*ssa.Return)
		if !ok {
			continue
		}
		last := rt.Results[len(rt.Results)-1]
		if !analysisutil.ImplementsError(last.Type()) {
			continue
		}
		if _, ok := last.(*ssa.Const); ok {
			// 型はerrorだがConst(nil)を返している
			continue
		}
		return true
	}
	return false
}

// isValidate returns true if fn is Validate()
func isValidate(fn *ssa.Function, validateMethods []Method) bool {
	if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg == nil {
		return false
	}
	path := fn.Pkg.Pkg.Path()

	for _, method := range validateMethods {
		if (path == method.packagePath || strings.HasSuffix(path, "vendor/"+method.packagePath)) && fn.Name() == method.name {
			return true
		}
	}
	return false
}

// scanVal scans value.
func scanVal(val ssa.Value) (*ssa.Function, error) {
	plugin := &visitorPlugin{}
	if err := ssawalk.Walk(ssawalk.NewDefaultVisitorWith(plugin.createOptions()...), val); err != nil {
		return nil, err
	}
	return plugin.fn, nil
}
