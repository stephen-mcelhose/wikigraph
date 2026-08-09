---
type: how-to
title: How to find a learning path through your wiki
status: draft
description: Use wikigraph goal to surface the pages closest to a target topic and visualise the route your wiki suggests for getting there.
tags: [how-to, goal, learning-path, mfpt, subgraph]
resource: cmd_goal.go
timestamp: 2026-08-09T07:31:56Z
---

> [!WARNING]
> **Prototype — exact slugs only. The interface will change.**
>
> `wikigraph goal` currently requires goals to be exact page slugs (e.g. `error-correction`).
> Natural-language goals like `"understand quantum error correction"` are **not yet supported**.
>
> Two open issues track the planned upgrade:
> - [#1 — semantic goal support and strategy goals](https://github.com/stephen-mcelhose/wikigraph/issues/1): natural-language `--goal` queries, semantic nearest-neighbour resolution, strategy goals
> - [#2 — `wikigraph vectorize`](https://github.com/stephen-mcelhose/wikigraph/issues/2): the local embedding layer that semantic search depends on
>
> See [[adr-001-embedding-layer]] for the architectural decision. **Do not rely on the current slug-only behaviour in production workflows — the interface will change when semantic search lands.**

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
- The quantum-go example wiki cloned locally:
  ```bash
  git clone https://github.com/stephen-mcelhose/quantum-go ~/quantum-go
  ```
- A wiki directory using the `[[page-slug]]` wikilink syntax (see [[adr-002-slug-resolution]])
- You know the **exact slug** (filename without `.md`) of your goal page(s)

## Steps

### 1. Find your goal slug

Goal pages must match an exact slug. If you are unsure of the slug, list
your wiki pages:

```bash
ls ~/quantum-go/wiki/*.md | xargs -I{} basename {} .md | sort
```

Pick the slug(s) you want to navigate toward.

### 2. Run with a single goal

```bash
wikigraph goal ~/quantum-go/wiki --goal error-correction
```

Output: `goal_graph.html` (default). Open it:

```bash
open goal_graph.html
```

You will see the 10 pages (default `--top 10`) closest to `error-correction`
by MFPT, connected by their actual wikilinks. The goal page itself is always
included regardless of `--top`.

### 3. Set multiple goals

If your target topic spans several pages, pass multiple `--goal` flags. Each
page's score is the *minimum* MFPT across all goals — it is pulled into the
subgraph if it is close to *any* of the goals.

```bash
wikigraph goal ~/quantum-go/wiki \
  --goal error-correction \
  --goal shors-algorithm \
  --goal grovers-algorithm
```

This is useful for topic clusters (e.g. "quantum algorithms") where
no single page fully represents the target.

### 4. Expand or shrink the subgraph

```bash
# Wider view — 20 closest pages
wikigraph goal ~/quantum-go/wiki --goal error-correction --top 20 -o error-correction-path.html

# Focused view — only the 5 closest
wikigraph goal ~/quantum-go/wiki --goal error-correction --top 5
```

If `--top` is larger than the number of reachable pages, wikigraph uses
whatever it can reach.

**Action:** Open the output file and confirm the node count in the title bar or page matches the `--top` value you passed.

### 5. Save to a specific file

```bash
wikigraph goal ~/quantum-go/wiki \
  --goal error-correction \
  --top 12 \
  -o /tmp/goal-path.html
```

**Action:** Confirm the file landed: `ls -l /tmp/goal-path.html`

### 6. Filter weak edges

```bash
wikigraph goal ~/quantum-go/wiki --goal error-correction --min-edge 0.02
```

Raises the edge visibility threshold (default: 0.005). Useful if the
subgraph is cluttered.

**Action:** Open `goal_graph.html` and confirm the graph has fewer visible edges than the default run.

## Verification

After running `wikigraph goal`, confirm:

1. stderr confirms the run completed. For the default `--top 10` against `~/quantum-go/wiki`:
   ```
   Pages: 36
   Written: goal_graph.html (10 nodes)
   ```
   If you used `-o`, the filename shown will match your argument.
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

- **Large nodes** — high [[stationary-distribution|stationary probability]] *within the subgraph*. These
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
| `unknown --goal slug: <slug>`                | Slug doesn't match any `.md` filename           | Check exact filename with `ls ~/quantum-go/wiki/*.md` |
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
