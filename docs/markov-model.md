---
type: concept
title: Markov Model of a Wiki
description: How wikigraph translates `[[wikilinks]]` into a row-stochastic transition matrix — the loadPages and buildAdjacency pipeline in wiki.go.
resource: wiki.go
tags: [markov, transition-matrix, wikilinks, adjacency, sink, teleportation]
timestamp: 2026-08-09T07:31:56Z
---

# Markov Model of a Wiki

wikigraph models a wiki as a **discrete-time Markov chain**: pages are
states, and outgoing wikilinks (`[[slug]]` syntax) define the transitions. A
[[random-walk|random walker]] follows links at random; the long-run fraction
of time spent on each page is its [[stationary-distribution|stationary distribution]]
π — used as centrality throughout the tool.

The full pipeline is in `wiki.go` and detailed in [[architecture]].

## Step 1 — loadPages

`os.ReadDir` lists all `.md` files in the wiki directory. Each file that
is not in the exclude set becomes a **slug** (filename without `.md`).
Slugs are sorted alphabetically; their sorted position is their matrix
index. Sorting gives a deterministic matrix regardless of filesystem order.
The rationale for the flat slug convention is in [[adr-002-slug-resolution]].

The exclude set is built from the `--exclude` flag (defaults: `index`,
`log`, `AGENTS`).

## Step 2 — Wikilink extraction (buildAdjacency)

Each page is read and matched against:

```
\[\[([A-Za-z][A-Za-z0-9-]*)(?:\|[^\]]+)?\]\]
```

- **Capture group 1** is the slug (the `[[slug|alias]]` form strips the alias)
- Match is case-insensitive: the slug is lowercased before index lookup
- Self-links (`j == i`) are discarded
- Duplicate links count as one (deduplicated via a `map[int]bool`)

## Step 3 — Adjacency matrix

A `gonum/mat.Dense` n×n matrix is filled with 1.0 for each outgoing link.
`catrace.NewRandomWalkKernel` then **row-normalises**: each row is divided
by its sum, producing a row-stochastic matrix where each entry Pᵢⱼ is the
probability of moving from page i to page j in one step.

## Sink handling — teleportation

A page with no outgoing links (a **[[sink-page|sink]]**) would produce a zero row —
making the matrix substochastic (row sum = 0, violating the stochastic property). To keep the
chain well-defined, `buildAdjacency` gives sinks a **uniform row**: equal
probability of jumping to any page. This is the same teleportation trick
used in PageRank.

Within each [[communicating-classes|communicating class]] the Kernel is irreducible, which
guarantees a unique stationary distribution per class. The Kernel is also aperiodic for any
class that contains at least one sink page, since the uniform teleportation row carries a
self-loop. Classes with no sinks may be periodic if the underlying link graph is bipartite.
See [[analyze]] for what communicating classes mean for wiki health.

## What the matrix represents

After normalisation, Pᵢⱼ is the **one-step transition probability** from
page i to page j. The stationary distribution π satisfies πP = π and gives
the long-run visit frequency for each page. Pages with high π are hubs —
many random walks converge on them. See [[mfpt]] for how mean first passage
time uses this matrix to measure structural distance between pages.

## Sources

- `wiki.go`
- `cmd_analyze.go`
