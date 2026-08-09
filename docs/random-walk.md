---
type: concept
title: Random Walk on Wiki Graphs
description: Mathematical model of a reader navigating a wiki by following wikilinks — the stochastic foundation behind all wikigraph subcommands.
tags: [random-walk, markov-chain, pagerank, stochastic-process]
timestamp: 2026-08-09T08:35:00Z
---

# Random Walk on Wiki Graphs

## Overview

A **random walk** on a wiki graph represents a reader who starts at a page $X_0$ and at each discrete step $t$ clicks one of the page's outgoing `[[wikilinks]]` with uniform probability.

## Key Properties

### Connection to transition matrix

The probability of moving from page $i$ to page $j$ in one step is:

$$P(X_{t+1} = j \mid X_t = i) = P_{ij}$$

The state distribution at step $t$ is given by $v^{(t)} = v^{(0)} P^t$.

### Long-run behavior

If the graph forms a single [[communicating-classes|recurrent class]], $v^{(t)}$ converges to a unique stationary distribution $\pi$ as $t \to \infty$, independent of starting page $v^{(0)}$:

$$\lim_{t \to \infty} P^t = \mathbf{1} \pi^T$$

## Related Concepts

- [[markov-model]] — Construction of transition matrix $P$
- [[stationary-distribution]] — Long-run visit frequency $\pi$
- [[entropy-rate]] — Step-by-step predictability

## Sources

- Aldous, D., & Fill, J. A. (2002). *Reversible Markov Chains and Random Walks on Graphs*.
- [[markov-model]]
