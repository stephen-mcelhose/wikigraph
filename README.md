# wikigraph

Interactive force-directed wikilink graph for Markdown wikis using `[[slug]]` wikilinks.

Visualises your wiki as a **Markov chain** — node size reflects how often a random
link-follower lands on that page (stationary distribution), not just raw in-degree.

## Install

```bash
go install github.com/stephen-mcelhose/wikigraph@latest
```

Requires Go 1.22+.

## Usage

```
wikigraph <subcommand> [flags] <wiki-dir>

Subcommands:
  graph    Generate a self-contained, interactive HTML graph
  goal     Subgraph highlighting the shortest paths to goal pages (MFPT)
  export   Export graph data as JSON, CSV, or DOT
  analyze  Print a wiki health report (orphans, sinks, central pages, suggestions)

Persistent flags (all subcommands):
  -e, --exclude strings   slugs to exclude (default [index, log, AGENTS])
  -r, --recursive         recursively scan subdirectories for Markdown pages
```

### graph

```bash
wikigraph graph <wiki-dir> [flags]

Flags:
  -o, --out string       output HTML file (default "wiki_graph.html")
  -t, --title string     browser tab title (default: <dir> wiki)
  -m, --min-edge float   omit edges below this transition probability (default 0.005)
  -s, --sed string       sed expression(s) applied to the HTML output (repeatable)
```

### goal

```bash
wikigraph goal <wiki-dir> --goal <slug> [--goal <slug> ...] [flags]

Flags:
  --goal strings   target page slug(s) (required, repeatable)
  --top int        number of nodes to include in the subgraph (default 10)
  -o, --out string output HTML file (default "wiki_goal.html")
```

### export

```bash
wikigraph export <wiki-dir> --format <json|csv|dot> [flags]

Flags:
  --format string    output format: json, csv, or dot (required)
  -o, --out string   output file base path (extension appended automatically)
  -m, --min-edge float  omit edges below this probability (default 0.005)
```

### analyze

```bash
wikigraph analyze <wiki-dir> [flags]

Flags:
  --orphan-pct float   show pages in the bottom N% by stationary probability (default 0.10)
  --suggest-top int    number of pages to compute missing-link suggestions for (default 5)
```

## Wiki format

- One `.md` file per page, named `<slug>.md`
- Subdirectories supported via `-r / --recursive` (e.g. Obsidian/Logseq vaults)
- Cross-references as `[[slug]]` or `[[slug|Display Name]]` wikilinks
- Compatible with: Obsidian, Foam, Logseq, Roam, and any hand-rolled wiki using the same convention

## How the graph works

wikigraph models your wiki as a **random walk**: a visitor lands on a page, picks a
link at random, and repeats. Pages with no outbound links teleport uniformly to any page.

- **Node size** — stationary distribution π(page): how often the random walker visits this
  page in the long run. Richer than in-degree: a page with one link from a highly-connected
  hub can outrank a page with ten links from leaf pages.
- **Node colour** — communicating class: pages that can all reach each other share a colour.
  Isolated clusters reveal wiki sections that aren't cross-linked.
- **Edge width** — transition probability: how likely the walker follows that specific link.

Open the HTML in any browser. No server needed. Drag, zoom, and pan are fully supported.

## Examples

```bash
# Basic graph of an Obsidian vault
wikigraph graph ~/notes/obsidian-vault

# Custom output file and title
wikigraph graph --out docs/graph.html --title "my notes" ~/notes

# Wiki health report — orphans, sinks, most central pages
wikigraph analyze ~/notes --suggest-top 5

# Learning path toward a goal — 12 closest pages by MFPT
wikigraph goal ~/notes --goal machine-learning --goal neural-networks --top 12 -o /tmp/goal.html

# Export for further analysis
wikigraph export ~/notes --format json -o /tmp/wiki
wikigraph export ~/notes --format dot  -o /tmp/wiki
```

## --sed flag (macOS / Linux only)

The `--sed` flag pipes the generated HTML through `sed` for custom post-processing.
It shells out via `exec.Command("sed", ...)` and is not available on Windows.
Windows users can use the built-in `--title` and `--min-edge` flags without `--sed`.

## Testing

See [docs/testing-runbook.md](docs/testing-runbook.md) for the full manual test plan
covering all subcommands and edge cases.

## Architecture

See [docs/adr-001-embedding-layer.md](docs/adr-001-embedding-layer.md) for the decision
record on the planned `wikigraph vectorize` semantic search feature.

## License

BSD 3-Clause
