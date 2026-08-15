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

var wikilinkRe = regexp.MustCompile(`\[\[([A-Za-z0-9][A-Za-z0-9-]*)(?:\|[^\]]+)?\]\]`)

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
	rootAbs, sourceDir := absRootAndSourceDir(rootDir, sourceFilePath)

	for _, m := range mdLinkRe.FindAllSubmatchIndex(raw, -1) {
		if isImageMarkdownLink(raw, m[0]) {
			continue
		}
		target, ok := normalizeMDLinkTarget(string(raw[m[4]:m[5]]))
		if !ok {
			continue
		}
		resolved := filepath.Clean(filepath.Join(sourceDir, target))
		if escapesWikiRoot(rootAbs, resolved) {
			fmt.Fprintf(warnOut, "warning: relative link %q in %s resolves outside the wiki root %q; ignoring\n", target, sourceFilePath, rootDir)
			continue
		}
		if slug, ok := lookupPathSlug(pathToSlug, resolved); ok {
			linked[slug] = true
		}
	}
	return linked
}

func absRootAndSourceDir(rootDir, sourceFilePath string) (rootAbs, sourceDir string) {
	rootAbs, err := filepath.Abs(rootDir)
	if err != nil {
		rootAbs = filepath.Clean(rootDir)
	}
	sourceAbs, err := filepath.Abs(sourceFilePath)
	if err != nil {
		sourceAbs = filepath.Clean(sourceFilePath)
	}
	return rootAbs, filepath.Dir(sourceAbs)
}

func isImageMarkdownLink(raw []byte, start int) bool {
	return start > 0 && raw[start-1] == '!'
}

// normalizeMDLinkTarget strips optional title / fragment / query and rejects
// absolute URLs, bare anchors, and empty targets.
func normalizeMDLinkTarget(target string) (string, bool) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "", false
	}
	if idx := strings.IndexAny(target, " \t"); idx >= 0 {
		target = target[:idx]
	}
	if isAbsoluteLinkTarget(target) {
		return "", false
	}
	if strings.HasPrefix(target, "#") {
		return "", false
	}
	if idx := strings.IndexAny(target, "#?"); idx >= 0 {
		target = target[:idx]
	}
	if target == "" {
		return "", false
	}
	return target, true
}

func escapesWikiRoot(rootAbs, resolved string) bool {
	rel, err := filepath.Rel(rootAbs, resolved)
	if err != nil {
		return false
	}
	return rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func lookupPathSlug(pathToSlug map[string]string, resolved string) (string, bool) {
	if slug, ok := pathToSlug[resolved]; ok {
		return slug, true
	}
	if !strings.HasSuffix(resolved, ".md") {
		if slug, ok := pathToSlug[resolved+".md"]; ok {
			return slug, true
		}
	}
	return "", false
}

// makeExcludeMap converts a slice of slugs to a set for O(1) lookup.
func makeExcludeMap(slugs []string) map[string]bool {
	m := make(map[string]bool, len(slugs))
	for _, s := range slugs {
		m[s] = true
	}
	return m
}

// isExcluded reports whether slug should be excluded.
// It checks for an exact match first, then a basename match so that
// --exclude index suppresses any */index.md at any directory depth.
func isExcluded(slug string, exclude map[string]bool) bool {
	return exclude[slug] || exclude[filepath.Base(slug)]
}

// loadPages reads all non-meta .md files from dir (optionally recursively)
// and returns sorted slugs, slug->index map, and slug->filePath map.
func loadPages(dir string, recursive bool, exclude map[string]bool) ([]string, map[string]int, map[string]string, error) {
	var pages []string
	var paths map[string]string
	var err error
	if recursive {
		pages, paths, err = loadPagesRecursive(dir, exclude)
	} else {
		pages, paths, err = loadPagesFlat(dir, exclude)
	}
	if err != nil {
		return nil, nil, nil, err
	}
	sort.Strings(pages)
	return pages, indexPages(pages), paths, nil
}

func indexPages(pages []string) map[string]int {
	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}
	return idx
}

func loadPagesFlat(dir string, exclude map[string]bool) ([]string, map[string]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("reading wiki dir %q: %w", dir, err)
	}
	paths := make(map[string]string)
	var pages []string
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
	return pages, paths, nil
}

func loadPagesRecursive(dir string, exclude map[string]bool) ([]string, map[string]string, error) {
	paths := make(map[string]string)
	var pages []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		return considerWikiFile(dir, path, d, err, exclude, paths, &pages)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("walking wiki dir %q: %w", dir, err)
	}
	return pages, paths, nil
}

func considerWikiFile(root, path string, d os.DirEntry, walkErr error, exclude map[string]bool, paths map[string]string, pages *[]string) error {
	if walkErr != nil {
		return walkErr
	}
	if d.IsDir() {
		if strings.HasPrefix(d.Name(), ".") && path != root {
			return filepath.SkipDir
		}
		return nil
	}
	if !strings.HasSuffix(d.Name(), ".md") {
		return nil
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}
	slug := strings.ReplaceAll(strings.TrimSuffix(rel, ".md"), string(filepath.Separator), "/")
	if isExcluded(slug, exclude) {
		return nil
	}
	if existingPath, ok := paths[slug]; ok {
		return fmt.Errorf("duplicate slug %q found in %q and %q", slug, existingPath, path)
	}
	paths[slug] = path
	*pages = append(*pages, slug)
	return nil
}

// defaultTeleportAlpha is the PageRank teleport probability (α). Link-following
// weight is (1−α). Matches catrace / firehose convention (not Google's d=0.85).
const defaultTeleportAlpha = 0.15

// buildAdjacency reads each page and extracts [[wikilinks]], returning a square
// raw adjacency matrix (real links only) and the slugs of sink pages.
// Sink rows are left as zeros; ergodicity comes from the teleporting kernel.
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
		raw, filePath, err := readPageRaw(paths, slug)
		if err != nil {
			return nil, nil, err
		}
		linked := resolveWikilinkTargets(raw, i, idx, filePath)
		if relativeLinks {
			mergeRelativeTargets(linked, raw, filePath, rootDir, pathToSlug, idx, i)
		}
		if len(linked) == 0 {
			sinks = append(sinks, slug)
			continue
		}
		for j := range linked {
			adj.Set(i, j, 1.0)
		}
	}
	return adj, sinks, nil
}

func readPageRaw(paths map[string]string, slug string) ([]byte, string, error) {
	filePath, ok := paths[slug]
	if !ok {
		return nil, "", fmt.Errorf("missing path for slug %q", slug)
	}
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("reading %s: %w", filePath, err)
	}
	return raw, filePath, nil
}

func resolveWikilinkTargets(raw []byte, self int, idx map[string]int, filePath string) map[int]bool {
	linked := map[int]bool{}
	for _, m := range wikilinkRe.FindAllSubmatch(raw, -1) {
		ref := strings.ToLower(string(m[1]))
		if j, ok := idx[ref]; ok && j != self {
			linked[j] = true
			continue
		}
		matches := lenientBasenameMatches(ref, self, idx)
		switch len(matches) {
		case 1:
			linked[matches[0]] = true
		default:
			if len(matches) > 1 {
				fmt.Fprintf(os.Stderr, "warning: [[%s]] in %s is ambiguous (%d matches); link dropped\n", ref, filePath, len(matches))
			}
		}
	}
	return linked
}

func lenientBasenameMatches(ref string, self int, idx map[string]int) []int {
	var matches []int
	for candidate, j := range idx {
		if filepath.Base(candidate) == ref && j != self {
			matches = append(matches, j)
		}
	}
	return matches
}

func mergeRelativeTargets(linked map[int]bool, raw []byte, filePath, rootDir string, pathToSlug map[string]string, idx map[string]int, self int) {
	for targetSlug := range resolveRelativeLinks(raw, filePath, rootDir, pathToSlug) {
		if j, ok := idx[targetSlug]; ok && j != self {
			linked[j] = true
		}
	}
}

// restartDistribution builds a probability vector over pages. With no seeds it
// is uniform (global PageRank). With seeds it concentrates mass equally on the
// named slugs (Personalized PageRank).
func restartDistribution(pages []string, seeds []string) ([]float64, error) {
	n := len(pages)
	v := make([]float64, n)
	if len(seeds) == 0 {
		u := 1.0 / float64(n)
		for i := range v {
			v[i] = u
		}
		return v, nil
	}
	idx := make(map[string]int, n)
	for i, p := range pages {
		idx[p] = i
	}
	seen := make(map[int]bool, len(seeds))
	var seedIdxs []int
	for _, s := range seeds {
		i, ok := idx[s]
		if !ok {
			return nil, fmt.Errorf("unknown --seed slug %q", s)
		}
		if !seen[i] {
			seen[i] = true
			seedIdxs = append(seedIdxs, i)
		}
	}
	mass := 1.0 / float64(len(seedIdxs))
	for _, i := range seedIdxs {
		v[i] = mass
	}
	return v, nil
}

// rawEdgeCount returns the number of directed wikilinks in the raw adjacency
// (entries > 0). Prefer this over counting nonzeros in the teleporting P.
func rawEdgeCount(adj *mat.Dense) int {
	r, c := adj.Dims()
	count := 0
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if adj.At(i, j) > 0 {
				count++
			}
		}
	}
	return count
}

// buildKernel builds a teleporting Markov kernel from the wiki directory using
// default α and a uniform restart (global PageRank).
// Returns the math kernel, raw adjacency, sorted page slugs, sink slugs, and any error.
func buildKernel(wikiDir string, recursive bool, exclude map[string]bool) (*catrace.Kernel, *mat.Dense, []string, []string, error) {
	return buildKernelWithOpts(wikiDir, recursive, exclude, false, defaultTeleportAlpha, nil)
}

// buildKernelWithOpts is like buildKernel but optionally enables parsing of
// standard Markdown [label](relative/path.md) links as edges, in addition to
// [[wikilinks]]. When relativeLinks is true, traversal is always recursive
// regardless of the recursive argument, since relative paths (e.g. ../sibling)
// commonly cross subdirectory boundaries.
//
// alpha is the teleport probability; seeds (page slugs) personalize the restart
// distribution. The returned adjacency is raw (no sink pre-fill); the kernel is
// NewTeleportingKernelFromAdj over that adj.
func buildKernelWithOpts(wikiDir string, recursive bool, exclude map[string]bool, relativeLinks bool, alpha float64, seeds []string) (*catrace.Kernel, *mat.Dense, []string, []string, error) {
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
	restart, err := restartDistribution(pages, seeds)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	k, err := catrace.NewTeleportingKernelFromAdj(adj, restart, alpha, pages)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("building teleporting kernel: %w", err)
	}
	return k, adj, pages, sinkPages, nil
}
