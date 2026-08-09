package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

var wikilinkRe = regexp.MustCompile(`\[\[([A-Za-z0-9][A-Za-z0-9/_-]*)(?:\|[^\]]+)?\]\]`)

// makeExcludeMap converts a slice of slugs to a set for O(1) lookup.
func makeExcludeMap(slugs []string) map[string]bool {
	m := make(map[string]bool, len(slugs))
	for _, s := range slugs {
		m[s] = true
	}
	return m
}

// loadPages reads all non-meta .md files from dir (recursively) and returns
// the sorted list of slugs and a slug→index map.
// Slugs are relative paths from dir without the .md suffix (e.g. "how-to/analyze").
// A slug is excluded if its full path OR its basename matches the exclude set.
func loadPages(dir string, exclude map[string]bool) ([]string, map[string]int, error) {
	var pages []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		slug := strings.TrimSuffix(rel, ".md")
		base := strings.TrimSuffix(d.Name(), ".md")
		if exclude[slug] || exclude[base] {
			return nil
		}
		pages = append(pages, slug)
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("reading wiki dir %q: %w", dir, err)
	}
	sort.Strings(pages)
	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}
	return pages, idx, nil
}

// buildAdjacency reads each page and extracts [[wikilinks]], returning a square
// adjacency matrix and the slugs of sink pages.
// Sink pages (no outgoing links) get uniform weight across all pages so the
// Markov kernel stays well-defined and the stationary distribution can be computed.
func buildAdjacency(dir string, pages []string, idx map[string]int) (*mat.Dense, []string, error) {
	n := len(pages)
	adj := mat.NewDense(n, n, nil)
	var sinks []string
	for i, slug := range pages {
		raw, err := os.ReadFile(filepath.Join(dir, slug+".md"))
		if err != nil {
			return nil, nil, fmt.Errorf("reading %s.md: %w", slug, err)
		}
		linked := map[int]bool{}
		for _, m := range wikilinkRe.FindAllSubmatch(raw, -1) {
			if j, ok := idx[strings.ToLower(string(m[1]))]; ok && j != i {
				linked[j] = true
			}
		}
		if len(linked) == 0 {
			// Sink node: teleport uniformly to avoid a zero row.
			sinks = append(sinks, slug)
			for j := 0; j < n; j++ {
				adj.Set(i, j, 1.0)
			}
		} else {
			for j := range linked {
				adj.Set(i, j, 1.0)
			}
		}
	}
	return adj, sinks, nil
}

// buildKernel builds a Markov kernel from the wiki directory.
// Returns the kernel, sorted page slugs, sink slugs (pages whose adjacency row
// was set to uniform teleportation), and any error.
func buildKernel(wikiDir string, exclude map[string]bool) (*catrace.Kernel, []string, []string, error) {
	pages, idx, err := loadPages(wikiDir, exclude)
	if err != nil {
		return nil, nil, nil, err
	}
	adj, sinkPages, err := buildAdjacency(wikiDir, pages, idx)
	if err != nil {
		return nil, nil, nil, err
	}
	k, err := catrace.NewRandomWalkKernel(adj, pages)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("building kernel: %w", err)
	}
	return k, pages, sinkPages, nil
}
