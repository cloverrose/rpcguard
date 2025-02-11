package rpcmethod

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"

	"github.com/gostaticanalysis/analysisutil"
)

type Checker struct {
	rpcTypes *rpcMethodTypes
	loader   *RPCTypesLoader
}

func BuildChecker(pass *analysis.Pass) *Checker {
	loader := &RPCTypesLoader{
		pass: pass,
	}
	rpcTypes, err := loader.load()
	if err != nil {
		return nil
	}
	return &Checker{
		rpcTypes: rpcTypes,
		loader:   loader,
	}
}

func (c *Checker) IsRPCMethod(fn *ssa.Function) bool {
	if len(fn.Params) != 3 {
		// params should be [receiver, ctx, req]
		return false
	}

	sig := fn.Signature
	if sig.Params().Len() != 2 || sig.Results().Len() != 2 {
		return false
	}
	if sig.Params().At(0).Type() != c.rpcTypes.ctxType {
		return false
	}
	if !c.loader.checkInnerType(sig.Params().At(1).Type(), c.rpcTypes.reqType) {
		return false
	}
	if !c.loader.checkInnerType(sig.Results().At(0).Type(), c.rpcTypes.resType) {
		return false
	}
	if !analysisutil.ImplementsError(sig.Results().At(1).Type()) {
		return false
	}
	return true
}

type rpcMethodTypes struct {
	ctxType types.Type
	reqType types.Type
	resType types.Type
}

type RPCTypesLoader struct {
	pass *analysis.Pass
}

// loadRPCMethodTypes loads types that are used in RPC method.
// If fail to load, this package doesn't have RPC methods. So we can skip this package.
func (l *RPCTypesLoader) load() (*rpcMethodTypes, error) {
	ctxType := l.getContextType()
	reqType := l.getInner(l.getRequestType())
	if reqType == nil {
		return nil, errors.New("no rpc method")
	}
	resType := l.getInner(l.getResponseType())
	if resType == nil {
		return nil, errors.New("no rpc method")
	}
	return &rpcMethodTypes{
		ctxType: ctxType,
		reqType: reqType,
		resType: resType,
	}, nil
}

func (l *RPCTypesLoader) getContextType() types.Type {
	return analysisutil.TypeOf(l.pass, "context", "Context")
}

func (l *RPCTypesLoader) getRequestType() types.Type {
	return analysisutil.TypeOf(l.pass, "connectrpc.com/connect", "*Request")
}

func (l *RPCTypesLoader) getResponseType() types.Type {
	return analysisutil.TypeOf(l.pass, "connectrpc.com/connect", "*Response")
}

// From *connect.Request[FooRequest], get connect.Request[T any]
func (l *RPCTypesLoader) getInner(tt types.Type) types.Type {
	ptr, ok := tt.(*types.Pointer)
	if !ok {
		return nil
	}

	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return nil
	}

	return named.Obj().Type()
}

func (l *RPCTypesLoader) checkInnerType(typ, wantType types.Type) bool {
	ityp := l.getInner(typ)
	if ityp == nil {
		return false
	}
	return ityp == wantType
}
