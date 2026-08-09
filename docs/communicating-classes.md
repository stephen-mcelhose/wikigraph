---
type: concept
title: Communicating Classes
description: A communicating class is a maximal set of Markov chain states that are mutually reachable — the structural unit that determines whether a wiki is connected or fragmented.
resource: cmd_analyze.go
tags: [communicating-classes, scc, tarjan, markov, graph-structure, reachability]
timestamp: 2026-08-09T07:31:56Z
---

# Communicating Classes

A **communicating class** is a maximal set of states in a Markov chain where
every state can reach every other state — possibly via intermediate states.
For a wiki, states are pages and transitions are `[[wikilinks]]`. Two pages
**communicate** if there is a directed path of links from A to B *and* from
B to A.

The class decomposition is the single most important structural fact about a
wiki: it tells you whether readers can navigate freely or whether the wiki is
fragmented into islands. See [[recurrent-class]] for what happens inside and
outside each class.

## How it is computed

wikigraph uses Kosaraju's strongly connected components (SCC) algorithm, via
`catrace.Kernel.Classes`. Kosaraju's algorithm runs in O(V + E) — linear in
pages and links — and produces a partition of all pages into SCCs. Each SCC
is a communicating class.

After SCC decomposition, each class is labelled **recurrent** or **transient**
based on whether any state in the class can escape to another class (transient)
or not (recurrent). See [[recurrent-class]].

## Reading the output

```
=== Communicating classes ===
Class 1 (recurrent): 11 page(s)
  adr-001-embedding-layer
  adr-002-slug-resolution
  analyze
  ...
```

| Field          | Meaning                                                              |
| -------------- | -------------------------------------------------------------------- |
| `Class N`      | Index — if there are multiple classes you see Class 1, Class 2, etc |
| `(recurrent)`  | All pages mutually reachable; stationary distribution defined here   |
| `(transient)`  | Pages escape to another class eventually; π = 0 for these pages     |
| `N page(s)`    | Size of the class                                                    |
| Page list      | Every slug in this class, in SCC discovery order (not alphabetical)  |

**One recurrent class containing all pages** is the healthy state: a random
walk started anywhere reaches everywhere. Multiple classes — or transient
classes — mean the wiki has structural gaps.

## Common failure modes

| Symptom | Cause | Fix |
| ------- | ----- | --- |
| `Class 2 (transient): 1 page` | A page has inbound links but no outbound link back to the main cluster | Add at least one `[[slug]]` to a page in the recurrent class |
| Two recurrent classes | Two clusters with no cross-links | Add links bridging the clusters in both directions |
| All pages transient | No cycles exist anywhere | Add cross-links; a linear chain has no recurrent class |
| Page missing from all classes | It has no `.md` extension, or it's in `--exclude` | Check filename and exclude list |

## Relationship to π and entropy

Only recurrent-class pages receive non-zero [[stationary-distribution|stationary probability]]
π — the long-run visit frequency of the [[markov-model]] [[random-walk|random walk]].
Transient pages are visited finitely often and then never again, so their π is
0 regardless of how many internal links they have.

[[entropy-rate|Entropy rate]] is computed over the recurrent class(es). A single large
recurrent class with varied out-degrees produces a healthy mid-range entropy
(wikigraph docs: 1.65 bits on 11 pages). See [[mfpt]] for how mean first
passage time uses the same Markov structure to measure distance between pages.

## Sources

- `cmd_analyze.go` — `kern.Classes(1e-10)`, recurrent/transient labelling, output formatting
- `catrace` library — Tarjan SCC via `Classes`
