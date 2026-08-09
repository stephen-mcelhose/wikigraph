---
type: concept
title: Wiki Entropy Rate
description: Entropy rate H of a random walk on a wiki graph — a global information-theoretic health signal measuring structural predictability and branching uniformity.
tags: [entropy, information-theory, markov-chain, health-signal]
timestamp: 2026-08-09T08:35:00Z
---

# Wiki Entropy Rate

## Overview

The **entropy rate** $H(X)$ measures the average information (in bits) per step of a random walker traversing the wiki graph:

$$H(X) = -\sum_{i} \pi_i \sum_{j} P_{ij} \log_2 P_{ij}$$

where $\pi$ is the [[stationary-distribution|stationary distribution]] and $P$ is the transition matrix.

## Key Properties

### What the range means

- **$H \approx 0$**: Highly deterministic graph (few branches, linear chains).
- **$H = \log_2 d$**: Uniform regular graph with node degree $d$.
- **Healthy wiki range**: Typically between $2.0$ and $4.0$ bits per step for a well-linked knowledge base of 20–100 pages.

### Why mid-range is healthy

- Too low ($< 1.5$ bits): Graph is overly rigid or linear.
- Too high ($> 5.0$ bits): Links are oversaturated and indiscriminate, offering no structural guidance.

## Related Concepts

- [[markov-model]] — Random walk kernel $P$
- [[stationary-distribution]] — Weighting factor $\pi_i$
- [[analyze]] — CLI health reporting

## Sources

- Cover, T. M., & Thomas, J. A. (2006). *Elements of Information Theory*. Wiley.
- [[catrace]]
