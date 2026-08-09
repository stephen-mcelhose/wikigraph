---
type: concept
title: wikigraph Architecture
description: Seven Go source files, the catrace dependency, and the data-flow pipeline from markdown wikilinks to Markov kernel to visual or numeric output.
tags: [architecture, go, cobra, catrace, pipeline]
timestamp: 2026-08-09T07:31:56Z
---

# wikigraph Architecture

wikigraph is ~850 lines of Go across seven files. The design is deliberately
thin: the tool's job is to translate a wiki's `[[wikilinks]]` into a Markov
kernel and then hand off to the `catrace` library for all the maths.

## File map

| File             | Role                                                                   |
| ---------------- | ---------------------------------------------------------------------- |
| `main.go`        | Cobra root command, `--exclude` persistent flag, subcommand wiring     |
| `wiki.go`        | `loadPages`, `buildAdjacency`, `buildKernel` — the translation layer  |
| `cmd_graph.go`   | `graph` subcommand — renders full kernel as interactive HTML           |
| `cmd_analyze.go` | `analyze` subcommand — prints health report with six sections          |
| `cmd_goal.go`    | `goal` subcommand — MFPT ranking + trace kernel for a target subgraph  |
| `cmd_export.go`  | `export` subcommand — serialises kernel as JSON, CSV, or DOT           |
| `sed.go`         | `applySed` — applies arbitrary sed expressions to HTML output          |

## The catrace dependency

All Markov maths lives in [[catrace|`github.com/stephen-mcelhose/catrace`]]. wikigraph
never implements its own linear algebra. The `catrace.Kernel` struct exposes:

- `P` (`mat.Dense`) — the n×n row-stochastic transition matrix
- `Stationary(tol, maxIter)` — power iteration to find [[stationary-distribution|π]]
- `Classes(tol)` — Kosaraju SCC decomposition → [[communicating-classes|recurrent vs transient sets]]
- `EntropyRate(base)` — [[entropy-rate|H]] = −Σᵢ πᵢ Σⱼ Pᵢⱼ log Pᵢⱼ
- `MeanFirstPassage(i, j)` — MFPT via fundamental matrix
- `CommuteTime(i, j)` — MFPT(i,j) + MFPT(j,i)
- `Trace(subset, tol)` — effective kernel on a subset of states
- `ToHTML(opts)` — D3-based force-directed graph

See [[markov-model]] for how the kernel is built; see [[mfpt]] for how MFPT
and commute time are used.

## Data-flow pipeline

```
docs/*.md
  │
  ▼ loadPages (wiki.go)
sorted []string slugs + slug→index map
  │
  ▼ buildAdjacency (wiki.go)
mat.Dense  n×n adjacency (0/1 + sink teleportation)
  │
  ▼ catrace.NewRandomWalkKernel
catrace.Kernel  (P is now row-stochastic)
  │
  ├──▶ cmd_graph.go   → k.ToHTML → wiki_graph.html
  ├──▶ cmd_analyze.go → Stationary + Classes + CommuteTime → text report
  ├──▶ cmd_goal.go    → MeanFirstPassage + Trace + ToHTML → goal_graph.html
  └──▶ cmd_export.go  → Stationary + Classes → JSON / CSV / DOT
```

## Persistent --exclude flag

`main.go` declares `--exclude` (`-e`) as a `PersistentFlag` on the root
command. Cobra makes it available to all subcommands. The default value is
`["index", "log", "AGENTS"]` — the three meta-files used by `llm-wiki`.
When a user passes `-e`, the default is replaced entirely, so they must
re-specify the defaults alongside any custom exclusions.

## Wiki conventions

The tool operates on flat wiki directories — all pages at one level, slugs
are bare filenames without path prefixes (e.g., `analyze`, not `how-to/analyze`).
The rationale is documented in [[adr-002-slug-resolution]].
The how-to documentation series that drove the flat structure was planned in
[[how-to-docs-plan]].

## Sources

- `main.go`
- `wiki.go`
- `cmd_graph.go`
- `cmd_analyze.go`
- `cmd_goal.go`
- `cmd_export.go`
- `sed.go`
