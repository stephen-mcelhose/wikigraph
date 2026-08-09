package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/catrace"
)

var (
	flagGoalGoals   []string
	flagGoalTop     int
	flagGoalOut     string
	flagGoalMinEdge float64
	flagGoalSed     []string
)

var goalCmd = &cobra.Command{
	Use:   "goal <wiki-dir>",
	Short: "Compute a learning-path subgraph toward goal pages [PROTOTYPE]",
	Long: `goal builds a subgraph focused on navigating toward one or more goal pages.

PROTOTYPE: goals must currently be exact page slugs. Semantic natural-language
goals ("understand quantum error correction") are planned — see
https://github.com/stephen-mcelhose/wikigraph/issues/1

For each page, mean first-passage time to any goal is computed. The top N
closest pages (by minimum MFPT across all goals) form a subgraph whose
effective (trace) kernel is rendered as an interactive HTML graph.

Node size reflects the stationary distribution of the trace kernel on the subset.`,

	Args: cobra.ExactArgs(1),
	RunE: runGoal,
}

func init() {
	goalCmd.Flags().StringArrayVar(&flagGoalGoals, "goal", nil, "goal page slug (repeatable, at least one required)")
	goalCmd.Flags().IntVar(&flagGoalTop, "top", 10, "number of pages in the subgraph")
	goalCmd.Flags().StringVarP(&flagGoalOut, "out", "o", "goal_graph.html", "output HTML file")
	goalCmd.Flags().Float64VarP(&flagGoalMinEdge, "min-edge", "m", 0.005, "omit edges below this transition probability")
	goalCmd.Flags().StringArrayVarP(&flagGoalSed, "sed", "s", nil, "sed expression(s) to apply to the HTML output (repeatable)")
}

func runGoal(cmd *cobra.Command, args []string) error {
	if len(flagGoalGoals) == 0 {
		return fmt.Errorf("at least one --goal slug is required")
	}

	wikiDir := args[0]
	exclude := makeExcludeMap(flagExclude)

	kern, pages, _, err := buildKernel(wikiDir, flagRecursive, exclude)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Pages: %d\n", len(pages))

	// slug → index map.
	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}

	// Validate --goal slugs.
	goalIdxs := make([]int, 0, len(flagGoalGoals))
	for _, g := range flagGoalGoals {
		i, ok := idx[g]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown --goal slug %q\nValid slugs:\n", g)
			for _, p := range pages {
				fmt.Fprintf(os.Stderr, "  %s\n", p)
			}
			return fmt.Errorf("unknown --goal slug: %s", g)
		}
		goalIdxs = append(goalIdxs, i)
	}

	n := len(pages)

	// For each page i, score = min over goals of MFPT(i, goal).
	const inf = 1e18
	scores := make([]float64, n)
	for i := range scores {
		scores[i] = inf
	}
	for _, gIdx := range goalIdxs {
		scores[gIdx] = 0
		for i := 0; i < n; i++ {
			if i == gIdx {
				continue
			}
			mfpt, err := kern.MeanFirstPassage(i, gIdx)
			if err != nil {
				continue // unreachable — leave at inf
			}
			if mfpt < scores[i] {
				scores[i] = mfpt
			}
		}
	}

	// Rank ascending; always include goal pages.
	type ranked struct {
		idx   int
		score float64
	}
	all := make([]ranked, n)
	for i := range all {
		all[i] = ranked{idx: i, score: scores[i]}
	}
	sort.Slice(all, func(a, b int) bool { return all[a].score < all[b].score })

	top := flagGoalTop
	if top > n {
		top = n
	}

	selected := make(map[int]bool)
	for _, g := range goalIdxs {
		selected[g] = true
	}
	for _, r := range all {
		if len(selected) >= top {
			break
		}
		selected[r.idx] = true
	}

	subset := make([]int, 0, len(selected))
	for i := range selected {
		subset = append(subset, i)
	}
	sort.Ints(subset)

	// Compute effective kernel on subset via trace.
	traceKern, err := kern.Trace(subset, 1e-9)
	if err != nil {
		return fmt.Errorf("trace failed (try increasing --top): %w", err)
	}

	title := filepath.Base(wikiDir) + " → " + flagGoalGoals[0]
	html, err := traceKern.ToHTML(&catrace.VisualiseOptions{
		Title:   title,
		MinEdge: flagGoalMinEdge,
		Width:   1400,
		Height:  900,
	})
	if err != nil {
		return fmt.Errorf("generating HTML: %w", err)
	}

	page := string(html)
	if len(flagGoalSed) > 0 {
		page, err = applySed(page, flagGoalSed)
		if err != nil {
			return err
		}
	}

	if err := os.WriteFile(flagGoalOut, []byte(page), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", flagGoalOut, err)
	}
	fmt.Fprintf(os.Stderr, "Written: %s (%d nodes)\n", flagGoalOut, len(subset))
	return nil
}
