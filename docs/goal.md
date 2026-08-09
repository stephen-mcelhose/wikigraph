---
type: how-to
title: How to find a learning path through your wiki
status: draft
description: Use wikigraph goal to surface the pages closest to a target topic and visualise the route your wiki suggests for getting there.
tags: [goal, learning-path, mfpt, subgraph]
resource: cmd_goal.go
timestamp: 2026-08-09T07:31:56Z
---

> **Draft** — `wikigraph goal` is a prototype. Goals must currently be exact
> page slugs. Natural-language goals ("understand quantum error correction")
> are planned. See [issue #1](https://github.com/stephen-mcelhose/wikigraph/issues/1)
> and [[adr-001-embedding-layer]] for the planned approach.

# How to find a learning path through your wiki

`wikigraph goal` answers the question: *given a topic I want to understand,
which pages in my wiki are the most structurally relevant stepping stones?*

It uses **[[mfpt|Mean First Passage Time]] (MFPT)** — a Markov chain metric that
measures how quickly a random walker starting at page X will reach your goal
page. The same random-walk model underlies [[analyze]]
(commute time suggestions) and [[graph]]
(node size and colour). Pages with a low MFPT are the ones your wiki's link structure naturally
funnels readers toward the goal through. These are your learning path.

The output is an interactive HTML subgraph of the N closest pages, rendered
the same way as `wikigraph graph`.

## Goal

An interactive HTML file showing the subgraph of pages closest to your target
topic, with edge widths reflecting how strongly each link contributes to the path.

## Prerequisites

- `wikigraph` installed and on your PATH
- A wiki directory with `[[slug]]` wikilinks
- You know the **exact slug** (filename without `.md`) of your goal page(s)

## Steps

### 1. Find your goal slug

Goal pages must match an exact slug. If you are unsure of the slug, list
your wiki pages:

```bash
ls ~/notes/*.md | xargs -I{} basename {} .md | sort
```

Pick the slug(s) you want to navigate toward.

### 2. Run with a single goal

```bash
wikigraph goal ~/notes --goal machine-learning
```

Output: `goal_graph.html` (default). Open it:

```bash
open goal_graph.html
```

You will see the 10 pages (default `--top 10`) closest to `machine-learning`
by MFPT, connected by their actual wikilinks. The goal page itself is always
included regardless of `--top`.

### 3. Set multiple goals

If your target topic spans several pages, pass multiple `--goal` flags. Each
page's score is the *minimum* MFPT across all goals — it is pulled into the
subgraph if it is close to *any* of the goals.

```bash
wikigraph goal ~/notes \
  --goal machine-learning \
  --goal neural-networks \
  --goal backpropagation
```

This is useful for topic clusters (e.g. "deep learning fundamentals") where
no single page fully represents the target.

### 4. Expand or shrink the subgraph

```bash
# Wider view — 20 closest pages
wikigraph goal ~/notes --goal machine-learning --top 20 -o learning-path.html

# Focused view — only the 5 closest
wikigraph goal ~/notes --goal machine-learning --top 5
```

If `--top` is larger than the number of reachable pages, wikigraph uses
whatever it can reach.

### 5. Save to a specific file

```bash
wikigraph goal ~/notes \
  --goal machine-learning \
  --top 12 \
  -o /tmp/ml-path.html
```

### 6. Filter weak edges

```bash
wikigraph goal ~/notes --goal machine-learning --min-edge 0.02
```

Raises the edge visibility threshold (default: 0.005). Useful if the
subgraph is cluttered.

## Interpreting the output

The subgraph uses the **trace kernel** on the selected subset — the effective
transition probabilities *within* those pages, accounting for paths that
leave and re-enter the subset.

- **Large nodes** — high stationary probability *within the subgraph*. These
  are the most-visited pages on the path to your goal.
- **Wide edges** — high probability that the random walker follows that link
  next. These are the canonical stepping stones.
- **The goal page(s)** — always present; often large, often central.

Pages that appear large and are not the goal are your **gateway concepts**:
mastering them first will naturally lead you to the goal.

## Limitations (prototype)

- Goals must be exact slugs — there is no fuzzy matching yet
- If a goal slug is unreachable (e.g. it is in a transient class with no
  inbound links), wikigraph will skip it and warn you
- Very large `--top` values on large wikis can be slow; the MFPT computation
  is O(pages²) per goal

## Troubleshooting

| Symptom                                      | Cause                                           | Fix                                              |
| -------------------------------------------- | ----------------------------------------------- | ------------------------------------------------ |
| `unknown --goal slug: <slug>`                | Slug doesn't match any `.md` filename           | Check exact filename with `ls ~/notes/*.md`      |
| `trace failed (try increasing --top)`        | Subset is not strongly connected                | Increase `--top` so more bridging pages included |
| Goal page missing from HTML                  | Shouldn't happen — goal is always included      | File a bug                                       |
| All nodes the same size                      | Single-page subset or uniform trace distribution | Increase `--top`                                |
| MFPT computation is very slow                | Large wiki, many goals                          | Reduce `--top` or number of `--goal` flags       |

## Sources

- `cmd_goal.go`
- https://github.com/stephen-mcelhose/wikigraph/issues/1
