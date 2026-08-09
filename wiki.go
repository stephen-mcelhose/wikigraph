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

var wikilinkRe = regexp.MustCompile(`\[\[([A-Za-z][A-Za-z0-9-]*)(?:\|[^\]]+)?\]\]`)

// makeExcludeMap converts a slice of slugs to a set for O(1) lookup.
func makeExcludeMap(slugs []string) map[string]bool {
	m := make(map[string]bool, len(slugs))
	for _, s := range slugs {
		m[s] = true
	}
	return m
}

// loadPages reads all non-meta .md files from dir (optionally recursively)
// and returns sorted slugs, slug->index map, and slug->filePath map.
func loadPages(dir string, recursive bool, exclude map[string]bool) ([]string, map[string]int, map[string]string, error) {
	paths := make(map[string]string)
	var pages []string

	if recursive {
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				// Skip hidden directories like .git, .obsidian, .trash
				if strings.HasPrefix(d.Name(), ".") && path != dir {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(d.Name(), ".md") {
				return nil
			}
			slug := strings.TrimSuffix(d.Name(), ".md")
			if exclude[slug] {
				return nil
			}
			if existingPath, ok := paths[slug]; ok {
				return fmt.Errorf("duplicate slug %q found in %q and %q", slug, existingPath, path)
			}
			paths[slug] = path
			pages = append(pages, slug)
			return nil
		})
		if err != nil {
			return nil, nil, nil, fmt.Errorf("walking wiki dir %q: %w", dir, err)
		}
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("reading wiki dir %q: %w", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			slug := strings.TrimSuffix(e.Name(), ".md")
			if exclude[slug] {
				continue
			}
			paths[slug] = filepath.Join(dir, e.Name())
			pages = append(pages, slug)
		}
	}

	sort.Strings(pages)
	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}
	return pages, idx, paths, nil
}

// buildAdjacency reads each page and extracts [[wikilinks]], returning a square
// adjacency matrix and the slugs of sink pages.
// Sink pages (no outgoing links) get uniform weight across all pages so the
// Markov kernel stays well-defined and the stationary distribution can be computed.
func buildAdjacency(pages []string, idx map[string]int, paths map[string]string) (*mat.Dense, []string, error) {
	n := len(pages)
	adj := mat.NewDense(n, n, nil)
	var sinks []string
	for i, slug := range pages {
		filePath, ok := paths[slug]
		if !ok {
			return nil, nil, fmt.Errorf("missing path for slug %q", slug)
		}
		raw, err := os.ReadFile(filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("reading %s: %w", filePath, err)
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
func buildKernel(wikiDir string, recursive bool, exclude map[string]bool) (*catrace.Kernel, []string, []string, error) {
	pages, idx, paths, err := loadPages(wikiDir, recursive, exclude)
	if err != nil {
		return nil, nil, nil, err
	}
	adj, sinkPages, err := buildAdjacency(pages, idx, paths)
	if err != nil {
		return nil, nil, nil, err
	}
	k, err := catrace.NewRandomWalkKernel(adj, pages)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("building kernel: %w", err)
	}
	return k, pages, sinkPages, nil
}
