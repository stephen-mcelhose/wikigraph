---
type: decision
title: Slug Resolution for Nested Wiki Directories
description: Chose a flat wiki layout with basename slugs over path-qualified slugs, keeping [[wikilinks]] natural and decoupled from filesystem structure.
tags: [adr, slug, wikilink, wiki-structure]
timestamp: 2026-08-09T07:19:55Z
status: accepted
---

# ADR-002 — Slug resolution for nested wiki directories

## Status

Accepted

## Context

After establishing `docs/` as the wiki root, pages were organised into
subdirectories (`how-to/`, `adr/`, `proposal/`) for filesystem clarity.
This forced a choice: what slug does `docs/how-to/analyze.md` get?
The decision affects every `[[wikilink]]` in the wiki — for example,
links to [[analyze]], [[goal]], and [[adr-001-embedding-layer]].

Two options were considered.

**Option A — Path-qualified slugs** (`how-to/analyze`)

Wikilinks must be written `[[how-to/analyze]]`. Globally unique by
construction. But: no mainstream wiki tool works this way (not Obsidian,
Roam, Logseq, or MediaWiki). Moving a file breaks all inbound links.
Authors must know and type the full directory path of every page they link
to. The wikilink regex must be extended to allow `/`, leaking filesystem
semantics into the link grammar.

**Option B — Flat layout with basename slugs** (`analyze`)

All pages live directly in `docs/` root. Slug = filename without `.md`.
Wikilinks are `[[analyze]]`. This matches every mainstream wiki convention.
Moving or reorganising files is impossible to do accidentally because there
are no subdirectories. Requires globally unique filenames — a mild
constraint that is easy to enforce and easy to verify.

## Decision

**Flat layout (Option B).**

The filesystem hierarchy is an organisational convenience for humans
browsing files. Wikilinks express conceptual relationships between pages.
Coupling the two — making authors write `[[how-to/analyze]]` — leaks
filesystem structure into the wiki's semantic layer. Organisation should
emerge from the graph of wikilinks, not from folders.

We also tried Option A in practice and found it created churn: every time
a page moved between subdirectories, all inbound wikilinks broke silently.
The regex extension needed to support `/` in slugs was a code smell
confirming the design was wrong.

## Implementation note — recursive scanning attempt

During implementation we attempted to support subdirectory-organised wikis by
replacing `os.ReadDir` with `filepath.Walk` in `loadPages`. This would have
allowed pages at arbitrary depths (e.g. `docs/how-to/analyze.md`) to be
discovered automatically.

**Reverted.** Recursive scanning directly violates the flat-layout standard:
it makes subdirectories a valid place to put pages, which re-introduces the
slug ambiguity problem (two files at different paths with the same basename),
undermines the "organisation emerges from wikilinks not folders" principle,
and silently breaks the moment any file moves. The standard is `os.ReadDir`
on a single flat directory — no recursion, ever.

This experience is why the standard is now **encoded** in `AGENTS.md` and
the `llm-wiki` skill: future LLM sessions must not re-introduce recursive
scanning as a "helpful improvement".

## Consequences

- All pages live at `docs/<slug>.md` — no subdirectories ever
- Slugs are kebab-case basenames: `analyze`, `adr-001-embedding-layer`, `testing-runbook`
- Wikilinks are `[[slug]]` — never `[[dir/slug]]`
- `AGENTS.md` and `llm-wiki` SKILL updated to enforce the flat rule
- `wiki.go` uses `os.ReadDir` (non-recursive) — `filepath.Walk` was tried and reverted
- The how-to guide series that triggered this decision is tracked in [[how-to-docs-plan]]

## Sources

- Obsidian wikilink resolution: https://help.obsidian.md/Linking+notes+and+files/Internal+links
- Karpathy LLM-wiki pattern: https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f
