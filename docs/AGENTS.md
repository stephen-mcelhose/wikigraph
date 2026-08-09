# Wiki Schema

This wiki is maintained by an LLM using the `llm-wiki` skill.
Wiki root: `docs/` in the wikigraph repository.

## Domain

Documentation for the `wikigraph` CLI — a tool that models a markdown wiki as
a Markov chain and exposes graph-theoretic analysis: stationary distribution,
communicating classes, MFPT-based learning paths, and multi-format export.

Topics in scope:
- How-to guides for each subcommand (`graph`, `analyze`, `goal`, `export`)
- Architecture decision records (ADRs)
- Operational runbooks (testing, releases)
- Proposals and design docs

## Conventions

- **Page slugs**: kebab-case, directory-prefixed where needed (e.g., `how-to/analyze`, `adr/001-embedding-layer`)
- **Frontmatter**: OKF — `type`, `title`, `description`, `timestamp` required; `resource`, `tags`, `status` optional
- **Types**: `concept` | `how-to` | `decision` | `runbook` | `proposal` | `spike`
- **Cross-references**: `[[Page Title]]` wikilinks — use the page's `title` frontmatter value
- **Sources section**: every page ends with `## Sources` listing URLs or file paths it was derived from
- **No raw/ folder**: sources are referenced inline via `resource:` frontmatter or `## Sources` sections

## Operations

Run these via the `llm-wiki` skill (wiki root = `docs/`):

- `ingest <source>` — read a new source (URL, file, GitHub issue), write or update a page, propagate wikilinks
- `query <question>` — synthesize an answer from wiki pages
- `lint` — audit for orphans, missing frontmatter, stale claims, broken wikilinks

## index.md

Structured catalog of all wiki pages. Updated on every write operation.

## log.md

Append-only chronological log. Format: `## [YYYY-MM-DD] operation | detail`
