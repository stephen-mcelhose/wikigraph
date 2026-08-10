---
type: how-to
title: How to analyse your wiki's health
description: Use wikigraph analyze to find orphaned pages, dead ends, structural clusters, and pages that would benefit most from new links.
tags: [how-to, analyze, health, orphans, sinks, markov, entropy]
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
- The example wiki cloned locally:
  ```bash
  git clone https://github.com/stephen-mcelhose/quantum-go ~/quantum-go
  ```
- A wiki directory: one `.md` file per page, named `<slug>.md`, with
  `[[slug]]` wikilinks in the body

## Steps

### 1. Run the basic health report

```bash
wikigraph analyze ~/quantum-go/wiki
```

You will see six sections printed to stdout.

### 2. Read the Overview section

> **Note:** The examples in steps 2–3 use an illustrative wiki with multiple
> communicating classes to show all output formats. Running against
> `~/quantum-go/wiki` produces `Pages: 36, Edges: 295, Entropy rate: 2.2003
> bits, Classes: 1` — a single healthy class with no transient pages.

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
| **Classes**    | Number of [[communicating-classes]] (see next section)                            |

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

A **[[recurrent-class|recurrent]]** class is a strongly connected component: every page in it can
reach every other page. The random walker eventually settles into recurrent
classes permanently.

A **[[recurrent-class|transient]]** class is a cluster that has no links back to the main graph.
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
  how-to-add-a-new-gate                     π=0.001603  → add inbound links
  gate-zoo                                  π=0.001749  → add inbound links
```

Orphans are pages in the **bottom N% by [[stationary-distribution|stationary probability]]**
— pages the [[random-walk|random walker]] rarely visits because few other pages link to them.

The default threshold is the bottom 10% (`--orphan-pct 0.10`). Widen it to
see more candidates:

```bash
wikigraph analyze ~/quantum-go/wiki --orphan-pct 0.20
```

**Action:** For each orphan, find pages in your wiki whose *content* is related
and add `[[orphan-slug]]` wikilinks from those pages to the orphan. The
section title tells you the exact remedy: "add inbound links".

> Note: orphans are identified by *low π*, not by zero in-degree. A page can
> have several inbound links and still have a low π if those links come only
> from other low-traffic pages. ADR and proposal pages are expected to score
> low by this metric — see [[adr-003-orphan-threshold]] for the rationale.

### 5. Read the Sink pages section

```
=== Sink pages (no outgoing links) ===
  algorithm-comparison                      → add outgoing links
  fuzz-testing                              → add outgoing links
```

[[sink-page|Sinks]] have **no outbound wikilinks at all**. In the [[markov-model]], a [[random-walk|random walker]]
landing on a sink teleports uniformly to any page. This is handled gracefully
by wikigraph, but sinks are usually an oversight — pages that were written
in isolation and never linked forward.

**Action:** Open each sink page and add at least one `[[slug]]` wikilink to a
related page. Even one link is enough to stop it being a sink.

### 6. Read the Most central section

```
=== Most central (top 5 by stationary distribution) ===
  1. composite-gates                        π=0.096414
  2. gate-application                       π=0.090011
  3. quantum-linear-algebra                 π=0.076762
  4. simulator-optimizations                π=0.067835
  5. quantum-dsl                            π=0.057854
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
  quantum-linear-algebra:
    → composite-gates                       (commute: 21.43)
    → quantum-dsl                           (commute: 31.80)
  simulator-optimizations:
    → quantum-dsl                           (commute: 36.06)
    → grovers-algorithm                     (commute: 39.46)
```

For each of the top-N pages (by centrality), this section lists unlinked
pairs sorted by **[[commute-time|commute time]]** — how quickly the [[random-walk|random walk]] bounces between
them. A low commute time means the two pages are already structurally close
(many indirect paths), so adding a direct link would be a natural improvement.

The default is the top 3 pages (`--suggest-top 3`). Increase it:

```bash
wikigraph analyze ~/quantum-go/wiki --suggest-top 10
```

Disable suggestions entirely (faster on large wikis):

```bash
wikigraph analyze ~/quantum-go/wiki --suggest-top 0
```

**Action:** For each suggestion, read both pages. If the topic connection is
genuine, add `[[target-slug]]` to the source page's body. If you want to
explore why two pages are close, see [[goal]]
for MFPT-based subgraph visualisation. The commute-time algorithm is explained in [[mfpt]].

### 8. Full example with all flags

```bash
wikigraph analyze ~/quantum-go/wiki \
  --orphan-pct 0.15 \
  --suggest-top 10 \
  --exclude index --exclude log --exclude AGENTS --exclude README
```

### 9. Summarise the report with an LLM

Capture the full output and paste it into your preferred LLM with this prompt:

```
wikigraph analyze ~/quantum-go/wiki > analyze-report.txt
```

**Suggested prompt:**

```
Below is the output of `wikigraph analyze` on my wiki. Please:

1. State in one sentence whether the wiki is structurally healthy or needs work,
   citing the entropy rate relative to log₂(pages) and the class count.
2. List the top 3 highest-impact actions, ordered by urgency. For each action,
   name the specific page slug(s) to act on and what to do (add inbound link,
   add outbound link, merge page, etc.).
3. From the Suggested missing links section, pick the single best link to add
   and explain in one sentence why that pair is structurally close.

<paste contents of analyze-report.txt here>
```

The prompt asks for slug-level specifics because vague advice ("improve your orphan pages") is not actionable. The entropy rate framing (see [[entropy-rate]]) gives the LLM the right benchmark: log₂(36) ≈ 5.17 bits is the theoretical maximum for a 36-page wiki; a healthy wiki sits above 70% of that.

## Verification

Run against a known wiki and check stdout directly:

```
wikigraph analyze ~/quantum-go/wiki
```

Expected stdout begins with:

```
=== Overview ===
Pages:        36
Edges:        295
Entropy rate: 2.2003 bits
Classes:      1
```

Confirm:
- All six sections (`Overview`, `Communicating classes`, `Orphan pages`, `Sink pages`, `Most central`, `Suggested missing links`) appear in stdout
- `Pages:` matches the number of `.md` files in your wiki
- The Communicating classes section lists at least one class

## Interpreting the numbers together

A healthy wiki typically looks like this:

| Signal                          | Healthy            | Needs work                      |
| ------------------------------- | ------------------ | ------------------------------- |
| Classes                         | 1 recurrent class  | Multiple classes, any transient |
| Orphan π values                 | > 0.005            | < 0.001                         |
| Sink count                      | 0                  | Many sinks                      |
| [[entropy-rate]] vs log₂(pages)               | > 70%  | < 50%                           |
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

- [`cmd_analyze.go`](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_analyze.go)
- [`wiki.go`](https://github.com/stephen-mcelhose/wikigraph/blob/main/wiki.go)
- [[architecture]] — full data-flow pipeline from wikilinks to Markov output
- [[how-to-docs-plan]] — proposal that drove the creation of these guides
- [[quickstart]] — getting-started guide covering installation through first analysis run
