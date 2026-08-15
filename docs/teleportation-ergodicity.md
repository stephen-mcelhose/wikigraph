---
type: concept
title: Teleportation and Ergodicity in Markov Wiki Graphs
description: Why α-damping PageRank makes the wiki walk ergodic; sink-only teleportation as historical pedagogy; raw structure vs teleporting math.
tags: [markov-chain, ergodicity, pagerank, teleportation, sink-page, stationary-distribution]
timestamp: 2026-08-15T21:30:00Z
resource: https://en.wikipedia.org/wiki/PageRank
---

# Teleportation and Ergodicity in Markov Wiki Graphs

## Overview

wikigraph models a wiki as a discrete-time Markov chain. For π, MFPT, and
related math the chain must be **ergodic**. A raw wiki digraph almost never is:
[[sink-page|sink]] pages produce zero rows, and disconnected clusters stay
unreachable.

**Teleportation** (PageRank α-damping) is the fix: at each step the surfer
follows a real link with probability $1-\alpha$ or jumps according to a restart
distribution $v$ with probability $\alpha$. Default $\alpha = 0.15$ (user-facing
`--alpha`); uniform $v$ → global PageRank; `--seed` → Personalized PageRank.

## Key Properties

### PageRank form (current default)

With raw adjacency $A$ and restart $v$:

$$
T_{ij} = \alpha v_j + (1-\alpha)\frac{A_{ij}}{\mathrm{rowsum}_i}
\quad\text{(non-sink)},\qquad
T_{ij} = v_j \quad\text{(sink)}.
$$

This is `catrace.NewTeleportingKernelFromAdj`. Naming: wikigraph's `--alpha` is
the **teleport** probability (firehose / catrace style), so link-following
weight is $1-\alpha$. Literature often writes damping $d \approx 0.85$ with
teleport $1-d$.

In the PDA / agent reading, $v$ is **intent** and $\alpha$ is **intentionality
weight** — the same dials as Personalized PageRank on a qualia or world kernel.
See [[knowledge-graph-to-pda-agents]].

### Sink-only teleportation (historical)

Earlier wikigraph pre-filled sink rows with $1/n$ and used
`NewRandomWalkKernel` with implicit $\alpha = 1$ on non-sinks. That is a valid
Markov construction and a useful teaching step ("dead ends break π"), but it is
**not** the product default anymore. See superseded
[[adr-011-sink-teleportation-vs-pagerank-damping]] and current
[[adr-012-teleporting-pagerank-default]].

### Raw structure vs teleporting math

| Signal | Where |
| ------ | ----- |
| Edge count, export edges, authoring SCCs | Raw digraph |
| π, entropy, commute, MFPT | Teleporting $T$ |
| Graph node sizes | `NodeMass` ← stationary of $T$ |
| Graph edges | Real links (α=0 base + `MinEdge`) |

Communicating classes on $T$ are nearly always one class (ergodic by
construction). `analyze` therefore reports **raw** SCCs as the fragmentation
signal, and lists sinks separately so dangling pages are not mistaken for
"fifty classes."

## Related Concepts

- [[sink-page]] — Structural sinks and restart handling
- [[stationary-distribution]] — PageRank π
- [[mfpt]] — Defined on the teleporting kernel
- [[communicating-classes]] — Raw digraph SCCs in analyze
- [[markov-model]] — Pipeline
- [[knowledge-graph-to-pda-agents]] — α / v as agent intent
- [[adr-012-teleporting-pagerank-default]] — Current decision
- [[adr-011-sink-teleportation-vs-pagerank-damping]] — Superseded
- [[pagerank-foundation-rewrite]] — Proposal that drove ADR-012
- [[catrace]] — TeleportingKernel / PPR APIs

## Sources

- Wikipedia: PageRank — https://en.wikipedia.org/wiki/PageRank
- Page, L., Brin, S., Motwani, R., & Winograd, T. (1999). *The PageRank Citation Ranking*. http://ilpubs.stanford.edu:8090/422/
- [[adr-012-teleporting-pagerank-default]]
- catrace PR #39 — `NewTeleportingKernelFromAdj`, `NodeMass`
