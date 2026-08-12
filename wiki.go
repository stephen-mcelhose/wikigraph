package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

var wikilinkRe = regexp.MustCompile(`\[\[([A-Za-z][A-Za-z0-9-]*)(?:\|[^\]]+)?\]\]`)

// mdLinkRe matches standard Markdown [label](target) links (but not [[wikilinks]]).
var mdLinkRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

// isAbsoluteLinkTarget reports whether a markdown link target points outside
// the wiki corpus (a URL), and should therefore be ignored as an edge source.
func isAbsoluteLinkTarget(target string) bool {
	return strings.HasPrefix(target, "http://") ||
		strings.HasPrefix(target, "https://") ||
		strings.HasPrefix(target, "//") ||
		strings.HasPrefix(target, "mailto:")
}

// invertPaths builds a slug->path map's inverse: absolute-path -> slug.
func invertPaths(paths map[string]string) map[string]string {
	inv := make(map[string]string, len(paths))
	for slug, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			abs = filepath.Clean(p)
		}
		inv[abs] = slug
	}
	return inv
}

// resolveRelativeLinks extracts [label](target) markdown links from raw whose
// target is a relative path (not an absolute URL), resolves each target
// relative to sourceFilePath's directory, and maps it to a known page slug.
// Targets that resolve outside rootDir produce a warning on stderr, since
// relative-link mode assumes all targets stay within the project root.
func resolveRelativeLinks(raw []byte, sourceFilePath, rootDir string, pathToSlug map[string]string) map[string]bool {
	return resolveRelativeLinksTo(os.Stderr, raw, sourceFilePath, rootDir, pathToSlug)
}

// resolveRelativeLinksTo is resolveRelativeLinks with an injectable warning
// writer, for testability.
func resolveRelativeLinksTo(warnOut io.Writer, raw []byte, sourceFilePath, rootDir string, pathToSlug map[string]string) map[string]bool {
	linked := map[string]bool{}

	rootAbs, err := filepath.Abs(rootDir)
	if err != nil {
		rootAbs = filepath.Clean(rootDir)
	}
	sourceAbs, err := filepath.Abs(sourceFilePath)
	if err != nil {
		sourceAbs = filepath.Clean(sourceFilePath)
	}
	sourceDir := filepath.Dir(sourceAbs)

	for _, m := range mdLinkRe.FindAllSubmatchIndex(raw, -1) {
		// Skip image links: ![alt](src)
		if m[0] > 0 && raw[m[0]-1] == '!' {
			continue
		}
		target := strings.TrimSpace(string(raw[m[4]:m[5]]))
		if target == "" {
			continue
		}
		// A link target may carry a trailing "title" after whitespace, e.g. (foo.md "Title").
		if idx := strings.IndexAny(target, " \t"); idx >= 0 {
			target = target[:idx]
		}
		if isAbsoluteLinkTarget(target) {
			continue
		}
		if strings.HasPrefix(target, "#") {
			// Same-page anchor only, not a cross-page edge.
			continue
		}
		// Strip fragment/query suffixes.
		if idx := strings.IndexAny(target, "#?"); idx >= 0 {
			target = target[:idx]
		}
		if target == "" {
			continue
		}

		resolved := filepath.Clean(filepath.Join(sourceDir, target))

		if rel, err := filepath.Rel(rootAbs, resolved); err == nil {
			if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				fmt.Fprintf(warnOut, "warning: relative link %q in %s resolves outside the wiki root %q; ignoring\n", target, sourceFilePath, rootDir)
				continue
			}
		}

		if slug, ok := pathToSlug[resolved]; ok {
			linked[slug] = true
			continue
		}
		if !strings.HasSuffix(resolved, ".md") {
			if slug, ok := pathToSlug[resolved+".md"]; ok {
				linked[slug] = true
			}
		}
	}
	return linked
}

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
	return buildAdjacencyWithOpts(pages, idx, paths, false, "")
}

// buildAdjacencyWithOpts is like buildAdjacency but optionally also parses
// standard Markdown [label](relative/path.md) links as edges, in addition to
// [[wikilinks]]. rootDir is used to warn when a relative link resolves
// outside the wiki root.
func buildAdjacencyWithOpts(pages []string, idx map[string]int, paths map[string]string, relativeLinks bool, rootDir string) (*mat.Dense, []string, error) {
	n := len(pages)
	if n == 0 {
		return nil, nil, fmt.Errorf("no .md pages found")
	}
	adj := mat.NewDense(n, n, nil)
	var sinks []string

	var pathToSlug map[string]string
	if relativeLinks {
		pathToSlug = invertPaths(paths)
	}

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
		if relativeLinks {
			for targetSlug := range resolveRelativeLinks(raw, filePath, rootDir, pathToSlug) {
				if j, ok := idx[targetSlug]; ok && j != i {
					linked[j] = true
				}
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
// Returns the kernel, transition matrix, sorted page slugs, sink slugs, and any error.
func buildKernel(wikiDir string, recursive bool, exclude map[string]bool) (*catrace.Kernel, *mat.Dense, []string, []string, error) {
	return buildKernelWithOpts(wikiDir, recursive, exclude, false)
}

// buildKernelWithOpts is like buildKernel but optionally enables parsing of
// standard Markdown [label](relative/path.md) links as edges, in addition to
// [[wikilinks]]. When relativeLinks is true, traversal is always recursive
// regardless of the recursive argument, since relative paths (e.g. ../sibling)
// commonly cross subdirectory boundaries.
func buildKernelWithOpts(wikiDir string, recursive bool, exclude map[string]bool, relativeLinks bool) (*catrace.Kernel, *mat.Dense, []string, []string, error) {
	if relativeLinks {
		recursive = true
	}
	pages, idx, paths, err := loadPages(wikiDir, recursive, exclude)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	adj, sinkPages, err := buildAdjacencyWithOpts(pages, idx, paths, relativeLinks, wikiDir)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	k, err := catrace.NewRandomWalkKernel(adj, pages)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("building kernel: %w", err)
	}
	// Row-normalize adj to create transition matrix P
	n := len(pages)
	P := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		rowSum := 0.0
		for j := 0; j < n; j++ {
			rowSum += adj.At(i, j)
		}
		if rowSum > 0 {
			for j := 0; j < n; j++ {
				P.Set(i, j, adj.At(i, j)/rowSum)
			}
		}
	}
	return k, P, pages, sinkPages, nil
}
