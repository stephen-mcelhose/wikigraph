package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestRawSCCs_TwoClosedComponents(t *testing.T) {
	// A↔B  and  C↔D  (no cross edges) — both closed.
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

// TestRawSCCs_GoldenOpenComponentHasOutboundEdge: A↔B with B→C makes {A,B}
// open (edge leaves the component) and {C} a closed sink singleton.
func TestRawSCCs_GoldenOpenComponentHasOutboundEdge(t *testing.T) {
	// A↔B→C
	adj := mat.NewDense(3, 3, []float64{
		0, 1, 0, // A → B
		1, 0, 1, // B → A, B → C
		0, 0, 0, // C sink
	})
	sccs := rawSCCs(adj)
	if len(sccs) != 2 {
		t.Fatalf("got %d SCCs, want 2", len(sccs))
	}

	var open, closed int
	for _, s := range sccs {
		if s.Closed {
			closed++
			if len(s.Members) != 1 {
				t.Errorf("closed sink SCC size %d, want 1 (C)", len(s.Members))
			}
		} else {
			open++
			if len(s.Members) != 2 {
				t.Errorf("open SCC size %d, want 2 (A,B)", len(s.Members))
			}
		}
	}
	if open != 1 || closed != 1 {
		t.Errorf("want 1 open + 1 closed SCC, got open=%d closed=%d", open, closed)
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
