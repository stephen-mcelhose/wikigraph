---
type: decision
title: "ADR-005: Page-Type Structural Conventions and Proposal/Spike Storage"
description: Standardises structural schemas across all six page types, adopts upfront Y-statements, drops redundant markdown status headings, and defers directory-based proposal storage.
status: accepted
resource: https://github.com/stephen-mcelhose/wikigraph/issues/24
tags: [adr, conventions, architecture, llm-wiki]
timestamp: 2026-08-09T22:35:00Z
---

# ADR-005: Page-Type Structural Conventions and Proposal/Spike Storage

## Context

`wikigraph` uses six wiki page types (`concept`, `how-to`, `decision`, `runbook`, `proposal`, `spike`), but lacked strict required schemas for each type. This led to inconsistent authoring across human and LLM contributions. Additionally, markdown `## Status` sections in ADRs duplicated frontmatter `status:` fields, decisions lacked uniform Y-statement structures, and proposals/spikes risked cluttering active documentation without clear retention/storage rules.

## Decision

In the context of standardising LLM-wiki authoring and page structures across wikigraph documentation, facing inconsistent schemas, redundant metadata, and retention ambiguity, we decided to mandate per-type required section schemas (including Y-statements for ADRs and minimal skeletons for concepts), drop redundant markdown `## Status` headings, and maintain proposals/spikes in the flat `docs/` directory for now, to achieve consistent automated linting and progressive agent discovery, accepting that moving proposals/spikes to dedicated code subdirectories is deferred until wikigraph tooling becomes directory-aware.

Key outcomes:
1. **Per-type schemas**: Mandatory required headings enforced via `llm-wiki lint` and `ingest` as specified in [[page-type-conventions]].
2. **ADR Format**: `## Decision` must start with Y-statement (`In the context of <X>, facing <Y>, we decided <Z>, to achieve <A>, accepting <B>.`). Markdown `## Status` heading is removed.
3. **`concept` Skeleton**: `Overview` → `Key Properties` → `Related Concepts` → `Sources`.
4. **Retention & Directory Layout**: Flat slug rules in `docs/` remain authoritative. Retention of completed proposals/spikes stays inside flat `docs/` until directory-aware tooling feature work is undertaken.

## Consequences

- All new and edited pages must strictly follow per-type required sections.
- `llm-wiki lint` fails if required sections are missing.
- `AGENTS.md` remains concise by delegating section schemas to [[page-type-conventions]].
- Tooling does not need to handle complex nested folder paths.

## Sources

- [GitHub issue #24 — docs: define per-type structural conventions for wiki pages](https://github.com/stephen-mcelhose/wikigraph/issues/24)
- [[page-type-conventions]]
