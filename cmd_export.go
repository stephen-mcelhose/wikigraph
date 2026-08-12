package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/catrace"
)

var (
	flagExportFormat  string
	flagExportOut     string
	flagExportMinEdge float64
)

var exportCmd = &cobra.Command{
	Use:   "export <wiki-dir>",
	Short: "Export the wiki graph as JSON, CSV, or DOT",
	Long: `export writes the wiki's Markov kernel to a file for use in external tools.

Formats:
  json (default)  node-link JSON compatible with D3/Observable
  csv             two files: <out>_nodes.csv and <out>_edges.csv
  dot             Graphviz DOT digraph`,

	Args: cobra.ExactArgs(1),
	RunE: runExport,
}

func init() {
	exportCmd.Flags().StringVar(&flagExportFormat, "format", "json", "output format: json, csv, or dot")
	exportCmd.Flags().StringVarP(&flagExportOut, "out", "o", "wiki_graph", "output file base name")
	exportCmd.Flags().Float64VarP(&flagExportMinEdge, "min-edge", "m", 0.005, "omit edges below this probability")
}

func runExport(cmd *cobra.Command, args []string) error {
	wikiDir := args[0]
	exclude := makeExcludeMap(flagExclude)

	kern, _, pages, _, err := buildKernelWithOpts(wikiDir, flagRecursive, exclude, flagRelativeLinks)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Pages: %d\n", len(pages))

	n := len(pages)

	// Stationary distribution (uniform fallback if chain is reducible).
	pi, err := kern.Stationary(1e-12, 5000)
	if err != nil {
		pi = make([]float64, n)
		u := 1.0 / float64(n)
		for i := range pi {
			pi[i] = u
		}
	}

	// Map each state to its recurrent class index; -1 means transient.
	classOf := make([]int, n)
	for i := range classOf {
		classOf[i] = -1
	}
	if cd, err := kern.Classes(1e-10); err == nil {
		for classNum, comp := range cd.Recurrent {
			for _, state := range comp {
				classOf[state] = classNum
			}
		}
	}

	switch flagExportFormat {
	case "json":
		return doExportJSON(kern, pages, pi, classOf, n)
	case "csv":
		return doExportCSV(kern, pages, pi, classOf, n)
	case "dot":
		return doExportDOT(kern, pages, pi, n)
	default:
		return fmt.Errorf("unknown format %q: choose json, csv, or dot", flagExportFormat)
	}
}

// --- JSON ---

type jsonNode struct {
	ID    string  `json:"id"`
	Pi    float64 `json:"pi"`
	Class int     `json:"class"`
}

type jsonLink struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Value  float64 `json:"value"`
}

type jsonGraph struct {
	Nodes []jsonNode `json:"nodes"`
	Links []jsonLink `json:"links"`
}

func doExportJSON(kern *catrace.Kernel, pages []string, pi []float64, classOf []int, n int) error {
	g := jsonGraph{
		Nodes: make([]jsonNode, n),
		Links: make([]jsonLink, 0),
	}
	for i, p := range pages {
		g.Nodes[i] = jsonNode{ID: p, Pi: pi[i], Class: classOf[i]}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := kern.P.At(i, j)
			if v < flagExportMinEdge {
				continue
			}
			g.Links = append(g.Links, jsonLink{Source: pages[i], Target: pages[j], Value: v})
		}
	}
	out, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling JSON: %w", err)
	}
	outFile := flagExportOut + ".json"
	if err := os.WriteFile(outFile, out, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", outFile, err)
	}
	fmt.Fprintf(os.Stderr, "Written: %s\n", outFile)
	return nil
}

// --- CSV ---

func doExportCSV(kern *catrace.Kernel, pages []string, pi []float64, classOf []int, n int) error {
	nodesFile := flagExportOut + "_nodes.csv"
	var nb strings.Builder
	nb.WriteString("slug,pi,class\n")
	for i, p := range pages {
		fmt.Fprintf(&nb, "%s,%.10f,%d\n", p, pi[i], classOf[i])
	}
	if err := os.WriteFile(nodesFile, []byte(nb.String()), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", nodesFile, err)
	}
	fmt.Fprintf(os.Stderr, "Written: %s\n", nodesFile)

	edgesFile := flagExportOut + "_edges.csv"
	var eb strings.Builder
	eb.WriteString("source,target,probability\n")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := kern.P.At(i, j)
			if v < flagExportMinEdge {
				continue
			}
			fmt.Fprintf(&eb, "%s,%s,%.10f\n", pages[i], pages[j], v)
		}
	}
	if err := os.WriteFile(edgesFile, []byte(eb.String()), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", edgesFile, err)
	}
	fmt.Fprintf(os.Stderr, "Written: %s\n", edgesFile)
	return nil
}

// --- DOT ---

func doExportDOT(kern *catrace.Kernel, pages []string, pi []float64, n int) error {
	var b strings.Builder
	b.WriteString("digraph wiki {\n")
	for i, p := range pages {
		fmt.Fprintf(&b, "  %q [weight=%.6f];\n", p, pi[i])
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := kern.P.At(i, j)
			if v < flagExportMinEdge {
				continue
			}
			fmt.Fprintf(&b, "  %q -> %q [weight=%.6f];\n", pages[i], pages[j], v)
		}
	}
	b.WriteString("}\n")

	outFile := flagExportOut + ".dot"
	if err := os.WriteFile(outFile, []byte(b.String()), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", outFile, err)
	}
	fmt.Fprintf(os.Stderr, "Written: %s\n", outFile)
	return nil
}
