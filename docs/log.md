# Wiki Log

<!-- Append-only. Never edit existing entries. -->

## [2026-08-09] init | wiki initialized at docs/ in wikigraph repository

## [2026-08-09] ingest | how-to/export, how-to/graph, how-to/analyze, how-to/goal — initial how-to guides written from source code

## [2026-08-09] lint | 11 pages checked, 4 issues found, 4 fixed
## [2026-08-09] edit | testing-runbook — added inline links at TC-02, TC-07, TC-11, TC-16 (first TC of each subcommand group)
## [2026-08-09] refactor | flatten wiki — remove subdirectories, use basename slugs
- docs/ is now flat: analyze, export, goal, graph, adr-001-embedding-layer, testing-runbook, how-to-docs-plan
- all [[wikilinks]] updated to [[slug]] format (slug = filename without .md)
- AGENTS.md updated: no-subdirectory rule, [[slug]] convention
- llm-wiki SKILL.md updated: flat page rule in conventions table and AGENTS.md template
- testing-runbook: all TCs updated to docs/ path and real output values (7 pages, verified)
- ADR-002: documents slug resolution decision (basename over path-qualified)
- Added how-to/index.md to master index (index gap)
- Added [[How to find a learning path through your wiki]] link to adr/001-embedding-layer (missing cross-reference)
- Added Completed Guides section to proposal/how-to-docs-plan (orphan — only linked from index)
- Added See Also section to testing-runbook (orphan — only linked from index)
