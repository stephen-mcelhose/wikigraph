---
type: concept
title: Stationary Distribution (π)
description: Long-run visit probability π — PageRank mass on a wiki walk, and the same fixed-point idea as long-run experience under PDA agent dynamics.
tags: [stationary-distribution, power-iteration, centrality, orphan-detection, pagerank, pda]
timestamp: 2026-08-15T21:38:39Z
---

# Stationary Distribution (π)

## Overview

The **stationary distribution** $\pi$ is an $n$-dimensional probability vector
satisfying $P^T \pi = \pi$ with $\sum_i \pi_i = 1$. It is the long-run fraction
of time a walker spends in each state.

In `wikigraph`, under the default teleporting kernel with uniform restart,
$\pi$ is **global PageRank** on pages — hub vs orphan signal in [[analyze]].
With `--seed`, the fixed point is **Personalized PageRank**: long-run mass
under persistent intent. That is the same $\pi$ vs PPR distinction catrace uses
for PDA agents (dynamics-only vs goal-directed experience). See
[[knowledge-graph-to-pda-agents]].

## Key Properties

### Power iteration algorithm

`catrace` computes $\pi$ using power iteration starting from a uniform vector $v^{(0)} = [1/n, \dots, 1/n]^T$:

$$v^{(k+1)} = P^T v^{(k)}$$

Iteration stops when $\|v^{(k+1)} - v^{(k)}\|_\infty < \text{tol}$ (default $10^{-8}$).
On a teleporting kernel this converges for $\alpha \in (0,1]$.

### Significance of π values

- **High $\pi$ (central hubs)**: Core pages (or agent states) that attract long-run mass.
- **Low $\pi$ (orphans)**: Rarely visited pages/states (bottom 10% in `wikigraph analyze`).
- **PPR vs $\pi$**: $\|\pi - \mathrm{ppr}(v,\alpha)\|$ measures how much intent moves the long run — on a wiki or on an agent kernel.

## Related Concepts

- [[knowledge-graph-to-pda-agents]] — π / PPR as dynamics vs intent
- [[random-walk]] — Convergence to stationary vector
- [[analyze]] — Orphan and central page sections
- [[adr-003-orphan-threshold]] — Low $\pi$ tolerance policy for governance pages
- [[teleportation-ergodicity]] — Why teleportation is required for π to be defined
- [[sink-page]] — Structural sinks; restart handling in the teleporting kernel
- [[adr-012-teleporting-pagerank-default]] — Default PageRank math
- [[catrace]] — Shared engine with PDA agents

## Sources

- Golub, G. H., & Van Loan, C. F. (2013). *Matrix Computations*. Johns Hopkins.
- [[catrace]]
- [[knowledge-graph-to-pda-agents]]
