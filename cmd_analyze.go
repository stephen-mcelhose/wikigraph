package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/spf13/cobra"
)

var (
	flagAnalyzeOrphanPct  float64
	flagAnalyzeSuggestTop int
	flagAnalyzeMinCommute float64
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze <wiki-dir>",
	Short: "Print a wiki health report",
	Long: `analyze prints six sections about the wiki's graph structure:

  1. Overview          — page count, edge count, entropy rate, class count
  2. Communicating classes — pages per class; transient classes flagged
  3. Orphan pages      — low stationary-distribution pages (add inbound links)
  4. Sink pages        — pages with no outgoing links (add outbound links)
  5. Most central      — top 5 by stationary distribution
  6. Suggested links   — unlinked pairs with low commute time (skip with --suggest-top 0)`,

	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	analyzeCmd.Flags().Float64Var(&flagAnalyzeOrphanPct, "orphan-pct", 0.10,
		"percentile threshold for orphan pages")
	analyzeCmd.Flags().IntVar(&flagAnalyzeSuggestTop, "suggest-top", 3,
		"suggested links per page, 0 to skip")
	analyzeCmd.Flags().Float64Var(&flagAnalyzeMinCommute, "min-commute", 2.0,
		"min commute time to surface a suggestion")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	wikiDir := args[0]
	exclude := makeExcludeMap(flagExclude)

	kern, pages, sinkPages, err := buildKernel(wikiDir, flagRecursive, exclude)
	if err != nil {
		return err
	}

	n := len(pages)

	// Stationary distribution.
	pi, statErr := kern.Stationary(1e-12, 5000)
	if statErr != nil {
		pi = make([]float64, n)
		u := 1.0 / float64(n)
		for i := range pi {
			pi[i] = u
		}
	}

	// Class decomposition.
	cd, err := kern.Classes(1e-10)
	if err != nil {
		return fmt.Errorf("class decomposition: %w", err)
	}

	// Edge count: all (i,j) pairs where p > 0.
	edgeCount := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if kern.P.At(i, j) > 1e-10 {
				edgeCount++
			}
		}
	}

	// Entropy rate (bits per step).
	entropyRate, _ := kern.EntropyRate(2)

	// Transient state set.
	transientSet := make(map[int]bool, len(cd.Transient))
	for _, t := range cd.Transient {
		transientSet[t] = true
	}

	// === 1. Overview ===
	fmt.Printf("=== Overview ===\n")
	fmt.Printf("Pages:        %d\n", n)
	fmt.Printf("Edges:        %d\n", edgeCount)
	fmt.Printf("Entropy rate: %.4f bits\n", entropyRate)
	fmt.Printf("Classes:      %d\n", len(cd.SCCs))
	fmt.Println()

	// === 2. Communicating classes ===
	fmt.Printf("=== Communicating classes ===\n")
	for i, comp := range cd.SCCs {
		recurrent := true
		for _, state := range comp {
			if transientSet[state] {
				recurrent = false
				break
			}
		}
		label := "recurrent"
		if !recurrent {
			label = "transient — add links out of this class"
		}
		fmt.Printf("Class %d (%s): %d page(s)\n", i+1, label, len(comp))
		for _, state := range comp {
			fmt.Printf("  %s\n", pages[state])
		}
	}
	fmt.Println()

	// === 3. Orphan pages ===
	piSorted := make([]float64, n)
	copy(piSorted, pi)
	sort.Float64s(piSorted)
	threshIdx := int(math.Floor(float64(n) * flagAnalyzeOrphanPct))
	if threshIdx >= n {
		threshIdx = n - 1
	}
	threshold := piSorted[threshIdx]

	type slugPi struct {
		slug string
		pi   float64
	}
	var orphans []slugPi
	for i, p := range pages {
		if pi[i] <= threshold {
			orphans = append(orphans, slugPi{slug: p, pi: pi[i]})
		}
	}
	sort.Slice(orphans, func(a, b int) bool { return orphans[a].pi < orphans[b].pi })

	fmt.Printf("=== Orphan pages (bottom %.0f%% by stationary distribution) ===\n",
		flagAnalyzeOrphanPct*100)
	if len(orphans) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, o := range orphans {
			fmt.Printf("  %-40s  π=%.6f  → add inbound links\n", o.slug, o.pi)
		}
	}
	fmt.Println()

	// === 4. Sink pages ===
	fmt.Printf("=== Sink pages (no outgoing links) ===\n")
	if len(sinkPages) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, s := range sinkPages {
			fmt.Printf("  %-40s  → add outgoing links\n", s)
		}
	}
	fmt.Println()

	// === 5. Most central ===
	type pageRank struct {
		slug string
		pi   float64
	}
	ranked := make([]pageRank, n)
	for i, p := range pages {
		ranked[i] = pageRank{slug: p, pi: pi[i]}
	}
	sort.Slice(ranked, func(a, b int) bool { return ranked[a].pi > ranked[b].pi })

	top5 := 5
	if top5 > n {
		top5 = n
	}
	fmt.Printf("=== Most central (top %d by stationary distribution) ===\n", top5)
	for i := 0; i < top5; i++ {
		fmt.Printf("  %d. %-40s  π=%.6f\n", i+1, ranked[i].slug, ranked[i].pi)
	}
	fmt.Println()

	// === 6. Suggested missing links ===
	if flagAnalyzeSuggestTop == 0 {
		return nil
	}

	fmt.Printf("=== Suggested missing links (lowest commute time, not yet linked) ===\n")
	type suggestion struct {
		target  string
		commute float64
	}
	anySuggestion := false
	for i, slug := range pages {
		var suggestions []suggestion
		for j := range pages {
			if i == j {
				continue
			}
			if kern.P.At(i, j) > 1e-10 {
				continue // already linked
			}
			ct, err := kern.CommuteTime(i, j)
			if err != nil || ct < flagAnalyzeMinCommute {
				continue
			}
			suggestions = append(suggestions, suggestion{target: pages[j], commute: ct})
		}
		sort.Slice(suggestions, func(a, b int) bool {
			return suggestions[a].commute < suggestions[b].commute
		})
		top := flagAnalyzeSuggestTop
		if top > len(suggestions) {
			top = len(suggestions)
		}
		if top == 0 {
			continue
		}
		anySuggestion = true
		fmt.Printf("  %s:\n", slug)
		for li := 0; li < top; li++ {
			fmt.Printf("    → %-38s  (commute: %.2f)\n",
				suggestions[li].target, suggestions[li].commute)
		}
	}
	if !anySuggestion {
		fmt.Println("  (none)")
	}
	fmt.Println()

	return nil
}
