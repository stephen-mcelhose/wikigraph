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
	kern, P, pages, _, err := buildKernel(tmpDir, false, exclude)
	if err != nil {
		t.Fatalf("buildKernel failed: %v", err)
	}

	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}

	n := len(pages)

	// 1. Union strategy test
	t.Run("Union Strategy", func(t *testing.T) {
		goals := []int{idx["d"]}
		res := selectUnion(kern, n, goals, 3)
		if len(res) != 3 {
			t.Fatalf("expected 3 nodes, got %d", len(res))
		}
		// Goal 'd' must be present
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

	// 3. Path strategy test
	t.Run("Path Strategy", func(t *testing.T) {
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

	// 4. Bottleneck strategy test
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
