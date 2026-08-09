---
type: concept
title: Sink Pages
description: A sink page has no outgoing wikilinks — its zero row in the adjacency matrix is fixed by uniform teleportation, the same trick used by PageRank.
resource: wiki.go
tags: [sink, teleportation, markov, pagerank, adjacency]
timestamp: 2026-08-09T08:05:55Z
---

# Sink Pages

A **sink page** is a wiki page with zero outgoing `[[wikilinks]]`. In the
adjacency matrix its row is all zeros — a [[random-walk]] arriving here has
nowhere to go next.

The [[analyze]] command reports sinks explicitly:

```
=== Sink pages (no outgoing links) ===
  appendix-a                                → add outgoing links
```

## Why zero rows break the Markov model

A row-stochastic [[markov-model|transition matrix]] requires every row to sum
to 1. A zero row is sub-stochastic: probability "leaks out" of the system.
Without a fix, the walk gets absorbed at sink pages and the
[[stationary-distribution]] is no longer well-defined for the rest of the
graph.

## The uniform teleportation fix

`buildAdjacency` in `wiki.go` detects zero rows and replaces them with
uniform probability:

```
Pᵢⱼ = 1/N   for all j   (if page i is a sink)
```

The walker landing on a sink teleports to any page with equal probability.
This is identical to the teleportation step in Google's PageRank and keeps
the chain ergodic. See [[markov-model]] for the full construction.

## Impact on the graph

Teleportation from sinks injects probability mass uniformly across all pages.
This can slightly inflate the [[stationary-distribution|π]] of otherwise
low-traffic pages. Once a real wikilink replaces the teleportation row, the
π distribution sharpens to reflect actual link structure.

A single wikilink is enough to stop a page being a sink — `buildAdjacency`
only triggers teleportation when the out-degree is exactly zero. See
[[architecture]] for the full `buildAdjacency` → [[catrace]] pipeline.

## How to fix a sink

Open the sink page and add at least one `[[slug]]` wikilink to a related page.
The [[analyze]] output tells you exactly which pages are sinks and labels each
with `→ add outgoing links`.

## Relationship to communicating classes

Sinks with teleportation are treated as if they link to all pages. They do not
form separate [[communicating-classes|communicating classes]] — the walk can
always escape via the teleportation rows. Without the fix, a sink would be an
absorbing state and would form its own trivial recurrent class.

## Sources

- [`wiki.go` — `buildAdjacency`, sink row detection](https://github.com/stephen-mcelhose/wikigraph/blob/main/wiki.go)
- [`cmd_analyze.go` — sink reporting section](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_analyze.go)
- [PageRank — Wikipedia](https://en.wikipedia.org/wiki/PageRank)
- [Absorbing Markov chain — Wikipedia](https://en.wikipedia.org/wiki/Absorbing_Markov_chain)
