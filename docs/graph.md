---
type: how-to
title: How to generate an interactive wiki graph
description: Turn your markdown wiki into a browsable force-directed graph that shows page centrality, link clusters, and transition probabilities.
tags: [how-to, graph, visualisation, html, force-directed]
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
- `quantum-go` cloned as the example wiki (see [[adr-004-quantum-go-example-wiki]]):

```bash
git clone https://github.com/stephen-mcelhose/quantum-go ~/quantum-go
```

- Or any wiki directory: one `.md` file per page, named `<slug>.md`, with
  `[[slug]]` wikilinks in the body

## What the graph encodes

| Visual property | What it means                                                                 |
| --------------- | ----------------------------------------------------------------------------- |
| **Node size**   | Stationary distribution π — how often a random walker lands on this page; a proxy for systemic importance |
| **Node colour** | Communicating class — pages that can all reach each other share a colour      |
| **Edge width**  | Transition probability — how likely the walker follows that specific link     |

Pages that are large and central are your wiki's hubs. Pages that are small
and isolated need more inbound links. Pages sharing a colour form a
self-contained cluster. For a full breakdown of what these properties mean and
how to act on them, see [[analyze]].

## Steps

### 1. Basic graph

```bash
wikigraph graph ~/quantum-go/wiki
```

Writes `wiki_graph.html` to the current directory. Open it:

```bash
open wiki_graph.html          # macOS
xdg-open wiki_graph.html      # Linux
```

### 2. Set a custom output path and title

```bash
wikigraph graph ~/quantum-go/wiki -o ~/quantum-go.html --title "quantum-go"
```

The title appears in the browser tab and as an overlay label in the top-left of the graph.

### 3. Reduce visual noise with `--min-edge`

Very small transition probabilities clutter the graph. The default threshold
is `0.005`. Raise it to show only the strongest links:

```bash
wikigraph graph ~/quantum-go/wiki --min-edge 0.05
```

Lower it to expose weak edges that the default filters out:

```bash
wikigraph graph ~/quantum-go/wiki --min-edge 0.001
```

Open `wiki_graph.html` — raising the threshold should visibly reduce the number of edges.

### 4. Post-process the HTML with `--sed` (macOS / Linux)

Inject custom CSS or replace text without modifying the source:

```bash
# Change the background colour
wikigraph graph ~/quantum-go/wiki -s 's/background: #1a1a2e/background: red/'

# Rename the title and enlarge the tooltip font
wikigraph graph ~/quantum-go/wiki \
  -s 's/wiki wiki/My Wiki/' \
  -s 's/font-size: 12px/font-size: 32px/'
```

> `wiki wiki` is the title when your wiki directory is named `wiki`. Use `--title "My Wiki"` for a reliable rename that works with any directory name.

`--sed` is not available on Windows. Use `--title` and `--min-edge` instead.

Open `wiki_graph.html` and confirm the substitution took effect (e.g. the
background colour or text changed as specified).

### 5. Exclude meta-pages

```bash
wikigraph graph ~/quantum-go/wiki --exclude index --exclude log --exclude README
```

`index`, `log`, and `AGENTS` are excluded by default. This replaces the
defaults, so re-add them if needed:

```bash
wikigraph graph ~/quantum-go/wiki -e index -e log -e AGENTS -e README
```

Open `wiki_graph.html` and confirm the excluded slugs no longer appear as nodes.

## Verification

After the run you should see on stderr:

```
Pages: 36
Written: wiki_graph.html
```

Open the HTML and confirm:
- Nodes are visible and draggable
- Multiple colours appear if your wiki has more than one communicating class (quantum-go has one class, so all nodes share a colour)
- Hovering a node shows its slug and π value

## Troubleshooting

| Symptom                              | Cause                                      | Fix                                                  |
| ------------------------------------ | ------------------------------------------ | ---------------------------------------------------- |
| `wikigraph: command not found`        | Binary not on PATH                         | Run `go install github.com/stephen-mcelhose/wikigraph@latest` or add its bin dir to PATH |
| No `Written:` line on stderr          | Bad output path or write-permission error  | Omit `-o` to write to the current directory, or check permissions on the target path |
| `Pages: 1` — almost empty graph      | Wrong directory, or only one `.md` file    | Check the path; ensure files use `.md` extension     |
| All nodes the same size              | All pages have near-equal π — happens when every page is a sink (uniform teleportation) or the link graph is highly symmetric | Add pages with varied out-degree; check for widespread sink pages with `wikigraph analyze` |
| Graph is a hairball of edges         | `--min-edge` too low                       | Raise to `0.02` or higher                            |
| `--sed` flag not recognised          | Running on Windows                         | Use `--title` / `--min-edge` instead                 |
| Output HTML is blank in browser      | Very large wiki (1000+ pages)              | Filter with `--min-edge 0.02` to reduce render load  |

## See also

- [[analyze]] — interpret the communicating classes and centrality shown in the graph
- [[testing-runbook]] — end-to-end verification steps for the `graph` subcommand

## Sources

- `cmd_graph.go`
- `wiki.go`
