---
type: decision
title: "ADR-012: Teleporting PageRank as Default Markov Math"
description: Supersede ADR-011. Default math is full α-damping via NewTeleportingKernelFromAdj; raw adjacency for structure/display; NodeMass for PageRank-sized graphs.
tags: [adr, pagerank, teleportation, ergodicity, visualization, markov-chain]
timestamp: 2026-08-15T21:30:00Z
status: accepted
---

# ADR-012 — Teleporting PageRank as Default Markov Math

## Context

[[adr-011-sink-teleportation-vs-pagerank-damping]] retained sink-only
teleportation (pre-fill zero rows with `1/n`) because catrace rejected zero
rows and lacked a clean display/math split. That kept π computable but:

1. Non-sink pages followed links with probability 1 (implicit α = 1) — not
   industry PageRank.
2. Sink teleport rows rendered as fake edges ("noisy sinks") in `graph` /
   `export`.
3. Personalization (`--seed` PPR) was awkward on a sink-prefilled adj.

catrace now provides `NewTeleportingKernelFromAdj`, `VisualiseOptions.NodeMass`,
and PPR helpers (catrace PR #39). [[pagerank-foundation-rewrite]] (epic #56)
proposed a foundation rewrite: raw adj everywhere, full α-damping as default
math, clean viz via NodeMass. Child issue #55 covers `graph` + shared kernel.

## Decision

In the context of wiki Markov analysis and knowledge-graph visualization,
facing sink-only teleportation vs full α-damping PageRank after catrace gained
teleport/NodeMass APIs, we decided to **adopt the teleporting kernel as the
shared default math and keep raw adjacency for structure and display**, to
achieve PageRank-compatible π, clean graphs, and natural `--seed` PPR,
accepting a one-time re-baseline of docs-wiki goldens and demoting sink-only
teleportation to a historical / pedagogical footnote.

Concretely:

1. **Raw adjacency** — `buildAdjacencyWithOpts` records real wikilinks only;
   sink rows stay zero; sinks are still reported as an authoring smell.
2. **Default math** — `NewTeleportingKernelFromAdj` with user-facing `--alpha`
   (default `0.15` = teleport probability; link weight `1−α`) and uniform
   restart → global PageRank for π, entropy, MFPT / commute, and `goal`.
3. **Personalization** — repeatable `--seed <slug>` builds a non-uniform
   restart → Personalized PageRank (used by `graph` now; optional `goal --rank
   ppr` remains issue #22).
4. **Display** — `graph` sizes nodes with `NodeMass` from the teleporting
   stationary distribution and draws edges from an α=0 base-link kernel with
   `MinEdge` high enough to hide sink restart rows (~0.02).
5. **Analyze / export structure** — Overview edge count = raw directed
   wikilink count; communicating classes = raw digraph SCCs (authoring
   signal). The teleporting chain is ergodic by construction and is not the
   health signal for fragmentation. Export emits raw edges + PageRank π.
6. **Hard cut** — no `--kernel sink-only` migration flag; ADR-011 is
   superseded.

## Consequences

- π, entropy, orphan thresholds, and commute suggestions change vs the
  sink-only era; [[testing-runbook]] and `wiki_integration_test.go` were
  re-golden'd for the docs wiki.
- `analyze` labels centrality as PageRank π; class section titles refer to the
  raw wikilink digraph.
- Spider traps and bored-surfer behaviour match peer PageRank tools.
- MFPT / `goal` remain first-class on the same teleporting kernel (PPR ranking
  is additive later via #22).
- Sink-only teleportation remains a valid Markov construction for teaching
  ("dead ends break π") but is not the product default.

## Options Considered

- **Incremental dual math** (PPR viz on `graph` only; sink-prefill elsewhere) —
  rejected as an end state (two π truths).
- **Legacy `--kernel sink-only` for one release** — rejected; hard cut with
  this ADR is enough.
- **Export teleporting `P`** — rejected for default export; raw edges + π
  match interop and clean viz. Math `P` remains available in-process via the
  kernel.

## Implementation Notes

- Shared builder: `buildKernelWithOpts(..., alpha, seeds)` in `wiki.go`.
- Raw SCCs: `rawSCCs` in `digraph.go` (Tarjan).
- Persistent flags: `--alpha`, `--seed`.

## Sources

- [[pagerank-foundation-rewrite]] — accepted proposal (epic #56)
- [[adr-011-sink-teleportation-vs-pagerank-damping]] — superseded predecessor
- [[teleportation-ergodicity]] — mathematical background
- https://github.com/stephen-mcelhose/wikigraph/issues/56
- https://github.com/stephen-mcelhose/wikigraph/issues/55
- https://github.com/stephen-mcelhose/catrace/pull/39
- Page, L., Brin, S., Motwani, R., & Winograd, T. (1999). *The PageRank
  Citation Ranking*. http://ilpubs.stanford.edu:8090/422/
