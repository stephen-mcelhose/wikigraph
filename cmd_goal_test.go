package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestGoalStrategies(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test graph:
	// A -> B -> C -> D
	// A -> X -> D
	// E (isolated island -> A)
	files := map[string]string{
		"a.md": "[[b]] [[x]]",
		"b.md": "[[c]]",
		"c.md": "[[d]]",
		"d.md": "[[a]]",
		"x.md": "[[d]]",
		"e.md": "[[a]]",
	}

	for fname, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, fname), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	exclude := makeExcludeMap(nil)
	kern, _, pages, _, err := buildKernel(tmpDir, false, exclude)
	if err != nil {
		t.Fatalf("buildKernel failed: %v", err)
	}

	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}

	n := len(pages)
	P := kern.P

	// 1. Union strategy test
	t.Run("Union Strategy", func(t *testing.T) {
		goals := []int{idx["d"]}
		res := selectUnion(kern, n, goals, 3)
		if len(res) != 3 {
			t.Fatalf("expected 3 nodes, got %d", len(res))
		}
		hasD := false
		for _, r := range res {
			if r == idx["d"] {
				hasD = true
			}
		}
		if !hasD {
			t.Errorf("expected goal 'd' in union result")
		}
	})

	// 2. Intersection strategy test
	t.Run("Intersection Strategy", func(t *testing.T) {
		goals := []int{idx["a"], idx["d"]}
		res := selectIntersection(kern, n, goals, 3)
		if len(res) != 3 {
			t.Fatalf("expected 3 nodes, got %d", len(res))
		}
		hasA, hasD := false, false
		for _, r := range res {
			if r == idx["a"] {
				hasA = true
			}
			if r == idx["d"] {
				hasD = true
			}
		}
		if !hasA || !hasD {
			t.Errorf("expected goals 'a' and 'd' in intersection result")
		}
	})

	// 3. Path strategy test (Shortest Path)
	t.Run("Path Strategy - Shortest Route", func(t *testing.T) {
		// Path from 'a' to 'd' -> shorter path is a -> x -> d (2 hops vs 3 hops)
		goals := []int{idx["a"], idx["d"]}
		res := selectPath(P, n, goals, 3)
		sort.Ints(res)
		expected := []int{idx["a"], idx["d"], idx["x"]}
		sort.Ints(expected)

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("expected path [a, d, x] (%v), got %v", expected, res)
		}
	})

	// 4. Path strategy test (Multi-goal Chain & Neighbor Expansion)
	t.Run("Path Strategy - Multi-goal Chain & Expansion", func(t *testing.T) {
		// Sequence: a -> b -> c -> d
		goals := []int{idx["a"], idx["b"], idx["c"], idx["d"]}
		res := selectPath(P, n, goals, 5)
		if len(res) != 5 {
			t.Fatalf("expected 5 nodes with expansion, got %d", len(res))
		}

		// All goal path nodes must be included
		for _, g := range []string{"a", "b", "c", "d"} {
			found := false
			for _, r := range res {
				if r == idx[g] {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected node %s in multi-goal path expansion", g)
			}
		}
	})

	// 5. Path strategy test (Unreachable Goal Pair Fallback)
	t.Run("Path Strategy - Unreachable Pair", func(t *testing.T) {
		// e -> a exists, but a -> e does NOT exist.
		goals := []int{idx["a"], idx["e"]}
		res := selectPath(P, n, goals, 2)
		// Should include 'a' and 'e' without crashing
		if len(res) != 2 {
			t.Fatalf("expected 2 goal nodes, got %d", len(res))
		}
	})

	// 6. Bottleneck strategy test
	t.Run("Bottleneck Strategy", func(t *testing.T) {
		goals := []int{idx["a"], idx["d"]}
		res, err := selectBottleneck(P, n, goals, 3)
		if err != nil {
			t.Fatalf("selectBottleneck failed: %v", err)
		}
		if len(res) != 3 {
			t.Fatalf("expected 3 nodes, got %d", len(res))
		}
		hasA, hasD := false, false
		for _, r := range res {
			if r == idx["a"] {
				hasA = true
			}
			if r == idx["d"] {
				hasD = true
			}
		}
		if !hasA || !hasD {
			t.Errorf("expected goals 'a' and 'd' in bottleneck result")
		}
	})
}

// TestBottleneckSelectsPathNode verifies that selectBottleneck identifies genuine
// chokepoint pages, not arbitrary or disconnected nodes.
//
// Graph: A → hub → B, hub → C, B → hub, C → hub
// hub is the only page connecting B and C; any walk between them must pass through hub.
// A is a dead-end tributary (reaches hub but not on any B↔C path).
// Goals: [B, C], top: 3 → expected result: {B, C, hub}, NOT {B, C, A}.
func TestBottleneckSelectsPathNode(t *testing.T) {
	tmpDir := t.TempDir()

	files := map[string]string{
		"hub.md": "[[b]] [[c]]",
		"b.md":   "[[hub]]",
		"c.md":   "[[hub]]",
		"a.md":   "[[hub]]", // tributary: reaches hub but not on B↔C path
	}
	for fname, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, fname), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	exclude := makeExcludeMap(nil)
	kern2, _, pages, _, err := buildKernel(tmpDir, false, exclude)
	if err != nil {
		t.Fatalf("buildKernel failed: %v", err)
	}
	P := kern2.P

	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}
	n := len(pages)

	goals := []int{idx["b"], idx["c"]}
	res, err := selectBottleneck(P, n, goals, 3)
	if err != nil {
		t.Fatalf("selectBottleneck failed: %v", err)
	}
	if len(res) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(res))
	}

	hasHub := false
	hasA := false
	for _, r := range res {
		if r == idx["hub"] {
			hasHub = true
		}
		if r == idx["a"] {
			hasA = true
		}
	}

	if !hasHub {
		t.Errorf("expected 'hub' (the chokepoint) in bottleneck result, got indices %v (pages: %v)", res, pages)
	}
	if hasA {
		t.Errorf("'a' (tributary, not on B↔C path) should NOT be in bottleneck result")
	}
}
