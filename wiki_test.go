package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gonum.org/v1/gonum/mat"
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

	// Recursive scan: should find page1 and sub/page2 (path-relative slug), skipping .hidden
	pagesRec, _, _, err := loadPages(tmpDir, true, exclude)
	if err != nil {
		t.Fatalf("recursive loadPages failed: %v", err)
	}
	if len(pagesRec) != 2 || pagesRec[0] != "page1" || pagesRec[1] != "sub/page2" {
		t.Errorf("expected [page1, sub/page2], got %v", pagesRec)
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
	kern, _, pages, _, err := buildKernelWithOpts(tmpDir, false, exclude, true, defaultTeleportAlpha, nil)
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

	src := idx["gate/01-discovery"]
	if kern.P.At(src, idx["gate/03-recommend"]) <= 0 {
		t.Errorf("expected edge gate/01-discovery -> gate/03-recommend (sibling relative link)")
	}
	if kern.P.At(src, idx["shared/notes"]) <= 0 {
		t.Errorf("expected edge gate/01-discovery -> shared/notes (../ relative traversal)")
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
	// relativeLinks=false: markdown links ignored; pages with only md links are sinks.
	adj, sinks, err := buildAdjacencyWithOpts(pages, idx, paths, false, tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if adj.At(idx["a"], idx["b"]) != 0 {
		t.Errorf("expected no raw edge for markdown link when relativeLinks=false, got adj[a][b]=%v", adj.At(idx["a"], idx["b"]))
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
	kern, _, pages, _, err := buildKernelWithOpts(tmpDir, false, exclude, true, defaultTeleportAlpha, nil)
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
	_, _, pages, _, err := buildKernelWithOpts(tmpDir, false, exclude, true, defaultTeleportAlpha, nil)
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

// TestNormalizeMDLinkTarget_GoldenAcceptAndRejectTable locks title/fragment/
// query stripping and absolute/anchor rejection used by relative-link parsing.
func TestNormalizeMDLinkTarget_GoldenAcceptAndRejectTable(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		want   string
		wantOK bool
	}{
		{name: "plain relative", in: "02-feasibility.md", want: "02-feasibility.md", wantOK: true},
		{name: "strips title suffix", in: `foo.md "Title"`, want: "foo.md", wantOK: true},
		{name: "strips fragment", in: "foo.md#section", want: "foo.md", wantOK: true},
		{name: "strips query", in: "foo.md?x=1", want: "foo.md", wantOK: true},
		{name: "reject empty", in: "   ", want: "", wantOK: false},
		{name: "reject bare anchor", in: "#section", want: "", wantOK: false},
		{name: "reject http", in: "https://example.com/a.md", want: "", wantOK: false},
		{name: "reject mailto", in: "mailto:a@b.com", want: "", wantOK: false},
		{name: "reject protocol-relative", in: "//cdn.example/x.md", want: "", wantOK: false},
		{name: "fragment-only after strip empty", in: "#", want: "", wantOK: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := normalizeMDLinkTarget(tc.in)
			if ok != tc.wantOK || got != tc.want {
				t.Errorf("normalizeMDLinkTarget(%q) = (%q, %v), want (%q, %v)",
					tc.in, got, ok, tc.want, tc.wantOK)
			}
		})
	}
}

func TestRelativeLinks_DanglingRelativeLinkDoesNotPanic(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.md"), []byte("[missing](does-not-exist.md)\n"), 0644)

	exclude := makeExcludeMap(nil)
	_, _, pages, _, err := buildKernelWithOpts(tmpDir, false, exclude, true, defaultTeleportAlpha, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %v", pages)
	}
}

func writePortfolioFixture(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	alpha := filepath.Join(tmpDir, "project-alpha")
	beta := filepath.Join(tmpDir, "project-beta")
	for _, d := range []string{alpha, beta} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}
	// Each project: 01-discovery <-> 02-feasibility via relative markdown links.
	os.WriteFile(filepath.Join(alpha, "01-discovery.md"), []byte("# Alpha Discovery\n\n[Feasibility](02-feasibility.md)\n"), 0644)
	os.WriteFile(filepath.Join(alpha, "02-feasibility.md"), []byte("# Alpha Feasibility\n\n[Discovery](01-discovery.md)\n"), 0644)
	os.WriteFile(filepath.Join(beta, "01-discovery.md"), []byte("# Beta Discovery\n\n[Feasibility](02-feasibility.md)\n"), 0644)
	os.WriteFile(filepath.Join(beta, "02-feasibility.md"), []byte("# Beta Feasibility\n\n[Discovery](01-discovery.md)\n"), 0644)
	return tmpDir
}

func loadPortfolioRaw(t *testing.T) (rawAdj *mat.Dense, pages []string, idx map[string]int) {
	t.Helper()
	tmpDir := writePortfolioFixture(t)
	_, rawAdj, pages, _, err := buildKernelWithOpts(tmpDir, false, makeExcludeMap(nil), true, defaultTeleportAlpha, nil)
	if err != nil {
		t.Fatalf("buildKernelWithOpts: %v", err)
	}
	if len(pages) != 4 {
		t.Fatalf("expected 4 pages, got %d (%v)", len(pages), pages)
	}
	idx = make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}
	for _, slug := range []string{
		"project-alpha/01-discovery",
		"project-alpha/02-feasibility",
		"project-beta/01-discovery",
		"project-beta/02-feasibility",
	} {
		if _, ok := idx[slug]; !ok {
			t.Fatalf("slug %q missing from pages %v", slug, pages)
		}
	}
	return rawAdj, pages, idx
}

func portfolioCrossProjectPairs() [][2]string {
	return [][2]string{
		{"project-alpha/01-discovery", "project-beta/01-discovery"},
		{"project-alpha/01-discovery", "project-beta/02-feasibility"},
		{"project-alpha/02-feasibility", "project-beta/01-discovery"},
		{"project-alpha/02-feasibility", "project-beta/02-feasibility"},
		{"project-beta/01-discovery", "project-alpha/01-discovery"},
		{"project-beta/01-discovery", "project-alpha/02-feasibility"},
		{"project-beta/02-feasibility", "project-alpha/01-discovery"},
		{"project-beta/02-feasibility", "project-alpha/02-feasibility"},
	}
}

func portfolioIntraProjectPairs() [][2]string {
	return [][2]string{
		{"project-alpha/01-discovery", "project-alpha/02-feasibility"},
		{"project-alpha/02-feasibility", "project-alpha/01-discovery"},
		{"project-beta/01-discovery", "project-beta/02-feasibility"},
		{"project-beta/02-feasibility", "project-beta/01-discovery"},
	}
}

// TC-27(a/b): portfolio wiki where every project subdirectory contains
// identically-named files. --relative-links must resolve sibling links within
// each project only. The plain wikilink path ([[02-feasibility]]) is ambiguous
// and covered separately in TestLenientWikilinkFallback_AmbiguousBasenameDropsLink.

func TestRelativeLinks_Portfolio_NoCrossProjectEdges(t *testing.T) {
	rawAdj, _, idx := loadPortfolioRaw(t)
	for _, e := range portfolioCrossProjectPairs() {
		if rawAdj.At(idx[e[0]], idx[e[1]]) > 0 {
			t.Errorf("cross-project edge %s -> %s must not exist (rawAdj=%.6f)",
				e[0], e[1], rawAdj.At(idx[e[0]], idx[e[1]]))
		}
	}
}

func TestRelativeLinks_Portfolio_IntraProjectEdges(t *testing.T) {
	rawAdj, _, idx := loadPortfolioRaw(t)
	for _, e := range portfolioIntraProjectPairs() {
		if rawAdj.At(idx[e[0]], idx[e[1]]) <= 0 {
			t.Errorf("intra-project edge %s -> %s must exist (rawAdj=%.6f)",
				e[0], e[1], rawAdj.At(idx[e[0]], idx[e[1]]))
		}
	}
}

func TestRelativeLinks_Portfolio_TwoRawSCCs(t *testing.T) {
	rawAdj, pages, _ := loadPortfolioRaw(t)
	sccs := rawSCCs(rawAdj)
	if len(sccs) != 2 {
		t.Errorf("expected 2 raw communicating classes (one per project), got %d", len(sccs))
	}
	for _, scc := range sccs {
		if len(scc.Members) != 2 {
			slugs := make([]string, len(scc.Members))
			for i, s := range scc.Members {
				slugs[i] = pages[s]
			}
			t.Errorf("expected each class to have 2 pages, got %d: %v", len(scc.Members), slugs)
		}
	}
}

func TestLoadPagesPathRelativeSlugs(t *testing.T) {
	// Two files with the same basename in different subdirectories must not
	// collide: recursive mode uses path-relative slugs.
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
	pages, _, _, err := loadPages(tmpDir, true, exclude)
	if err != nil {
		t.Fatalf("unexpected error with path-relative slugs: %v", err)
	}
	// Sorted: "folder/note" < "note"
	if len(pages) != 2 || pages[0] != "folder/note" || pages[1] != "note" {
		t.Errorf("expected [folder/note, note], got %v", pages)
	}
}

func TestIsExcluded(t *testing.T) {
	exclude := makeExcludeMap([]string{"index", "log"})

	// Exact match on bare stem (flat-mode slug).
	if !isExcluded("index", exclude) {
		t.Error("expected 'index' to be excluded (exact match)")
	}
	// Basename match on path-relative slug (recursive-mode slug).
	if !isExcluded("subdir/index", exclude) {
		t.Error("expected 'subdir/index' to be excluded (basename match)")
	}
	if !isExcluded("a/b/c/log", exclude) {
		t.Error("expected 'a/b/c/log' to be excluded (deep basename match)")
	}
	// Non-excluded slug.
	if isExcluded("my-page", exclude) {
		t.Error("expected 'my-page' not to be excluded")
	}
	if isExcluded("subdir/my-page", exclude) {
		t.Error("expected 'subdir/my-page' not to be excluded")
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

	// Path-relative slugs, sorted lexicographically.
	expectedPages := []string{"archive/2025/old-notes", "projects/quantum/quantum-intro", "root_notes/index-note"}
	for i, p := range pages {
		if p != expectedPages[i] {
			t.Errorf("expected page %d to be %s, got %s", i, expectedPages[i], p)
		}
	}

	pi, err := kern.Stationary(1e-12, 5000)
	if err != nil {
		t.Fatalf("stationary distribution failed: %v", err)
	}

	// quantum-intro (index 1) has inbound links from both index-note and old-notes, so it should have highest pi.
	// [[wikilinks]] in fixtures resolve via the lenient basename fallback.
	if pi[1] <= pi[0] || pi[1] <= pi[2] {
		t.Errorf("expected projects/quantum/quantum-intro to be most central, got pi=%v", pi)
	}
}

func TestWikilinkRe_DigitLeadingSlug(t *testing.T) {
	// wikilinkRe must match slugs that start with a digit, e.g. [[02-feasibility]].
	cases := []struct {
		input string
		want  string // empty string means no match expected
	}{
		{"[[02-feasibility]]", "02-feasibility"},
		{"[[1-intro]]", "1-intro"},
		{"[[my-note]]", "my-note"},                 // letter-leading still works
		{"[[note2]]", "note2"},                     // digit in body still works
		{"[[02-feasibility|Alias]]", "02-feasibility"}, // alias form
	}
	for _, tc := range cases {
		m := wikilinkRe.FindStringSubmatch(tc.input)
		if tc.want == "" {
			if m != nil {
				t.Errorf("input %q: expected no match, got %v", tc.input, m)
			}
			continue
		}
		if m == nil {
			t.Errorf("input %q: expected match %q, got no match", tc.input, tc.want)
			continue
		}
		if m[1] != tc.want {
			t.Errorf("input %q: got group[1]=%q, want %q", tc.input, m[1], tc.want)
		}
	}
}

func TestLenientWikilinkFallback_UniqueBasename(t *testing.T) {
	// [[02-note]] (digit-leading) in a recursive vault where only one file
	// matches the basename should resolve via the lenient fallback.
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	// root.md links to [[02-note]], which lives at sub/02-note.md (slug: sub/02-note)
	os.WriteFile(filepath.Join(tmpDir, "root.md"), []byte("[[02-note]]"), 0644)
	os.WriteFile(filepath.Join(subDir, "02-note.md"), []byte("# Note"), 0644)

	exclude := makeExcludeMap(nil)
	kern, _, pages, _, err := buildKernel(tmpDir, true, exclude)
	if err != nil {
		t.Fatalf("buildKernel failed: %v", err)
	}
	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}
	if kern.P.At(idx["root"], idx["sub/02-note"]) <= 0 {
		t.Errorf("expected lenient fallback to resolve [[02-note]] -> sub/02-note; P=%v", kern.P)
	}
}

func TestLenientWikilinkFallback_AmbiguousBasenameDropsLink(t *testing.T) {
	// [[note]] in a recursive vault where two files share the basename
	// must be dropped (not mis-resolved) and a warning printed to stderr.
	tmpDir := t.TempDir()
	aDir := filepath.Join(tmpDir, "a")
	bDir := filepath.Join(tmpDir, "b")
	for _, d := range []string{aDir, bDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}
	// root links to [[note]]; both a/note.md and b/note.md exist — ambiguous.
	os.WriteFile(filepath.Join(tmpDir, "root.md"), []byte("[[note]]"), 0644)
	os.WriteFile(filepath.Join(aDir, "note.md"), []byte("# A"), 0644)
	os.WriteFile(filepath.Join(bDir, "note.md"), []byte("# B"), 0644)

	exclude := makeExcludeMap(nil)
	kern, _, pages, sinks, err := buildKernel(tmpDir, true, exclude)
	if err != nil {
		t.Fatalf("buildKernel failed: %v", err)
	}
	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}
	// root should be a sink (no resolved links) because the only wikilink was ambiguous.
	found := false
	for _, s := range sinks {
		if s == "root" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'root' to be a sink when its only wikilink is ambiguous; sinks=%v", sinks)
	}
	// Neither a/note nor b/note should have an explicit raw edge from root.
	// Sink rows in the teleporting kernel collapse to the uniform restart (1/n).
	if kern.P.At(idx["root"], idx["a/note"]) > 0 && kern.P.At(idx["root"], idx["b/note"]) > 0 {
		n := float64(len(pages))
		if kern.P.At(idx["root"], idx["a/note"])*n != 1.0 {
			t.Errorf("expected uniform restart row for sink root, got P[root][a/note]=%v", kern.P.At(idx["root"], idx["a/note"]))
		}
	}
}
