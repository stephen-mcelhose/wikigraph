---
type: index
title: wikigraph Docs Index
description: Master index of all documentation pages in the wikigraph wiki.
timestamp: 2026-08-09T06:54:46Z
---

# wikigraph Docs Index

| Page                         | Type       | Description                                                    |
| ---------------------------- | ---------- | -------------------------------------------------------------- |
| [[quickstart]]               | `how-to`   | Quickstart guide: setting up llm-wiki, visualising, analysing, and linting |
| [[llm-wiki-pattern]]         | `concept`  | Compounding knowledge base architecture (Karpathy pattern)     |
| [[export]]                   | `how-to`   | Export graph data as JSON, CSV, or DOT for external tools      |
| [[graph]]                    | `how-to`   | Generate an interactive force-directed HTML graph              |
| [[analyze]]                  | `how-to`   | Wiki health report — orphans, sinks, hubs, link suggestions    |
| [[goal]]                     | `how-to`   | Learning paths via MFPT — draft                                |
| [[adr-001-embedding-layer]]  | `decision` | ADR: chromem-go + Ollama chosen for semantic goal resolution   |
| [[adr-002-slug-resolution]]  | `decision` | ADR: flat wiki layout with basename slugs (superseded by ADR-006) |
| [[adr-006-recursive-vault-traversal]] | `decision` | ADR: recursive vault traversal and basename slug resolution |
| [[adr-003-orphan-threshold]] | `decision` | ADR: accept low π for ADR/proposal pages — not a defect       |
| [[adr-004-quantum-go-example-wiki]] | `decision` | ADR: Use quantum-go as the canonical example wiki in how-to docs |
| [[adr-005-page-type-conventions-and-proposal-storage]] | `decision` | ADR: Page-type structural conventions and proposal/spike storage strategy |
| [[testing-runbook]]          | `runbook`  | Manual test plan covering all subcommands and edge cases       |
| [[how-to-docs-plan]]         | `proposal` | Proposal tracking the how-to docs initiative (GitHub issue #3) |
| [[page-type-conventions]]    | `proposal` | Proposal to define required sections for each wiki page type   |
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
| [[adr-007-subgraph-partitioning-and-path-strategies]] | `decision` | ADR: Subgraph partitioning and path strategies for goal subcommand |
| [[adr-008-prototype-math-strategies]]                 | `decision` | ADR: Accept in-process math for prototype path/bottleneck strategies |
| [[absorbing-markov-chain]]   | `concept`  | Absorbing states, block canonical form, and fundamental matrix N      |
| [[bottleneck-centrality]]    | `concept`  | Random-walk betweenness metric for gatekeeper/chokepoint pages        |
| [[path-sequence]]            | `concept`  | Dijkstra negative log-likelihood transition sequence strategy          |
| [[kernel-identifiability]]   | `concept`  | Which outputs (π, commute times, export, Trace) recover P, and when  |
| [[graph-topologies]]         | `concept`  | Reference catalog of named graph topologies and their Markov walk properties |
| [[graph-models]]             | `concept`  | Reference catalog of random graph generative models (ER, BA, WS, SBM, LFR)  |
| [[adr-009-wiki-gen-make-vs-buy]] | `decision` | ADR: hybrid make+buy — keep wiki-gen for hierarchical, add nx-to-wiki for named graphs |
| [[adr-010-path-relative-slugs]]       | `decision` | ADR: path-relative slugs in recursive mode; lenient wikilink basename fallback for structured folder wikis |
| [[adr-011-sink-teleportation-vs-pagerank-damping]] | `decision` | ADR: retain sink-only teleportation; separate display adj from math adj; defer full α-damping |
| [[teleportation-ergodicity]]          | `concept`  | Why sink-only teleportation is legitimate; PageRank damping factor comparison; when to reconsider |

## Sources

Generated and maintained by the `llm-wiki` skill from the pages in `docs/`.
