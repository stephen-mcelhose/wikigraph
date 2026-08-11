---
type: concept
title: Bottleneck Centrality Metric
description: Random-walk betweenness metric used by wikigraph goal --strategy bottleneck to identify gatekeeper pages controlling flow between goal pairs.
tags: [bottleneck, betweenness, random-walk, markov-chain, fundamental-matrix, effective-resistance]
timestamp: 2026-08-10T01:15:00Z
---

# Bottleneck Centrality Metric

> [!WARNING]
> **Prototype strategy.** `--strategy bottleneck` computes the fundamental matrix $N = (I-Q)^{-1}$ directly in `cmd_goal.go` using `gonum/mat`, bypassing the catrace library. This deviates from wikigraph's standard practice of delegating all Markov chain mathematics to catrace. Accepted under [[adr-008-prototype-math-strategies|ADR-008]] while the strategy is validated; a catrace API is the intended end-state.

## Overview

In `wikigraph goal --strategy bottleneck`, the **bottleneck score** $B_{s,t}(k)$ measures how critical page $k$ is as a gateway/chokepoint when navigating from source goal $s$ to target goal $t$.

It is based on **Random Walk Betweenness Centrality** (or effective current flow in electrical network theory).

## Mathematical Definition

For a pair of goal nodes $(s, t)$, we construct an [[absorbing-markov-chain]] where the target goal $t$ is designated as an absorbing state (once the random walk hits $t$, it stays there).

Let $Q$ be the sub-transition matrix for all non-$t$ transient nodes. The **fundamental matrix** $N$ is:

$$N = (I - Q)^{-1}$$

The entry $N_{s, k}$ gives the expected number of times a random walk starting at $s$ visits node $k$ before being absorbed at target $t$.

### Multi-Goal Pair Aggregation

When multiple goals $\{g_1, g_2, \dots, g_M\}$ are provided, the total bottleneck score for intermediate page $k$ is the sum across all ordered goal pairs:

$$\text{BottleneckScore}(k) = \sum_{a \neq b} N_{g_a, k}^{(g_b)}$$

where $N^{(g_b)}$ is the fundamental matrix with goal $g_b$ set as the absorbing state.

## Interpretation in Subgraph Selection

- **High Bottleneck Score**: Page $k$ sits on almost all structural pathways between $s$ and $t$. If a user wants to travel conceptually from $s$ to $t$, mastering page $k$ is an unavoidable prerequisite.
- **Low Bottleneck Score**: Page $k$ is either off the path entirely or sits on one of many parallel redundant pathways.

## Related Concepts

- [[absorbing-markov-chain]] — Foundational state absorption theory and fundamental matrix $N$
- [[commute-time]] — Symmetric resistance distance metric
- [[mfpt]] — Mean First Passage Time foundation
- [[goal]] — Command surface utilizing bottleneck scoring
- [[catrace]] — Markov chain library (trace kernel and MFPT; fundamental matrix is computed directly via `gonum/mat` in `cmd_goal.go`)
- [[adr-007-subgraph-partitioning-and-path-strategies]] — Decision record for the four goal strategies

## Sources

- Newman, M. E. J. (2005). *A measure of betweenness centrality based on random walks*. Social Networks, 27(1), 39–54. [DOI: 10.1016/j.socnet.2004.11.009](https://doi.org/10.1016/j.socnet.2004.11.009) | [arXiv:cond-mat/0309045](https://arxiv.org/abs/cond-mat/0309045)
