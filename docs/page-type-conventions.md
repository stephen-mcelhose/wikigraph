---
type: proposal
title: Per-Type Structural Conventions for Wiki Pages
description: Proposal defining required sections for each wiki page type, progressive discovery via AGENTS.md, and proposal/spike retention strategy.
resource: https://github.com/stephen-mcelhose/wikigraph/issues/24
tags: [documentation, conventions, wiki, llm-wiki]
timestamp: 2026-08-09T22:30:00Z
---

# Per-Type Structural Conventions for Wiki Pages

The wiki uses six page types (`concept`, `how-to`, `decision`, `runbook`, `proposal`, `spike`), but structural requirements were previously informal. This page defines exact section schemas, progressive discovery mechanics, and proposal/spike retention policies.

## Decision Summary

1. **Progressive Discovery**: `AGENTS.md` remains lean and points to `[[page-type-conventions]]` (or `[[adr-005-page-type-conventions-and-proposal-storage]]`). `llm-wiki` checks required section rules prior to page creation/editing.
2. **Redundant Heading Removal**: `## Status` in markdown body is deprecated/dropped for `decision` pages because OKF frontmatter `status` (`proposed|accepted|deprecated`) is the single source of truth.
3. **Upfront Y-Statement**: `decision` pages MUST begin `## Decision` with the standard Y-Statement pattern.
4. **Minimal `concept` Skeleton**: Enforces structure while retaining authoring flow.
5. **Proposal/Spike History Retention**: Deferred moving proposals/spikes into code subdirectories for now (would break flat `docs/*.md` layout invariant). Captured via [[adr-005-page-type-conventions-and-proposal-storage]]; to be revisited if tooling expands directory awareness.

## Required Sections by Page Type

| Page Type | Required Sections (Canonical Order) | Key Requirements |
| :--- | :--- | :--- |
| **`decision`** | 1. `## Context`<br>2. `## Decision`<br>3. `## Consequences`<br>4. `## Sources` | • Markdown `## Status` omitted (frontmatter `status:` is canonical)<br>• `## Decision` MUST start with Y-Statement:<br>`In the context of <X>, facing <Y>, we decided <Z>, to achieve <A>, accepting <B>.` |
| **`how-to`** | 1. `## Goal`<br>2. `## Prerequisites`<br>3. `## Steps`<br>4. `## Verification`<br>5. `## Troubleshooting`<br>6. `## Sources` | • Strictly complies with 24-point how-to rubric |
| **`concept`** | 1. `## Overview`<br>2. `## Key Properties`<br>3. `## Related Concepts`<br>4. `## Sources` | • Minimal skeleton to ensure graph interlinking without restricting mathematical flow |
| **`proposal`** | 1. `## Problem`<br>2. `## Proposed Solution`<br>3. `## Alternatives Considered`<br>4. `## Open Questions`<br>5. `## Sources` | • Value: Architectural design decisions before consensus |
| **`spike`** | 1. `## Question`<br>2. `## Method`<br>3. `## Findings`<br>4. `## Conclusions`<br>5. `## Sources` | • Value: Timeboxed technical investigation with empirical code/data findings |
| **`runbook`** | 1. `## Overview`<br>2. `## Prerequisites`<br>3. `## Test Cases`<br>4. `## Sources` | • Operational test procedures and assertions |

## Enforcement

- **`llm-wiki lint`**: Blocking enforcement. Audits page type against required top-level section headings and reports violations as errors.
- **`llm-wiki ingest`**: Blocking enforcement. Validates generated/edited pages before writing.

## Sources

- [GitHub issue #24 — docs: define per-type structural conventions for wiki pages](https://github.com/stephen-mcelhose/wikigraph/issues/24)
- [[adr-005-page-type-conventions-and-proposal-storage]]
- [[how-to-docs-plan]]
- [[index]]
