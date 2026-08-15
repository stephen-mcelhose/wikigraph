package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

// goalChainFixture builds:
//
//	A -> B -> C -> D
//	A -> X -> D
//	E -> A
func goalChainFixture(t *testing.T) (kern *catrace.Kernel, P *mat.Dense, pages []string, idx map[string]int, n int) {
	t.Helper()
	tmpDir := t.TempDir()
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
	kern, _, pages, _, err := buildKernel(tmpDir, false, makeExcludeMap(nil))
	if err != nil {
		t.Fatalf("buildKernel failed: %v", err)
	}
	idx = make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}
	return kern, kern.P, pages, idx, len(pages)
}

// mustEqualSlugSet asserts got indices match wantSlugs exactly (order-independent).
func mustEqualSlugSet(t *testing.T, pages []string, got []int, wantSlugs ...string) {
	t.Helper()
	gotSlugs := make([]string, len(got))
	for i, g := range got {
		gotSlugs[i] = pages[g]
	}
	sort.Strings(gotSlugs)
	want := append([]string(nil), wantSlugs...)
	sort.Strings(want)
	if !reflect.DeepEqual(gotSlugs, want) {
		t.Errorf("selected set = %v, want golden %v", gotSlugs, want)
	}
}

// TestSelectUnion_GoldenTop3NearGoalD locks the OR-neighborhood around d.
// Fixture MFPT ranking: c and x are the two closest non-goal pages.
func TestSelectUnion_GoldenTop3NearGoalD(t *testing.T) {
	kern, _, pages, idx, n := goalChainFixture(t)
	res := selectUnion(kern, n, []int{idx["d"]}, 3)
	mustEqualSlugSet(t, pages, res, "c", "d", "x")
}

// TestSelectIntersection_GoldenTop3ForGoalsAD locks AND-prerequisites for {a,d}.
func TestSelectIntersection_GoldenTop3ForGoalsAD(t *testing.T) {
	kern, _, pages, idx, n := goalChainFixture(t)
	res := selectIntersection(kern, n, []int{idx["a"], idx["d"]}, 3)
	mustEqualSlugSet(t, pages, res, "a", "c", "d")
}

// TestSelectPath_GoldenShortestRouteAXD: a→d prefers a→x→d over a→b→c→d.
func TestSelectPath_GoldenShortestRouteAXD(t *testing.T) {
	_, P, pages, idx, n := goalChainFixture(t)
	res := selectPath(P, n, []int{idx["a"], idx["d"]}, 3)
	mustEqualSlugSet(t, pages, res, "a", "d", "x")
}

// TestSelectPath_GoldenMultiGoalChainExpandsNeighborE: goals a,b,c,d fill 4 of
// top 5; neighbor expansion must pick e (only remaining bidirectional neighbor).
func TestSelectPath_GoldenMultiGoalChainExpandsNeighborE(t *testing.T) {
	_, P, pages, idx, n := goalChainFixture(t)
	res := selectPath(P, n, []int{idx["a"], idx["b"], idx["c"], idx["d"]}, 5)
	mustEqualSlugSet(t, pages, res, "a", "b", "c", "d", "e")
}

// TestSelectPath_GoldenUnreachablePairKeepsGoalsOnly: a↛e has no path; result
// is exactly the two seeded goals (no phantom intermediates).
func TestSelectPath_GoldenUnreachablePairKeepsGoalsOnly(t *testing.T) {
	_, P, pages, idx, n := goalChainFixture(t)
	res := selectPath(P, n, []int{idx["a"], idx["e"]}, 2)
	mustEqualSlugSet(t, pages, res, "a", "e")
}

// TestSelectBottleneck_GoldenTop3ForGoalsAD locks betweenness fill for {a,d}.
func TestSelectBottleneck_GoldenTop3ForGoalsAD(t *testing.T) {
	_, P, pages, idx, n := goalChainFixture(t)
	res, err := selectBottleneck(P, n, []int{idx["a"], idx["d"]}, 3)
	if err != nil {
		t.Fatalf("selectBottleneck failed: %v", err)
	}
	mustEqualSlugSet(t, pages, res, "a", "c", "d")
}

// TestSelectBottleneck_GoldenSingleGoalPairsAllSources covers bottleneckPairs
// when len(goals)==1 (every non-goal paired as source → goal).
func TestSelectBottleneck_GoldenSingleGoalPairsAllSources(t *testing.T) {
	_, P, pages, idx, n := goalChainFixture(t)
	res, err := selectBottleneck(P, n, []int{idx["d"]}, 3)
	if err != nil {
		t.Fatalf("selectBottleneck failed: %v", err)
	}
	mustEqualSlugSet(t, pages, res, "a", "c", "d")
}

// TestSelectBottleneck_GoldenStarChokepointExactSet is the automated form of
// runbook TC-07e: hub is the only B↔C gatekeeper; tributary a must be absent.
//
// Graph: A → hub → B, hub → C, B → hub, C → hub
// Goals: [B, C], top: 3 → golden {b, c, hub}.
func TestSelectBottleneck_GoldenStarChokepointExactSet(t *testing.T) {
	tmpDir := t.TempDir()
	files := map[string]string{
		"hub.md": "[[b]] [[c]]",
		"b.md":   "[[hub]]",
		"c.md":   "[[hub]]",
		"a.md":   "[[hub]]",
	}
	for fname, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, fname), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	kern, _, pages, _, err := buildKernel(tmpDir, false, makeExcludeMap(nil))
	if err != nil {
		t.Fatalf("buildKernel failed: %v", err)
	}
	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}

	res, err := selectBottleneck(kern.P, len(pages), []int{idx["b"], idx["c"]}, 3)
	if err != nil {
		t.Fatalf("selectBottleneck failed: %v", err)
	}
	mustEqualSlugSet(t, pages, res, "b", "c", "hub")
}
