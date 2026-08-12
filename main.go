// wikigraph is a suite of tools for analysing a wiki's internal link structure.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var defaultExcludes = []string{"index", "log", "AGENTS"}

// flagExclude is a persistent flag inherited by all subcommands.
var flagExclude []string

// flagRecursive controls whether subdirectories are scanned recursively.
var flagRecursive bool

// flagRelativeLinks enables parsing of standard Markdown [label](path.md)
// links (relative paths only) as graph edges, in addition to [[wikilinks]].
// Implies recursive traversal, since relative paths often cross subdirectories.
var flagRelativeLinks bool

var rootCmd = &cobra.Command{
	Use:   "wikigraph",
	Short: "Interactive wiki link graph and analysis tools.",
	Long: `wikigraph is a suite of tools for analysing a wiki's internal link structure.

Wiki format expected:
  - One Markdown file per page, named <slug>.md (e.g. grovers-algorithm.md)
  - Cross-references written as [[slug]] wikilinks anywhere in the body
  - With --relative-links, standard Markdown [label](relative/path.md) links
    are also treated as edges (absolute URLs are ignored); this mode is
    always recursive and warns if a link resolves outside the wiki root
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
	rootCmd.PersistentFlags().BoolVarP(&flagRecursive, "recursive", "r", false,
		"recursively scan subdirectories for Markdown pages")
	rootCmd.PersistentFlags().BoolVar(&flagRelativeLinks, "relative-links", false,
		"also parse standard Markdown [label](relative/path.md) links as edges (absolute URLs ignored); implies --recursive")
	rootCmd.AddCommand(graphCmd, goalCmd, exportCmd, analyzeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
