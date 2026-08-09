---
type: concept
title: Stationary Distribution (π)
description: Long-run visit probability π finding central hub pages and orphan pages in a wiki random walk.
tags: [stationary-distribution, power-iteration, centrality, orphan-detection]
timestamp: 2026-08-09T08:35:00Z
---

# Stationary Distribution (π)

## Overview

The **stationary distribution** $\pi$ is an $n$-dimensional probability vector satisfying $P^T \pi = \pi$ with $\sum_i \pi_i = 1$. It represents the long-run fraction of time a random walker spends on each page.

## Key Properties

### Power iteration algorithm

`catrace` computes $\pi$ using power iteration starting from a uniform vector $v^{(0)} = [1/n, \dots, 1/n]^T$:

$$v^{(k+1)} = P^T v^{(k)}$$

Iteration stops when $\|v^{(k+1)} - v^{(k)}\|_\infty < \text{tol}$ (default $10^{-8}$).

### Significance of π values

- **High $\pi$ (central hubs)**: Core pages frequently linked from across the wiki (e.g. [[analyze]], [[random-walk]]).
- **Low $\pi$ (orphans)**: Pages rarely reached by following links (bottom 10% flagged in `wikigraph analyze`).

## Related Concepts

- [[random-walk]] — Convergence to stationary vector
- [[analyze]] — Orphan and central page sections
- [[adr-003-orphan-threshold]] — Low $\pi$ tolerance policy for governance pages

## Sources

- Golub, G. H., & Van Loan, C. F. (2013). *Matrix Computations*. Johns Hopkins.
- [[catrace]]
