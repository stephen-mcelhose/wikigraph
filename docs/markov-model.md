---
type: concept
title: Wikilink Markov Model
description: How wikilinks are extracted and converted into an n×n row-stochastic transition matrix P and catrace.Kernel struct.
tags: [markov-chain, transition-matrix, adjacency, parser]
timestamp: 2026-08-09T08:35:00Z
---

# Wikilink Markov Model

## Overview

`wikigraph` models a markdown wiki as a discrete-time Markov chain $X_t$ on $n$ states, where states are pages and transitions correspond to a reader clicking `[[wikilinks]]`.

## Key Properties

### Pipeline steps (`wiki.go`)

1. **`loadPages`**: Discovers markdown files in `docs/`, sorts slugs alphabetically, builds `slug -> index` map.
2. **`buildAdjacency`**: Parses `[[slug]]` and `[[slug|alias]]` wikilinks from page bodies using regular expressions (`\[\[([A-Za-z][A-Za-z0-9-]*)(?:\|[^\]]+)?\]\]`).
3. **Teleportation for sinks**: If a page has zero outgoing links ([[sink-page]]), a uniform row $P_{ij} = 1/n$ is inserted.
4. **Row-stochastic normalization**: Row $i$ with $k_i > 0$ outgoing links gets uniform probability $P_{ij} = 1/k_i$ for linked targets $j$.

### Transition Matrix $P$

The resulting matrix $P$ satisfies:

$$\sum_{j=1}^n P_{ij} = 1 \quad \forall i$$

## Related Concepts

- [[architecture]] — Data-flow pipeline
- [[sink-page]] — Sink handling and uniform teleportation
- [[catrace]] — `catrace.Kernel` wrapper
- [[stationary-distribution]] — Power iteration on $P$

## Sources

- `wiki.go`
- Norris, J. R. (1998). *Markov Chains*. Cambridge University Press.
