---
type: how-to
title: How to analyse your wiki's health
description: Use wikigraph analyze to find orphaned pages, dead ends, structural clusters, and pages that would benefit most from new links.
tags: [analyze, health, orphans, sinks, markov, entropy]
resource: cmd_analyze.go
timestamp: 2026-08-09T07:31:56Z
---

# How to analyse your wiki's health

`wikigraph analyze` prints a six-section report that tells you which pages
nobody visits, which pages lead nowhere, how your wiki clusters into separate
regions, and where to add links to knit it together.

It is the most information-dense subcommand. This guide works through each
section so you know exactly what to act on. The same structure is visualised
interactively by [[graph]], and the data
can be exported for external tools via [[export]].

## Goal

A terminal report that tells you:
1. How big and connected your wiki is
2. Which pages are isolated clusters
3. Which pages need inbound links (orphans)
4. Which pages need outbound links (sinks)
5. Which pages are the most important hubs
6. Specific pairs of pages that should probably be linked

## Prerequisites

- `wikigraph` installed and on your PATH
- A wiki directory: one `.md` file per page, named `<slug>.md`, with
  `[[slug]]` wikilinks in the body

## Steps

### 1. Run the basic health report

```bash
wikigraph analyze ~/notes
```

You will see six sections printed to stdout. stderr shows `Pages: N` as a
sanity check.

### 2. Read the Overview section

```
=== Overview ===
Pages:        47
Edges:        183
Entropy rate: 3.1204 bits
Classes:      3
```

| Field          | What it means                                                                      |
| -------------- | ---------------------------------------------------------------------------------- |
| **Pages**      | Total number of `.md` files loaded (after exclusions)                              |
| **Edges**      | Number of directed links with probability > 1e-10                                 |
| **Entropy rate** | Bits per step of the random walk. Higher = more evenly spread link structure.   |
| **Classes**    | Number of communicating classes (see next section)                                |

A high entropy rate (close to log₂(Pages)) means the walker has genuine choice
at each step — a sign of a well-connected wiki. A low entropy rate means a few
pages dominate.

### 3. Read the Communicating classes section

```
=== Communicating classes ===
Class 1 (recurrent): 44 page(s)
  machine-learning
  neural-networks
  ...
Class 2 (transient — add links out of this class): 2 page(s)
  quantum-draft
  error-correction
Class 3 (recurrent): 1 page(s)
  appendix-a
```

A **recurrent** class is a strongly connected component: every page in it can
reach every other page. The random walker eventually settles into recurrent
classes permanently.

A **transient** class is a cluster that has no links back to the main graph.
Once the random walker leaves, it never returns. Pages in transient classes
are structurally isolated — add at least one outbound link from the cluster
to the main graph to fix this.

A single-page recurrent class that is also isolated from the rest of the wiki
indicates a completely orphaned page (no inbound *and* no outbound links
within a larger component).

**Action:** For each transient class, open those pages and add wikilinks
pointing into the main (largest) recurrent class.

### 4. Read the Orphan pages section

```
=== Orphan pages (bottom 10% by stationary distribution) ===
  quantum-draft                             π=0.000312  → add inbound links
  error-correction                          π=0.000418  → add inbound links
```

Orphans are pages in the **bottom N% by stationary probability** — pages the
random walker rarely visits because few other pages link to them.

The default threshold is the bottom 10% (`--orphan-pct 0.10`). Widen it to
see more candidates:

```bash
wikigraph analyze ~/notes --orphan-pct 0.20
```

**Action:** For each orphan, find pages in your wiki whose *content* is related
and add `[[orphan-slug]]` wikilinks from those pages to the orphan. The
section title tells you the exact remedy: "add inbound links".

> Note: orphans are identified by *low π*, not by zero in-degree. A page can
> have several inbound links and still have a low π if those links come only
> from other low-traffic pages.

### 5. Read the Sink pages section

```
=== Sink pages (no outgoing links) ===
  appendix-a                                → add outgoing links
  glossary-terms                            → add outgoing links
```

Sinks have **no outbound wikilinks at all**. In the [[markov-model]], a visitor
landing on a sink teleports uniformly to any page. This is handled gracefully
by wikigraph, but sinks are usually an oversight — pages that were written
in isolation and never linked forward.

**Action:** Open each sink page and add at least one `[[slug]]` wikilink to a
related page. Even one link is enough to stop it being a sink.

### 6. Read the Most central section

```
=== Most central (top 5 by stationary distribution) ===
  1. machine-learning                       π=0.084231
  2. neural-networks                        π=0.071509
  3. linear-algebra                         π=0.063812
  4. probability                            π=0.058201
  5. calculus                               π=0.047388
```

These are your **hubs**: pages that a random walker visits most often. They
are the most influential pages for spreading information across your wiki.

Use this list to:
- Ensure hubs are high quality and up to date
- Add links *from* hubs to content you want readers to discover
- Check that hubs genuinely deserve their centrality (are they linked from
  many pages because they are important, or because they are unavoidable?)

### 7. Read the Suggested missing links section

```
=== Suggested missing links (lowest commute time, not yet linked) ===
  machine-learning:
    → linear-algebra                        (commute: 4.21)
    → probability                           (commute: 5.83)
  neural-networks:
    → backpropagation                       (commute: 6.10)
```

For each of the top-N pages (by centrality), this section lists unlinked
pairs sorted by **commute time** — how quickly the random walk bounces between
them. A low commute time means the two pages are already structurally close
(many indirect paths), so adding a direct link would be a natural improvement.

The default is the top 5 pages (`--suggest-top 5`). Increase it:

```bash
wikigraph analyze ~/notes --suggest-top 10
```

Disable suggestions entirely (faster on large wikis):

```bash
wikigraph analyze ~/notes --suggest-top 0
```

**Action:** For each suggestion, read both pages. If the topic connection is
genuine, add `[[target-slug]]` to the source page's body. If you want to
explore why two pages are close, see [[goal]]
for MFPT-based subgraph visualisation. The commute-time algorithm is explained in [[mfpt]].

### 8. Full example with all flags

```bash
wikigraph analyze ~/notes \
  --orphan-pct 0.15 \
  --suggest-top 10 \
  --exclude index --exclude log --exclude AGENTS --exclude README
```

## Verification

- Stderr shows `Pages: N` matching the number of `.md` files you expect
- Each of the six sections appears in stdout
- The Communicating classes section shows at least one class

## Interpreting the numbers together

A healthy wiki typically looks like this:

| Signal                          | Healthy            | Needs work                      |
| ------------------------------- | ------------------ | ------------------------------- |
| Classes                         | 1 recurrent class  | Multiple classes, any transient |
| Orphan π values                 | > 0.005            | < 0.001                         |
| Sink count                      | 0                  | Many sinks                      |
| Entropy rate vs log₂(pages)     | > 70%              | < 50%                           |
| Top hub π                       | < 0.15             | > 0.25 (one page dominates)     |

## Troubleshooting

| Symptom                                      | Cause                                        | Fix                                              |
| -------------------------------------------- | -------------------------------------------- | ------------------------------------------------ |
| All pages are in one transient class         | No page links back to another                | Add cross-links; check wikilink syntax `[[slug]]` |
| Suggestions section is very slow             | Large wiki with many pages                   | Use `--suggest-top 0` or `--suggest-top 3`       |
| Expected page is missing from the report     | Wrong directory or not a `.md` file          | Check file extension and directory path           |
| π values for all pages are identical         | Wiki has only one page, or all sinks         | Add real wikilinks between pages                  |
| `commute: +Inf` in suggestions               | Page is unreachable from target (transient)  | Fix the transient class first                     |

## Sources

- `cmd_analyze.go`
- `wiki.go`
