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

## [2026-08-10] lint | 30 pages checked, 0 orphans, index updated with 4 new pages (adr-007, absorbing-markov-chain, bottleneck-centrality, path-sequence)

## [2026-08-10] lint | 30 pages checked, 7 issues found, 7 fixed
- Stale claim: architecture.md line count updated ~850 → ~1,567; cmd_goal.go file-map entry updated to list all 4 strategies
- Missing xref: added [[adr-007]] link to bottleneck-centrality.md, path-sequence.md, absorbing-markov-chain.md
- Missing xref: added [[absorbing-markov-chain]] link to adr-007 bottleneck decision
- Missing xref: added [[adr-007]] to goal.md See also section
- Low-π orphan: quickstart (π=0.005) resolved — added inbound link from analyze.md (π=0.107)
- adr-007 (π=0.004) resolved — 3 new inbound links from the new concept pages
- Cannot fix without new source: llm-wiki-pattern (π=0.003), path-sequence (π=0.009) remain low — linked only from low-traffic pages; subagent GAN review pending

## [2026-08-10] decision | ADR-008 — prototype math strategies accepted
- path and bottleneck marked [PROTOTYPE] in CLI --help, --strategy flag description, and stderr on invocation
- bottleneck-centrality.md and path-sequence.md each received a [!WARNING] callout citing ADR-008
- ADR-008 written: accepts in-process math as time-boxed; bottleneck exception temporary (catrace API planned), path exception permanent (Dijkstra ≠ Markov math)
- index.md updated with ADR-008

## [2026-08-10] lint | 31 pages checked, 2 issues found, 1 fixed, 1 accepted
- Missing xref: added [[adr-008]] to architecture.md Related Concepts (π: 0.004 → 0.008)
- Low-π accepted: llm-wiki-pattern, adr-005, adr-004 — all meta/ADR pages, inherently low-traffic, no forced links added
- OKF frontmatter: all 31 pages pass; 1 recurrent class, 0 sinks, entropy 2.34 bits
- Added adr-007, absorbing-markov-chain, bottleneck-centrality, path-sequence to index.md
- Verified OKF frontmatter across all 33 pages
- Updated log.md

## [2026-08-11] lint | 35 pages checked (post-rebase on main), 14 issues found, 14 fixed
## [2026-08-13] amend | adr-011-sink-teleportation-vs-pagerank-damping
- Revised decision: accept noisy sink visualization (star pattern) as known consequence of ergodicity
- Removed "separate display adj from math adj" — not viable: catrace rejects zero rows, lazy π, self-loops create non-ergodic display chain
- Updated [[sink-page]] "noisy sink visualization" section to match revised decision

## [2026-08-13] ingest | teleportation-ergodicity + adr-011-sink-teleportation-vs-pagerank-damping
- New concept page: [[teleportation-ergodicity]] — sink-only teleportation vs full PageRank α-damping; canonical reference https://en.wikipedia.org/wiki/PageRank
- New decision: [[adr-011-sink-teleportation-vs-pagerank-damping]] — retain sink-only teleportation; separate display adj from math adj; revisit if spider traps, π asymmetries, or catrace α-damping support emerge
- Propagated to [[sink-page]]: added display vs math adjacency section, ADR-011 and teleportation-ergodicity cross-refs
- Propagated to [[stationary-distribution]]: added teleportation-ergodicity and sink-page cross-refs
- index.md updated; 38 pages total

## [2026-08-13] ingest | adr-010-path-relative-slugs — path-relative slugs in recursive mode (issue #53, PR #54)
- Renamed from adr-007-path-relative-slugs to adr-010 (ADR-007 number already taken by subgraph-partitioning)
- Decision: recursive mode uses path-relative slugs (e.g. subdir/page) instead of bare stems; flat mode unchanged
- Lenient wikilink fallback: [[basename]] resolves when unique across all slugs; warns and drops when ambiguous
- --exclude now matches by basename at any depth in recursive mode
- wikilinkRe broadened to accept digit-leading slugs ([[02-feasibility]] etc.)
- Cross-references added: architecture, markov-model, adr-006 all link to [[adr-010-path-relative-slugs]]
- testing-runbook: 3 broken [[adr-007-path-relative-slugs]] wikilinks fixed to [[adr-010-path-relative-slugs]]
- testing-runbook: TCs 23-27 added covering recursive mode, path-relative slugs, --relative-links, and portfolio wiki
- wiki_integration_test.go added as automated guard for docs/ wiki graph invariants
- index.md updated with ADR-010 row; 36 pages total
- OKF frontmatter missing: added to graph-topologies.md, graph-models.md, graphs/caveman.md, graphs/barbell.md, graphs/karate-club.md, graphs/krackhardt-kite.md
- Section rename: `## Links` → `## Sources` in graph-topologies, graph-models; `## References` → `## Sources` in all four graph write-ups
- Index gap: added graph-topologies, graph-models, adr-009-wiki-gen-make-vs-buy to index.md
- Transient class: graph-topologies, graph-models, adr-009 formed isolated transient class — fixed by adding `## Related Concepts` wikilinks in both concept pages and inbound links from architecture.md
- Broken wikilink: `[[wiki-gen]]` in adr-009 Sources replaced with plain text (page does not exist yet)
- Stale sources: adr-009 text references to research/ files upgraded to `[[graph-topologies]]` and `[[graph-models]]` wikilinks
- Testing-runbook: updated to 35 pages, 176 edges, entropy 2.41 bits, new orphan/central values
- Advisory (not fixed): graphs/*.md in subdirectory violate flat layout — provisional pending #42; low π on graph-topologies (0.005), graph-models (0.004), adr-009 (0.002) accepted per ADR-003
- Final state: 35 pages, 176 edges, 1 recurrent class, 0 sinks, entropy 2.41 bits
