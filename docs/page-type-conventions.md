---
type: proposal
title: Per-Type Structural Conventions for Wiki Pages
description: Proposal defining required sections, optional sections, tone/audience guidance, frontmatter schemas, progressive discovery via AGENTS.md, and proposal/spike retention strategy.
resource: https://github.com/stephen-mcelhose/wikigraph/issues/24
tags: [documentation, conventions, wiki, llm-wiki]
timestamp: 2026-08-09T23:00:00Z
---

# Per-Type Structural Conventions for Wiki Pages

The wiki uses six page types (`concept`, `how-to`, `decision`, `runbook`, `proposal`, `spike`), but structural requirements were previously informal. This page defines exact section schemas, optional sections, tone/audience guidelines, required frontmatter fields, progressive discovery mechanics, and retention policies.

## Decision Summary

1. **Progressive Discovery**: `AGENTS.md` remains lean and points to `[[page-type-conventions]]` (or `[[adr-005-page-type-conventions-and-proposal-storage]]`). `llm-wiki` checks required section rules prior to page creation/editing.
2. **Redundant Heading Removal**: `## Status` in markdown body is deprecated/dropped for `decision` pages because OKF frontmatter `status` (`proposed|accepted|deprecated`) is the single source of truth.
3. **Upfront Y-Statement**: `decision` pages MUST begin `## Decision` with the standard Y-Statement pattern.
4. **Minimal `concept` Skeleton**: Enforces structure while retaining authoring flow.
5. **Proposal/Spike Lifecycle**: Proposals and spikes are immutable history records. Upon decision, a proposal links to its accepted/rejected `decision` ADR rather than changing its own type to `decision`.
6. **Proposal/Spike History Retention**: Deferred moving proposals/spikes into code subdirectories for now (would break flat `docs/*.md` layout invariant). Captured via [[adr-005-page-type-conventions-and-proposal-storage]]; to be revisited if tooling expands directory awareness.

## Frontmatter Requirements by Page Type

All pages require `type`, `title`, `description`, `timestamp`. Type-specific additions:

| Page Type | Required Frontmatter Fields |
| :--- | :--- |
| `decision` | `status: proposed\|accepted\|deprecated\|superseded` |
| `how-to` | `resource:` (CLI subcommand or guide topic) |
| `proposal` | `resource:` (GitHub issue URL) |
| `spike` | `resource:` (GitHub issue URL or research topic) |
| `runbook` | `tags:` (must include `runbook` or `qa`) |
| `concept` | `tags:` |

## Required & Optional Sections with Tone Guidance

| Page Type | Required Sections (Canonical Order) | Optional Sections | Tone & Audience |
| :--- | :--- | :--- | :--- |
| **`decision`** | 1. `## Context`<br>2. `## Decision`<br>3. `## Consequences`<br>4. `## Sources` | • `## Options Considered`<br>• `## Implementation Notes` | **Objective, analytical**. Audience: developers/architects seeking rationale. Omit markdown `## Status`. `## Decision` starts with Y-statement (`In the context of <X>, facing <Y>, we decided <Z>, to achieve <A>, accepting <B>`). |
| **`how-to`** | 1. `## Goal`<br>2. `## Prerequisites`<br>3. `## Steps`<br>4. `## Verification`<br>5. `## Troubleshooting`<br>6. `## Sources` | • `## Examples`<br>• `## Related Guides` | **Imperative, task-focused**. Audience: end users executing a task. Must satisfy 24-point how-to rubric. |
| **`concept`** | 1. `## Overview`<br>2. `## Key Properties`<br>3. `## Related Concepts`<br>4. `## Sources` | • `## Mathematical Formalism`<br>• `## Code Implementation`<br>• `## Examples` | **Explanatory, domain-oriented**. Audience: reader learning Markov chains or graph concepts. Narrative freedom within minimal skeleton. |
| **`proposal`** | 1. `## Problem`<br>2. `## Proposed Solution`<br>3. `## Alternatives Considered`<br>4. `## Open Questions`<br>5. `## Sources` | • `## Implementation Notes`<br>• `## Related Proposals and Decisions` | **Persuasive, trade-off oriented**. Audience: team evaluating architectural options before consensus. |
| **`spike`** | 1. `## Question`<br>2. `## Method`<br>3. `## Findings`<br>4. `## Conclusions`<br>5. `## Sources` | • `## Benchmark Results`<br>• `## Code Samples`<br>• `## Open Questions` | **Empirical, evidence-based**. Audience: engineers reviewing test data, POC code results, or performance benchmarks. |
| **`runbook`** | 1. `## Overview`<br>2. `## Prerequisites`<br>3. `## Test Cases`<br>4. `## Sources` | • `## Expected Outputs`<br>• `## Known Edge Cases` | **Procedural, assertion-heavy**. Audience: QA/maintainers executing manual or automated verification runs. |

## Enforcement

- **`llm-wiki lint`**: Blocking enforcement. Audits page type against required top-level section headings and reports violations as errors.
- **`llm-wiki ingest`**: Blocking enforcement. Validates generated/edited pages before writing.

## Sources

- [GitHub issue #24 — docs: define per-type structural conventions for wiki pages](https://github.com/stephen-mcelhose/wikigraph/issues/24)
- [[adr-005-page-type-conventions-and-proposal-storage]]
- [[how-to-docs-plan]]
- [[testing-runbook]]
- [[index]]
