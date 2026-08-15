---
type: concept
title: Mean First Passage Time (MFPT)
description: Expected steps M(i,j) for a random walk to reach target page j from start page i — foundation for wikigraph goal learning paths.
tags: [mfpt, markov-chain, fundamental-matrix, goal-path]
timestamp: 2026-08-09T08:35:00Z
---

# Mean First Passage Time (MFPT)

## Overview

The **Mean First Passage Time (MFPT)** $M(i,j)$ is the expected number of random-walk steps required to reach target page $j$ for the first time starting from page $i$.

## Key Properties

### Computation via fundamental matrix

For an irreducible ergodic Markov chain with transition matrix $P$ and stationary vector $\pi$, the fundamental matrix $Z$ is defined as:

$$Z = (I - P + W)^{-1}$$

where $W = \mathbf{1} \pi^T$. The MFPT matrix $M$ is given by:

$$M(i,j) = \frac{Z_{jj} - Z_{ij}}{\pi_j}$$

### Applications in wikigraph

1. **`wikigraph goal --goal <slug>`**: Ranks all candidate prerequisite pages $i$ by ascending $M(i, \text{goal})$ to generate an optimal learning path.
2. **`wikigraph analyze`**: Forms symmetric [[commute-time]] $K(i,j) = M(i,j) + M(j,i)$ for link recommendations.

On PDA agent kernels the same $M(i,j)$ is **access cost** to a target
experience or task-complete state — complementary to PPR occupancy under
intent. See [[knowledge-graph-to-pda-agents]].

## Related Concepts

- [[commute-time]] — Symmetric $M(i,j) + M(j,i)$ metric
- [[goal]] — CLI subcommand utilizing MFPT
- [[catrace]] — Linear algebra implementation of $Z$ matrix
- [[knowledge-graph-to-pda-agents]] — MFPT as agent latency / access cost
- [[stationary-distribution]] — Complementary long-run mass

## Sources

- Kemeny, J. G., & Snell, J. L. (1976). *Finite Markov Chains*. Springer-Verlag.
- [[catrace]]
- [[knowledge-graph-to-pda-agents]]
