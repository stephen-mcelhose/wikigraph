---
type: decision
title: "ADR-011: Sink-Only Teleportation vs Full PageRank Damping"
description: Retain sink-only teleportation for ergodicity. Accept noisy sink visualization as the cost of a correct Markov chain. Defer full α-damping as a future option.
tags: [adr, ergodicity, teleportation, pagerank, markov-chain, visualization]
timestamp: 2026-08-13T00:52:04Z
status: accepted
---

# ADR-011 — Sink-Only Teleportation vs Full PageRank Damping

## Context

wikigraph computes stationary distribution (π), mean first passage time (MFPT),
and communicating classes on a Markov chain derived from wiki links. These
require an ergodic (irreducible + aperiodic) chain. Raw wiki graphs are never
ergodic: sink pages — pages with no outgoing wikilinks or relative-link edges —
produce zero rows in the transition matrix, making π undefined.

The current fix: `buildAdjacencyWithOpts` fills zero rows for sinks with `1.0`
in every column (uniform teleportation). This achieves ergodicity but has a
visual side-effect: teleportation rows appear as real edges in `graph` and
`export` output — a sink page appears to link to every other page ("noisy
sink").

We investigated separating the display adjacency (raw links only) from the math
adjacency (with teleportation). Reading the `catrace` source revealed this is
not viable without deeper changes:

- `catrace.NewRandomWalkKernel` explicitly rejects zero rows with an error.
- `Stationary` is computed lazily from `k.P` on every call; there is no cached
  result to preserve if we zero out sink rows post-construction.
- Replacing teleportation with self-loops for the display kernel creates
  absorbing states, making the display chain non-ergodic and causing
  `Stationary` to fail, which collapses all node sizes to equal in the
  visualization.

The question raised: is the sink-only approach mathematically sound, or should
we adopt full PageRank-style damping (α ≈ 0.85 applied to all pages)?

See [[teleportation-ergodicity]] for the mathematical background.

## Decision

In the context of achieving ergodicity for Markov chain analysis of wiki link
graphs, facing the choice between sink-only teleportation and full PageRank
α-damping:

1. **Retain sink-only teleportation and accept the noisy sink visualization.**
   Sink pages teleport uniformly to all n pages. This is mathematically
   legitimate and is the direct analogue of how PageRank handles dangling nodes
   (Page et al., 1999). The "star pattern" a sink page produces in the graph
   is an accurate representation of its role in the ergodic chain: from a sink,
   a random reader can reach any page. It is noisy, not wrong.

2. **Do not attempt to separate display adjacency from math adjacency.**
   Doing so requires either modifying catrace (out of scope) or working around
   its validation and lazy computation in fragile ways. The ergodicity
   requirement is non-negotiable for correct math; the display is secondary.

3. **Defer full α-damping.** Full PageRank damping (`P = α·A + (1−α)·J/n`)
   is more theoretically principled — it guarantees ergodicity unconditionally,
   eliminates the sink/non-sink asymmetry, and models realistic browsing
   behaviour. However, adopting it would shift π values for all pages (not just
   sinks), breaking existing runbook expectations, and would require exposing α
   as a user-facing parameter. The benefit does not outweigh the disruption at
   this stage.

## Consequences

- Sink pages produce a star pattern in `graph` and `export` output. This is a
  known, accepted property of the ergodic chain — not a rendering bug.
- `analyze` reports sink pages explicitly so authors know which pages to fix.
- The stationary distribution, MFPT, and class computation are correct.
- **Revisit if any of the following occur:**
  - Spider traps cause π to concentrate undesirably on a subset of pages.
  - Sink/non-sink asymmetry produces misleading π rankings in production wikis.
  - catrace adds native support for α-damping or a display-only rendering path.
  - A user-facing damping factor (`--alpha`) is requested for PageRank
    compatibility.

## Supersession candidate

[[pagerank-foundation-rewrite]] (epic
https://github.com/stephen-mcelhose/wikigraph/issues/56) proposes replacing
this ADR’s defaults with full α-damping PageRank, raw-edge display via
`NodeMass`, and an `analyze` section revisit. Until that proposal is accepted
and a successor ADR is written, this decision remains in force.

## Sources

- [[teleportation-ergodicity]] — mathematical background
- [[sink-page]] — sink detection and reporting
- [[stationary-distribution]] — requires ergodic chain
- [[pagerank-foundation-rewrite]] — proposed supersession (epic #56)
- Wikipedia: PageRank — https://en.wikipedia.org/wiki/PageRank
- Page, L., Brin, S., Motwani, R., & Winograd, T. (1999). *The PageRank
  Citation Ranking: Bringing Order to the Web*. Stanford InfoLab. http://ilpubs.stanford.edu:8090/422/
