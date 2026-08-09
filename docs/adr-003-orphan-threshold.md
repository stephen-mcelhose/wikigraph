---
type: decision
title: Accept Low Stationary Probability for ADR and Proposal Pages
description: ADR and proposal pages will structurally fall in the bottom 10% orphan band in a small wiki; this is expected and should not be fixed by adding artificial cross-links.
tags: [adr, orphan, stationary-distribution, wiki-governance]
timestamp: 2026-08-09T08:21:01Z
status: accepted
---

# ADR-003 — Accept low stationary probability for ADR and proposal pages

## Context

During a lint pass on the 19-page wiki, `wikigraph analyze` flagged governance pages (`how-to-docs-plan`, `adr-001-embedding-layer`) as orphans in the bottom 10% by [[stationary-distribution|stationary distribution]].

Attempting to raise their π by adding cross-links shifted the orphan band to neighbouring pages (a geometric tail effect). Governance pages document decisions rather than concepts, so lower visit frequency is structurally expected.

Two options were considered:
- **Option A**: Treat ADR/proposal orphan warnings as defects and artificially inflate cross-links.
- **Option B**: Accept low π as structurally correct for governance pages and refrain from adding mechanical links.

## Decision

In the context of health analysis via `wikigraph analyze`, facing false-positive orphan warnings on decision and proposal pages, we decided to accept low stationary probability for governance pages, to achieve clean conceptual cross-linking without noise, accepting that governance pages will naturally cluster in the bottom stationary distribution band.

## Consequences

- **Lint policy**: ADR and proposal pages in the bottom 10% band are not acted upon unless they have zero inbound links.
- **Real orphan signal**: A governance page is genuinely orphaned only if `grep -r '\[\[<slug>\]\]' docs/*.md` returns no results outside `index.md` and `log.md`.
- **Threshold guidance**: Use `wikigraph analyze --orphan-pct 0.05` to tighten the band when governance pages predominate.

## Sources

- [[analyze]] — `--orphan-pct` flag documentation
- [[stationary-distribution]] — why π measures structural centrality
- [[how-to-docs-plan]] — example of proposal page with low π
- [[adr-001-embedding-layer]] — example of decision page with low π
