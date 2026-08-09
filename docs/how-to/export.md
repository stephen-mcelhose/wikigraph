---
type: how-to
title: How to export your wiki graph for external tools
description: Export the Markov kernel of your wiki as JSON, CSV, or DOT so you can analyse or visualise it outside wikigraph.
tags: [export, json, csv, dot, d3, graphviz]
---

# How to export your wiki graph for external tools

Use `wikigraph export` when you want to take the graph data into another tool —
D3/Observable, Gephi, Graphviz, a spreadsheet, or your own script.

## Goal

At the end of this guide you will have one or more files on disk containing
every page's centrality score and every link's transition probability, filtered
to the edges that matter.

## Prerequisites

- `wikigraph` installed and on your PATH (`wikigraph --version` to check)
- A wiki directory: one `.md` file per page, named `<slug>.md`, with
  `[[slug]]` wikilinks in the body

## Choosing a format

| Format | Output file(s)                         | Good for                              |
| ------ | -------------------------------------- | ------------------------------------- |
| `json` | `<out>.json`                           | D3, Observable, custom scripts        |
| `csv`  | `<out>_nodes.csv`, `<out>_edges.csv`   | Spreadsheets, Gephi, pandas           |
| `dot`  | `<out>.dot`                            | Graphviz (`dot`, `neato`, `fdp`)      |

## Steps

### 1. Export as JSON (default)

```bash
wikigraph export ~/notes --format json -o /tmp/wiki
```

Output: `/tmp/wiki.json`

The JSON is D3 node-link format:

```json
{
  "nodes": [
    { "id": "machine-learning", "pi": 0.042, "class": 0 },
    { "id": "neural-networks",  "pi": 0.031, "class": 0 }
  ],
  "links": [
    { "source": "machine-learning", "target": "neural-networks", "value": 0.25 }
  ]
}
```

- `pi` — stationary probability (centrality): how often a random walker visits this page
- `class` — communicating class index (`-1` means the page is transient / isolated)
- `value` — transition probability on that edge

### 2. Export as CSV

```bash
wikigraph export ~/notes --format csv -o /tmp/wiki
```

Output: two files.

`/tmp/wiki_nodes.csv`:
```
slug,pi,class
machine-learning,0.0420000000,0
neural-networks,0.0310000000,0
quantum-computing,0.0180000000,1
```

`/tmp/wiki_edges.csv`:
```
source,target,probability
machine-learning,neural-networks,0.2500000000
neural-networks,machine-learning,0.1200000000
```

### 3. Export as DOT (Graphviz)

```bash
wikigraph export ~/notes --format dot -o /tmp/wiki
```

Output: `/tmp/wiki.dot`

Render it immediately:

```bash
dot -Tpng /tmp/wiki.dot -o /tmp/wiki.png
open /tmp/wiki.png
```

### 4. Filter out weak edges

By default, edges with a transition probability below `0.005` are omitted.
Raise the threshold to focus only on the strongest links:

```bash
wikigraph export ~/notes --format json -o /tmp/wiki --min-edge 0.05
```

Lower it (or set to `0`) to include everything:

```bash
wikigraph export ~/notes --format dot -o /tmp/wiki --min-edge 0
```

### 5. Exclude meta-pages

`index`, `log`, and `AGENTS` are excluded by default. Add your own:

```bash
wikigraph export ~/notes --format json -o /tmp/wiki --exclude index --exclude log --exclude README
```

## Verification

Check the output landed:

```bash
# JSON
cat /tmp/wiki.json | python3 -m json.tool | head -20

# CSV
head /tmp/wiki_nodes.csv
wc -l /tmp/wiki_edges.csv

# DOT
head /tmp/wiki.dot
```

stderr during the run prints `Pages: N` and `Written: <file>` — if you see those, it worked.

## Troubleshooting

| Symptom                              | Cause                                          | Fix                                              |
| ------------------------------------ | ---------------------------------------------- | ------------------------------------------------ |
| `unknown format "xlsx"`              | Unsupported format                             | Use `json`, `csv`, or `dot`                      |
| Output file has 0 edges              | `--min-edge` too high for your wiki            | Lower with `--min-edge 0.001` or `--min-edge 0`  |
| A page you expect is missing         | It has no `.md` extension, or it's in `--exclude` | Check filename and exclude list              |
| `Pages: 1` on stderr                 | Only one page found — likely wrong directory   | Double-check the path to your wiki              |
