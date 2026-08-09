# Wiki Log

<!-- Append-only. Never edit existing entries. -->

## [2026-08-09] init | wiki initialized at docs/ in wikigraph repository

## [2026-08-09] ingest | how-to/export, how-to/graph, how-to/analyze, how-to/goal — initial how-to guides written from source code

## [2026-08-09] lint | 11 pages checked, 4 issues found, 4 fixed
## [2026-08-09] edit | testing-runbook — added inline links at TC-02, TC-07, TC-11, TC-16 (first TC of each subcommand group)
## [2026-08-09] refactor | migrate runbook to use docs/ as wiki under test
- wiki.go: recursive directory scanning via filepath.WalkDir; slugs are relative paths (e.g. how-to/analyze); exclusions match on full slug OR basename; wikilink regex extended to allow / and leading digits
- all [[Title style]] wikilinks updated to [[slug/format]] across 7 docs files
- runbook: all 22 TCs updated — commands, page counts, slugs, pass criteria — verified against real output
- Added how-to/index.md to master index (index gap)
- Added [[How to find a learning path through your wiki]] link to adr/001-embedding-layer (missing cross-reference)
- Added Completed Guides section to proposal/how-to-docs-plan (orphan — only linked from index)
- Added See Also section to testing-runbook (orphan — only linked from index)
