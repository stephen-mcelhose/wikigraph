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

It uses **[[mfpt|MFPT]]** to score pages by structural proximity to your goal. The same model underlies [[analyze]] and [[graph]].

The output is an interactive HTML subgraph of the N closest pages, rendered
the same way as `wikigraph graph`.

## Goal

An interactive HTML file showing the subgraph of pages closest to your target
topic, with edge widths reflecting how strongly each link contributes to the path.

## Prerequisites

- `wikigraph` installed and on your PATH
- A wiki directory using the `[[page-slug]]` wikilink syntax (see [[adr-002-slug-resolution]])
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

## Verification

After running `wikigraph goal`, confirm:

1. stderr prints a confirmation line:
   ```
   Written: goal_graph.html (N nodes)
   ```
   If you used `-o`, the filename will match your argument.
2. The output file exists:
   ```bash
   ls -l goal_graph.html
   ```
3. The graph contains your goal page — open the file and confirm a node is labelled with your goal slug:
   ```bash
   open goal_graph.html
   ```

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
- If the goal page is in a transient class (no inbound links from other pages),
  most pages will score as unreachable — this surfaces as `trace failed`;
  add wikilinks pointing to the goal page from related pages to resolve it
- Very large `--top` values on large wikis can be slow; the MFPT computation
  is O(pages²) per goal

## Troubleshooting

| Symptom                                      | Cause                                           | Fix                                              |
| -------------------------------------------- | ----------------------------------------------- | ------------------------------------------------ |
| `unknown --goal slug: <slug>`                | Slug doesn't match any `.md` filename           | Check exact filename with `ls ~/notes/*.md`      |
| `trace failed (try increasing --top)`        | Subset is not strongly connected                | Increase `--top` so more bridging pages included |
| `trace failed` and increasing `--top` doesn't help | Goal page is in a transient class — no other pages link to it | Add wikilinks pointing to the goal page from related pages |
| Goal page missing from HTML                  | Shouldn't happen — goal is always included      | File a bug                                       |
| All nodes the same size                      | Single-page subset or uniform trace distribution | Increase `--top`                                |
| MFPT computation is very slow                | Large wiki, many goals                          | Reduce `--top` or number of `--goal` flags       |

## See also

- [[mfpt]] — the metric that powers goal ranking
- [[analyze]] — surface under-linked and over-linked pages in your wiki
- [[graph]] — render your entire wiki as an interactive graph

## Sources

- `cmd_goal.go`
- https://github.com/stephen-mcelhose/wikigraph/issues/1
