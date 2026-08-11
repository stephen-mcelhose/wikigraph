---
type: concept
title: Kernel Identifiability — Recovering P from Analysis Outputs
description: What each wikigraph output signal tells you about the transition matrix P, and when recovery is lossless vs. lossy.
tags: [markov-chain, transition-matrix, identifiability, catrace, export, kernel-recovery]
timestamp: 2026-08-11T00:00:00Z
---

# Kernel Identifiability — Recovering P from Analysis Outputs

## Overview

The **transition matrix** $P$ is the complete description of a Markov chain. Every analysis
`wikigraph` produces — stationary distribution $\pi$, communicating classes, commute times —
is derived from $P$. The question of **kernel identifiability** asks the reverse: given one of
these outputs, can you reconstruct $P$?

The answer depends on which output you have.

## Signal hierarchy

| Observable                                     | Recovers $P$?       | Notes                                                                              |
| ---------------------------------------------- | ------------------- | ---------------------------------------------------------------------------------- |
| $\pi$ alone                                    | No                  | Many distinct graphs share the same $\pi$                                          |
| Commute times for a subset of pairs            | No                  | Top-5 per node (from `analyze --suggest-top 5`) is far too sparse                  |
| Full $N \times N$ commute time matrix          | No (general graphs) | Recovers $P$ only for trees; non-trees have multiple realizations                  |
| `catrace.Trace(S)` on a subset $S$             | Yes, for $S$        | Lossless collapsed kernel $P_S$ — see below                                        |
| `catrace.Trace(S = \text{all states})`         | Yes                 | Trivially returns $P$ itself                                                       |
| `wikigraph export --min-edge 0` (edges CSV)    | Yes                 | Every $(i,j)$ entry of $P$ is in the file, including structural zeros              |
| `wikigraph export` (default `--min-edge 0.005`)| Yes (for nx-to-wiki graphs) | All non-zero edges exceed 0.005; sparse format encodes $P$ completely      |
| The `.md` files themselves                     | Yes (uniform walks) | $P_{ij} = 1/\deg(i)$ if `[[j]]` appears in `i.md`; files are a lossless encoding |

## Why π alone is insufficient

$\pi$ satisfies $\pi P = \pi$. For a given $\pi$, there are infinitely many row-stochastic
matrices $P$ that satisfy this equation. Knowing that `off-16` has $\pi = 0.109$ tells you
nothing about which other nodes `off-16` links to.

## Why the full commute time matrix is insufficient (in general)

The commute time $K(i,j) = 2|E| \cdot R_{ij}$ where $R_{ij}$ is the effective resistance.
For a **tree**, effective resistance uniquely determines the graph — there is only one spanning
tree. For a **general graph**, multiple distinct adjacency structures can produce identical
effective resistance matrices. Commute times are a spectral fingerprint of the graph, not a
reconstruction of it.

## catrace.Trace — lossless collapse onto a subset

`catrace.Trace(S, tol)` computes the **effective trace kernel** $P_S$: a $|S| \times |S|$
row-stochastic matrix where $P_S(i,j)$ is the probability that a walk starting at $i \in S$
hits $j \in S$ before returning to $i$, integrated over all paths through the excluded states
$S^c$. This is sometimes called the **censored chain** on $S$.

$P_S$ is lossless with respect to hitting probabilities within $S$: any quantity computable
from $P$ restricted to first-passage behaviour among states in $S$ can be computed from $P_S$
instead. When $S$ is all states, $P_S = P$ exactly.

`wikigraph goal` uses this to build meaningful subgraphs: it extracts the effective kernel
over the goal-relevant subset rather than the full $N \times N$ matrix.

## For nx-to-wiki uniform walks: the .md files ARE P

Every graph produced by `nx-to-wiki` runs `G.to_directed()` on a connected undirected graph.
Because every node in the source graph has at least one neighbour, every generated `.md` file
has at least one `[[wikilink]]` — there are **no sink pages**. This matters because `wiki.go`
handles sink pages differently:

```
// wiki.go — buildAdjacency
if len(linked) == 0 {
    // Sink node: teleport uniformly to avoid a zero row.
    for j := 0; j < n; j++ { adj.Set(i, j, 1.0) }
}
```

A sink page's row is set to $P_{ij} = 1/n$ for all $j$ — uniform teleportation, not derivable
from the file's wikilinks alone. **For sink-free wikis only**, the transition matrix is a
uniform random walk:

$$P_{ij} = \frac{1}{\deg(i)} \quad \text{if } \texttt{[[j]]} \text{ appears in } \texttt{i.md} \text{ (self-links excluded)}, \quad \text{else } 0$$

All nx-to-wiki generated wikis are sink-free by construction, so for these graphs the `.md`
files fully determine $P$ without running `catrace`. For a general user wiki — which may contain
sink pages — the `.md` files alone are insufficient; you must also know which pages are sinks
and apply the $1/n$ teleportation rule.

Consequently, the `wikigraph analyze` text output is **lossy** relative to $P$: it reports $\pi$
and top-5 commute times per node, but not the individual $P_{ij}$ values (156 entries for the
karate club, scaling as $O(n \cdot \bar{d})$ for general graphs).

## Exporting P

`wikigraph export` writes $P$ to a file. The edges output IS the transition matrix:

```bash
# Sparse P — non-zero entries only (sufficient for nx-to-wiki graphs, all edges > 0.005)
wikigraph export /tmp/nxwiki-karate --format csv -o /tmp/karate-export

# Dense P — full N×N including structural zeros (explicit lossless record)
wikigraph export /tmp/nxwiki-karate --format csv -o /tmp/karate-export --min-edge 0
```

For the karate club (34 nodes, 156 directed edges):

| Export flag     | Edge rows | What you get                           |
| --------------- | --------- | -------------------------------------- |
| default         | 156       | Sparse $P$: non-zero entries only       |
| `--min-edge 0`  | 1,156     | Dense $P$: full $34 \times 34$ matrix   |

Both are lossless for this graph (all non-zero $P_{ij} \geq 1/17 \approx 0.059 \gg 0.005$).

## Summary

- `wikigraph analyze` is **intentionally lossy** — it surfaces the high-signal summaries ($\pi$,
  top commute-time pairs) and discards the bulk of $P$.
- `wikigraph export` is **lossless** — it exposes the full kernel.
- `catrace.Trace(S)` is **lossless for subset $S$** — the canonical collapse operation when you
  want a blackbox kernel over a specific set of states.
- For uniform walks, the source `.md` files and the export edges file are equivalent lossless
  representations of $P$.

## Related Concepts

- [[markov-model]] — How $P$ is constructed from wikilinks
- [[catrace]] — `catrace.Kernel` and the `Trace` method
- [[commute-time]] — Why $K(i,j)$ is a graph metric but not a $P$-recovery tool
- [[stationary-distribution]] — Why $\pi$ alone does not identify $P$
- [[export]] — `wikigraph export` reference and format details
- [[absorbing-markov-chain]] — Fundamental matrix $N = (I-Q)^{-1}$; related collapse for absorbing chains

## Sources

- Kemeny, J. G., & Snell, J. L. (1960). *Finite Markov Chains*. D. Van Nostrand Company.
- Chandra et al. (1996). The Electrical Resistancy of a Graph along with its Applications to Random Walks.
- [[catrace]] — `Trace(subset, tol)` implementation
