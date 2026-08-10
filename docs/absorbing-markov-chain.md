---
type: concept
title: Absorbing Markov Chain
description: A Markov chain with absorbing states that cannot be left once entered, used in wikigraph for bottleneck centrality and random walk reachability analysis.
tags: [absorbing-markov-chain, fundamental-matrix, transient-states, markov-chain, bottleneck]
timestamp: 2026-08-10T01:25:00Z
---

# Absorbing Markov Chain

## Overview

An **absorbing Markov chain** is a Markov chain containing at least one **absorbing state** — a state that, once entered, cannot be left ($P_{ii} = 1$) — such that every non-absorbing (transient) state can reach at least one absorbing state in a finite number of steps.

In `wikigraph`, absorbing Markov chains model goal-directed navigation and calculate [[bottleneck-centrality]] between pairs of pages.

## Canonical Form & Fundamental Matrix

By reordering states into transient ($T$) and absorbing ($A$) sets, the transition matrix $P$ takes the canonical block form:

$$P = \begin{pmatrix} Q & R \\ 0 & I \end{pmatrix}$$

where:
- $Q$ is the $t \times t$ transition matrix between transient states.
- $R$ is the $t \times r$ transition matrix from transient states to absorbing states.
- $I$ is the $r \times r$ identity matrix (absorbing states remain in place).
- $0$ is the $r \times t$ zero matrix.

### Fundamental Matrix $N$

The **fundamental matrix** $N$ captures the expected number of visits to transient state $j$ starting from transient state $i$ before absorption:

$$N = I + Q + Q^2 + Q^3 + \dots = \sum_{k=0}^{\infty} Q^k = (I - Q)^{-1}$$

## Usage in WikiGraph

1. **Bottleneck Strategy (`wikigraph goal --strategy bottleneck`)**: Set target goal $t$ as the absorbing state. Entry $N_{s, k}$ gives the expected visits to page $k$ on walks starting at source $s$ before reaching goal $t$.
2. **Trace Kernel Subgraphing (`catrace`)**: Computes effective transition matrices on subsets of pages by absorbing transitions outside the subset and tracking re-entry probabilities.

## Related Concepts

- [[bottleneck-centrality]] — Bottleneck scoring algorithm built on $N = (I-Q)^{-1}$
- [[recurrent-class]] — Classification of states and communicating classes
- [[mfpt]] — Expected steps to reach target states
- [[catrace]] — Fundamental matrix computation engine

## Sources

- Kemeny, J. G., & Snell, J. L. (1976). *Finite Markov Chains*. Springer-Verlag. [DOI: 10.1007/978-0-387-90192-3](https://doi.org/10.1007/978-0-387-90192-3) | [Springer Link](https://link.springer.com/book/10.1007/978-0-387-90192-3)
