package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestRawSCCs_TwoComponents(t *testing.T) {
	// A↔B  and  C↔D  (no cross edges)
	adj := mat.NewDense(4, 4, []float64{
		0, 1, 0, 0,
		1, 0, 0, 0,
		0, 0, 0, 1,
		0, 0, 1, 0,
	})
	sccs := rawSCCs(adj)
	if len(sccs) != 2 {
		t.Fatalf("got %d SCCs, want 2", len(sccs))
	}
	for _, s := range sccs {
		if len(s.Members) != 2 {
			t.Errorf("SCC size %d, want 2", len(s.Members))
		}
		if !s.Closed {
			t.Errorf("expected closed SCC, got open: %v", s.Members)
		}
	}
}

func TestRawEdgeCount(t *testing.T) {
	adj := mat.NewDense(3, 3, []float64{
		0, 1, 1,
		0, 0, 1,
		0, 0, 0,
	})
	if got := rawEdgeCount(adj); got != 3 {
		t.Errorf("rawEdgeCount=%d, want 3", got)
	}
}
