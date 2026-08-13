---
type: decision
title: Recursive Vault Traversal and Slug Resolution
description: Supersedes ADR-002. Introduces opt-in recursive vault directory walking via -r/--recursive flag while preserving basename slug matching and prohibiting duplicate basenames.
tags: [adr, slug, recursive, obsidian, wikilink, wiki-structure]
timestamp: 2026-08-09T16:00:00Z
status: accepted
supersedes: adr-002-slug-resolution
superseded_by: adr-007-path-relative-slugs
---

# ADR-006 — Recursive Vault Traversal and Slug Resolution

## Context

ADR-002 established a flat wiki layout requirement (`docs/<slug>.md`) for `wikigraph`'s internal documentation and core assumption in `wiki.go`. However, real-world Knowledge Bases (PKM) such as Obsidian, Logseq, and Foam frequently organize notes into nested directory hierarchies (e.g. `projects/`, `areas/`, `archive/`).

Requiring users to flatten nested vaults before running `wikigraph` limits utility. Conversely, forcing path-qualified slugs (e.g. `projects/my-note`) makes `[[wikilink]]` references brittle when files move between directories.

## Decision

In the context of analyzing nested Markdown vaults (e.g., Obsidian/Logseq), facing the need to preserve simple `[[slug]]` wikilinks without requiring flat filesystem layouts:

1. **Supersede ADR-002**: Replace ADR-002's flat-directory mandate with a flexible traversal model.
2. **Opt-in Recursive Traversal**: Introduce a persistent `-r, --recursive` CLI flag across all `wikigraph` subcommands (`graph`, `goal`, `analyze`, `export`).
3. **Basename Slug Matching**: Maintain basename slug lookup (`slug = filename` without `.md`), matching default Obsidian behavior. `[[my-note]]` resolves to `my-note.md` regardless of directory depth.
4. **Ignored Directories**: Automatically ignore hidden directories (`.git`, `.obsidian`, `.logseq`, `.trash`, `.vscode`) and assets during recursive walking.
5. **Collision Policy**: Enforce strict uniqueness on file basenames across nested subdirectories. If `work/idea.md` and `archive/idea.md` both exist, `wikigraph` halts with a clear collision error: `duplicate slug "idea" found in "work/idea.md" and "archive/idea.md"`.

## Consequences

- Existing flat wikis continue to work without flags or behavioral changes.
- Complex Obsidian and Logseq vaults can be parsed via `-r / --recursive`.
- `buildKernel`, `loadPages`, and `buildAdjacency` in `wiki.go` now map `slug -> relativeFilePath`.
- Slugs remain globally unique across subdirectories.

## Sources

- [Obsidian internal link resolution](https://help.obsidian.md/Linking+notes+and+files/Internal+links)
- [[adr-002-slug-resolution]] (superseded)
- GitHub Issue #26
