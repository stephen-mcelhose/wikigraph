---
type: concept
title: Commute Time Metric
description: Symmetric distance metric K(i,j) = MFPT(i,j) + MFPT(j,i) used by wikigraph analyze to discover missing link suggestions between related pages.
tags: [commute-time, mfpt, distance-metric, link-suggestions, markov-chain]
timestamp: 2026-08-09T08:35:00Z
---

# Commute Time Metric

## Overview

The **commute time** $K(i,j)$ between page $i$ and page $j$ is the expected number of random-walk steps to travel from $i$ to $j$ and return back to $i$:

$$K(i,j) = M(i,j) + M(j,i)$$

where $M(i,j)$ is the [[mfpt|Mean First Passage Time]].

## Key Properties

### Why symmetry matters

Unlike MFPT, which is asymmetric ($M(i,j) \neq M(j,i)$), commute time is strictly symmetric ($K(i,j) = K(j,i)$). This makes $K(i,j)$ a true graph metric (satisfying triangle inequality).

### Relationship to graph resistance distance

Commute time is directly proportional to **effective resistance** $R_{ij}$ in an electrical network where edges represent unit resistors:

$$K(i,j) = 2m \cdot R_{ij}$$

Two pages have a small commute time if there are many short, parallel paths connecting them in the wiki.

### Reading Suggested Missing Links

In `wikigraph analyze docs/`:

```
=== Suggested missing links (lowest commute time, not yet linked) ===
analyze:
  → catrace (commute: 27.76)
```

Low commute time between unlinked pages indicates strong implicit connectivity — ideal candidates for explicit `[[wikilinks]]`.

## Related Concepts

- [[mfpt]] — Asymmetric mean first passage time foundation
- [[analyze]] — How `wikigraph analyze` formats link suggestions
- [[catrace]] — Fundamental matrix linear algebra implementation

## Sources

- Chandra et al. (1996). *The Electrical Resistancy of a Graph along with its Applications to Random Walks*.
- [[mfpt]]
