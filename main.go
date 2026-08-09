// wikigraph is a suite of tools for analysing a wiki's internal link structure.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var defaultExcludes = []string{"index", "log", "AGENTS"}

// flagExclude is a persistent flag inherited by all subcommands.
var flagExclude []string

var rootCmd = &cobra.Command{
	Use:   "wikigraph",
	Short: "Interactive wiki link graph and analysis tools.",
	Long: `wikigraph is a suite of tools for analysing a wiki's internal link structure.

Wiki format expected:
  - One Markdown file per page, named <slug>.md (e.g. grovers-algorithm.md)
  - Cross-references written as [[slug]] wikilinks anywhere in the body
  - Meta-files (index, log, AGENTS) are excluded automatically via --exclude
  - All other .md files become nodes in the graph

Subcommands:
  graph    Generate an interactive force-directed HTML graph
  goal     Compute a learning-path subgraph toward one or more goal pages
  export   Export the wiki graph as JSON, CSV, or DOT
  analyze  Print a wiki health report`,
}

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&flagExclude, "exclude", "e", defaultExcludes,
		"slugs to exclude from all subcommands (meta-pages, changelogs, etc.)")
	rootCmd.AddCommand(graphCmd, goalCmd, exportCmd, analyzeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
