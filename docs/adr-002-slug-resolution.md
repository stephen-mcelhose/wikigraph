---
type: decision
title: Slug Resolution for Nested Wiki Directories
description: Chose a flat wiki layout with basename slugs over path-qualified slugs, keeping [[wikilinks]] natural and decoupled from filesystem structure.
tags: [adr, slug, wikilink, wiki-structure]
timestamp: 2026-08-09T07:19:55Z
status: superseded
superseded_by: adr-006-recursive-vault-traversal
---

# ADR-002 — Slug resolution for nested wiki directories

## Context

After establishing `docs/` as the wiki root, pages were organised into subdirectories (`how-to/`, `adr/`, `proposal/`) for filesystem clarity. This forced a choice: what slug does `docs/how-to/analyze.md` get? The decision affects every `[[wikilink]]` in the wiki — for example, links to [[analyze]], [[goal]], and [[adr-001-embedding-layer]].

Two options were considered:
- **Option A (Path-qualified slugs)**: `how-to/analyze`. Leaks filesystem structure into link grammar and creates churn whenever files move.
- **Option B (Flat layout with basename slugs)**: All pages live directly in `docs/`. Slug = filename without `.md`. Matches mainstream wiki conventions.

## Decision

In the context of resolving wikilinks across docs, facing breaking links whenever files move between subdirectories, we decided to adopt a flat layout with basename slugs directly in `docs/`, to achieve simple, natural wikilinks decoupled from filesystem structure, accepting that filenames must remain globally unique across the wiki.

## Consequences

- All pages live at `docs/<slug>.md` — no subdirectories ever.
- Slugs are kebab-case basenames (`analyze`, `adr-001-embedding-layer`, `testing-runbook`).
- Wikilinks are `[[slug]]` — never `[[dir/slug]]`.
- `AGENTS.md` and `llm-wiki` SKILL enforce the flat rule.
- `wiki.go` uses `os.ReadDir` (non-recursive).

## Sources

- [Obsidian wikilink resolution](https://help.obsidian.md/Linking+notes+and+files/Internal+links)
- [Karpathy LLM-wiki pattern](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- [[how-to-docs-plan]]
