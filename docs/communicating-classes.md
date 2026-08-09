---
type: concept
title: Communicating Classes and Recurrent Sets
description: Maximal mutually-reachable sets of wiki pages — how Kosaraju's SCC algorithm partitions a wiki graph into recurrent and transient components.
tags: [markov-chain, graph-theory, scc, recurrent, transient]
timestamp: 2026-08-09T08:35:00Z
---

# Communicating Classes and Recurrent Sets

## Overview

A **communicating class** is a maximal set of states (pages) where every page in the set can reach every other page in the set following directed wikilinks. In graph theory, communicating classes correspond to **strongly connected components (SCCs)** on non-zero transitions ($P_{ij} > 0$).

## Key Properties

### How it is computed

`catrace` computes communicating classes via Kosaraju's two-pass DFS algorithm on the adjacency graph:

1. **First DFS pass**: Computes post-order traversal on $P$.
2. **Second DFS pass**: Runs DFS on transpose $P^T$ in reverse post-order to extract components.
3. **Classification**:
   - **Recurrent class**: A component with no outgoing transitions to other components. Once a random walker enters, it never leaves.
   - **Transient class**: A component with outgoing transitions to another component. A random walker will eventually leave and never return.

### Reading the output

In `wikigraph analyze docs/`:

```
=== Communicating classes ===
Class 1 (recurrent): 23 page(s)
  analyze
  architecture
  ...
```

A healthy wiki forms **a single recurrent class** containing all pages.

### Common failure modes

- **Disconnected islands**: Subgraphs linked internally but unreachable from the main wiki.
- **Trap subgraphs / Sink pages**: Pages or clusters with no path back to main content.

## Related Concepts

- [[recurrent-class]] — Detailed distinction between recurrent and transient sets
- [[sink-page]] — Pages with zero outgoing links that create dead ends
- [[stationary-distribution]] — How transient pages lose long-run probability $\pi = 0$
- [[catrace]] — Kosaraju SCC algorithm implementation

## Sources

- Feller, W. (1968). *An Introduction to Probability Theory and Its Applications*.
- `catrace` package documentation
