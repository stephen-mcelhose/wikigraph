---
type: proposal
title: Per-Type Structural Conventions for Wiki Pages
description: Proposal to define required and optional sections for each of the six wiki page types, making authorship consistent for both humans and the LLM.
resource: https://github.com/stephen-mcelhose/wikigraph/issues/24
tags: [documentation, conventions, wiki, llm-wiki]
timestamp: 2026-08-09T21:00:52Z
---

# Per-Type Structural Conventions for Wiki Pages

The wiki uses six page types — `concept`, `how-to`, `decision`, `runbook`, `proposal`, `spike` — but only defines them by name. No structural conventions exist: no required sections, no ordering, no tone guidance. This proposal defines the work needed to fill that gap.

## Current state

Structure has emerged inconsistently from practice:

| Type       | Observed sections                                                              | Formalised? |
| ---------- | ------------------------------------------------------------------------------ | ----------- |
| `how-to`   | Goal → Prerequisites → Steps → Verification → Troubleshooting → Sources       | Partially — rubric exists in [[how-to-docs-plan]] but not in [[AGENTS.md]] |
| `decision` | Varies by ADR — no consistent Status / Context / Decision / Consequences shape | No |
| `runbook`  | Prerequisites → TC-## test cases (inferred from [[testing-runbook]])           | No |
| `proposal` | Problem → Proposed Solution → Alternatives → Open Questions → Sources          | No |
| `concept`  | Fully freeform — each page organises itself differently                        | No |
| `spike`    | No pages exist yet                                                             | N/A |

## What needs to be decided

For each type, agree on:

1. **Required sections** — must be present; lint should flag absence
2. **Optional sections** — expected but not mandatory
3. **Section order** — canonical top-to-bottom sequence
4. **Tone and scope** — e.g. "concept pages explain the idea to a reader unfamiliar with Markov chains"

## Open questions

- **`decision`**: standard ADR format (Status / Context / Decision / Consequences) or something lighter? Standard ADR has Status as the first field — useful for scanning open vs. accepted decisions.
- **`concept`**: is freeform acceptable, or should we enforce a minimal skeleton (Definition → Properties → Related concepts)? Risk of over-constraining pages that have natural narrative flow.
- **`spike`**: what distinguishes a spike from a proposal? A proposal is a design question seeking a decision; a spike is an investigation with findings. The natural shape is Question → Method → Findings → Conclusions → Open questions.
- **`how-to`**: the existing rubric (issue #11) should be the source of truth — does it need to be resolved into this work or tracked separately?
- **Where do conventions live?** Inline in `AGENTS.md` (one authoritative place, always read by the LLM), or a separate `conventions.md` page that `AGENTS.md` references?

## Out of scope

- Rewriting existing pages to conform (follow-on task, tracked separately)
- Changing OKF frontmatter fields

## Sources

- [GitHub issue #24 — docs: define per-type structural conventions for wiki pages](https://github.com/stephen-mcelhose/wikigraph/issues/24)
