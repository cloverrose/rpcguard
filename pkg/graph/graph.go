package graph

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
)

// https://www.logarithmic.net/pfh/blog/01208083168

type Graph[T comparable] struct {
	Vertices []T
	Edges    map[T][]T
}

// NewGraph returns new graph.
func NewGraph[T comparable]() *Graph[T] {
	return &Graph[T]{
		Edges: make(map[T][]T),
	}
}

func (g *Graph[T]) AddEdge(start, end T) {
	if idx := slices.Index(g.Vertices, start); idx == -1 {
		g.Vertices = append(g.Vertices, start)
	}
	if idx := slices.Index(g.Vertices, end); idx == -1 {
		g.Vertices = append(g.Vertices, end)
	}
	if idx := slices.Index(g.Edges[start], end); idx == -1 {
		g.Edges[start] = append(g.Edges[start], end)
	}
}

func (g *Graph[T]) LogValue() slog.Value {
	coreData := make(map[string][]string)
	for v, wlist := range g.Edges {
		coreData2 := make([]string, 0, len(wlist))
		for _, w := range wlist {
			coreData2 = append(coreData2, fmt.Sprintf("%v", w))
		}
		coreData[fmt.Sprintf("%v", v)] = coreData2
	}
	jsonBytes, err := json.Marshal(coreData)
	if err != nil {
		return slog.Value{}
	}
	return slog.StringValue(string(jsonBytes))
}

func Decomposition[T comparable](g *Graph[T]) [][]T {
	index := make(map[T]int)
	lowLink := make(map[T]int)
	stack := make([]T, 0)
	onStack := make(map[T]bool)
	var currentIndex int

	var result [][]T
	var strongConnect func(v T)
	strongConnect = func(v T) {
		index[v] = currentIndex
		lowLink[v] = currentIndex
		currentIndex++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range g.Edges[v] {
			if _, ok := index[w]; !ok {
				strongConnect(w)
				lowLink[v] = min(lowLink[v], lowLink[w])
			} else if onStack[w] {
				lowLink[v] = min(lowLink[v], lowLink[w])
			}
		}

		if lowLink[v] == index[v] {
			var w T
			var scc []T
			for {
				w, stack = stack[len(stack)-1], stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}
			result = append(result, scc)
		}
	}

	for _, v := range g.Vertices {
		if _, ok := index[v]; !ok {
			strongConnect(v)
		}
	}

	return result
}
