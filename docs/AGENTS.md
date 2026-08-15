# Wiki Schema

This wiki is maintained by an LLM using the `llm-wiki` skill.
Wiki root: `docs/` in the wikigraph repository.

## Domain

Documentation for the `wikigraph` CLI — a tool that models a markdown wiki as
a Markov chain and exposes graph-theoretic analysis: stationary distribution,
communicating classes, MFPT-based learning paths, and multi-format export.

The same Markov craft (PageRank π, Personalized PageRank via `--seed`/`--alpha`,
MFPT, entropy) is the bridge from knowledge-graph analysis to PDA agents and
multi-agent networks in [[catrace]] — see [[knowledge-graph-to-pda-agents]].

Topics in scope:
- How-to guides for each subcommand (`graph`, `analyze`, `goal`, `export`)
- Architecture decision records (ADRs)
- Operational runbooks (testing, releases)
- Proposals and design docs
- Concept pages for Markov / PageRank foundations and the KG→PDA agent leap

## Conventions

- **Page slugs**: kebab-case, flat — all pages live directly in `docs/`, never in subdirectories (e.g., `analyze.md`, `adr-001-embedding-layer.md`)
- **Frontmatter**: OKF — `type`, `title`, `description`, `timestamp` required; `resource`, `tags`, `status` optional
- **Types**: `concept` | `how-to` | `decision` | `runbook` | `proposal` | `spike` (for required section schemas by page type, see [[page-type-conventions]])
- **Cross-references**: `[[slug]]` wikilinks — the slug is the filename without `.md`
- **Sources section**: every page ends with `## Sources` listing URLs or file paths it was derived from
- **No subdirectories**: filesystem layout is flat. Organisation emerges from wikilinks and the graph, not from folders.

## Operations

Run these via the `llm-wiki` skill (wiki root = `docs/`):

- `ingest <source>` — read a new source (URL, file, GitHub issue), write or update a page, propagate wikilinks
- `query <question>` — synthesize an answer from wiki pages
- `lint` — audit for orphans, missing frontmatter, stale claims, broken wikilinks

## index.md

Structured catalog of all wiki pages. Updated on every write operation.

## log.md

Append-only chronological log. Format: `## [YYYY-MM-DD] operation | detail`

## Verification gate

**After any code or docs/ change, run the integration test before considering work done:**

```bash
cd ~/repos/wikigraph
go test -run TestDocsWiki_Integration -v
```

`TestDocsWiki_Integration` in `wiki_integration_test.go` calls `buildKernel("docs/", ...)` directly and asserts page count, edge count, class count, most central page, and lowest-π page. It is the automated guard for the manual runbook. A failure means either:
- The constants in `wiki_integration_test.go` need updating (wiki changed intentionally), or
- A code change broke the graph analysis (regression).

**Never ignore a failing integration test.** Fix it or update the constants with the correct values, then update the runbook to match.

## testing-runbook.md

**Always update `testing-runbook.md` whenever pages are added or removed from the wiki.**

The runbook hardcodes expected `wikigraph analyze` output values — page counts, edge counts, entropy rate, orphan π values, and the full page list. These go stale the moment the wiki changes. After any ingest, lint, or structural edit that changes the page count:

1. Run `go test -run TestDocsWiki_Integration` — if it fails, update the constants in `wiki_integration_test.go` first.
2. Run `wikigraph analyze docs/` to get current human-readable values.
3. Update the header block (pages, edges, entropy, "Last verified").
4. Update the full page list in the Prerequisites section.
5. Update every hardcoded `Pages: N` expected value in the TC pass criteria.
6. Update the TC-16 spot-check table (Overview row, Classes row, Orphans row, Most central row).
7. Update TC-18 (lowest-π page identity and π value) and TC-19 (count equals N).
8. Update TC-20 exclusion expectations (`Pages: N-1` after excluding one page).
9. Update TC-12 CSV node row count.
