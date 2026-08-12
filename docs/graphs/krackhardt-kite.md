---
type: spike
title: Krackhardt Kite — wikigraph Analysis
description: Empirical wikigraph analysis of the Krackhardt Kite graph (10-node multi-centrality teaching benchmark) as a Markov chain. Path provisional — will move to flat docs/ when #42 resolves.
tags: [benchmark, krackhardt-kite, networkx, markov-chain, spike, draft]
status: draft
timestamp: 2026-08-11T00:00:00Z
---

<!-- Path provisional — will be updated when #42 resolves -->

# Krackhardt Kite — multi-centrality teaching benchmark

> ⚠️ **DRAFT** — AI-assisted write-up, not yet verified by human analysis. Treat all findings as provisional.

> **nx-to-wiki flag**: `--graph krackhardt-kite --role-names`
> **Nodes**: 10 · **Directed edges**: 36 · **Naming**: tier2-structural-role-walk
> **Source**: empirical — fixed

---

## Background

The Krackhardt Kite is a 10-node, 18-edge social network constructed by David Krackhardt
specifically to demonstrate that "most central" is not a single well-defined property — different
centrality measures crown different nodes. Node 3 has the highest **degree** (6 ties), nodes 5
and 6 have the highest **closeness** (they can reach every other node in the fewest average
hops), and node 7 has the highest **betweenness** (it sits on the most shortest paths between
other pairs, acting as the sole bridge between the dense "kite body" and the sparse "tail").
The graph is drawn to resemble a kite: a diamond-shaped cluster of six densely interconnected
nodes (the body), a single bridging node (the string), and a two-node pendant tail trailing off
to an isolated leaf. It is a standard teaching example in social network analysis courses
precisely because it separates degree, betweenness, and closeness centrality into three
different winners on a graph small enough to compute by hand.

**Origin**: Krackhardt, D. (1990). Assessing the Political Landscape: Structure, Cognition,
and Power in Organizations. *Administrative Science Quarterly*, 35(2), 342–369.
https://doi.org/10.2307/2393394

*(Note: the issue's originally listed DOI, 10.1177/0149206390015001001, resolves to a 404 —
it is not a valid DOI for this paper. The correct DOI for Krackhardt (1990), verified above, is
10.2307/2393394.)*

---

## Graph Properties

| Property                | Value / Description                                                     |
| ------------------------ | ------------------------------------------------------------------------ |
| Nodes                    | 10                                                                       |
| Undirected edges         | 18                                                                       |
| Directed edges (×2)      | 36                                                                       |
| Degree distribution      | Right-skewed; min 1 (node 9, the pendant leaf), max 6 (node 3, the degree hub) |
| Community structure      | None — a single connected community with distinct structural zones (body / bridge / tail) |
| Diameter                 | 4                                                                        |
| Average clustering coef  | 0.6444                                                                   |

**Key parameters**: None — this is a fixed empirical dataset with no generative knobs.

**Structural zones** (from node-id numbering, verified against degree/betweenness/closeness):

| Zone         | Node ids   | Description                                                        |
| ------------ | ---------- | -------------------------------------------------------------------- |
| Kite body    | 0, 1, 2, 3, 4, 5, 6 | Densely interconnected core; node 3 (degree 6) sits at its centre |
| Bridge       | 7          | Sole connector between the kite body and the tail; highest betweenness (0.3889) |
| Tail         | 8          | Degree-2 pendant hanging off the bridge                              |
| Leaf         | 9          | Degree-1 leaf at the very end of the tail                           |

---

## Transition Matrix Structure

For every graph produced by `nx-to-wiki`, the random walk is **uniform**: each step chooses a
neighbour with equal probability. The transition matrix $P$ is fully determined by the
adjacency structure of the `.md` files:

$$P_{ij} = \frac{1}{\deg(i)} \quad \text{if } [[j]] \text{ appears in i.md}, \quad \text{else } 0$$

Representative rows (degree-hub, betweenness-broker, mid-degree, leaf):

| Slug          | Node id | Degree | $P_{ij}$ (non-zero) | Structural note                                            |
| ------------- | ------- | ------ | -------------------- | ------------------------------------------------------------ |
| node-03       | 3       | 6      | 1/6 ≈ 0.166667        | Degree hub — highest degree in the graph, but only moderate betweenness |
| connector-00  | 7       | 3      | 1/3 ≈ 0.333333        | Betweenness broker — sole bridge between kite body and tail    |
| hub-00        | 5       | 5      | 1/5 = 0.200000        | Closeness leader — high degree and high betweenness           |
| leaf-00       | 9       | 1      | 1/1 = 1.000000        | Degree-1 pendant; all weight on its single edge (→ node-05 / node 8) |

**The `.md` files are a lossless encoding of $P$.** The minimum non-zero $P_{ij}$ in this graph
is $1/6 \approx 0.167$ (`node-03`'s row), well above the `--min-edge` default filter of 0.005.

### Export

```bash
# Sparse P — non-zero entries only (36 edge rows; every edge exceeds --min-edge 0.005)
wikigraph export /tmp/nxwiki-kite --format csv -o /tmp/kite-export

# Dense P — full 10×10 matrix including structural zeros (100 edge rows)
wikigraph export /tmp/nxwiki-kite --format csv -o /tmp/kite-export --min-edge 0
```

Both commands write two files:

- `kite-export_nodes.csv` — 10 rows: `slug, π, class`
- `kite-export_edges.csv` — 36 rows (sparse) or 90 rows (dense): `source, target, probability`

No `--min-edge 0` warning applies: minimum $P_{ij} = 1/6 \approx 0.167 \gg 0.005$, so the
default sparse export is lossless for this graph.

---

## Slug Naming

**Naming tier**: Tier 2 — structural-role walk (`--role-names` required; the Krackhardt kite
has no built-in node labels, so `nx-to-wiki` falls through to `assign_role_slugs()` in
`tools/nx-to-wiki/main.py`).

**Assignment algorithm** (`assign_role_slugs`, verbatim from source):

1. Compute `degree = dict(G.degree())` and `betweenness = nx.betweenness_centrality(G)` for
   every node.
2. Compute the 75th percentile of the degree distribution (`deg_p75`) and of the betweenness
   distribution (`bet_p75`) across all nodes, using `statistics.quantiles(values, n=100,
   method="inclusive")`.
3. Classify each node `n` with degree `d` and betweenness `b`:
   - `d == 1` → role `leaf`
   - `d >= deg_p75 and b >= bet_p75` → role `hub`
   - `b >= bet_p75` (and not already a leaf/hub) → role `connector`
   - otherwise → role `node`
4. Group nodes by role. Within each role group, sort node ids ascending and assign
   `{role}-{index:02d}` (zero-padded to the width of the group size, minimum width 2).

**Why this ordering**: For this graph, `deg_p75 = 4.75` and `bet_p75 ≈ 0.2292`. Node 3 has the
single highest degree (6) but a betweenness of only 0.1019 — below `bet_p75` — so it does
**not** qualify as `hub`; it falls into the generic `node` group instead, becoming `node-03`
(the 4th of 6 nodes in the `node` group, sorted by id: 0, 1, 2, 3, 4, 8). Nodes 5 and 6 both
have degree 5 (≥ 4.75) and betweenness 0.2315 (≥ 0.2292), so both qualify as `hub`, becoming
`hub-00` (id 5) and `hub-01` (id 6). Node 7 has degree 3 (below `deg_p75`) but the highest
betweenness in the graph (0.3889 ≥ 0.2292), so it becomes the sole `connector`, `connector-00`.
Node 9 has degree 1, so the `leaf` rule fires first regardless of betweenness, giving
`leaf-00`. The remaining five nodes (0, 1, 2, 4, 8) join node 3 in the generic `node` group,
sorted by id (0, 1, 2, 3, 4, 8) and zero-padded to width 2 since the group has 6 members.

Full node-id → slug mapping:

| Slug          | NetworkX node id | Degree | Betweenness | Role rule applied                          |
| ------------- | ----------------- | ------ | ----------- | -------------------------------------------- |
| node-00       | 0                  | 4      | 0.0231      | Below both percentiles → `node`              |
| node-01       | 1                  | 4      | 0.0231      | Below both percentiles → `node`              |
| node-02       | 2                  | 3      | 0.0000      | Below both percentiles → `node`              |
| node-03       | 3                  | 6      | 0.1019      | Degree ≥ p75 but betweenness < p75 → `node`  |
| node-04       | 4                  | 3      | 0.0000      | Below both percentiles → `node`              |
| hub-00        | 5                  | 5      | 0.2315      | Degree ≥ p75 and betweenness ≥ p75 → `hub`   |
| hub-01        | 6                  | 5      | 0.2315      | Degree ≥ p75 and betweenness ≥ p75 → `hub`   |
| connector-00  | 7                  | 3      | 0.3889      | Degree < p75, betweenness ≥ p75 → `connector`|
| node-05       | 8                  | 2      | 0.2222      | Betweenness (0.2222) just below p75 (0.2292) → `node` |
| leaf-00       | 9                  | 1      | 0.0000      | Degree == 1 → `leaf` (checked first)         |

Note: `node-05` (node id 8) narrowly misses the `connector` classification — its betweenness
(0.2222) falls just short of the 75th-percentile threshold (0.2292) — despite sitting directly
on the tail between the bridge (`connector-00`) and the leaf (`leaf-00`). This is a case where
the Tier 2 percentile cutoff produces a counter-intuitive result: visually and structurally,
node 8 is part of the bridge/tail chain, but the role algorithm classifies it as a generic
`node`.

---

## How to Generate

```bash
# Minimal — role-names flag is required (no built-in labels for this graph)
python3 tools/nx-to-wiki/main.py --graph krackhardt-kite --role-names --out /tmp/nxwiki-kite
wikigraph graph /tmp/nxwiki-kite -o /tmp/nx-kite.html --title "Krackhardt Kite"
wikigraph analyze /tmp/nxwiki-kite --suggest-top 5
```

No parameter variants apply — `--graph krackhardt-kite` always produces the same 10-node,
18-edge graph; there are no generative knobs to sweep. Omitting `--role-names` falls through
to Tier 3 fallback naming (`page-00` … `page-09`), which loses all structural information —
always pass `--role-names` for this graph.

**What the output directory contains:**

- 10 `.md` files, one per node
- Each file links to all its undirected neighbours (symmetric — every link is bidirectional)
- File layout: `# {slug}\n\n[[neighbour-1]] [[neighbour-2]] ...`
- Why directed edges = 2× undirected: `G.to_directed()` adds both u→v and v→u for every edge

---

## wikigraph Analysis

> **⚠️ Communicating classes are always trivial here.** All wikis produced by `nx-to-wiki` use
> `G.to_directed()` on a connected undirected graph, which yields a **strongly connected**
> directed graph. `wikigraph analyze` will report **one communicating class containing all
> nodes**. This is expected and correct — not a bug.
>
> This graph has no community attribute; the meaningful structural distinction is **zone**
> (kite body / bridge / tail / isolate), not community. The cross-community section of the
> template is replaced below with a **Centrality Comparison** section, since `wikigraph` has
> no betweenness or closeness metric of its own — those are computed directly against the
> underlying NetworkX graph and compared to π.

### Raw analyze output

```
=== Overview ===
Pages:        10
Edges:        36
Entropy rate: 1.9720 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 10 page(s)
  connector-00
  hub-00
  hub-01
  leaf-00
  node-00
  node-01
  node-02
  node-03
  node-04
  node-05

=== Orphan pages (bottom 10% by stationary distribution) ===
  leaf-00                                   π=0.027778  → add inbound links
  node-05                                   π=0.055556  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. node-03                                   π=0.166667
  2. hub-00                                    π=0.138889
  3. hub-01                                    π=0.138889
  4. node-01                                   π=0.111111
  5. node-00                                   π=0.111111

=== Suggested missing links (lowest commute time, not yet linked) ===
  connector-00:
    → node-03                                 (commute: 27.00)
    → node-01                                 (commute: 31.14)
    → node-00                                 (commute: 31.14)
    → node-04                                 (commute: 34.05)
    → node-02                                 (commute: 34.05)
  hub-00:
    → node-01                                 (commute: 18.49)
    → node-04                                 (commute: 22.38)
    → node-05                                 (commute: 57.41)
    → leaf-00                                 (commute: 93.41)
  hub-01:
    → node-00                                 (commute: 18.49)
    → node-02                                 (commute: 22.38)
    → node-05                                 (commute: 57.41)
    → leaf-00                                 (commute: 93.41)
  leaf-00:
    → connector-00                            (commute: 72.00)
    → hub-01                                  (commute: 93.41)
    → hub-00                                  (commute: 93.41)
    → node-03                                 (commute: 99.00)
    → node-00                                 (commute: 103.14)
  node-00:
    → hub-01                                  (commute: 18.49)
    → node-04                                 (commute: 23.59)
    → connector-00                            (commute: 31.14)
    → node-05                                 (commute: 67.14)
    → leaf-00                                 (commute: 103.14)
  node-01:
    → hub-00                                  (commute: 18.49)
    → node-02                                 (commute: 23.59)
    → connector-00                            (commute: 31.14)
    → node-05                                 (commute: 67.14)
    → leaf-00                                 (commute: 103.14)
  node-02:
    → hub-01                                  (commute: 22.38)
    → node-01                                 (commute: 23.59)
    → node-04                                 (commute: 28.22)
    → connector-00                            (commute: 34.05)
    → node-05                                 (commute: 70.05)
  node-03:
    → connector-00                            (commute: 27.00)
    → node-05                                 (commute: 63.00)
    → leaf-00                                 (commute: 99.00)
  node-04:
    → hub-00                                  (commute: 22.38)
    → node-00                                 (commute: 23.59)
    → node-02                                 (commute: 28.22)
    → connector-00                            (commute: 34.05)
    → node-05                                 (commute: 70.05)
  node-05:
    → hub-00                                  (commute: 57.41)
    → hub-01                                  (commute: 57.41)
    → node-03                                 (commute: 63.00)
    → node-00                                 (commute: 67.14)
    → node-01                                 (commute: 67.14)
```

### Stationary distribution (π)

| Rank | Slug     | π value  | Expected structural role                                       |
| ---- | -------- | -------- | ------------------------------------------------------------------ |
| 1    | node-03  | 0.166667 | Degree hub (node 3, degree 6, the highest in the graph)             |
| 2    | hub-00   | 0.138889 | Closeness/degree co-leader (node 5, degree 5)                       |
| 3    | hub-01   | 0.138889 | Closeness/degree co-leader (node 6, degree 5)                       |
| 4    | node-01  | 0.111111 | Kite-body member (node 1, degree 4)                                 |
| 5    | node-00  | 0.111111 | Kite-body member (node 0, degree 4)                                 |

**Finding**: `node-03` (node 3, the highest-degree node) leads π at 0.166667. π is a
degree-driven measure — it does **not** identify the betweenness broker (node 7,
`connector-00`) as central at all: `connector-00` does not even appear in the top 5, despite
having the highest betweenness centrality (0.3889) in the graph. This is the graph's core
pedagogical point made numerically: degree centrality and betweenness centrality identify
completely different nodes.

### Suggested links

Top-2 π nodes (`node-03`, tied `hub-00`/`hub-01`) and the lowest-π node (`leaf-00`,
π=0.027778):

| From      | To          | Commute time | Passes through the broker (`connector-00`)? |
| --------- | ----------- | ------------ | --------------------------------------------- |
| node-03   | connector-00 | 27.00        | Is the broker itself                           |
| node-03   | node-05      | 63.00        | Yes — only path from node-03 to node-05 crosses connector-00 |
| node-03   | leaf-00      | 99.00        | Yes — only path crosses connector-00 and node-05 |
| hub-00    | node-01      | 18.49        | No — within kite body                          |
| hub-00    | node-04      | 22.38        | No — within kite body                          |
| hub-00    | node-05      | 57.41        | Yes — crosses connector-00                     |
| hub-00    | leaf-00      | 93.41        | Yes — crosses connector-00 and node-05         |
| leaf-00   | connector-00 | 72.00        | Is the broker itself                           |
| leaf-00   | hub-01       | 93.41        | Yes — crosses connector-00                     |
| leaf-00   | hub-00       | 93.41        | Yes — crosses connector-00                     |
| leaf-00   | node-03      | 99.00        | Yes — crosses connector-00                     |
| leaf-00   | node-00      | 103.14       | Yes — crosses connector-00                     |

**Finding**: Every suggestion for `node-03` and every suggestion for `leaf-00` that targets a
kite-body node necessarily routes through `connector-00` (node 7) — the graph has only one
path between the kite body/tail-adjacent nodes and the tail, so `connector-00` mediates every
long-range commute. For `hub-00`, 2 of 4 suggestions stay within the kite body (short commute
times 18.49–22.38) while the other 2 reach into the tail and must cross `connector-00`. The
remaining 6 nodes not shown above (`connector-00`, `hub-01`, `node-00`, `node-02`, `node-04`)
follow the same pattern: any suggestion targeting `node-05` or `leaf-00` crosses
`connector-00`, and any suggestion targeting a kite-body node does not.

### Centrality Comparison (replaces cross-community section — no community attribute exists)

| Node id | Slug          | Degree | Betweenness | Closeness | π        |
| ------- | ------------- | ------ | ----------- | --------- | -------- |
| 3       | node-03       | 6 (max)| 0.1019      | 0.6000    | 0.166667 (max) |
| 5       | hub-00        | 5      | 0.2315      | 0.6429 (max) | 0.138889 |
| 6       | hub-01        | 5      | 0.2315      | 0.6429 (max) | 0.138889 |
| 7       | connector-00  | 3      | 0.3889 (max)| 0.6000    | 0.083333 |
| 9       | leaf-00       | 1 (min)| 0.0000      | 0.3103 (min) | 0.027778 (min) |

**Finding**: Three different nodes win three different centrality measures — exactly the
graph's designed teaching point. Degree picks node 3 (`node-03`); closeness picks nodes 5 and 6
(`hub-00`/`hub-01`, tied); betweenness picks node 7 (`connector-00`). π (the random-walk
stationary distribution) agrees with **degree**, not with betweenness or closeness, because π
is mathematically proportional to degree for a uniform random walk (π_i = deg_i / (2|E|)).
`connector-00`, despite having the highest betweenness in the graph, ranks only 6th in π
(0.083333) — its low degree (3) suppresses its stationary visitation probability even though
it is structurally indispensable for connecting the tail to the body.

### Stationary distribution spread

- **Max π**: 0.166667 (node-03)
- **Min π**: 0.027778 (leaf-00)
- **Ratio max/min**: 6.0 — **Hub-dominated** (falls in the 5–20 tier)

**Finding**: The ratio of 6.0 equals `node-03`'s degree exactly (degree 6), consistent with
π_i ∝ deg_i for a uniform random walk on an undirected graph. `leaf-00` (node 9, the sole
degree-1 pendant) sits at the theoretical minimum, π = 1/(2·18) = 0.027778.

### Visual observations

> *The following description is inferred from the graph's structural properties and stationary
> distribution. The force-directed layout can be verified by opening `/tmp/nx-kite.html`.*

Given the kite's designed shape, the force-directed layout should show a dense hexagonal
cluster on one side comprising `node-00` through `node-04` plus `hub-00`/`hub-01`, with
`node-03` visually largest (highest π=0.167, sitting near the cluster's geometric centre with
6 edges radiating outward). `connector-00` should appear as a single narrow bridge node
positioned between this dense cluster and a short two-node tail (`node-05` then `leaf-00`),
visibly smaller than the kite-body nodes despite its structural importance — its low π
(0.083333) means the layout's node-sizing (if driven by π) will under-represent its true
significance as the sole cut vertex separating the tail from the body. `leaf-00` should appear
as the smallest node in the entire graph, a single pendant at the extreme end of the tail with
one thin edge.

### Markov questions — answered

- **Does π rank the known degree-hub (node 3) highest?**
  Yes — `node-03` (node 3, degree 6) has π=0.166667, the highest in the graph.

- **Does `suggest` recommend links that pass through the known broker node (node 7)?**
  Yes — every suggestion in `node-03`'s and `leaf-00`'s top-5 lists that targets a node on the
  opposite side of the graph (tail vs. body) necessarily has its shortest random-walk path
  routed through `connector-00` (node 7), since it is the sole cut vertex connecting the two
  regions. `node-03`'s 2nd and 3rd suggestions (`node-05` at commute 63.00, `leaf-00` at commute
  99.00) both cross `connector-00`.

- **Is π hub-dominated or flat (max/min ratio)?**
  Hub-dominated — max/min ratio = 6.0 (0.166667 / 0.027778), which falls in the 5–20 tier.

- **Does the degree-hub (node 3) lead π while the broker (node 7) shows up in `suggest`
  recommendations rather than in π rank?**
  Yes, precisely. `node-03` leads π at 0.166667 (rank 1 of 10), while `connector-00` ranks only
  6th in π (0.083333) and does not appear in the top-5 π list at all. But `connector-00`
  appears as the first suggested link for `node-03` itself (commute time 27.00, the lowest
  commute time in `node-03`'s entire suggestion list) and mediates 2 of `node-03`'s other 2
  suggestions (`node-05`, `leaf-00`) by lying on their only path — confirming the broker's
  importance surfaces through the `suggest` mechanism (path-mediation) rather than through π
  (visitation frequency).

---

## Sources

- Krackhardt, D. (1990). Assessing the Political Landscape: Structure, Cognition, and Power in
  Organizations. *Administrative Science Quarterly*, 35(2), 342–369.
  https://doi.org/10.2307/2393394
- NetworkX: https://networkx.org/documentation/stable/reference/generated/networkx.generators.small.krackhardt_kite_graph.html
- Graph topology notes: Star Graph / Markov walk section from graph-topologies research notes

---

## Definition of Done

- [x] All sections above filled — no placeholder text remaining
- [x] Slug ordering algorithm documented with a node-id mapping table (not just the pattern)
- [x] `wikigraph analyze --suggest-top 5` output pasted verbatim
- [x] Cross-community edge count computed and recorded (replaced with Centrality Comparison
      section — no community attribute exists in this graph)
- [x] All Markov questions answered with actual numbers from analyze output
- [x] At least one variant command shown, or a note explaining why none apply (fixed graph —
      no generative knobs; `--role-names` is required, not optional)
- [x] Every reference link manually verified to resolve (correct DOI 302→200 replaces the
      issue's broken DOI; NetworkX docs 200 at the `generators.small` path)
- [x] File committed to branch — path to be updated once #42 resolves
