---
type: decision
title: Subgraph Partitioning and Path Strategy Algorithms for Goal Subcommand
description: Decision to support multiple goal strategies (union, intersection, path, bottleneck) and the specific design for path traversal ordering and expansion.
resource: cmd_goal.go
tags: [adr, goal, path-strategy, graph-partitioning, dijkstra, mfpt]
timestamp: 2026-08-10T01:00:00Z
status: accepted
---

# ADR-007 — Subgraph Partitioning and Path Strategy Algorithms for Goal Subcommand

## Context

`wikigraph goal` initially used a single hardcoded scoring function: $\text{Score}(i) = \min_{g} \text{MFPT}(i, g)$ (`union`). While effective at finding the local union of page neighborhoods around goal nodes, it failed to address other primary structural questions:

1. **Prerequisite Intersection (`intersection`)**: Finding shared pages that sit close to *all* specified goals simultaneously.
2. **Sequential Stepping-Stone Traversal (`path`)**: Finding the optimal sequential reading order from Goal A to Goal B (and through additional goal milestones in sequence).
3. **Chokepoint Identification (`bottleneck`)**: Finding gatekeeper pages that control transition paths between goal pairs.

Specifically for the **`path` strategy**, three critical design challenges arose:
- How to evaluate multi-goal inputs (e.g. `--goal A --goal B --goal C`).
- How to map transition probabilities $P_{ij}$ to shortest-path edge weights.
- How to expand the subgraph when the path length $K$ is smaller than `--top N`.

## Options Considered

### For Multi-Goal `path` Evaluation
- **Option A (Pairwise All-to-All)**: Evaluate all $O(M^2)$ goal combinations ($A \leftrightarrow B, B \leftrightarrow C, A \leftrightarrow C$) and combine all paths.
  *Rejected*: Loses the directional, sequential intent of ordered goal flags and creates cluttered subgraphs.
- **Option B (Ordered Sequential Chain $A \rightarrow B \rightarrow C$)**: Treat `--goal` flag order as an explicit sequential learning pipeline, computing shortest transition paths between adjacent pairs ($A \rightarrow B$, then $B \rightarrow C$).
  *Selected*: Preserves user intent for curriculum generation while remaining deterministic and clean.

### For Path Edge Weighting
- **Option A (Unweighted Hop Count)**: Standard BFS shortest path.
  *Rejected*: Ignores link structure and transition probabilities ($P_{ij}$).
- **Option B (Negative Log Probability Weights $w_{ij} = -\log P_{ij}$)**: Run Dijkstra's algorithm on $w_{ij}$.
  *Selected*: Minimizing $\sum -\log P_{ij}$ mathematically maximizes the cumulative transition probability product $\prod P_{ij}$ along the path.

### For `--top N` Subgraph Expansion
- **Option A (Truncate or Fail)**: Render only exact path nodes, ignoring `--top N`.
  *Rejected*: Leaves graphs sparse and isolated without local context.
- **Option B ($k$-Hop Probabilistic Neighbor Expansion)**: Force-include path nodes first, then select highest transition probability 1-hop outgoing and incoming neighbors until $N$ nodes are reached.
  *Selected*: Provides critical context around path stepping stones while preserving full control over `--top N`.

## Decision

We decided to extend `wikigraph goal` with a `--strategy` flag supporting `union`, `intersection`, `path`, and `bottleneck`:

1. **`union`** *(default)*: Ranks nodes by $\min_g \text{MFPT}(i, g)$.
2. **`intersection`**: Ranks nodes by $\max_g \text{MFPT}(i, g)$ (requires reachability to all goals).
3. **`path`**:
   - Evaluates goals in flag order ($g_1 \rightarrow g_2 \rightarrow \dots \rightarrow g_M$).
   - Uses Dijkstra's algorithm on edge weights $w_{ij} = -\log P_{ij}$ (where $P_{ij} > 0$).
   - Expands undersized paths up to `--top N` using 1-hop probabilistic neighbors.
4. **`bottleneck`**: Scores nodes by random walk betweenness centrality across goal pairs.

## Consequences

- **Flexibility**: Users can customize goal visualizations according to their learning intent (curriculum path vs prerequisite intersection vs neighborhood exploration).
- **Determinism**: Sequential flag evaluation gives deterministic results for multi-goal path generation.
- **Backwards Compatibility**: Default remains `union`, ensuring zero breaking changes for existing workflows.

## Sources

- [[goal]] — Command surface using these strategies
- [[path-sequence]] — Path algorithm details
- [[bottleneck-centrality]] — Bottleneck centrality algorithm details
- [Issue #29: feat(goal): alternative partitioning and scoring strategies for multi-goal subgraphs](https://github.com/stephen-mcelhose/wikigraph/issues/29)
- `cmd_goal.go`
