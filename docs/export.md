---
type: how-to
title: How to export your wiki graph for external tools
description: Export the Markov kernel of your wiki as JSON, CSV, or DOT so you can analyse or visualise it outside wikigraph.
tags: [how-to, export, json, csv, dot, d3, graphviz]
resource: cmd_export.go
timestamp: 2026-08-09T07:31:56Z
---

# How to export your wiki graph for external tools

Use `wikigraph export` when you want to take the graph data into another tool —
D3/Observable, Gephi, Graphviz, a spreadsheet, or your own script.

## Goal

At the end of this guide you will have one or more files on disk containing
every page's PageRank π and every **raw** wikilink edge (not teleportation mass),
filtered by `--min-edge` on edge weight.

## Prerequisites

- `wikigraph` installed and on your PATH (`wikigraph --version` to check)
- The quantum-go example wiki cloned locally:
  ```bash
  git clone https://github.com/stephen-mcelhose/quantum-go ~/quantum-go
  ```
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
wikigraph export ~/quantum-go/wiki --format json -o /tmp/wiki
```

Expected stderr:

```
Pages: 36
Written: /tmp/wiki.json
```

The JSON is D3 node-link format:

```json
{
  "nodes": [
    { "id": "algorithm-comparison", "pi": 0.025630, "class": 0 },
    { "id": "arithmetic-gates",     "pi": 0.041796, "class": 0 }
  ],
  "links": [
    { "source": "algorithm-comparison", "target": "arithmetic-gates", "value": 0.027778 }
  ]
}
```

- `pi` — PageRank π (teleporting kernel; see [[stationary-distribution]])
- `class` — raw digraph SCC index (see [[analyze]])
- `value` — raw directed wikilink weight (typically `1.0`), not teleporting math $P$

### 2. Export as CSV

```bash
wikigraph export ~/quantum-go/wiki --format csv -o /tmp/wiki
```

Expected stderr:

```
Pages: 36
Written: /tmp/wiki_nodes.csv
Written: /tmp/wiki_edges.csv
```

`/tmp/wiki_nodes.csv` (36 rows + header):
```
slug,pi,class
algorithm-comparison,0.0256299258,0
arithmetic-gates,0.0417959521,0
bb84-qkd,0.0144622871,0
```

`/tmp/wiki_edges.csv` (header `source,target,weight`; one row per raw wikilink):
```
source,target,weight
algorithm-comparison,arithmetic-gates,1.0000000000
arithmetic-gates,bb84-qkd,1.0000000000
```

### 3. Export as DOT (Graphviz)

```bash
wikigraph export ~/quantum-go/wiki --format dot -o /tmp/wiki
```

Expected stderr:

```
Pages: 36
Written: /tmp/wiki.dot
```

`/tmp/wiki.dot` begins:

```dot
digraph wiki {
  "algorithm-comparison" [weight=0.025630];
  "arithmetic-gates" [weight=0.041796];
  "bb84-qkd" [weight=0.014462];
  ...
  "algorithm-comparison" -> "arithmetic-gates" [weight=1.000000];
```

Render it immediately:

```bash
dot -Tpng /tmp/wiki.dot -o /tmp/wiki.png
open /tmp/wiki.png
```

### 4. Filter out weak edges

Raw export edges are weight `1.0` by default, so the historical `--min-edge 0.005`
threshold keeps all real wikilinks. Raise above `1` to drop edges, or use
`graph --min-edge` on the base-link kernel when visualising:

```bash
wikigraph export ~/quantum-go/wiki --format json -o /tmp/wiki --min-edge 0.05
```

Lower it (or set to `0`) to include everything:

```bash
wikigraph export ~/quantum-go/wiki --format dot -o /tmp/wiki --min-edge 0
```

**Action:** Confirm the link count changes — `cat /tmp/wiki.json | jq '.links | length'` before and after adjusting `--min-edge` to see the difference.

### 5. Exclude meta-pages

`index`, `log`, and `AGENTS` are excluded by default. Add your own:

```bash
wikigraph export ~/quantum-go/wiki --format json -o /tmp/wiki --exclude index --exclude log --exclude README
```

**Action:** Open `/tmp/wiki.json` and confirm the excluded slugs do not appear in the
`nodes` array.

## Verification

Run all three formats against the canonical example wiki:

```bash
wikigraph export ~/quantum-go/wiki --format json -o /tmp/wiki
wikigraph export ~/quantum-go/wiki --format csv  -o /tmp/wiki
wikigraph export ~/quantum-go/wiki --format dot  -o /tmp/wiki
```

Expected stderr for each run:

```
Pages: 36
Written: /tmp/wiki.json

Pages: 36
Written: /tmp/wiki_nodes.csv
Written: /tmp/wiki_edges.csv

Pages: 36
Written: /tmp/wiki.dot
```

Spot-check the files:

```bash
# JSON — two top nodes by pi
cat /tmp/wiki.json | jq '.nodes[:2]'

# CSV — header + first three node rows
head -4 /tmp/wiki_nodes.csv

# DOT — opening stanza
head -5 /tmp/wiki.dot
```

If `Pages: 36` appears on stderr for each run and all four output files exist (`wiki.json`, `wiki_nodes.csv`, `wiki_edges.csv`, `wiki.dot`), the export succeeded.

## Troubleshooting

| Symptom                              | Cause                                          | Fix                                              |
| ------------------------------------ | ---------------------------------------------- | ------------------------------------------------ |
| `unknown format "xlsx"`              | Unsupported format                             | Use `json`, `csv`, or `dot`                      |
| Output file has 0 edges              | `--min-edge` too high for your wiki            | Lower with `--min-edge 0.001` or `--min-edge 0`  |
| A page you expect is missing         | It has no `.md` extension, or it's in `--exclude` | Check filename and exclude list              |
| `Pages: 1` on stderr                 | Only one page found — likely wrong directory   | Double-check the path to your wiki               |
| JSON contains `\uXXXX` escapes or garbled node IDs | Page filenames contain non-ASCII characters | Rename files to ASCII slugs, or set `LANG=en_US.UTF-8` before running |

## Sources

- [`cmd_export.go`](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_export.go)
- [`wiki.go`](https://github.com/stephen-mcelhose/wikigraph/blob/main/wiki.go)
- [[architecture]] — full data-flow pipeline from wikilinks to Markov output
