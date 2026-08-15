---
type: decision
title: Path-Relative Slugs in Recursive Mode
description: Supersedes ADR-006. In recursive mode, slugs are derived from the path relative to the wiki root rather than the bare filename stem, eliminating collision errors in structured folder wikis.
tags: [adr, slug, recursive, wikilink, wiki-structure]
timestamp: 2026-08-12T23:38:21Z
status: accepted
supersedes: adr-006-recursive-vault-traversal
---

# ADR-010 — Path-Relative Slugs in Recursive Mode

## Context

ADR-006 introduced `-r/--recursive` traversal with basename slug matching (`slug = filename` without `.md`) and a strict collision policy: if two files share a basename across directories, `wikigraph` halts with an error.

This works for Obsidian-style vaults where filenames are globally unique by convention. It breaks for structured folder wikis — documentation repos, POC portfolios, onboarding guides — where each project directory legitimately contains files with identical names (`01-discovery.md`, `02-feasibility.md`, `index.md`). These wikis cannot be analysed at all with `--recursive`.

The `--relative-links` feature (issue #51) compounds the problem: it forces `recursive = true` and builds a path-to-slug reverse map, but since slugs were still bare stems the collision error fired before any link resolution could occur.

## Decision

In the context of supporting structured folder wikis with repeated filenames, facing the need to make `--recursive` and `--relative-links` viable for non-Obsidian repositories:

1. **Path-relative slugs in recursive mode**: When `-r/--recursive` is active, slugs are derived from the file's path relative to the wiki root directory, with the `.md` extension stripped and path separators normalised to `/`. Example: `firehose/canada-grower360/01-discovery.md` → slug `canada-grower360/01-discovery`.

2. **Flat mode unchanged**: Non-recursive mode continues to use bare filename stems. No existing users of flat mode are affected.

3. **Suffix-match exclude**: `--exclude <name>` suppresses any page whose path-relative slug ends with `/<name>` or equals `<name>` exactly. This preserves the pre-existing behaviour where `--exclude index` suppresses every `index.md` at any directory depth.

4. **Lenient wikilink resolution**: `[[basename]]` references are resolved with a two-step lookup. First, an exact match against the slug index is attempted (covers flat mode and path-relative slugs authored in full). If no exact match is found, the index is searched for any slug whose basename component equals the reference. If exactly one slug matches, the link resolves. If multiple slugs share the basename (ambiguous), a warning is printed to stderr and the link is dropped.

5. **`[[subdir/page]]` syntax deferred**: The `wikilinkRe` regex is not changed. Authoring `[[canada-grower360/01-discovery]]` in wikilink form is not yet supported. This can be addressed in a follow-up if demand arises.

## Consequences

- Structured folder wikis with repeated filenames can be analysed with `--recursive` and `--relative-links` without pre-processing.
- Output (graph node labels, `analyze` report, `export` CSV/JSON/DOT) shows path-relative slugs in recursive mode (e.g. `canada-grower360/01-discovery`).
- `--goal` requires full path-relative slugs in recursive mode (e.g. `--goal canada-grower360/01-discovery`).
- Existing `[[basename]]` wikilinks in Obsidian-style vaults continue to resolve via the lenient fallback, as long as each basename is globally unique within the vault. Ambiguous basenames produce a stderr warning rather than a silent wrong edge.
- `--relative-links` (`buildAdjacencyWithOpts`) benefits automatically: its `pathToSlug` inverse map is now keyed by path-relative slugs, which is exactly what relative-path resolution produces.

## Sources

- [[adr-006-recursive-vault-traversal]] (superseded)
- GitHub Issue #53
