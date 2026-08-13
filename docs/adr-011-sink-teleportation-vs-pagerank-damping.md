---
type: decision
title: "ADR-011: Sink-Only Teleportation vs Full PageRank Damping"
description: Retain sink-only teleportation for ergodicity. Separate the math adjacency (with teleportation) from the display adjacency (real links only). Defer full α-damping as a future option.
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
side-effect: the modified adjacency matrix is passed directly to `catrace` and
then to the graph renderer and export formats. Teleportation rows appear as
real edges, producing misleading star patterns in the visualization — a sink
page appears to link to every other page.

The question raised: is the sink-only approach mathematically sound, or should
we adopt full PageRank-style damping (α ≈ 0.85 applied to all pages)?

See [[teleportation-ergodicity]] for the mathematical background.

## Decision

In the context of achieving ergodicity for Markov chain analysis of wiki link
graphs, facing the choice between sink-only teleportation and full PageRank
α-damping:

1. **Retain sink-only teleportation for math.** Sink pages teleport uniformly
   to all n pages. Non-sink pages follow real links with probability 1. This is
   mathematically legitimate: it produces a valid Markov chain with a
   well-defined stationary distribution. The math is not faked — it is a
   deliberate and standard modelling choice (see PageRank's original treatment
   of dangling nodes in Page et al., 1999).

2. **Separate display adjacency from math adjacency.** The root cause of the
   visualisation bug is that the same adjacency matrix was used for both the
   Markov math and the graph renderer. Fix: `buildAdjacencyWithOpts` returns
   the raw adjacency (real links only, zero rows for sinks). `buildKernelWithOpts`
   adds teleportation rows to a copy before passing to catrace. The graph and
   export commands use the raw adjacency for display — sinks appear as
   dead-ends, which is the truth about the wiki structure.

3. **Defer full α-damping.** Full PageRank damping (`P = α·A + (1−α)·J/n`)
   is more theoretically principled — it guarantees ergodicity unconditionally,
   eliminates the sink/non-sink asymmetry, and models realistic browsing
   behaviour. However, adopting it would shift π values for all pages (not just
   sinks), breaking existing runbook expectations, and would require exposing α
   as a user-facing parameter. The benefit does not outweigh the disruption at
   this stage.

## Consequences

- Sink pages are correctly shown as dead-ends in `graph` and `export` output.
- The stationary distribution, MFPT, and class computation are unchanged —
  they continue to use the ergodic (teleportation-filled) adjacency.
- `analyze`'s sink detection continues to report sink pages accurately.
- **Revisit if any of the following occur:**
  - Spider traps (small strongly-connected components with no outgoing edges)
    cause π to concentrate undesirably on a subset of pages.
  - Sink/non-sink asymmetry produces misleading π rankings in production wikis.
  - A user-facing damping factor (`--alpha`) is requested for PageRank
    compatibility.
  - catrace adds native support for α-damping.

## Sources

- [[teleportation-ergodicity]] — mathematical background
- [[sink-page]] — sink detection and reporting
- [[stationary-distribution]] — requires ergodic chain
- Wikipedia: PageRank — https://en.wikipedia.org/wiki/PageRank
- Page, L., Brin, S., Motwani, R., & Winograd, T. (1999). *The PageRank
  Citation Ranking: Bringing Order to the Web*. Stanford InfoLab. http://ilpubs.stanford.edu:8090/422/
