---
type: concept
title: Wikilink Markov Model
description: How wikilinks become raw adjacency and a teleporting PageRank kernel (π, MFPT) in wikigraph.
tags: [markov-chain, transition-matrix, adjacency, parser, pagerank]
timestamp: 2026-08-15T21:30:00Z
---

# Wikilink Markov Model

## Overview

`wikigraph` models a markdown wiki as a discrete-time Markov chain on $n$
pages. **Structure** is the raw digraph of `[[wikilinks]]`. **Math** (π, MFPT,
entropy) uses a teleporting / PageRank kernel over that digraph.

## Key Properties

### Pipeline steps (`wiki.go`)

1. **`loadPages`**: Discovers markdown files, sorts slugs, builds
   `slug → index` (path-relative slugs in recursive mode; see
   [[adr-010-path-relative-slugs]]).
2. **`buildAdjacency`**: Parses `[[slug]]` / `[[slug|alias]]` (and optional
   relative Markdown links). Real links only; [[sink-page|sink]] rows stay
   zero.
3. **`NewTeleportingKernelFromAdj`**: Builds
   $T_{ij} = \alpha v_j + (1-\alpha) A_{ij}/\mathrm{rowsum}_i$ (sinks → $v$).
   Default $\alpha = 0.15$; $v$ uniform unless `--seed` (PPR).

### What is raw vs math

| Artifact | Source |
| -------- | ------ |
| Edge count, export edges, raw SCCs | Raw adjacency $A$ |
| π, entropy, commute, MFPT / `goal` | Teleporting $T$ |
| `graph` node size | Stationary of $T$ via `NodeMass` |
| `graph` edges | Base-link kernel ($\alpha=0$) + `MinEdge` |

See [[adr-012-teleporting-pagerank-default]].

## Related Concepts

- [[architecture]] — Data-flow pipeline
- [[sink-page]] — Structural sinks vs restart handling
- [[catrace]] — `NewTeleportingKernelFromAdj`, `NodeMass`
- [[stationary-distribution]] — PageRank π
- [[teleportation-ergodicity]] — Why α-damping guarantees ergodicity
- [[kernel-identifiability]] — Recovering structure from outputs

## Sources

- `wiki.go`
- [[adr-012-teleporting-pagerank-default]]
