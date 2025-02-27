package callgraph

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"golang.org/x/tools/go/ssa"

	"github.com/cloverrose/rpcguard/pkg/errtrace/ssawalk"
	"github.com/cloverrose/rpcguard/pkg/graph"

	"github.com/cloverrose/rpcguard/passes/wraperr/rtn"
	"github.com/cloverrose/rpcguard/passes/wraperr/visitors/norm"
)

type CallGraph struct {
	indicesFunc func(fn *ssa.Function) []int
	order       []*ssa.Function
	data        map[*ssa.Function]*FuncInfo
}

func New(indicesFunc func(fn *ssa.Function) []int) *CallGraph {
	return &CallGraph{
		indicesFunc: indicesFunc,
		data:        make(map[*ssa.Function]*FuncInfo),
	}
}

func (cg *CallGraph) Scan(srcFunc *ssa.Function) error {
	indices := cg.indicesFunc(srcFunc)
	if len(indices) == 0 {
		return fmt.Errorf("no indices for srcFunc: %s", srcFunc)
	}
	cg.order = append(cg.order, srcFunc)
	info, err := scanFunc(srcFunc, indices)
	if err != nil {
		return err
	}

	cg.data[srcFunc] = info
	return nil
}

// scanFunc scans single function and returns all returnInfo.
func scanFunc(fn *ssa.Function, indices []int) (*FuncInfo, error) {
	data := make(map[*ssa.Return]*returnInfo)
	for _, val := range rtn.GetReturnsAt(fn, indices) {
		info, err := scanVal(val.Value, getTargetIndex(val.Value))
		if err != nil {
			return nil, err
		}
		data[val.Return] = info
	}
	return &FuncInfo{data: data}, nil
}

func getTargetIndex(val ssa.Value) []int {
	switch val := val.(type) {
	case *ssa.Extract:
		return []int{val.Index}
	default:
		return []int{0}
	}
}

// scanVal scans single return value and returns returnInfo.
func scanVal(val ssa.Value, indices []int) (*returnInfo, error) {
	plugin := &visitorPlugin{
		normalizeFunc: norm.NewNormalizeFunc(indices),
	}
	if err := ssawalk.Walk(ssawalk.NewDefaultVisitorWith(plugin.createOptions()...), val); err != nil {
		return nil, err
	}
	return &returnInfo{
		toFuncs:        plugin.toFuncs,
		isObviouslyBad: plugin.isBad,
	}, nil
}

type scanPluginFunc func(fn *ssa.Function, indicesFunc func(fn *ssa.Function) []int) (map[*ssa.Return]map[*ssa.Function][]*ssa.Function, error)

// ScanWithPlugin scans fn with plugin and updates call graph.
func (cg *CallGraph) ScanWithPlugin(plugin scanPluginFunc, fn *ssa.Function) error {
	if _, ok := cg.data[fn]; !ok {
		return fmt.Errorf("unexpected: not found for fn: %s", fn)
	}

	rtnToReplacements, err := plugin(fn, cg.indicesFunc)
	if err != nil {
		return err
	}

	for rt, replacement := range rtnToReplacements {
		if _, ok := cg.data[fn].data[rt]; !ok {
			return fmt.Errorf("unexpected: not found for rt: %s", rt)
		}
		if len(replacement) == 0 {
			continue
		}
		newToFuncs := make([]*ssa.Function, 0, len(cg.data[fn].data[rt].toFuncs))
		for _, toFunc := range cg.data[fn].data[rt].toFuncs {
			converts, ok := replacement[toFunc]
			if ok {
				newToFuncs = append(newToFuncs, converts...)
			} else {
				newToFuncs = append(newToFuncs, toFunc)
			}
		}
		cg.data[fn].data[rt].toFuncs = newToFuncs
	}
	return nil
}

// GetReturnInfo returns FuncInfo for the given f.
func (cg *CallGraph) GetReturnInfo(f *ssa.Function) *FuncInfo {
	info, ok := cg.data[f]
	if !ok {
		return nil
	}
	return info
}

// Convert converts CallGraph to Graph.
func (cg *CallGraph) Convert() *graph.Graph[*ssa.Function] {
	g := graph.NewGraph[*ssa.Function]()
	for srcFunc, info := range cg.data {
		for _, toFunc := range info.GetAllToFuncs() {
			g.AddEdge(srcFunc, toFunc)
		}
	}
	return g
}

func (cg *CallGraph) LogValue() slog.Value {
	coreData := make(map[string][]string, len(cg.data))
	for srcFunc, info := range cg.data {
		toFuncs := info.GetAllToFuncs()
		str := make([]string, 0, len(toFuncs))
		for _, toFunc := range toFuncs {
			str = append(str, toFunc.String())
		}
		coreData[srcFunc.String()] = str
	}
	jsonBytes, err := json.Marshal(coreData)
	if err != nil {
		return slog.Value{}
	}
	return slog.StringValue(string(jsonBytes))
}

// returnInfo holds information for single return value.
type returnInfo struct {
	toFuncs        []*ssa.Function // source of this return value's functions.
	isObviouslyBad bool            // if this return is obviously bad or not. E.g. return returns non function (alloc etc).
}

// FuncInfo holds information for single function.
// Single function has several returns, data holds mapping from return to returnInfo.
type FuncInfo struct {
	data map[*ssa.Return]*returnInfo
}

// GetReturns returns all returns of single function.
func (i *FuncInfo) GetReturns() []*ssa.Return {
	return slices.Collect(maps.Keys(i.data))
}

// GetAllToFuncs returns all functions that are called from this function.
func (i *FuncInfo) GetAllToFuncs() []*ssa.Function {
	ret := make([]*ssa.Function, 0, len(i.data))
	for _, info := range i.data {
		ret = append(ret, info.toFuncs...)
	}
	return ret
}

// IsObviouslyBad returns true if this function is obviously bad.
func (i *FuncInfo) IsObviouslyBad() bool {
	for _, info := range i.data {
		if info.isObviouslyBad {
			return true
		}
	}
	return false
}

// IsObviouslyOK returns true if this function is obviously ok.
func (i *FuncInfo) IsObviouslyOK() bool {
	return !i.IsObviouslyBad() && len(i.GetAllToFuncs()) == 0
}

// IsObviouslyBadReturn returns true if given rtn is obviously bad.
func (i *FuncInfo) IsObviouslyBadReturn(rtn *ssa.Return) bool {
	info, ok := i.data[rtn]
	if !ok {
		return false
	}
	return info.isObviouslyBad
}

// GetToFuncs returns toFuncs for the given rtn.
func (i *FuncInfo) GetToFuncs(rtn *ssa.Return) []*ssa.Function {
	info, ok := i.data[rtn]
	if !ok {
		return nil
	}
	return info.toFuncs
}
