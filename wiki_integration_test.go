package main

// TestDocsWiki_Integration is the automated guard for the manual testing runbook.
// It runs the same analysis as `wikigraph analyze docs/` and asserts the key
// invariants from the TC-16 spot-check table in docs/testing-runbook.md.
//
// When docs/ changes (page added/removed, links added/removed):
//  1. Run `wikigraph analyze docs/` to get the new values.
//  2. Update the constants below to match.
//  3. Update the matching values in docs/testing-runbook.md
//     (header, prerequisites list, TC-16 table, TC-18, TC-19, TC-20, TC-12).
//
// The test is intentionally strict: a failing constant is a signal to update
// the runbook, not to loosen the assertion.
//
// Run standalone: go test -run TestDocsWiki_Integration
import (
	"testing"
)

// docsWikiExpected holds the expected values for the docs/ wiki.
// Update these whenever a page is added, removed, or significantly relinked.
const (
	docsWantPages   = 39
	docsWantEdges   = 201 // P[i][j] > 1e-10 entries in the transition matrix
	docsWantClasses = 1
	docsWantCentral = "analyze"            // slug with the highest stationary probability
	docsWantLowestPi = "adr-009-wiki-gen-make-vs-buy" // slug with the lowest stationary probability
)

func TestDocsWiki_Integration(t *testing.T) {
	exclude := makeExcludeMap([]string{"index", "log", "AGENTS"})
	kern, _, pages, sinks, err := buildKernel("docs/", false, exclude)
	if err != nil {
		t.Fatalf("buildKernel(docs/): %v — is the binary being run from the repo root?", err)
	}

	// Page count.
	if len(pages) != docsWantPages {
		t.Errorf("page count: got %d, want %d\n\tUpdate docsWantPages and docs/testing-runbook.md",
			len(pages), docsWantPages)
	}

	// Edge count (same threshold as cmd_analyze.go).
	n := len(pages)
	edgeCount := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if kern.P.At(i, j) > 1e-10 {
				edgeCount++
			}
		}
	}
	if edgeCount != docsWantEdges {
		t.Errorf("edge count: got %d, want %d\n\tUpdate docsWantEdges and docs/testing-runbook.md",
			edgeCount, docsWantEdges)
	}

	// No sink pages.
	if len(sinks) != 0 {
		t.Errorf("sinks: expected none, got %v", sinks)
	}

	// Stationary distribution.
	pi, err := kern.Stationary(1e-12, 5000)
	if err != nil {
		t.Fatalf("stationary distribution: %v", err)
	}

	// Most central page.
	maxIdx := 0
	for i := 1; i < n; i++ {
		if pi[i] > pi[maxIdx] {
			maxIdx = i
		}
	}
	if pages[maxIdx] != docsWantCentral {
		t.Errorf("most central: got %q (π=%.6f), want %q\n\tUpdate docsWantCentral and docs/testing-runbook.md",
			pages[maxIdx], pi[maxIdx], docsWantCentral)
	}

	// Lowest-π page (TC-18 spot-check).
	minIdx := 0
	for i := 1; i < n; i++ {
		if pi[i] < pi[minIdx] {
			minIdx = i
		}
	}
	if pages[minIdx] != docsWantLowestPi {
		t.Errorf("lowest-π page: got %q (π=%.6f), want %q\n\tUpdate docsWantLowestPi and docs/testing-runbook.md",
			pages[minIdx], pi[minIdx], docsWantLowestPi)
	}

	// Communicating classes.
	cd, err := kern.Classes(1e-10)
	if err != nil {
		t.Fatalf("class decomposition: %v", err)
	}
	if len(cd.SCCs) != docsWantClasses {
		t.Errorf("class count: got %d, want %d\n\tUpdate docsWantClasses and docs/testing-runbook.md",
			len(cd.SCCs), docsWantClasses)
	}
}
