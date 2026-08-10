# Wiki Log

<!-- Append-only. Never edit existing entries. -->

## [2026-08-09] init | wiki initialized at docs/ in wikigraph repository

## [2026-08-09] ingest | how-to/export, how-to/graph, how-to/analyze, how-to/goal — initial how-to guides written from source code

## [2026-08-09] lint | 11 pages checked, 4 issues found, 4 fixed
## [2026-08-09] edit | testing-runbook — added inline links at TC-02, TC-07, TC-11, TC-16 (first TC of each subcommand group)
## [2026-08-09] ingest | wikigraph source code — 3 concept pages (architecture, markov-model, mfpt)
- Lint: added resource: + ## Sources to all 4 how-to pages; fixed both orphans (testing-runbook, how-to-docs-plan)
- architecture: 7-file map, catrace API, data-flow pipeline, --exclude behaviour
- markov-model: loadPages → buildAdjacency → buildKernel; wikilink regex; sink teleportation
- mfpt: MFPT definition; use in goal (ranking + trace kernel); use in analyze (commute time suggestions)
- Result: 11 pages, 36 edges, 1 recurrent class (was 3 classes), entropy 1.65 bits
## [2026-08-09] refactor | flatten wiki — remove subdirectories, use basename slugs
- docs/ is now flat: analyze, export, goal, graph, adr-001-embedding-layer, testing-runbook, how-to-docs-plan
- all [[wikilinks]] updated to [[slug]] format (slug = filename without .md)
- AGENTS.md updated: no-subdirectory rule, [[slug]] convention
- llm-wiki SKILL.md updated: flat page rule in conventions table and AGENTS.md template
- testing-runbook: all TCs updated to docs/ path and real output values (7 pages, verified)
- ADR-002: documents slug resolution decision (basename over path-qualified)
- Added how-to/index.md to master index (index gap)
- Added [[How to find a learning path through your wiki]] link to [[adr-001-embedding-layer]] (missing cross-reference)
- Added Completed Guides section to proposal/how-to-docs-plan (orphan — only linked from index)
- Added See Also section to testing-runbook (orphan — only linked from index)
## [2026-08-09] ingest | communicating-classes, recurrent-class concept pages
- communicating-classes: Tarjan SCC, recurrent/transient labelling, how to read analyze output, failure modes
- recurrent-class: recurrent vs transient distinction, π implications, practical fix guide
- cross-refs added: analyze, architecture, markov-model all link to new pages
## [2026-08-09] lint | 19 pages checked, 8 issues found, 8 fixed; 1 risk accepted
- Missing ## Sources: added to adr-001-embedding-layer, how-to-docs-plan, testing-runbook
- Stale content: testing-runbook page counts and expected output values updated (was 8 pages, now 19)
- Orphans: resolved how-to-docs-plan and adr-002-slug-resolution via targeted inbound links from higher-π pages
- Added [[catrace]] wikilinks in mfpt, stationary-distribution, commute-time, sink-page
- Added [[architecture]] to analyze Sources section
- Sources sections in analyze.md upgraded from bare filenames to GitHub URLs
- Risk accepted: ADR/proposal pages structurally land in bottom 10% by π — documented in [[adr-003-orphan-threshold]]

## [2026-08-09] ingest | adr-003-orphan-threshold — accept low π for governance pages
- Decision: low π on ADR/proposal pages is expected, not a defect requiring artificial cross-links
- Triggered by lint pass finding how-to-docs-plan and adr-001-embedding-layer in orphan band
- Final state: 20 pages, 125 edges, 1 recurrent class, 0 sinks

## [2026-08-09] ingest | concept pages for issue #6

Created six new concept pages: [[random-walk]], [[stationary-distribution]],
[[entropy-rate]], [[sink-page]], [[commute-time]], [[catrace]]. Updated
cross-references in markov-model, mfpt, communicating-classes, recurrent-class,
architecture, and analyze. Corrected SCC algorithm attribution from Tarjan to
Kosaraju in communicating-classes and architecture. Registered all six pages
in [[index]]. Closes https://github.com/stephen-mcelhose/wikigraph/issues/6.

## [2026-08-09] lint | 21 pages checked, 4 issues found, 4 fixed
- adr-004-quantum-go-example-wiki: Changed non-standard type to decision, added missing description field, and resolved sink status by linking to graph, analyze, goal, and export
- goal: Adjusted tags list to start with how-to per rubric
- index: Fixed index gap by registering adr-004-quantum-go-example-wiki
- testing-runbook: Updated hardcoded page counts, edge counts, and spot-check values to stay synchronized with the expanded 21-page wiki state

## [2026-08-09] ingest | issue #24 — per-type structural conventions proposal → page-type-conventions.md
- adr-005: Added ADR-005 establishing per-type required section schemas, Y-statements, redundant status heading deprecation, and deferred proposal storage subdirectories.
- AGENTS.md: Updated conventions section to point to [[page-type-conventions]] for structural schemas.

## [2026-08-09] review | 7 pages critically reviewed — 9 issues fixed, 3 advisory
- goal: Added [[stationary-distribution]] link at "large nodes" description
- mfpt: Removed backticks from inside [[catrace]] wikilink alias; added [[recurrent-class]] and [[communicating-classes]] links in Infinite MFPT section
- stationary-distribution: Removed backticks from inside [[catrace]] wikilink alias; removed stray backtick after closing bracket
- architecture: Removed backticks from inside [[catrace]] wikilink alias
- analyze: Fixed [[entropy-rate\|Entropy rate]] → [[entropy-rate]] in table (escaped pipe broke wikilink alias parsing); added [[random-walk]] link in sink section
- markov-model: Added [[catrace]] link on catrace.NewRandomWalkKernel
- Advisory (not fixed): index.md could group pages by type; architecture.md ASCII diagram subcommand files could be linked; stationary-distribution.md Sources has hardcoded GitHub URLs

## [2026-08-09] lint | 23 pages checked, 11 concept pages migrated to required section schema (Overview -> Key Properties -> Related Concepts -> Sources); 0 remaining schema violations.
- goal: Fixed unclosed code block at section 5 (Action text was inside the fenced block)
- architecture: Added [[commute-time]] link (prose said "commute time" without linking); added [[sink-page|sink]] link in data-flow pipeline
- stationary-distribution: Added [[commute-time]] link in command table
- [[how-to/analyze]] in adr-002-slug-resolution: intentional (illustrates rejected Option A) — left as-is

## [2026-08-09] lint | 26 pages checked, 0 orphans, 0 sinks, 1 recurrent class, entropy 2.27 bits
- Added quickstart and llm-wiki-pattern pages
- Verified OKF frontmatter across all 26 pages
- Fixed cross-references and updated index.md and testing-runbook.md
