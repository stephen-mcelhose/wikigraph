---
type: concept
title: The catrace Go Package
description: Overview of github.com/stephen-mcelhose/catrace — the Go library providing all Markov chain mathematics, graph algorithms, and D3 visualization for wikigraph.
tags: [catrace, go, markov-chain, linear-algebra, gonum]
timestamp: 2026-08-09T08:35:00Z
---

# The catrace Go Package

## Overview

`catrace` (`github.com/stephen-mcelhose/catrace`) is an open-source Go library providing Markov chain analysis, graph decomposition, and D3 force-directed graph rendering. `wikigraph` imports `catrace` and delegates all linear algebra to its `Kernel` struct.

### Repository and version

- **Import path:** `github.com/stephen-mcelhose/catrace`
- **Primary type:** `catrace.Kernel`
- **Dependencies:** `gonum.org/v1/gonum/mat`

## Key Properties

### The Kernel struct

A `Kernel` wraps an $n \times n$ row-stochastic transition matrix $P$ and an array of state labels (slugs).

### Analysis methods used by wikigraph

- **`Stationary(tol, maxIter)`**: Computes stationary distribution $\pi$ via power iteration ($P^T \pi = \pi$). Used by `analyze` for centrality and orphan detection ([[stationary-distribution]]).
- **`Classes(tol)`**: Computes communicating classes via Kosaraju's SCC algorithm on non-zero transitions ($P_{ij} > 0$). Used by `analyze` ([[communicating-classes]]).
- **`EntropyRate(base)`**: Computes entropy rate $H = -\sum_i \pi_i \sum_j P_{ij} \log P_{ij}$. Used by `analyze` ([[entropy-rate]]).
- **`MeanFirstPassage(i, j)`**: Computes expected steps from state $i$ to state $j$ using the fundamental matrix $Z = (I - P + W)^{-1}$. Used by `goal` ([[mfpt]]).
- **`CommuteTime(i, j)`**: Computes $K(i,j) = M(i,j) + M(j,i)$. Used by `analyze` for missing link suggestions ([[commute-time]]).
- **`Trace(subset, tol)`**: Computes effective trace kernel $P_S$ on a state subset $S$, accounting for paths through excluded states $S^c$. Used by `goal` subgraphs.
- **`ToHTML(opts)`**: Renders force-directed graph HTML using D3.js. Used by `graph` and `goal`.

## Related Concepts

- [[architecture]] — How `wikigraph` wires `catrace` into subcommand handlers
- [[markov-model]] — Construction of transition matrix $P$ from wikilinks
- [[stationary-distribution]] — Power iteration implementation
- [[communicating-classes]] — Strongly connected components
- [[mfpt]] — Fundamental matrix calculation

## Sources

- [catrace repository](https://github.com/stephen-mcelhose/catrace)
- `architecture.md`
