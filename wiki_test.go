package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPagesFlatAndRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	// Create structure:
	// tmpDir/
	//   page1.md
	//   sub/
	//     page2.md
	//   .hidden/
	//     page3.md
	if err := os.WriteFile(filepath.Join(tmpDir, "page1.md"), []byte("[[page2]]"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "page2.md"), []byte("[[page1]]"), 0644); err != nil {
		t.Fatal(err)
	}
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.Mkdir(hiddenDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hiddenDir, "page3.md"), []byte("[[page1]]"), 0644); err != nil {
		t.Fatal(err)
	}

	exclude := makeExcludeMap(nil)

	// Flat scan: should only find page1
	pagesFlat, _, _, err := loadPages(tmpDir, false, exclude)
	if err != nil {
		t.Fatalf("flat loadPages failed: %v", err)
	}
	if len(pagesFlat) != 1 || pagesFlat[0] != "page1" {
		t.Errorf("expected [page1], got %v", pagesFlat)
	}

	// Recursive scan: should find page1 and page2, skipping .hidden
	pagesRec, _, _, err := loadPages(tmpDir, true, exclude)
	if err != nil {
		t.Fatalf("recursive loadPages failed: %v", err)
	}
	if len(pagesRec) != 2 || pagesRec[0] != "page1" || pagesRec[1] != "page2" {
		t.Errorf("expected [page1, page2], got %v", pagesRec)
	}
}

// --- Relative markdown link edges (issue #51) ---

func TestRelativeLinks_SiblingAndTraversalAndAbsoluteURL(t *testing.T) {
	tmpDir := t.TempDir()

	// Layout:
	// tmpDir/
	//   gate/03-recommend.md
	//   gate/01-discovery.md   links: sibling "03-recommend.md", absolute URL, ../shared/notes.md
	//   shared/notes.md
	gateDir := filepath.Join(tmpDir, "gate")
	sharedDir := filepath.Join(tmpDir, "shared")
	for _, d := range []string{gateDir, sharedDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}
	os.WriteFile(filepath.Join(gateDir, "03-recommend.md"), []byte("# Gate 03\n"), 0644)
	os.WriteFile(filepath.Join(sharedDir, "notes.md"), []byte("# Notes\n"), 0644)
	os.WriteFile(filepath.Join(gateDir, "01-discovery.md"), []byte(strings.Join([]string{
		"# Gate 01",
		"",
		"[Gate 03](03-recommend.md)",
		"[shared notes](../shared/notes.md)",
		"[external](https://example.com)",
		"",
	}, "\n")), 0644)

	exclude := makeExcludeMap(nil)
	kern, _, pages, _, err := buildKernelWithOpts(tmpDir, false, exclude, true)
	if err != nil {
		t.Fatalf("buildKernelWithOpts failed: %v", err)
	}
	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d (%v)", len(pages), pages)
	}

	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}

	src := idx["01-discovery"]
	if kern.P.At(src, idx["03-recommend"]) <= 0 {
		t.Errorf("expected edge 01-discovery -> 03-recommend (sibling relative link)")
	}
	if kern.P.At(src, idx["notes"]) <= 0 {
		t.Errorf("expected edge 01-discovery -> notes (../ relative traversal)")
	}
	if len(pages) != 3 {
		t.Errorf("absolute URL should not add a node; got pages=%v", pages)
	}
}

func TestRelativeLinks_RedundantLabelMatchesTarget(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "stakeholder-notes.md"), []byte("# Stakeholder Notes\n"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("[stakeholder-notes](stakeholder-notes.md)\n"), 0644)

	exclude := makeExcludeMap(nil) // explicit nil: default excludes would drop 'index'
	pages, idx, paths, err := loadPages(tmpDir, false, exclude)
	if err != nil {
		t.Fatal(err)
	}
	adj, _, err := buildAdjacencyWithOpts(pages, idx, paths, true, tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if adj.At(idx["index"], idx["stakeholder-notes"]) <= 0 {
		t.Errorf("expected edge index -> stakeholder-notes")
	}
}

func TestRelativeLinks_DisabledByDefault_NonBreaking(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.md"), []byte("[b](b.md)\n"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.md"), []byte("# B\n"), 0644)

	exclude := makeExcludeMap(nil)
	pages, idx, paths, err := loadPages(tmpDir, false, exclude)
	if err != nil {
		t.Fatal(err)
	}
	// relativeLinks=false: identical behaviour to pre-#51 buildAdjacency.
	adj, sinks, err := buildAdjacencyWithOpts(pages, idx, paths, false, tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if adj.At(idx["a"], idx["b"]) != 1.0 {
		t.Errorf("expected uniform sink teleport row for 'a' (markdown link ignored), got adj[a][b]=%v", adj.At(idx["a"], idx["b"]))
	}
	if len(sinks) != 2 {
		t.Errorf("expected both pages to be sinks without relative-link parsing, got sinks=%v", sinks)
	}
}

func TestRelativeLinks_WikilinksStillWorkWhenEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.md"), []byte("[[b]] and [[b|Display Name]]\n"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.md"), []byte("# B\n"), 0644)

	exclude := makeExcludeMap(nil)
	kern, _, pages, _, err := buildKernelWithOpts(tmpDir, false, exclude, true)
	if err != nil {
		t.Fatal(err)
	}
	idx := map[string]int{}
	for i, p := range pages {
		idx[p] = i
	}
	if kern.P.At(idx["a"], idx["b"]) <= 0 {
		t.Errorf("expected [[wikilink]] edge to still work with --relative-links enabled")
	}
}

func TestRelativeLinks_ImpliesRecursive(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(tmpDir, "root.md"), []byte("[child](sub/child.md)\n"), 0644)
	os.WriteFile(filepath.Join(subDir, "child.md"), []byte("# Child\n"), 0644)

	exclude := makeExcludeMap(nil)
	// recursive=false explicitly, but relativeLinks=true should force recursion.
	_, _, pages, _, err := buildKernelWithOpts(tmpDir, false, exclude, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 2 {
		t.Fatalf("expected relative-links mode to imply recursive traversal, got pages=%v", pages)
	}
}

func TestRelativeLinks_WarnsOnEscapeAboveRoot(t *testing.T) {
	tmpDir := t.TempDir()
	projectRoot := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectRoot, 0755)
	outsideDir := filepath.Join(tmpDir, "outside")
	os.MkdirAll(outsideDir, 0755)
	os.WriteFile(filepath.Join(outsideDir, "escaped.md"), []byte("# Escaped\n"), 0644)
	os.WriteFile(filepath.Join(projectRoot, "page.md"), []byte("[escaped](../outside/escaped.md)\n"), 0644)

	raw, err := os.ReadFile(filepath.Join(projectRoot, "page.md"))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	linked := resolveRelativeLinksTo(&buf, raw, filepath.Join(projectRoot, "page.md"), projectRoot, map[string]string{})
	if len(linked) != 0 {
		t.Errorf("expected no edges for a path escaping the project root, got %v", linked)
	}
	if !strings.Contains(buf.String(), "warning:") || !strings.Contains(buf.String(), "outside the wiki root") {
		t.Errorf("expected a warning about escaping the wiki root, got: %q", buf.String())
	}
}

func TestRelativeLinks_DanglingRelativeLinkDoesNotPanic(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.md"), []byte("[missing](does-not-exist.md)\n"), 0644)

	exclude := makeExcludeMap(nil)
	_, _, pages, _, err := buildKernelWithOpts(tmpDir, false, exclude, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %v", pages)
	}
}

func TestLoadPagesDuplicateSlugError(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "note.md"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "folder")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "note.md"), []byte("world"), 0644); err != nil {
		t.Fatal(err)
	}

	exclude := makeExcludeMap(nil)
	_, _, _, err := loadPages(tmpDir, true, exclude)
	if err == nil {
		t.Fatalf("expected error on duplicate slug, got nil")
	}
}

func TestBuildKernelSampleNestedVault(t *testing.T) {
	tmpDir := t.TempDir()

	// Build sample vault layout
	// tmpDir/
	//   root_notes/index-note.md
	//   projects/quantum/quantum-intro.md
	//   archive/2025/old-notes.md
	//   .obsidian/config-note.md
	dirs := []string{
		filepath.Join(tmpDir, "root_notes"),
		filepath.Join(tmpDir, "projects", "quantum"),
		filepath.Join(tmpDir, "archive", "2025"),
		filepath.Join(tmpDir, ".obsidian"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}

	os.WriteFile(filepath.Join(tmpDir, "root_notes", "index-note.md"), []byte("[[quantum-intro]] and [[old-notes]]"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "projects", "quantum", "quantum-intro.md"), []byte("[[index-note]] and [[old-notes]]"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "archive", "2025", "old-notes.md"), []byte("[[quantum-intro]]"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".obsidian", "config-note.md"), []byte("[[index-note]]"), 0644)

	exclude := makeExcludeMap(nil)

	// Flat scan on root -> no pages error
	_, _, _, _, err := buildKernel(tmpDir, false, exclude)
	if err == nil {
		t.Fatalf("expected error on flat scan with no root .md files, got nil")
	}

	// Recursive scan -> 3 pages found, .obsidian ignored
	kern, _, pages, _, err := buildKernel(tmpDir, true, exclude)
	if err != nil {
		t.Fatalf("recursive buildKernel failed: %v", err)
	}

	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d (%v)", len(pages), pages)
	}

	expectedPages := []string{"index-note", "old-notes", "quantum-intro"}
	for i, p := range pages {
		if p != expectedPages[i] {
			t.Errorf("expected page %d to be %s, got %s", i, expectedPages[i], p)
		}
	}

	pi, err := kern.Stationary(1e-12, 5000)
	if err != nil {
		t.Fatalf("stationary distribution failed: %v", err)
	}

	// quantum-intro has inbound links from both index-note and old-notes, so it should have highest pi
	if pi[2] <= pi[0] || pi[2] <= pi[1] {
		t.Errorf("expected quantum-intro to be most central, got pi=%v", pi)
	}
}
