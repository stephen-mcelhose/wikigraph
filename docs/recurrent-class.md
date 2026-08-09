---
type: concept
title: Recurrent vs Transient Classes
description: Classification of communicating classes into recurrent (closed) and transient (leaky) sets in a wiki Markov chain.
tags: [recurrent, transient, markov-chain, scc, graph-theory]
timestamp: 2026-08-09T08:35:00Z
---

# Recurrent vs Transient Classes

## Overview

In a Markov chain, communicating classes are categorized as either **recurrent** (absorbing sets from which escape is impossible) or **transient** (leaky sets that eventually lose all probability mass).

## Key Properties

### Recurrent Class

A class $C$ is recurrent if $P_{ij} = 0$ for all $i \in C, j \notin C$. Once a random walker enters $C$, they remain in $C$ forever.

### Transient Class

A class $C$ is transient if there exists at least one path from $C$ to another class $C'$. As $t \to \infty$, the probability of being in $C$ drops to zero ($\pi_i = 0$ for all $i \in C$).

### Practical impact on wikis

If a wiki has a transient class, readers who navigate out of those pages can never navigate back. In `wikigraph analyze`, transient pages are flagged so authors can add return links.

## Related Concepts

- [[communicating-classes]] — Kosaraju SCC decomposition
- [[stationary-distribution]] — Stationary probability $\pi_i = 0$ on transient nodes
- [[sink-page]] — Extreme 1-page transient cases

## Sources

- [[communicating-classes]]
- Feller, W. (1968). *An Introduction to Probability Theory and Its Applications*.
