---
type: how-to
title: How to generate an interactive wiki graph
description: Turn your markdown wiki into a browsable force-directed graph that shows page centrality, link clusters, and transition probabilities.
tags: [graph, visualisation, html, force-directed]
resource: cmd_graph.go
timestamp: 2026-08-09T07:31:56Z
---

# How to generate an interactive wiki graph

`wikigraph graph` reads your wiki and produces a self-contained HTML file you
open in any browser — no server, no dependencies. Drag nodes, zoom, and pan
to explore your wiki's link topology.

## Goal

A single `.html` file that shows every page as a node and every wikilink as a
directed edge, with size and colour encoding graph-theoretic properties.

## Prerequisites

- `wikigraph` installed and on your PATH
- A wiki directory: one `.md` file per page, named `<slug>.md`, with
  `[[slug]]` wikilinks in the body

## What the graph encodes

| Visual property | What it means                                                                 |
| --------------- | ----------------------------------------------------------------------------- |
| **Node size**   | Stationary distribution π — how often a random walker lands on this page     |
| **Node colour** | Communicating class — pages that can all reach each other share a colour      |
| **Edge width**  | Transition probability — how likely the walker follows that specific link     |

Pages that are large and central are your wiki's hubs. Pages that are small
and isolated need more inbound links. Pages sharing a colour form a
self-contained cluster. For a full breakdown of what these properties mean and
how to act on them, see [[analyze]].

## Steps

### 1. Basic graph

```bash
wikigraph graph ~/notes
```

Opens a file called `wiki_graph.html` in the current directory. Open it:

```bash
open wiki_graph.html          # macOS
xdg-open wiki_graph.html      # Linux
```

### 2. Set a custom output path and title

```bash
wikigraph graph ~/notes -o docs/graph.html --title "my notes"
```

The title appears in the browser tab and as an `<h1>` in the page.

### 3. Reduce visual noise with `--min-edge`

Very small transition probabilities clutter the graph. The default threshold
is `0.005`. Raise it to show only the strongest links:

```bash
wikigraph graph ~/notes --min-edge 0.05
```

Lower it to see every edge (including tenuous ones from teleportation):

```bash
wikigraph graph ~/notes --min-edge 0.001
```

### 4. Post-process the HTML with `--sed` (macOS / Linux)

Inject custom CSS or replace text without modifying the source:

```bash
# Change the background colour
wikigraph graph ~/notes -s 's/background:#1a1a2e/background:#0f172a/'

# Apply multiple expressions
wikigraph graph ~/notes \
  -s 's/wiki_graph/My Wiki/' \
  -s 's/font-size:12px/font-size:14px/'
```

`--sed` is not available on Windows. Use `--title` and `--min-edge` instead.

### 5. Exclude meta-pages

```bash
wikigraph graph ~/notes --exclude index --exclude log --exclude README
```

`index`, `log`, and `AGENTS` are excluded by default. This replaces the
defaults, so re-add them if needed:

```bash
wikigraph graph ~/notes -e index -e log -e AGENTS -e README
```

## Verification

After the run you should see on stderr:

```
Pages: 42
Written: wiki_graph.html
```

Open the HTML and confirm:
- Nodes are visible and draggable
- Multiple colours appear if your wiki has separate clusters
- Hovering a node shows its slug and π value

## Troubleshooting

| Symptom                              | Cause                                      | Fix                                                  |
| ------------------------------------ | ------------------------------------------ | ---------------------------------------------------- |
| `Pages: 1` — almost empty graph      | Wrong directory, or only one `.md` file    | Check the path; ensure files use `.md` extension     |
| All nodes the same size              | Wiki may be a single communicating class   | Expected — add more cross-links to differentiate     |
| Graph is a hairball of edges         | `--min-edge` too low                       | Raise to `0.02` or higher                            |
| `--sed` flag not recognised          | Running on Windows                         | Use `--title` / `--min-edge` instead                 |
| Output HTML is blank in browser      | Very large wiki (1000+ pages)              | Filter with `--min-edge 0.02` to reduce render load  |

## See also

- [[analyze]] — interpret the communicating classes and centrality shown in the graph
- [[testing-runbook]] — end-to-end verification steps for the `graph` subcommand

## Sources

- `cmd_graph.go`
- `wiki.go`
