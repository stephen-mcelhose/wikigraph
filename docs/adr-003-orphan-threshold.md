---
type: decision
title: Accept Low Stationary Probability for ADR and Proposal Pages
description: ADR and proposal pages will structurally fall in the bottom 10% orphan band in a small wiki; this is expected and should not be fixed by adding artificial cross-links.
tags: [adr, orphan, stationary-distribution, wiki-governance]
timestamp: 2026-08-09T08:21:01Z
---

# ADR-003 — Accept low stationary probability for ADR and proposal pages

**Date:** 2026-08-09  
**Status:** Accepted

---

## Context

During a lint pass ([[log]]) on the 19-page wiki, `wikigraph analyze` flagged two
pages as orphans (bottom 10% by [[stationary-distribution|stationary distribution]]):

```
how-to-docs-plan    π=0.019279  → add inbound links
adr-001-embedding-layer  π=0.019290  → add inbound links
```

Both pages have multiple real inbound links (e.g. `analyze` → `how-to-docs-plan`,
`adr-001` → `adr-002`, `architecture` → `how-to-docs-plan`). Attempting to raise
their π by adding more cross-links caused the orphan band to shift: each time a link
was added to resolve one page, a neighbouring page fell below the threshold instead.
This is a geometric tail problem — in a 19-page wiki, a 10% orphan threshold will
always flag `floor(19 × 0.10) ≈ 2` pages, regardless of how many links are added.

**Two options were considered.**

**Option A — Treat ADR/proposal orphan warnings as defects to fix**

Chase the bottom-10% band by adding cross-links until every page exceeds some
minimum π. In a small wiki this is effectively impossible: every link added shifts
the threshold, and the pages that fall naturally low (governance, decisions,
proposals) are low for structural reasons — they are not frequently cited in
conceptual explanations.

**Option B — Accept low π as the correct state for governance pages**

ADR and proposal pages (`type: decision`, `type: proposal`) are governance
artefacts. They document why a decision was made; they are not concepts that
other pages cite regularly. Their low stationary probability is structurally
correct — a random walker following wikilinks spends less time on ADRs than on
core concepts like [[markov-model]] or [[stationary-distribution]]. The orphan
warning is a false positive for this page type.

## Decision

**Option B — accept low π for governance pages.**

The 10% orphan threshold in `wikigraph analyze` is calibrated for content pages
(concepts, how-tos, runbooks) where a low π genuinely signals poor integration.
ADR and proposal pages have a different function: they record decisions and
rationale, not concepts to be traversed. Their π is expected to be low, and
artificially inflating it by adding mechanical cross-links degrades the quality
of the links that already exist.

## Consequences

- **Lint policy:** During wiki lint passes, ADR and proposal pages in the bottom
  10% orphan band should be noted but not acted upon, unless the page has
  **zero** inbound links from non-excluded pages.
- **Real orphan signal:** A governance page is genuinely orphaned only if
  `grep -r '\[\[<slug>\]\]' docs/*.md` returns no results outside `index.md` and
  `log.md`. A non-zero inbound link count is sufficient.
- **Threshold guidance:** For wikis with a significant proportion of governance
  pages, consider running `wikigraph analyze --orphan-pct 0.05` to tighten the
  band and avoid noisy false positives.
- **No change to tooling:** `wikigraph analyze` is correct. The 10% threshold is
  a sensible default for content-heavy wikis; the interpretation of results
  requires human judgement for governance page types.

## Sources

- [[analyze]] — `--orphan-pct` flag documentation
- [[stationary-distribution]] — why π measures structural centrality, not importance
- [[how-to-docs-plan]] — example of a proposal page with low but acceptable π
- [[adr-001-embedding-layer]] — example of a decision page with low but acceptable π
