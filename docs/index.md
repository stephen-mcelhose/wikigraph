---
type: index
title: wikigraph Docs Index
description: Master index of all documentation pages in the wikigraph wiki.
timestamp: 2026-08-09T06:54:46Z
---

# wikigraph Docs Index

| Page                         | Type       | Description                                                    |
| ---------------------------- | ---------- | -------------------------------------------------------------- |
| [[export]]                   | `how-to`   | Export graph data as JSON, CSV, or DOT for external tools      |
| [[graph]]                    | `how-to`   | Generate an interactive force-directed HTML graph              |
| [[analyze]]                  | `how-to`   | Wiki health report — orphans, sinks, hubs, link suggestions    |
| [[goal]]                     | `how-to`   | Learning paths via MFPT — draft                                |
| [[adr-001-embedding-layer]]  | `decision` | ADR: chromem-go + Ollama chosen for semantic goal resolution   |
| [[adr-002-slug-resolution]]  | `decision` | ADR: flat wiki layout with basename slugs over path-qualified  |
| [[adr-003-orphan-threshold]] | `decision` | ADR: accept low π for ADR/proposal pages — not a defect       |
| [[testing-runbook]]          | `runbook`  | Manual test plan covering all subcommands and edge cases       |
| [[how-to-docs-plan]]         | `proposal` | Proposal tracking the how-to docs initiative (GitHub issue #3) |
| [[architecture]]             | `concept`  | Seven Go files, catrace dependency, and data-flow pipeline     |
| [[markov-model]]             | `concept`  | How wikilinks become a row-stochastic Markov kernel            |
| [[mfpt]]                     | `concept`  | Mean first passage time — used by goal and analyze             |
| [[communicating-classes]]    | `concept`  | Maximal mutually-reachable page sets; the wiki connectivity test |
| [[recurrent-class]]          | `concept`  | Recurrent vs transient classes — what gets π, what doesn't     |
| [[random-walk]]              | `concept`  | Random walk on a graph — the foundational model behind all analysis   |
| [[stationary-distribution]]  | `concept`  | π — long-run visit frequency, centrality, orphan detection            |
| [[entropy-rate]]             | `concept`  | Bits per step of the random walk — a wiki health signal               |
| [[sink-page]]                | `concept`  | Zero-outgoing-link pages and the teleportation fix                    |
| [[commute-time]]             | `concept`  | Symmetric MFPT-based distance metric for link suggestions             |
| [[catrace]]                  | `concept`  | The Go library providing all Markov chain mathematics                 |
