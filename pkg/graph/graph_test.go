package graph

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecomposition(t *testing.T) {
	t.Parallel()
	type testCase[T comparable] struct {
		name string
		g    *Graph[T]
		want [][]T
	}
	tests := []testCase[int]{
		{
			name: "No SCC",
			g: func() *Graph[int] {
				g := NewGraph[int]()
				g.AddEdge(0, 1)
				return g
			}(),
			want: [][]int{{1}, {0}},
		},
		{
			name: "SCC",
			g: func() *Graph[int] {
				g := NewGraph[int]()
				g.AddEdge(0, 1)
				g.AddEdge(1, 0)
				return g
			}(),
			want: [][]int{{1, 0}},
		},
		{
			name: "2 SCC",
			g: func() *Graph[int] {
				g := NewGraph[int]()
				g.AddEdge(0, 1)
				g.AddEdge(1, 0)
				g.AddEdge(1, 2)
				return g
			}(),
			want: [][]int{{2}, {1, 0}},
		},
		{
			name: "2 SCC",
			g: func() *Graph[int] {
				g := NewGraph[int]()
				g.AddEdge(0, 1)
				g.AddEdge(1, 0)
				g.AddEdge(1, 2)
				g.AddEdge(2, 3)
				g.AddEdge(3, 2)
				return g
			}(),
			want: [][]int{{3, 2}, {1, 0}},
		},
		{
			name: "Complex", //nolint:usestdlibvars // this is not related to constant.Complex.
			g: func() *Graph[int] {
				g := NewGraph[int]()
				g.AddEdge(0, 2)
				g.AddEdge(0, 3)
				g.AddEdge(1, 5)
				g.AddEdge(2, 1)
				g.AddEdge(2, 3)
				g.AddEdge(3, 0)
				g.AddEdge(3, 1)
				g.AddEdge(3, 7)
				g.AddEdge(4, 6)
				g.AddEdge(4, 7)
				g.AddEdge(5, 1)
				g.AddEdge(5, 8)
				g.AddEdge(6, 9)
				g.AddEdge(7, 8)
				g.AddEdge(7, 9)
				g.AddEdge(8, 10)
				g.AddEdge(9, 4)
				return g
			}(),
			want: [][]int{{10}, {8}, {5, 1}, {6, 4, 9, 7}, {3, 2, 0}},
		},
		{
			name: "Complex 2",
			g: func() *Graph[int] {
				g := NewGraph[int]()
				g.AddEdge(7, 8)
				g.AddEdge(7, 9)
				g.AddEdge(8, 10)
				g.AddEdge(3, 1)
				g.AddEdge(3, 7)
				g.AddEdge(4, 6)
				g.AddEdge(4, 7)
				g.AddEdge(5, 1)
				g.AddEdge(5, 8)
				g.AddEdge(0, 2)
				g.AddEdge(0, 3)
				g.AddEdge(1, 5)
				g.AddEdge(2, 1)
				g.AddEdge(2, 3)
				g.AddEdge(3, 0)
				g.AddEdge(6, 9)
				g.AddEdge(9, 4)
				return g
			}(),
			want: [][]int{{10}, {8}, {6, 4, 9, 7}, {5, 1}, {2, 0, 3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Decomposition(tt.g)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Decomposition() diff (-want,+got) %s", diff)
			}
		})
	}
}
