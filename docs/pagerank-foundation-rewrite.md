---
type: proposal
title: PageRank Foundation Rewrite — Teleporting Kernel as Default Math
description: Propose replacing sink-only teleportation with full α-damping PageRank as the shared Markov foundation, clean raw-edge visualization via NodeMass, and a planned analyze section revisit — tracked by epic #56.
resource: https://github.com/stephen-mcelhose/wikigraph/issues/56
tags: [proposal, pagerank, personalized-pagerank, teleportation, visualization, analyze, adr-011]
timestamp: 2026-08-15T20:50:08Z
---

# PageRank Foundation Rewrite — Teleporting Kernel as Default Math

## Problem

wikigraph today builds a Markov chain by pre-filling sink rows with uniform
`1/n` mass (`buildAdjacencyWithOpts`), then runs `NewRandomWalkKernel`. That
achieves ergodicity but:

1. **Visualization is noisy** — teleportation rows render as fake edges
   (sink → everywhere “star”), which fights the product goal of clean knowledge
   graphs.
2. **Centrality is not PageRank** — non-sink pages follow links with probability 1
   (implicit α = 1). Industry tools (Neo4j GDS, NetworkX/igraph tutorials,
   Obsidian graph plugins) treat PageRank / Personalized PageRank as the
   default importance lens, with real links as the only drawn edges.
3. **[[adr-011-sink-teleportation-vs-pagerank-damping]] deferred full α-damping**
   largely because catrace lacked a clean display/math split. catrace now has
   `NewTeleportingKernelFromAdj` and `VisualiseOptions.NodeMass` (catrace PR
   #39) — the constraint that drove ADR-011 is lifting.
4. An **incremental** path (PPR only on `graph`, leave `analyze`/`export`/`goal`
   on sink-prefill) creates two competing π definitions and postpones a
   re-baseline we already accept we will do.

Reference implementation of the desired viz pattern:
`eis-intake-firehose/tools/firehose-graph` — PPR for node mass; base-link kernel
(`alpha=0`) + `MinEdge` for edges.

## Proposed Solution

Adopt a **foundation rewrite (options B/C)**, not a permanent opt-in dual
world:

### Shared model

1. **Raw adjacency everywhere** — real links only; zero rows for sinks; no
   `1/n` pre-fill in the shared builder.
2. **Default math** — teleporting kernel with uniform restart and user-facing
   `--alpha` (default `0.15`) → **global PageRank** for π, orphans, “most
   central,” and entropy rate.
3. **Personalization** — repeatable `--seed <slug>` builds a non-uniform
   restart → **Personalized PageRank** (PPR) for `graph` (and later optional
   `goal --rank ppr` per issue #22).
4. **Display** — always raw / base-link edges + `NodeMass` from the teleporting
   walk’s stationary distribution; `MinEdge` high enough to suppress artificial
   sink uniform edges when a base kernel is used for rendering.
5. **MFPT / commute / `goal`** — remain first-class, computed on the same
   teleporting kernel. PPR does not replace MFPT; they answer different questions
   (visit mass vs hitting time).

### Pedagogical stance (sink-only)

Sink-only teleportation is a **valid Markov construction** and a useful
*stepping stone* (“dead ends break π”), but it is **not** the long-term
definition of wiki centrality:

| Intent | Sink-only | Full α-damping |
| ------ | --------- | -------------- |
| Fix dangling nodes | Yes | Yes (via α term) |
| Teach “what is π?” | Works | Works |
| Match PageRank / clean KG viz | Incomplete | Canonical |
| Spider traps | Can trap mass | Softened by α |
| Natural `--seed` PPR | Awkward | Natural |

**Preserve as structure, not as default walk:**

- **Sink pages** = zero out-degree (authoring smell) — still reported from raw
  adjacency.
- **Dangling-node handling** documented *inside* the PageRank / Google-matrix
  story, not as a competing product default.
- Sink-only may remain a footnote: “α → 1 and teleport only on dangling nodes ≈
  the pre-rewrite model.”

### Analyze section revisit (required planning)

Changing the underlying kernel means revisiting which `analyze` sections speak
about the **raw graph** vs the **teleporting walk**:

| Section | Under PageRank default | Revisit? |
| ------- | ---------------------- | -------- |
| Sink pages | Raw out-degree | No — keep structural |
| Most central | Becomes PageRank π | Yes — clarify label |
| Orphans (bottom % of π) | Membership / thresholds shift | Yes — re-baseline; maybe retune `--orphan-pct` |
| Entropy rate | Different P → different H | Yes — numbers; meaning OK |
| Edge count | Today counts math `P` (includes teleport mass) | Yes — prefer **raw** edge count for health |
| Communicating classes | Teleporting chain → often one class | Yes — consider raw-graph classes as authoring signal |
| Suggested links (commute) | Values shift | Yes — re-check on docs wiki |

This is a **section semantics** pass, not a rewrite of the subcommand.

### ADR path

- Do **not** silently amend [[adr-011-sink-teleportation-vs-pagerank-damping]]
  mid-implementation.
- When this proposal is accepted, write a new `decision` (likely ADR-012) that
  **supersedes** ADR-011: default full α-damping; affirm raw display vs
  teleporting math separation; sink-only demoted to historical / pedagogical
  footnote.
- Until then ADR-011 stays `accepted` with this proposal as the supersession
  candidate.

### Delivery sequencing (blockers)

1. Merge catrace PR #39 (`NewTeleportingKernelFromAdj`, `NodeMass`, PPR helpers)
   — verified 2026-08-15: APIs on feature branch, **not** yet on `main`.
2. Clean open wikigraph PRs first (especially #54 touching `wiki.go`).
3. Estimate planning for analyze raw-vs-math section split + docs/runbook /
   integration-test re-golden (explicitly accepted).
4. Child issues remain separate under epic #56: #55 (clean `graph` + seeds),
   #22 (`goal --rank ppr`).

### Locked product defaults (from planning)

- Default restart: **uniform** (global PageRank).
- Personalization: **explicit `--seed` slugs**.
- Default `analyze` “most central”: **PageRank** (not sink-only π).
- Re-baseline docs wiki golden numbers: **yes**, pending estimate planning.

## Alternatives Considered

### A — Incremental opt-in only

Ship clean PPR viz on `graph` only; leave `analyze`/`export`/`goal` on
sink-prefill until a later “swap defaults” gate.

- **Rejected as the end state** — two π truths, delayed cost, ADR-011 limbo.
- May still appear as a *thin* intermediate PR if sequencing requires it, but
  the epic goal is the foundation rewrite, not permanent dual math.

### Keep sink-only forever; only hide edges in viz

Display/math split without adopting α-damping.

- Rejected for centrality and PPR ergonomics; spider-trap and “bored surfer”
  behaviour stay wrong relative to peer tools.

### Full product rewrite (drop MFPT / `goal`)

Literature treats PPR and MFPT as complementary, not substitutes. Rejected.

## Open Questions

1. **Communicating classes in `analyze`:** report SCCs on the **raw** graph
   (authoring signal), the **teleporting** chain (math), or both?
2. **Edge count:** confirm raw directed wikilink count as the Overview metric
   (recommended).
3. **Legacy flag:** is a `--kernel sink-only` (or similar) needed for one
   release, or is a hard cut with ADR-012 enough?
4. **Default α:** `0.15` (teleport probability, firehose / catrace style) vs
   damping `d = 0.85` naming — pick one user-facing convention and document it
   everywhere.
5. **Export formats:** should JSON/CSV/DOT export raw adj, teleporting `P`,
   or both (with clear field names)?
6. **Estimate:** sizing the analyze semantics pass + runbook/integration
   re-golden before coding starts.

## Out of Scope

- Semantic goal resolution (GitHub issues #1 / #2)
- Dropping MFPT / making `goal` PPR-only (PPR is additive via issue #22)
- Permanent dual-math (clean `graph` only; sink-prefill forever for
  `analyze` / `export`) — that was the old incremental “pass 1” idea; **not**
  this proposal. Changing default adjacency / π for `analyze` and `export`
  **is in scope**.

## Implementation Notes

- Epic: https://github.com/stephen-mcelhose/wikigraph/issues/56
- Children: https://github.com/stephen-mcelhose/wikigraph/issues/55,
  https://github.com/stephen-mcelhose/wikigraph/issues/22
- Firehose reference: `eis-intake-firehose/tools/firehose-graph`
- External pattern sources informing this proposal: Neo4j PageRank /
  Personalized PageRank docs; igraph personalized PageRank tutorial; Obsidian
  Advanced Graph View / Cartographer (PageRank-sized nodes, real edges).

## Related Proposals and Decisions

- [[adr-011-sink-teleportation-vs-pagerank-damping]] — current accepted decision;
  supersession target
- [[teleportation-ergodicity]] — mathematical background
- [[sink-page]] — structural sink reporting
- [[stationary-distribution]] — π semantics
- [[analyze]] — health report whose section semantics need planning
- [[page-type-conventions]] — this page is a `proposal` pending a future ADR

## Sources

- https://github.com/stephen-mcelhose/wikigraph/issues/56
- https://github.com/stephen-mcelhose/wikigraph/issues/55
- https://github.com/stephen-mcelhose/wikigraph/issues/22
- https://github.com/stephen-mcelhose/catrace/pull/39
- [[adr-011-sink-teleportation-vs-pagerank-damping]]
- [[teleportation-ergodicity]]
- Neo4j GDS PageRank — https://neo4j.com/docs/graph-data-science/current/algorithms/page-rank/
- igraph personalized PageRank tutorial — https://python.igraph.org/en/stable/tutorials/personalized_pagerank.html
- Page, L., Brin, S., Motwani, R., & Winograd, T. (1999). *The PageRank Citation Ranking*. http://ilpubs.stanford.edu:8090/422/
