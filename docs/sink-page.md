---
type: concept
title: Sink Pages and Teleportation
description: Handling pages with zero outgoing wikilinks via uniform teleportation to maintain row-stochastic transition matrices.
tags: [sink-page, teleportation, markov-chain, dead-end]
timestamp: 2026-08-09T08:35:00Z
---

# Sink Pages and Teleportation

## Overview

A **sink page** is a page with zero outgoing `[[wikilinks]]`. In the raw adjacency matrix, a sink creates an all-zero row, violating the row-stochastic requirement ($\sum_j P_{ij} = 1$).

## Key Properties

### The uniform teleportation fix

When `wikigraph` encounters a sink row $i$ in `wiki.go`, it replaces the all-zero row with a uniform probability distribution over all $n$ pages:

$$P_{ij} = \frac{1}{n} \quad \forall j$$

This simulates a reader who reaches a dead end and opens a random page in the wiki.

### How to fix a sink

In `wikigraph analyze docs/`:

```
=== Sink pages (no outgoing links) ===
  some-page -> add outgoing links
```

Authors fix sinks by adding relevant outgoing `[[wikilinks]]`.

### Display vs math adjacency

Teleportation rows must **not** be rendered as real graph edges — they are a
mathematical construct, not links authored in the wiki. The display adjacency
(used by `graph` and `export`) contains only real links; zero rows for sinks
are left as-is so sinks appear as dead-ends. The math adjacency (passed to
catrace) has teleportation rows added before kernel construction. See
[[adr-011-sink-teleportation-vs-pagerank-damping]] and
[[teleportation-ergodicity]].

## Related Concepts

- [[markov-model]] — Transition matrix construction
- [[communicating-classes]] — Sink components
- [[analyze]] — Sink reporting section
- [[teleportation-ergodicity]] — Mathematical justification and PageRank comparison
- [[adr-011-sink-teleportation-vs-pagerank-damping]] — Decision: retain sink-only teleportation, defer full α-damping

## Sources

- `wiki.go` — `buildAdjacency` function
- Page, L., et al. (1999). *The PageRank Citation Ranking: Bringing Order to the Web*.
- Wikipedia: PageRank — https://en.wikipedia.org/wiki/PageRank
