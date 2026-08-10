---
type: concept
title: Path Strategy Sequence Metric
description: Path-finding strategy using negative log-likelihood transition weights and Dijkstra search to determine sequential reading order through wiki goals.
tags: [path-strategy, dijkstra, transition-probability, log-likelihood, curriculum]
timestamp: 2026-08-10T01:15:00Z
---

# Path Strategy Sequence Metric

## Overview

`wikigraph goal --strategy path` generates a sequential stepping-stone reading path connecting goal pages in the exact order specified by the user (`--goal A --goal B --goal C`).

Unlike distance-based metric subgraphs, `path` surfaces the single most probable sequential chain of wikilinks connecting each adjacent goal pair $g_m \rightarrow g_{m+1}$.

## Mathematical Formulation

Let $P_{ij}$ be the transition probability from page $i$ to page $j$ in the [[markov-model|Markov transition matrix]].

To find the path maximizing cumulative transition probability:

$$\max \prod_{(i,j) \in \text{Path}} P_{ij}$$

We transform edge weights into additive non-negative distances using the **negative log-likelihood**:

$$w_{ij} = -\log P_{ij} \quad (P_{ij} > 0)$$

Since $0 < P_{ij} \le 1$, $w_{ij} \ge 0$. Minimizing total path weight $\sum w_{ij}$ via **Dijkstra's Shortest Path Algorithm** is mathematically equivalent to maximizing the product of transition probabilities along the route.

## Subgraph Expansion

If the core path sequence contains $K < N$ pages (where $N = \text{--top}$):
1. All path sequence nodes are included first.
2. Remaining $N - K$ slots are filled by selecting 1-hop outgoing ($P_{\text{path}, j}$) and incoming ($P_{i, \text{path}}$) neighbor nodes ranked by transition probability.

## Related Concepts

- [[goal]] — `wikigraph goal` command documentation
- [[markov-model]] — Foundations of transition probabilities $P_{ij}$
- [[bottleneck-centrality]] — Alternative strategy for identifying chokepoints
- [[absorbing-markov-chain]] — State absorption mechanics

## Sources

- Cormen, T. H., Leiserson, C. E., Rivest, R. L., & Stein, C. (2009). *Introduction to Algorithms* (3rd ed.). MIT Press. ISBN: 978-0-262-03384-8.
