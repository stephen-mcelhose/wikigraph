---
type: decision
title: Accepting In-Process Math for Prototype Goal Strategies
description: Decision to permit path and bottleneck strategies to implement their own linear algebra directly in cmd_goal.go, breaking the standard catrace delegation pattern, while they are validated as exploratory features.
resource: cmd_goal.go
tags: [adr, goal, prototype, catrace, architecture, path-strategy, bottleneck]
timestamp: 2026-08-10T02:41:28Z
status: accepted
---

# ADR-008 — Accepting In-Process Math for Prototype Goal Strategies

## Context

wikigraph's [[architecture|architectural principle]] is that all Markov chain mathematics is delegated to the [[catrace]] library. The `union` and `intersection` goal strategies follow this fully — they receive a `*catrace.Kernel` and call `MeanFirstPassage` on it.

Two exploratory strategies added in [[adr-007-subgraph-partitioning-and-path-strategies|ADR-007]] break this pattern:

- **`path`** — runs Dijkstra's algorithm on the raw `*mat.Dense` transition matrix directly in `cmd_goal.go`. Dijkstra is a general graph algorithm, not Markov-chain-specific math, so it is not a natural fit for catrace.
- **`bottleneck`** — computes the fundamental matrix $N = (I-Q)^{-1}$ directly in `cmd_goal.go` using `gonum/mat`. This _is_ stochastic-matrix-specific math and is therefore a natural future catrace API — but catrace does not expose it today.

Both strategies are exploratory. Their value as learning-path tools has not yet been validated against real wikis.

## Options Considered

### Option A — Require catrace APIs before shipping

Add `FundamentalMatrix(absorbing int)` to catrace, and accept that `path` cannot follow the pattern (Dijkstra doesn't belong there). Ship only after catrace is updated.

*Rejected:* Blocks exploration and couples two unrelated release cycles. catrace is a separate library; a new API requires its own design, testing, and release. Waiting prevents us from learning whether these strategies are useful at all.

### Option B — Accept in-process math as a time-boxed prototype

Ship both strategies with explicit prototype labelling in the CLI help text, strategy flag description, stderr warnings, and wiki documentation. Accept the architecture deviation as a conscious, documented risk. Prioritise catrace API work for `bottleneck` (where it belongs) once the strategy proves useful.

*Selected.*

## Decision

We accept in-process math for `path` and `bottleneck` under the following conditions:

1. **Labelled at every surface.** Both strategies are marked `[PROTOTYPE]` in CLI `--help` output, the `--strategy` flag description, and stderr when invoked.
2. **Documented in the wiki.** [[path-sequence]] and [[bottleneck-centrality]] each carry a visible callout explaining the deviation and citing this ADR.
3. **Time-boxed.** Once either strategy is validated as genuinely useful on real wikis, the architecture debt is paid:
   - `bottleneck`: add `FundamentalMatrix(absorbing int) (*mat.Dense, error)` to catrace and refactor `selectBottleneck` to use it.
   - `path`: document the intentional exception — Dijkstra on a transition matrix is not stochastic-chain math and does not belong in catrace; `selectPath` remains in wikigraph permanently.
4. **No production reliance.** The prototype label signals that the interface and scoring semantics may change.

## Consequences

- Users see clear warnings when using either prototype strategy — no silent surprises.
- The architecture deviation is discoverable via `wikigraph goal --help`, stderr, and the wiki.
- catrace remains clean; no speculative API is added before we know the strategy is worth it.
- The `path` exception is explicitly accepted as permanent (Dijkstra ≠ Markov math); `bottleneck` exception is explicitly temporary.

## Sources

- [[architecture]] — catrace delegation principle
- [[adr-007-subgraph-partitioning-and-path-strategies]] — original decision to add path and bottleneck strategies
- [[bottleneck-centrality]] — fundamental matrix math and prototype callout
- [[path-sequence]] — Dijkstra approach and prototype callout
- [[catrace]] — library that owns Markov chain linear algebra
- `cmd_goal.go`
