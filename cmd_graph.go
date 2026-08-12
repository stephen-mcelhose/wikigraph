package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/catrace"
)

var (
	flagGraphOut     string
	flagGraphTitle   string
	flagGraphMinEdge float64
	flagGraphSed     []string
)

var graphCmd = &cobra.Command{
	Use:   "graph <wiki-dir>",
	Short: "Generate an interactive force-directed graph HTML",
	Long: `graph generates an interactive force-directed graph of a wiki's internal
link structure and writes it to a self-contained HTML file.

  Node size   ∝  stationary distribution (centrality — how often a random walk lands here)
  Node colour =  communicating class (cluster of pages that can all reach each other)
  Edge width  ∝  transition probability (how likely a random walk follows that link)
  Drag, zoom, and pan are fully supported.`,

	Args: cobra.ExactArgs(1),
	RunE: runGraph,
}

func init() {
	graphCmd.Flags().StringVarP(&flagGraphOut, "out", "o", "wiki_graph.html", "output HTML file")
	graphCmd.Flags().StringVarP(&flagGraphTitle, "title", "t", "", "graph title shown in browser tab (default: <wiki-dir> wiki)")
	graphCmd.Flags().Float64VarP(&flagGraphMinEdge, "min-edge", "m", 0.005, "omit edges below this transition probability")
	graphCmd.Flags().StringArrayVarP(&flagGraphSed, "sed", "s", nil, "sed expression(s) to apply to the HTML output (repeatable, e.g. -s 's/foo/bar/')")
}

func runGraph(cmd *cobra.Command, args []string) error {
	wikiDir := args[0]

	title := flagGraphTitle
	if title == "" {
		title = filepath.Base(wikiDir) + " wiki"
	}

	exclude := makeExcludeMap(flagExclude)

	recursive := flagRecursive || flagRelativeLinks
	pages, idx, paths, err := loadPages(wikiDir, recursive, exclude)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Pages: %d\n", len(pages))

	adj, _, err := buildAdjacencyWithOpts(pages, idx, paths, flagRelativeLinks, wikiDir)
	if err != nil {
		return err
	}

	k, err := catrace.NewRandomWalkKernel(adj, pages)
	if err != nil {
		return fmt.Errorf("building kernel: %w", err)
	}

	html, err := k.ToHTML(&catrace.VisualiseOptions{
		Title:   title,
		MinEdge: flagGraphMinEdge,
		Width:   1400,
		Height:  900,
	})
	if err != nil {
		return fmt.Errorf("generating HTML: %w", err)
	}

	page := string(html)

	if len(flagGraphSed) > 0 {
		page, err = applySed(page, flagGraphSed)
		if err != nil {
			return err
		}
	}

	if err := os.WriteFile(flagGraphOut, []byte(page), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", flagGraphOut, err)
	}
	fmt.Fprintf(os.Stderr, "Written: %s\n", flagGraphOut)
	return nil
}
