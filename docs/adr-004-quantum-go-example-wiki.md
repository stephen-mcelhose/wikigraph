---
type: decision
title: "ADR-004: Use quantum-go as the canonical example wiki in how-to docs"
description: Use quantum-go as the canonical example wiki in all how-to guides to make verification reproducible and realistic.
status: accepted
timestamp: 2026-08-09T09:00:00Z
tags: [adr, documentation, examples]
---

# ADR-004: Use quantum-go as the canonical example wiki in how-to docs

## Context

The how-to guides for `wikigraph` ([[graph]], [[analyze]], [[goal]], [[export]]) need a concrete wiki to use as the running example throughout steps and verification sections.

Options considered:
- **`~/notes`**: Generic placeholder. Unclonable and unverifiable against real CLI execution.
- **`./docs`**: Repo's own wiki. Creates self-referential documentation.
- **`~/quantum-go/wiki`**: A real third-party Go project wiki.

## Decision

In the context of authoring how-to guides for CLI subcommands, facing unverifiable example outputs from generic path placeholders, we decided to use `~/quantum-go/wiki` as the canonical example wiki, to achieve reproducible, real-world step verification in documentation, accepting that readers must clone `github.com/stephen-mcelhose/quantum-go` before following guides.

## Consequences

- **Verifiable**: Readers can clone quantum-go and reproduce exact CLI outputs shown in docs.
- **Tool-agnostic**: Demonstrates wikigraph working on an independent project wiki.
- **Maintenance**: Verification numbers (Pages, Edges, Entropy rate) are re-checked if quantum-go grows.

## Sources

- [quantum-go repository](https://github.com/stephen-mcelhose/quantum-go)
- [[how-to-docs-plan]]
- [[graph]]
- [[analyze]]
- [[goal]]
- [[export]]
