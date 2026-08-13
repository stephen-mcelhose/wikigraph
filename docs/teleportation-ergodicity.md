---
type: concept
title: Teleportation and Ergodicity in Markov Wiki Graphs
description: Why sink-only teleportation is mathematically legitimate, how it relates to PageRank's damping factor, and when to reconsider full α-damping.
tags: [markov-chain, ergodicity, pagerank, teleportation, sink-page, stationary-distribution]
timestamp: 2026-08-13T00:52:04Z
resource: https://en.wikipedia.org/wiki/PageRank
---

# Teleportation and Ergodicity in Markov Wiki Graphs

## Overview

wikigraph models a wiki as a discrete-time Markov chain. For the core math to
work — stationary distribution π, MFPT, communicating classes — the chain must
be **ergodic**: irreducible (every state reachable from every other) and
aperiodic. A raw wiki graph almost never satisfies this: sink pages (no
outgoing links) produce zero rows in the transition matrix, making π undefined.

**Teleportation** is the standard fix: a surfer who lands on a sink page is
assumed to jump uniformly to any page at random. This is not a mathematical
hack — it is the same mechanism Google's PageRank uses to handle dangling
nodes.

## Key Properties

### PageRank's damping factor

Page, Brin, Motwani, and Winograd (1998) introduced the **damping factor** `d`
(≈ 0.85) to model a surfer who, at each step, either follows a link with
probability `d` or teleports to a uniformly random page with probability `1−d`:

```
PR(pᵢ) = (1−d)/N + d · Σⱼ PR(pⱼ) / L(pⱼ)
```

Applied to **all** pages (not just sinks), this guarantees ergodicity and
models the "bored surfer" who occasionally types a new URL. The stationary
distribution of this modified chain is PageRank.

### Sink-only teleportation (wikigraph's approach)

wikigraph applies teleportation only to sink pages — pages with zero outgoing
links after wikilink and relative-link parsing. Non-sink pages follow their
real links with probability 1. This is a strict subset of full PageRank
damping:

| | wikigraph | Full PageRank |
|---|---|---|
| Sink pages | uniform jump to all n pages | uniform jump (via 1−d term) |
| Non-sink pages | pure link-following | d · link + (1−d) · random jump |
| Damping parameter | none (implicit α = 1 for non-sinks) | α ≈ 0.85 |
| Ergodicity guarantee | only if no isolated non-sink components | always |

Sink-only teleportation is a valid Markov chain with a well-defined stationary
distribution. The math is not faked. The stationary distribution reflects the
genuine link structure — sink pages act as uniform redistributors in the chain,
which is a reasonable model of a reader who gets stuck and picks a random next
page.

### Why this matters for visualization

The teleportation rows in the transition matrix P should **not** be rendered as
graph edges — they are a mathematical construct, not real links. Rendering them
produces misleading star patterns (a sink appears to link to every page). The
correct approach is to build two separate adjacency representations:

- **Raw adjacency** — real links only, zero rows for sinks; used for display
- **Math adjacency** — teleportation rows added for sinks; passed to catrace for π, MFPT, and class computation

See [[adr-011-sink-teleportation-vs-pagerank-damping]] for the architectural decision.

## Related Concepts

- [[sink-page]] — Zero-outgoing-link pages and how they are detected
- [[stationary-distribution]] — π requires an ergodic chain
- [[mfpt]] — Mean first passage time is undefined on non-ergodic chains
- [[communicating-classes]] — Sink pages form transient classes without teleportation
- [[markov-model]] — The full pipeline from wikilinks to P
- [[adr-011-sink-teleportation-vs-pagerank-damping]] — Decision record

## Sources

- Wikipedia: PageRank — https://en.wikipedia.org/wiki/PageRank
- Page, L., Brin, S., Motwani, R., & Winograd, T. (1999). *The PageRank Citation Ranking: Bringing Order to the Web*. Stanford InfoLab Technical Report. http://ilpubs.stanford.edu:8090/422/
- Brin, S., & Page, L. (1998). *The Anatomy of a Large-Scale Hypertextual Web Search Engine*. Computer Networks and ISDN Systems, 30(1–7), 107–117. https://doi.org/10.1016/S0169-7552(98)00110-X
