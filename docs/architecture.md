---
type: concept
title: wikigraph Architecture
description: Seven Go source files, the catrace dependency, and the data-flow pipeline from markdown wikilinks to Markov kernel to visual or numeric output.
tags: [architecture, go, cobra, catrace, pipeline]
timestamp: 2026-08-09T07:31:56Z
---

# wikigraph Architecture

## Overview

wikigraph is ~1,567 lines of Go across seven files. The design is deliberately thin: the tool's job is to translate a wiki's `[[wikilinks]]` into a Markov kernel and then hand off to the `catrace` library for linear algebra and graph analysis. The `union` and `intersection` goal strategies delegate fully to catrace; `path` (Dijkstra on the raw transition matrix) and `bottleneck` (fundamental matrix via `gonum/mat`) implement their own math directly — both are candidates for future catrace APIs where the math is stochastic-matrix-specific.

### File map

| File | Role |
| --- | --- |
| `main.go` | Cobra root command, `--exclude` and `-r/--recursive` persistent flags, subcommand wiring |
| `wiki.go` | `loadPages`, `buildAdjacency`, `buildKernel` — the translation layer |
| `wiki_test.go` | Unit tests for flat/recursive page loading, duplicate slug validation, and kernel construction |
| `cmd_graph.go` | `graph` subcommand — renders full kernel as interactive HTML |
| `cmd_analyze.go` | `analyze` subcommand — prints health report with six sections |
| `cmd_goal.go` | `goal` subcommand — MFPT ranking + trace kernel; four strategies: `union`, `intersection`, `path`, `bottleneck` (see [[adr-007-subgraph-partitioning-and-path-strategies]]) |
| `cmd_export.go` | `export` subcommand — serialises kernel as JSON, CSV, or DOT |
| `sed.go` | `applySed` — applies arbitrary sed expressions to HTML output |

## Key Properties

### The catrace dependency

All Markov maths lives in [[catrace|github.com/stephen-mcelhose/catrace]]. wikigraph never implements its own linear algebra. The `catrace.Kernel` struct exposes:

- `P` (`mat.Dense`) — the n×n row-stochastic transition matrix
- `Stationary(tol, maxIter)` — power iteration to find [[stationary-distribution|π]]
- `Classes(tol)` — Kosaraju SCC decomposition → [[communicating-classes|recurrent vs transient sets]]
- `EntropyRate(base)` — [[entropy-rate|H]] = −Σᵢ πᵢ Σⱼ Pᵢⱼ log Pᵢⱼ
- `MeanFirstPassage(i, j)` — MFPT via fundamental matrix
- `CommuteTime(i, j)` — MFPT(i,j) + MFPT(j,i)
- `Trace(subset, tol)` — effective kernel on a subset of states
- `ToHTML(opts)` — D3-based force-directed graph

### Data-flow pipeline

```
docs/*.md
  │
  ▼ loadPages (wiki.go)
sorted []string slugs + slug→index map
  │
  ▼ buildAdjacency (wiki.go)
mat.Dense  n×n adjacency (0/1 + [[sink-page|sink]] teleportation)
  │
  ▼ catrace.NewRandomWalkKernel
catrace.Kernel  (P is now row-stochastic)
  │
  ├──▶ cmd_graph.go   → k.ToHTML → wiki_graph.html
  ├──▶ cmd_analyze.go → Stationary + Classes + CommuteTime → text report
  ├──▶ cmd_goal.go    → MeanFirstPassage + Trace + ToHTML → goal_graph.html
  └──▶ cmd_export.go  → Stationary + Classes → JSON / CSV / DOT
```

### Persistent --exclude flag

`main.go` declares `--exclude` (`-e`) as a `PersistentFlag` on the root command. Cobra makes it available to all subcommands. The default value is `["index", "log", "AGENTS"]` — the three meta-files used by `llm-wiki`.

## Related Concepts

- [[markov-model]] — How wikilinks are parsed into a row-stochastic Markov kernel
- [[catrace]] — Go package providing Markov chain linear algebra
- [[mfpt]] — Mean first passage time calculations for goal paths
- [[commute-time]] — Symmetric distance metric used for link suggestions
- [[adr-006-recursive-vault-traversal]] — Rationale for recursive vault traversal and slug resolution
- [[adr-007-subgraph-partitioning-and-path-strategies]] — Subgraph strategy partitioning decisions for goal subcommand
- [[adr-008-prototype-math-strategies]] — Accepted deviation: path and bottleneck implement math directly pending catrace APIs
- [[how-to-docs-plan]] — Documentation initiative driving subcommand interface

## Sources

- `main.go`
- `wiki.go`
- `cmd_graph.go`
- `cmd_analyze.go`
- `cmd_goal.go`
- `cmd_export.go`
- `sed.go`
