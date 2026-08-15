---
type: concept
title: Sink Pages and Teleportation
description: Pages with zero outgoing wikilinks — structural smell on raw adjacency; handled in the teleporting PageRank kernel via the restart distribution.
tags: [sink-page, teleportation, markov-chain, dead-end, pagerank]
timestamp: 2026-08-15T21:30:00Z
---

# Sink Pages and Teleportation

## Overview

A **sink page** is a page with zero outgoing `[[wikilinks]]` (and no relative-link
edges when that mode is on). On the **raw** adjacency matrix the sink row is all
zeros. That is an authoring smell — a dead end for readers — and `analyze`
reports it in the Sink pages section.

Ergodicity for π / MFPT does **not** come from pre-filling that row with fake
`1/n` wikilinks. It comes from the **teleporting kernel** (PageRank α-damping):
sink rows collapse to the restart distribution. See [[adr-012-teleporting-pagerank-default]].

## Key Properties

### Detection (raw graph)

`buildAdjacencyWithOpts` leaves sink rows as zeros and appends the slug to the
sink list. Overview edge counts and export edges ignore sinks' nonexistent
out-links.

### Math (teleporting kernel)

With `NewTeleportingKernelFromAdj(adj, restart, alpha, …)`:

- Non-sink row: $T_{ij} = \alpha v_j + (1-\alpha)\, A_{ij}/\mathrm{rowsum}_i$
- Sink row: $T_{ij} = v_j$ (full restart; usually uniform → $1/n$)

Default `--alpha` is `0.15` (teleport probability). `--seed` personalizes $v$.

### How to fix a sink

In `wikigraph analyze docs/`:

```
=== Sink pages (no outgoing links) ===
  some-page → add outgoing links
```

Authors fix sinks by adding relevant outgoing `[[wikilinks]]`.

### Visualization

`graph` draws base-link edges (α=0 display kernel + `MinEdge`) and sizes nodes
with PageRank `NodeMass`. Sink restart mass is not drawn as a star of fake
wikilinks.

## Related Concepts

- [[adr-012-teleporting-pagerank-default]] — current decision
- [[teleportation-ergodicity]] — why teleportation is needed
- [[markov-model]] — raw adj + teleporting kernel pipeline
- [[stationary-distribution]] — PageRank π
- [[adr-011-sink-teleportation-vs-pagerank-damping]] — superseded sink-only ADR

## Sources

- `wiki.go` — `buildAdjacencyWithOpts`, `buildKernelWithOpts`
- [[adr-012-teleporting-pagerank-default]]
- catrace `NewTeleportingKernelFromAdj`
