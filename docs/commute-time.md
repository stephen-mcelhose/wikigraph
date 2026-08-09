---
type: concept
title: Commute Time
description: CT(i,j) = MFPT(i,j) + MFPT(j,i) — a symmetric graph distance metric used by wikigraph to suggest missing links between structurally close pages.
resource: cmd_analyze.go
tags: [commute-time, mfpt, graph-distance, link-suggestion, random-walk]
timestamp: 2026-08-09T08:05:55Z
---

# Commute Time

**Commute time** CT(i, j) is the expected number of steps for a [[random-walk]]
to travel from page i to page j *and back* to page i:

```
CT(i, j) = MFPT(i, j) + MFPT(j, i)
```

where [[mfpt|MFPT(i, j)]] is the mean first passage time from i to j. See
[[mfpt]] for how MFPT is computed via the fundamental matrix.

## Why symmetry matters

MFPT alone is asymmetric — MFPT(i, j) ≠ MFPT(j, i) in general, because the
link structure in each direction may differ. Commute time adds both directions,
making it **symmetric**: CT(i, j) = CT(j, i).

This symmetry is what makes commute time a proper distance metric. You can
compare any two page pairs on equal footing and rank them by structural
proximity without worrying about which direction is easier.

## Relationship to graph resistance distance

In an undirected graph, commute time is exactly proportional to the
**effective resistance** between two nodes in an equivalent electrical network
(each edge is a unit resistor). For directed wiki graphs the analogy is
approximate but the intuition holds: pages connected by many indirect paths
have low effective resistance and low commute time.

## Reading the Suggested missing links section

[[analyze]] uses commute time to recommend where to add wikilinks:

```
=== Suggested missing links (lowest commute time, not yet linked) ===
  machine-learning:
    → linear-algebra    (commute: 4.21)
    → probability       (commute: 5.83)
```

Low commute time means the pages are already structurally close — the walk
bounces between them quickly via indirect paths. Adding a direct `[[slug]]`
wikilink formalises that relationship for readers.

Pairs are ranked ascending by commute time; the top `--suggest-top` (default 5)
are shown. Pass `--suggest-top 0` to skip the section and save computation
on large wikis.

## Infinite commute time

CT(i, j) is infinite when i and j are in different
[[communicating-classes|communicating classes]] with no return path. Such
pairs are skipped by [[analyze]] — fixing the class structure (see
[[recurrent-class]]) is required before suggestions become meaningful.

## Relationship to MFPT

Commute time is a derived concept built entirely on [[mfpt]]. The
[[mfpt]] page covers the fundamental matrix computation; this page covers the
symmetric combination and its role as a link-suggestion metric.

## Sources

- [[catrace]] — `Kernel.CommuteTime(i, j int)`, `Kernel.MeanFirstPassage(i, j int)`
- [`cmd_analyze.go` — `kern.CommuteTime`, suggested missing links](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_analyze.go)
- [Commute time — Wikipedia](https://en.wikipedia.org/wiki/Hitting_time#Commute_time)
- [Effective resistance — Wikipedia](https://en.wikipedia.org/wiki/Resistance_distance)
