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

That thin slice is intentional pedagogy as well as product: the same kernels, π / PPR, MFPT, and entropy you measure on a wiki are the craft for PDA agents and multi-agent networks in catrace — see [[knowledge-graph-to-pda-agents]].

### File map

| File | Role |
| --- | --- |
| `main.go` | Cobra root command, persistent flags (`--exclude`, `-r`, `--relative-links`, `--alpha`, `--seed`), subcommand wiring |
| `wiki.go` | `loadPages`, `buildAdjacency`, `buildKernel` — raw adj + teleporting kernel |
| `digraph.go` | Raw digraph SCCs (`rawSCCs`) and `rawEdgeCount` for authoring signals |
| `wiki_test.go` | Unit tests for flat/recursive page loading, duplicate slug validation, and kernel construction |
| `cmd_graph.go` | `graph` subcommand — PageRank `NodeMass` + base-link HTML |
| `cmd_analyze.go` | `analyze` subcommand — health report (raw structure + PageRank π) |
| `cmd_goal.go` | `goal` subcommand — MFPT ranking + trace kernel; four strategies: `union`, `intersection`, `path`, `bottleneck` (see [[adr-007-subgraph-partitioning-and-path-strategies]]) |
| `cmd_export.go` | `export` subcommand — raw edges + PageRank π as JSON, CSV, or DOT |
| `sed.go` | `applySed` — applies arbitrary sed expressions to HTML output |

## Key Properties

### The catrace dependency

All Markov maths lives in [[catrace|github.com/stephen-mcelhose/catrace]]. wikigraph never implements its own linear algebra. The `catrace.Kernel` struct exposes:

- `NewTeleportingKernelFromAdj`, `NodeMass`, PPR helpers
- `Stationary(tol, maxIter)` — power iteration to find [[stationary-distribution|π]]
- `Classes(tol)` — Kosaraju SCC on math $P$ (usually one class under α-damping; analyze uses raw SCCs instead)
- `EntropyRate(base)` — [[entropy-rate|H]] = −Σᵢ πᵢ Σⱼ Pᵢⱼ log Pᵢⱼ
- `MeanFirstPassage(i, j)` — MFPT via fundamental matrix
- `CommuteTime(i, j)` — MFPT(i,j) + MFPT(j,i)
- `Trace(subset, tol)` — effective kernel on a subset of states
- `ToHTML(opts)` — D3-based force-directed graph (`NodeMass` optional)

### Data-flow pipeline

```
docs/*.md
  │
  ▼ loadPages (wiki.go)
sorted []string slugs + slug→index map
  │
  ▼ buildAdjacency (wiki.go)
mat.Dense  n×n raw adjacency (0/1; [[sink-page|sink]] rows stay zero)
  │
  ▼ catrace.NewTeleportingKernelFromAdj (α, restart)
catrace.Kernel  (teleporting / PageRank P)
  │
  ├──▶ cmd_graph.go   → NodeMass(π) + α=0 base edges → wiki_graph.html
  ├──▶ cmd_analyze.go → PageRank π + raw SCCs/edges + CommuteTime → text report
  ├──▶ cmd_goal.go    → MeanFirstPassage + Trace + ToHTML → goal_graph.html
  └──▶ cmd_export.go  → raw edges + PageRank π → JSON / CSV / DOT
```

See [[adr-012-teleporting-pagerank-default]].

### Persistent flags

`main.go` declares persistent flags on the root command (available to all
subcommands). Defaults include `--exclude index,log,AGENTS`, `--alpha 0.15`,
and optional `--seed` for Personalized PageRank.

## Related Concepts

- [[markov-model]] — How wikilinks are parsed into a row-stochastic Markov kernel
- [[catrace]] — Go package providing Markov chain linear algebra
- [[mfpt]] — Mean first passage time calculations for goal paths
- [[commute-time]] — Symmetric distance metric used for link suggestions
- [[adr-006-recursive-vault-traversal]] — Rationale for recursive vault traversal and slug resolution
- [[adr-010-path-relative-slugs]] — Path-relative slug scheme and lenient wikilink fallback for structured folder wikis
- [[adr-007-subgraph-partitioning-and-path-strategies]] — Subgraph strategy partitioning decisions for goal subcommand
- [[adr-008-prototype-math-strategies]] — Accepted deviation: path and bottleneck implement math directly pending catrace APIs
- [[adr-012-teleporting-pagerank-default]] — Default PageRank / teleporting kernel
- [[knowledge-graph-to-pda-agents]] — Shared operators from wiki walks to PDA / multi-agent kernels
- [[how-to-docs-plan]] — Documentation initiative driving subcommand interface
- [[adr-009-wiki-gen-make-vs-buy]] — Decision: nx-to-wiki Python converter for named NetworkX benchmark graphs
- [[graph-topologies]] — Named topology catalog (barbell, caveman, WS, SBM, …) used in benchmark experiments
- [[graph-models]] — Random graph model catalog (ER, BA, LFR) used in benchmark experiments

## Sources

- `main.go`
- `wiki.go`
- `digraph.go`
- `cmd_graph.go`
- `cmd_analyze.go`
- `cmd_goal.go`
- `cmd_export.go`
- `sed.go`
- [[adr-012-teleporting-pagerank-default]]
