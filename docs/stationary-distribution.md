---
type: concept
title: Stationary Distribution
description: The stationary distribution π is the long-run visit frequency of each page under a random walk — used as centrality throughout wikigraph.
resource: cmd_analyze.go
tags: [stationary-distribution, pi, centrality, power-iteration, orphan]
timestamp: 2026-08-09T08:05:55Z
---

# Stationary Distribution

The **stationary distribution** π is a probability vector satisfying:

```
πP = π    and    Σᵢ π(i) = 1
```

It gives the fraction of time a [[random-walk]] spends on each page in the
long run. Pages with high π are structural hubs that the walk gravitates
towards; pages with low π are peripheral. wikigraph uses π as the primary
centrality measure throughout [[analyze]], [[export]], and [[graph]].

## Power iteration algorithm

[[catrace|catrace.Kernel.Stationary]](tol, maxIter) computes π iteratively:

1. Start with a uniform vector v⁰ = [1/N, 1/N, …, 1/N]
2. Multiply: vᵏ⁺¹ = vᵏ P
3. Repeat until ‖vᵏ⁺¹ − vᵏ‖₁ < `tol`

The `tol` parameter controls convergence precision (wikigraph uses 1e-9).
`maxIter` caps iterations to prevent infinite loops on near-degenerate chains.
In practice, well-connected wikis converge in tens of iterations.

## Convergence conditions

Convergence to a unique π is guaranteed when the [[markov-model|Markov chain]]
is **irreducible** (all states communicate) and **aperiodic** (no fixed cycles
force the walk to return only at multiples of some period d > 1).

Within each [[recurrent-class|recurrent class]], wikigraph's
[[sink-page|sink teleportation]] ensures aperiodicity. If the wiki has
multiple [[communicating-classes|communicating classes]], each recurrent class
gets its own π (transient classes get π = 0).

## What high and low π means

| π value   | Interpretation                                                            |
| --------- | ------------------------------------------------------------------------- |
| High      | Hub page — many walks converge here via dense incoming link structure     |
| Low       | Peripheral — few indirect paths lead here; candidate for more inbound links |
| Zero      | Page is in a [[communicating-classes|transient class]] — cut off from recurrent core |

Low π ≠ zero in-degree. A page can have many inbound links and still score
low if those links come only from other low-traffic pages.

## Orphan detection

The [[analyze]] command's **Orphan pages** section lists pages in the bottom
N% by π (default: bottom 10%, controlled by `--orphan-pct`). These pages are
rarely visited by the random walk and typically need more inbound wikilinks
from well-connected pages.

## Where π appears in wikigraph output

| Command     | Use of π                                                         |
| ----------- | ---------------------------------------------------------------- |
| [[analyze]] | Orphan detection, most-central ranking, [[commute-time]] weighting  |
| [[export]]  | `pi` field in JSON/CSV per page                                 |
| [[graph]]   | Node size in the force-directed visualisation encodes π         |

## Sources

- [`cmd_analyze.go` — `kern.Stationary`, orphan section](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_analyze.go)
- [`cmd_export.go` — `pi` field in export](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_export.go)
- [Stationary distribution — Wikipedia](https://en.wikipedia.org/wiki/Stationary_distribution)
- [Power iteration — Wikipedia](https://en.wikipedia.org/wiki/Power_iteration)
