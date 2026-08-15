---
type: concept
title: The catrace Go Package
description: Shared Markov engine for wikigraph — kernels, PageRank/PPR, MFPT, viz — and the PDA agent formalism those tools also serve.
tags: [catrace, go, markov-chain, linear-algebra, gonum, pda, agents]
timestamp: 2026-08-15T21:38:39Z
---

# The catrace Go Package

## Overview

`catrace` (`github.com/stephen-mcelhose/catrace`) is an open-source Go library for
finite-state Markov models: analysis (stationary / PageRank, MFPT, entropy,
classes, trace) and D3 force-directed rendering. `wikigraph` imports `catrace`
and delegates linear algebra to its `Kernel` type.

The same library implements the **PDA triplet** (Perception / Decision / Action)
for single- and multi-agent models. Wikigraph is the knowledge-graph *slice* of
that craft — see [[knowledge-graph-to-pda-agents]] for the leap from wiki walks
to agents and networks of agents.

### Repository and version

- **Import path:** `github.com/stephen-mcelhose/catrace`
- **Primary types:** `catrace.Kernel`, `catrace.Agent` (PDA)
- **Dependencies:** `gonum.org/v1/gonum/mat`

## Key Properties

### The Kernel struct

A `Kernel` wraps an $n \times n$ row-stochastic transition matrix $P$ and an array of state labels (slugs).

### Analysis methods used by wikigraph

- **`NewTeleportingKernelFromAdj(adj, restart, alpha, names)`**: Builds the
  PageRank / teleporting kernel from raw adjacency (sinks → restart). Primary
  entry point after [[adr-012-teleporting-pagerank-default]]. Restart + $\alpha$
  are the same **intent** dials used for goal-directed agents in catrace.
- **`Stationary(tol, maxIter)`**: Computes π via power iteration. Used for
  PageRank centrality / orphans ([[stationary-distribution]]).
- **`Classes(tol)`**: SCCs on math $P$ (usually one class under α-damping).
  `analyze` reports **raw** digraph SCCs instead ([[communicating-classes]]).
- **`EntropyRate(base)`**: Entropy rate of the teleporting walk ([[entropy-rate]]).
- **`MeanFirstPassage(i, j)`**: MFPT via fundamental matrix. Used by `goal` ([[mfpt]]).
- **`CommuteTime(i, j)`**: $K(i,j) = M(i,j) + M(j,i)$. Used by `analyze` suggestions ([[commute-time]]).
- **`Trace(subset, tol)`**: Effective kernel on a subset. Used by `goal`.
- **`ToHTML(opts)`**: D3 force graph; `VisualiseOptions.NodeMass` sizes nodes
  independently of drawn edges. Used by `graph` and `goal`.

### PDA (beyond wikigraph’s CLI)

In catrace, an `Agent` is three kernels **P** ($W\to X$), **D** ($X\to G$),
**A** ($G\to W$). Compositions $Q = DAP$, $S = APD$, $W = PDA$ are the closed
perceive→decide→act loops. Personalized PageRank on those kernels is
goal-directed long-run mass — the agent analogue of `wikigraph graph --seed`.

You do not need to build a full PDA agent to get value from wikigraph; you *do*
reuse the same measurements when you step up to agents. That bridge is
[[knowledge-graph-to-pda-agents]].

## Related Concepts

- [[knowledge-graph-to-pda-agents]] — Shared operators: wiki walk ↔ PDA / multi-agent
- [[architecture]] — How `wikigraph` wires `catrace` into subcommand handlers
- [[markov-model]] — Raw adjacency + teleporting kernel
- [[adr-012-teleporting-pagerank-default]] — Current default math
- [[stationary-distribution]] — Power iteration / PageRank π
- [[teleportation-ergodicity]] — α-damping and intent
- [[communicating-classes]] — Strongly connected components
- [[mfpt]] — Fundamental matrix calculation
- [[kernel-identifiability]] — What each analysis output reveals about $P$

## Sources

- [catrace repository](https://github.com/stephen-mcelhose/catrace)
- catrace wiki: PDA Triplet Model; Personalized PageRank and Agent Modeling
- [[knowledge-graph-to-pda-agents]]
- `architecture.md`
