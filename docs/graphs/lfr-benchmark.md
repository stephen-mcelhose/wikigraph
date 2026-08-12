<!-- Path provisional — will be updated when #42 resolves -->

# LFR Benchmark Graph — power-law community-detection gold standard

> **nx-to-wiki flag**: `--graph lfr --n [nodes] --tau1 [degree-exponent] --tau2 [community-exponent] --mu [mixing]`
> **Nodes**: n (parameterised) · **Directed edges**: variable (parameterised; see self-loop note)
> **Naming**: tier1-lfr-community-attr
> **Source**: generative — parameterised

---

## Background

The Lancichinetti–Fortunato–Radicchi (LFR) benchmark is the standard random-graph model used
to evaluate community-detection algorithms, because — unlike the planted partition / SBM model,
which produces equal-sized communities and a uniform degree distribution — LFR reproduces two
properties observed in real-world networks: a **power-law degree distribution** (exponent
`τ1`) and a **power-law community-size distribution** (exponent `τ2`). Every node's edges are
split between intra-community ties and inter-community ties according to the **mixing
parameter μ**: each node keeps a fraction `1-μ` of its edges inside its own community and
sends the remaining fraction `μ` to nodes in other communities. `μ` is the LFR analogue of the
planted-partition `p_out/p_in` ratio, but applied per-node rather than per-edge-probability,
and it is the primary experimental variable used to stress-test detection algorithms:
**μ=0.5 is the conventional detection boundary** — below it, a majority of each node's edges
stay within its own community, so community structure is (in principle) still recoverable;
above it, most of each node's edges leave its community, and the community signal degrades
toward what a detection algorithm would see in a random graph with no real community structure.

**Origin**: Lancichinetti, A., Fortunato, S., & Radicchi, F. (2008). Benchmark graphs for
testing community detection algorithms. *Physical Review E*, 78, 046110.
https://doi.org/10.1103/PhysRevE.78.046110

---

## Graph Properties

**Key parameters**:

| Flag              | Default | Range      | What it controls                                             |
| ----------------- | ------- | ---------- | ---------------------------------------------------------------- |
| `--n`             | 50      | ≥ 1        | Total node count                                                  |
| `--tau1`          | 2.5     | > 1        | Power-law exponent for the degree distribution                      |
| `--tau2`          | 1.5     | > 1        | Power-law exponent for the community-size distribution               |
| `--mu`            | 0.3     | [0.0, 1.0] | Fraction of each node's edges that cross community boundaries         |
| `--average-degree`| none (defaults to 5.0 if neither this nor `--min-degree` set) | > 0 | Target mean degree across the graph |
| `--min-degree`    | none    | ≥ 1        | Minimum node degree (alternative to `--average-degree`)              |
| `--max-degree`    | none    | ≥ min       | Maximum node degree                                                 |
| `--min-community` | none (defaults to `max(10, n // 10)`) | ≥ 1 | Minimum community size |
| `--max-community` | none    | ≥ min       | Maximum community size                                              |
| `--seed`          | none    | any int     | RNG seed for reproducible generation                                |

Three parameter sets (as directed by the issue), all with `n=50`, `tau1=2.5`, `tau2=1.5`,
`seed=42`. **No generation failures occurred** for any of the three μ values tested — all three
commands ran successfully on the first attempt with `nx-to-wiki`'s implicit
`--average-degree 5.0` default (applied automatically since neither `--average-degree` nor
`--min-degree` was passed) and `--min-community 10` default (`max(10, 50 // 10) = 10`), so the
`--average-degree 5` fallback mentioned in the issue's guidance was not needed:

| Parameter set | μ   | Nodes | Communities | Undirected edges | Self-loops | Directed edges (self-loops excluded) |
| ------------- | --- | ----- | ------------- | ------------------ | ------------ | ---------------------------------------- |
| μ=0.1         | 0.1 | 50    | 4              | 105                 | 13            | 184                                        |
| μ=0.3         | 0.3 | 50    | 4              | 109                 | 8             | 202                                        |
| μ=0.5         | 0.5 | 50    | 4              | 114                 | 1             | 226                                        |

**Important deviation from the naive edge-count expectation**: `LFR_benchmark_graph` can
produce self-loops (NetworkX's own docstring notes this explicitly). `nx-to-wiki`'s
`write_wiki` function excludes self-loops when building each page's outbound link list (via
`out_links = sorted({slugs[v] for v in D.successors(n) if v != n})`), so the directed edge
count reported by `wikigraph analyze` (184 / 202 / 226) is **lower** than `2 ×
undirected_edges` would suggest (210 / 218 / 228) — the difference in each case equals `2 ×
self-loop count`. This is a structural quirk specific to LFR among the seven graphs in this
series; every other graph in this documentation set has zero self-loops.

With `min_community=10` and `n=50`, the realised community-size distribution is `[16, 14, 10,
10]` — 4 communities, identical across all three μ values (community membership is fixed once
generated with a given seed; only which edges are intra- vs. inter-community changes with μ).

---

## Transition Matrix Structure

For every graph produced by `nx-to-wiki`, the random walk is **uniform over the retained
(non-self-loop) edges**: each step chooses a neighbour with equal probability. The transition
matrix $P$ is fully determined by the adjacency structure of the `.md` files:

$$P_{ij} = \frac{1}{\deg(i)} \quad \text{if } [[j]] \text{ appears in i.md}, \quad \text{else } 0$$

where $\deg(i)$ here is the *retained* (self-loop-excluded) out-degree recorded in the `.md`
file, not NetworkX's raw `G.degree()` (which would double-count a self-loop).

### μ=0.1

| Slug   | Retained degree | $P_{ij}$ (non-zero) | Structural note                                       |
| ------ | ----------------- | -------------------- | -------------------------------------------------------- |
| c0-03  | 14                 | 1/14 ≈ 0.071429       | Highest-degree/highest-π node; NetworkX node id 10, raw degree 16 (includes 1 self-loop) |
| c0-08  | 9                  | 1/9 ≈ 0.111111        | Second-highest π                                          |
| c3-04  | 1                  | 1/1 = 1.000000        | Lowest-π node — a near-isolated leaf                       |

### μ=0.3

| Slug   | Retained degree | $P_{ij}$ (non-zero) | Structural note                                       |
| ------ | ----------------- | -------------------- | -------------------------------------------------------- |
| c1-03  | 16                 | 1/16 = 0.062500       | Highest-degree/highest-π node                              |
| c2-00  | 8                  | 1/8 = 0.125000        | Mid-tier node                                              |
| c0-02  | 2                  | 1/2 = 0.500000        | Lowest-degree node                                         |

### μ=0.5

| Slug   | Retained degree | $P_{ij}$ (non-zero) | Structural note                                       |
| ------ | ----------------- | -------------------- | -------------------------------------------------------- |
| c2-02  | 18                 | 1/18 ≈ 0.055556       | Highest-degree/highest-π node — grew further as μ rose      |
| c0-07  | 10                 | 1/10 = 0.100000       | Mid-tier node                                              |
| c0-13  | 2                  | 1/2 = 0.500000        | Lowest-degree tier                                          |

**The `.md` files are a lossless encoding of $P$ for all three parameter sets** (over the
retained, self-loop-free edge set). Minimum non-zero $P_{ij}$ across all three sets is
$1/18 \approx 0.0556$ (μ=0.5's highest-degree node) — well above the `--min-edge` default
filter of 0.005.

### Export

```bash
wikigraph export /tmp/nxwiki-lfr-01 --format csv -o /tmp/lfr-01-export
wikigraph export /tmp/nxwiki-lfr-03 --format csv -o /tmp/lfr-03-export
wikigraph export /tmp/nxwiki-lfr-05 --format csv -o /tmp/lfr-05-export
```

Sparse edge-row counts: 184 (μ=0.1), 202 (μ=0.3), 226 (μ=0.5) — matching the self-loop-excluded
directed-edge counts reported by `wikigraph analyze`, not `2 × G.number_of_edges()`. No
`--min-edge 0` warning applies to any of the three sets tested.

---

## Slug Naming

**Naming tier**: Tier 1

**Assignment algorithm** (`_slugs_from_community_attr(G, attr="community")`, verbatim from
source — the same helper used for planted-partition's `block` attribute, but here the
attribute is **set-valued**, not scalar): `LFR_benchmark_graph` assigns every node a
`community` attribute that is a **Python `set` of node ids** — every member of a node's own
community, including itself (not an integer community index). `nx-to-wiki`:

1. Reads `raw[n] = G.nodes[n]["community"]` — a `set` of co-member node ids for every node.
2. Canonicalises: `canon_keys[n] = frozenset(raw[n])` — converts each node's community set to
   an immutable `frozenset` so it can be used as a dictionary key and compared for equality.
   **Two nodes belong to the same community if and only if their `community` sets, converted
   to `frozenset`, are exactly equal** — this is the crucial distinguishing detail versus
   scalar-attribute graphs (karate club's `club`, planted-partition's `block`): there is no
   pre-existing integer community id anywhere in the LFR node attributes; community identity
   is entirely reconstructed from set equality.
3. Computes `unique = sorted(set(canon_keys.values()), key=lambda k: sorted(k))` — deduplicates
   the frozensets and sorts them by their own sorted member-list (so the community containing
   the lowest node id sorts first).
4. Assigns `community_id[key] = i` for `i, key` in `enumerate(unique)` — giving each unique
   frozenset a 0-indexed integer id in that sorted order.
5. Groups nodes by their `community_id`; within each group, sorts node ids ascending and
   assigns `c{cid}-{index:02d}` (zero-padded to `max(2, len(str(group_size - 1)))` digits).

**Why this ordering**: Because `sorted(unique, key=lambda k: sorted(k))` sorts communities by
their own lowest member's node id, `c0` is always the community that contains node 0 (assuming
node 0 exists, which it always does for `n ≥ 1`), `c1` is the community containing the
lowest-numbered node not in `c0`, and so on. This means — unlike the barbell or
planted-partition graphs, where slug assignment is a pure arithmetic function of node id and
parameters — **the LFR node-id→slug mapping cannot be predicted without first knowing which
nodes the generator grouped together**, since LFR's community assignment is itself an output
of the (seeded) random generation process, not an input parameter.

Full node-id → slug mapping for μ=0.1 (seed=42; identical for μ=0.3 and μ=0.5 since community
membership, unlike edge placement, is stable across the μ sweep for a fixed seed):

| Slug   | NetworkX node id | Community (frozenset, truncated) | Community size |
| ------ | ------------------ | ----------------------------------- | ----------------- |
| c0-00  | 0                   | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-01  | 3                   | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-02  | 9                   | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-03  | 10                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-04  | 12                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-05  | 26                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-06  | 29                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-07  | 31                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-08  | 35                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-09  | 36                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-10  | 38                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-11  | 40                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-12  | 44                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-13  | 45                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-14  | 46                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c0-15  | 49                  | {0, 3, 9, 10, 12, ...}               | 16                 |
| c1-00  | 1                   | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-01  | 7                   | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-02  | 16                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-03  | 17                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-04  | 18                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-05  | 21                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-06  | 23                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-07  | 25                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-08  | 27                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-09  | 28                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-10  | 32                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-11  | 33                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-12  | 42                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c1-13  | 43                  | {1, 7, 16, 21, 23, ...}              | 14                 |
| c2-00  | 2                   | {2, 4, 6, 13, 20, ...}               | 10                 |
| c2-01  | 4                   | {2, 4, 6, 13, 20, ...}               | 10                 |
| c2-02  | 6                   | {2, 4, 6, 13, 20, ...}               | 10                 |
| c2-03  | 13                  | {2, 4, 6, 13, 20, ...}               | 10                 |
| c2-04  | 20                  | {2, 4, 6, 13, 20, ...}               | 10                 |
| c2-05  | 22                  | {2, 4, 6, 13, 20, ...}               | 10                 |
| c2-06  | 34                  | {2, 4, 6, 13, 20, ...}               | 10                 |
| c2-07  | 39                  | {2, 4, 6, 13, 20, ...}               | 10                 |
| c2-08  | 41                  | {2, 4, 6, 13, 20, ...}               | 10                 |
| c2-09  | 48                  | {2, 4, 6, 13, 20, ...}               | 10                 |
| c3-00  | 5                   | {5, 8, 11, 14, 15, ...}              | 10                 |
| c3-01  | 8                   | {5, 8, 11, 14, 15, ...}              | 10                 |
| c3-02  | 11                  | {5, 8, 11, 14, 15, ...}              | 10                 |
| c3-03  | 14                  | {5, 8, 11, 14, 15, ...}              | 10                 |
| c3-04  | 15                  | {5, 8, 11, 14, 15, ...}              | 10                 |
| c3-05  | 19                  | {5, 8, 11, 14, 15, ...}              | 10                 |
| c3-06  | 24                  | {5, 8, 11, 14, 15, ...}              | 10                 |
| c3-07  | 30                  | {5, 8, 11, 14, 15, ...}              | 10                 |
| c3-08  | 37                  | {5, 8, 11, 14, 15, ...}              | 10                 |
| c3-09  | 47                  | {5, 8, 11, 14, 15, ...}              | 10                 |

`c0` contains node 0 (16 members), `c1` contains node 1 (14 members), `c2` contains node 2 (10
members), `c3` contains node 5 (10 members) — confirming the sort-by-lowest-member rule above.

---

## How to Generate

```bash
python3 tools/nx-to-wiki/main.py --graph lfr --n 50 --tau1 2.5 --tau2 1.5 --mu 0.1 --seed 42 --out /tmp/nxwiki-lfr-01
python3 tools/nx-to-wiki/main.py --graph lfr --n 50 --tau1 2.5 --tau2 1.5 --mu 0.3 --seed 42 --out /tmp/nxwiki-lfr-03
python3 tools/nx-to-wiki/main.py --graph lfr --n 50 --tau1 2.5 --tau2 1.5 --mu 0.5 --seed 42 --out /tmp/nxwiki-lfr-05

wikigraph graph /tmp/nxwiki-lfr-01 -o /tmp/nx-lfr-01.html --title "LFR mu=0.1"
wikigraph graph /tmp/nxwiki-lfr-03 -o /tmp/nx-lfr-03.html --title "LFR mu=0.3"
wikigraph graph /tmp/nxwiki-lfr-05 -o /tmp/nx-lfr-05.html --title "LFR mu=0.5"

wikigraph analyze /tmp/nxwiki-lfr-01 --suggest-top 5
wikigraph analyze /tmp/nxwiki-lfr-03 --suggest-top 5
wikigraph analyze /tmp/nxwiki-lfr-05 --suggest-top 5
```

**No deviation from the issue's commands was required** — all three ran successfully without
needing `--average-degree 5` as an explicit fallback (it is already `nx-to-wiki`'s implicit
default when neither `--average-degree` nor `--min-degree` is passed).

**Constraint to be aware of**: `min_community` must exceed `min_degree` (or the implicit
average-degree-derived minimum) or `LFR_benchmark_graph` raises an exception internally. With
`n=50`, the default `min_community = max(10, 50 // 10) = 10` comfortably exceeds the low
degrees produced by `average_degree=5.0`, so no constraint violation occurred in this sweep.

**What each output directory contains:**

- 50 `.md` files, one per node (fixed across all three parameter sets since `n` is fixed)
- Each file links to all its **non-self-loop** undirected neighbours (symmetric — every link
  is bidirectional); self-loops present in the underlying NetworkX graph are silently dropped
  by `write_wiki`'s `if v != n` filter
- File layout: `# {slug}\n\n[[neighbour-1]] [[neighbour-2]] ...`
- Directed-edge count is `2 × (undirected edges − self-loops)`, not `2 × undirected edges`

---

## wikigraph Analysis

> **⚠️ Communicating classes are always trivial here — do not treat class colouring as a
> ground-truth-tracking signal.** All wikis produced by `nx-to-wiki` use `G.to_directed()` on a
> connected undirected graph, which yields a **strongly connected** directed graph;
> `wikigraph analyze` reports **1 communicating class of 50 pages** for every μ value tested,
> by construction. This is expected and is not reported further per the issue's explicit
> guidance not to frame class colouring as tracking the planted communities. The meaningful
> signal is π and the fraction of `suggest` recommendations that are inter-community.

### Comparison across parameter sets

| Parameter set | μ   | Nodes | Communities | Directed edges | Entropy rate | Max π    | Min π    | Ratio | Cross-community % |
| ------------- | --- | ----- | ------------- | ---------------- | -------------- | -------- | -------- | ----- | -------------------- |
| μ=0.1         | 0.1 | 50    | 4              | 184               | 2.0903 bits    | 0.076087 | 0.005435 | 14.0  | 22 of 184 (12.0%)     |
| μ=0.3         | 0.3 | 50    | 4              | 202               | 2.2020 bits    | 0.079208 | 0.009901 | 8.0   | 86 of 202 (42.6%)     |
| μ=0.5         | 0.5 | 50    | 4              | 226               | 2.4113 bits    | 0.079646 | 0.008850 | 9.0   | 144 of 226 (63.7%)    |

**Trend**: Cross-community edge percentage rises sharply and monotonically with μ (12.0% →
42.6% → 63.7%), directly tracking the mixing parameter's definition. Entropy rate also rises
monotonically (2.09 → 2.20 → 2.41 bits) as the graph gains more total edges and becomes more
uniformly connected. The max/min π ratio, notably, does **not** trend monotonically with μ — it
peaks at μ=0.1 (14.0) and *falls* at higher μ (8.0, then 9.0). This is because π's ratio here
is dominated by each set's single highest-degree node (whose degree happens to be largest, by
chance of the power-law degree distribution, at μ=0.1) rather than by the overall mixing level
— the power-law degree distribution's inherent variability, not μ, is the primary driver of π
skew in this graph family, unlike the caveman/planted-partition graphs where π skew stayed
roughly constant across their sweeps.

### μ=0.1

#### Raw analyze output

```
=== Overview ===
Pages:        50
Edges:        184
Entropy rate: 2.0903 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 50 page(s)
  c0-00
  c0-01
  c0-02
  c0-03
  c0-04
  c0-05
  c0-06
  c0-07
  c0-08
  c0-09
  c0-10
  c0-11
  c0-12
  c0-13
  c0-14
  c0-15
  c1-00
  c1-01
  c1-02
  c1-03
  c1-04
  c1-05
  c1-06
  c1-07
  c1-08
  c1-09
  c1-10
  c1-11
  c1-12
  c1-13
  c2-00
  c2-01
  c2-02
  c2-03
  c2-04
  c2-05
  c2-06
  c2-07
  c2-08
  c2-09
  c3-00
  c3-01
  c3-02
  c3-03
  c3-04
  c3-05
  c3-06
  c3-07
  c3-08
  c3-09

=== Orphan pages (bottom 10% by stationary distribution) ===
  c3-04                                     π=0.005435  → add inbound links
  c3-05                                     π=0.005435  → add inbound links
  c2-01                                     π=0.010870  → add inbound links
  c2-06                                     π=0.010870  → add inbound links
  c2-04                                     π=0.010870  → add inbound links
  c2-03                                     π=0.010870  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. c0-03                                     π=0.076087
  2. c0-08                                     π=0.048913
  3. c3-06                                     π=0.043478
  4. c3-01                                     π=0.032609
  5. c2-07                                     π=0.032609

=== Suggested missing links (lowest commute time, not yet linked) ===
  c0-00:
    → c0-08                                   (commute: 103.44)
    → c0-09                                   (commute: 112.54)
    → c0-07                                   (commute: 122.82)
    → c0-12                                   (commute: 128.21)
    → c0-05                                   (commute: 138.69)
  c0-01:
    → c0-08                                   (commute: 110.73)
    → c0-02                                   (commute: 127.58)
    → c0-12                                   (commute: 127.58)
    → c0-07                                   (commute: 130.59)
    → c0-11                                   (commute: 139.50)
  c0-02:
    → c0-08                                   (commute: 84.40)
    → c0-07                                   (commute: 101.20)
    → c0-12                                   (commute: 108.99)
    → c3-01                                   (commute: 111.37)
    → c3-06                                   (commute: 118.20)
  c0-03:
    → c3-01                                   (commute: 87.80)
    → c0-05                                   (commute: 92.85)
    → c3-07                                   (commute: 97.14)
    → c3-06                                   (commute: 97.75)
    → c0-15                                   (commute: 100.90)
  c0-04:
    → c0-08                                   (commute: 98.43)
    → c0-07                                   (commute: 112.98)
    → c0-13                                   (commute: 118.59)
    → c0-11                                   (commute: 133.21)
    → c0-02                                   (commute: 136.61)
  c0-05:
    → c0-03                                   (commute: 92.85)
    → c0-12                                   (commute: 119.03)
    → c0-09                                   (commute: 125.58)
    → c0-13                                   (commute: 125.72)
    → c0-11                                   (commute: 131.75)
  c0-06:
    → c0-07                                   (commute: 137.64)
    → c0-12                                   (commute: 141.95)
    → c0-13                                   (commute: 149.98)
    → c0-11                                   (commute: 154.06)
    → c0-02                                   (commute: 157.49)
  c0-07:
    → c0-13                                   (commute: 99.66)
    → c0-02                                   (commute: 101.20)
    → c0-09                                   (commute: 104.74)
    → c0-10                                   (commute: 109.85)
    → c0-04                                   (commute: 112.98)
  c0-08:
    → c0-11                                   (commute: 77.98)
    → c0-02                                   (commute: 84.40)
    → c0-09                                   (commute: 86.57)
    → c0-10                                   (commute: 86.64)
    → c0-04                                   (commute: 98.43)
  c0-09:
    → c0-08                                   (commute: 86.57)
    → c0-12                                   (commute: 102.83)
    → c0-07                                   (commute: 104.74)
    → c0-00                                   (commute: 112.54)
    → c0-11                                   (commute: 123.56)
  c0-10:
    → c0-08                                   (commute: 86.64)
    → c0-07                                   (commute: 109.85)
    → c0-12                                   (commute: 112.75)
    → c0-11                                   (commute: 114.31)
    → c0-13                                   (commute: 123.98)
  c0-11:
    → c0-08                                   (commute: 77.98)
    → c0-12                                   (commute: 102.51)
    → c0-10                                   (commute: 114.31)
    → c0-13                                   (commute: 117.17)
    → c0-09                                   (commute: 123.56)
  c0-12:
    → c0-11                                   (commute: 102.51)
    → c0-13                                   (commute: 102.69)
    → c0-09                                   (commute: 102.83)
    → c0-02                                   (commute: 108.99)
    → c0-10                                   (commute: 112.75)
  c0-13:
    → c0-07                                   (commute: 99.66)
    → c0-12                                   (commute: 102.69)
    → c0-11                                   (commute: 117.17)
    → c0-04                                   (commute: 118.59)
    → c0-15                                   (commute: 118.78)
  c0-14:
    → c0-03                                   (commute: 123.51)
    → c0-07                                   (commute: 153.33)
    → c0-02                                   (commute: 156.23)
    → c0-12                                   (commute: 158.05)
    → c0-11                                   (commute: 166.49)
  c0-15:
    → c0-03                                   (commute: 100.90)
    → c0-08                                   (commute: 110.99)
    → c0-13                                   (commute: 118.78)
    → c0-07                                   (commute: 125.07)
    → c0-12                                   (commute: 136.15)
  c1-00:
    → c1-06                                   (commute: 203.82)
    → c1-11                                   (commute: 214.88)
    → c1-02                                   (commute: 216.19)
    → c1-04                                   (commute: 218.45)
    → c1-05                                   (commute: 224.26)
  c1-01:
    → c0-03                                   (commute: 151.84)
    → c2-05                                   (commute: 160.09)
    → c0-08                                   (commute: 162.41)
    → c2-07                                   (commute: 176.02)
    → c0-07                                   (commute: 176.16)
  c1-02:
    → c1-04                                   (commute: 179.43)
    → c0-03                                   (commute: 184.73)
    → c0-08                                   (commute: 204.23)
    → c1-03                                   (commute: 214.40)
    → c1-00                                   (commute: 216.19)
  c1-03:
    → c0-03                                   (commute: 103.10)
    → c0-07                                   (commute: 124.12)
    → c0-12                                   (commute: 136.52)
    → c0-02                                   (commute: 140.48)
    → c0-13                                   (commute: 148.98)
  c1-04:
    → c1-12                                   (commute: 160.00)
    → c0-03                                   (commute: 172.09)
    → c0-08                                   (commute: 178.16)
    → c1-02                                   (commute: 179.43)
    → c1-05                                   (commute: 181.81)
  c1-05:
    → c1-06                                   (commute: 133.66)
    → c0-08                                   (commute: 156.67)
    → c0-07                                   (commute: 175.93)
    → c0-12                                   (commute: 178.74)
    → c0-11                                   (commute: 179.11)
  c1-06:
    → c1-05                                   (commute: 133.66)
    → c0-03                                   (commute: 171.96)
    → c1-03                                   (commute: 182.87)
    → c0-08                                   (commute: 186.67)
    → c1-11                                   (commute: 193.29)
  c1-07:
    → c1-12                                   (commute: 229.75)
    → c1-02                                   (commute: 236.90)
    → c1-09                                   (commute: 246.11)
    → c1-04                                   (commute: 249.84)
    → c1-05                                   (commute: 252.67)
  c1-08:
    → c1-06                                   (commute: 223.09)
    → c1-10                                   (commute: 238.61)
    → c1-03                                   (commute: 244.51)
    → c1-12                                   (commute: 249.33)
    → c1-02                                   (commute: 275.13)
  c1-09:
    → c2-00                                   (commute: 212.21)
    → c2-05                                   (commute: 215.61)
    → c2-08                                   (commute: 232.37)
    → c0-03                                   (commute: 245.58)
    → c1-07                                   (commute: 246.11)
  c1-10:
    → c1-04                                   (commute: 200.33)
    → c1-12                                   (commute: 214.88)
    → c1-06                                   (commute: 236.84)
    → c1-08                                   (commute: 238.61)
    → c1-02                                   (commute: 266.20)
  c1-11:
    → c1-06                                   (commute: 193.29)
    → c1-12                                   (commute: 207.47)
    → c1-00                                   (commute: 214.88)
    → c1-03                                   (commute: 223.27)
    → c1-02                                   (commute: 239.66)
  c1-12:
    → c1-04                                   (commute: 160.00)
    → c0-03                                   (commute: 175.95)
    → c0-08                                   (commute: 194.59)
    → c1-03                                   (commute: 201.33)
    → c1-11                                   (commute: 207.47)
  c1-13:
    → c1-06                                   (commute: 246.11)
    → c2-07                                   (commute: 246.11)
    → c0-03                                   (commute: 281.98)
    → c2-05                                   (commute: 287.46)
    → c2-00                                   (commute: 288.27)
  c2-00:
    → c2-08                                   (commute: 119.61)
    → c0-03                                   (commute: 138.85)
    → c2-09                                   (commute: 141.29)
    → c2-03                                   (commute: 146.93)
    → c0-01                                   (commute: 148.92)
  c2-01:
    → c2-07                                   (commute: 156.09)
    → c2-05                                   (commute: 176.74)
    → c2-02                                   (commute: 204.92)
    → c2-08                                   (commute: 219.00)
    → c2-06                                   (commute: 229.69)
  c2-02:
    → c2-07                                   (commute: 113.80)
    → c2-08                                   (commute: 156.18)
    → c2-06                                   (commute: 168.11)
    → c2-09                                   (commute: 174.73)
    → c0-03                                   (commute: 185.53)
  c2-03:
    → c2-00                                   (commute: 146.93)
    → c2-05                                   (commute: 151.48)
    → c2-08                                   (commute: 191.61)
    → c2-09                                   (commute: 207.27)
    → c2-06                                   (commute: 225.01)
  c2-04:
    → c2-00                                   (commute: 156.09)
    → c2-05                                   (commute: 182.75)
    → c2-02                                   (commute: 212.34)
    → c2-08                                   (commute: 212.45)
    → c2-09                                   (commute: 228.18)
  c2-05:
    → c2-07                                   (commute: 82.81)
    → c0-03                                   (commute: 124.38)
    → c0-08                                   (commute: 140.55)
    → c0-12                                   (commute: 146.42)
    → c3-06                                   (commute: 148.48)
  c2-06:
    → c2-07                                   (commute: 154.82)
    → c2-02                                   (commute: 168.11)
    → c2-08                                   (commute: 182.80)
    → c2-09                                   (commute: 202.80)
    → c0-03                                   (commute: 207.42)
  c2-07:
    → c2-05                                   (commute: 82.81)
    → c2-02                                   (commute: 113.80)
    → c0-03                                   (commute: 148.23)
    → c2-06                                   (commute: 154.82)
    → c2-01                                   (commute: 156.09)
  c2-08:
    → c2-00                                   (commute: 119.61)
    → c0-03                                   (commute: 140.74)
    → c0-08                                   (commute: 154.14)
    → c2-02                                   (commute: 156.18)
    → c3-01                                   (commute: 156.46)
  c2-09:
    → c2-00                                   (commute: 141.29)
    → c2-08                                   (commute: 168.46)
    → c2-02                                   (commute: 174.73)
    → c2-06                                   (commute: 202.80)
    → c2-03                                   (commute: 207.27)
  c3-00:
    → c3-07                                   (commute: 78.47)
    → c0-03                                   (commute: 108.36)
    → c0-08                                   (commute: 115.96)
    → c0-02                                   (commute: 130.39)
    → c3-09                                   (commute: 140.92)
  c3-01:
    → c3-06                                   (commute: 69.12)
    → c3-08                                   (commute: 84.20)
    → c0-03                                   (commute: 87.80)
    → c0-02                                   (commute: 111.37)
    → c0-07                                   (commute: 120.97)
  c3-02:
    → c3-07                                   (commute: 91.62)
    → c3-03                                   (commute: 98.07)
    → c0-08                                   (commute: 102.77)
    → c3-08                                   (commute: 105.96)
    → c0-02                                   (commute: 121.83)
  c3-03:
    → c3-02                                   (commute: 98.07)
    → c3-08                                   (commute: 99.07)
    → c0-03                                   (commute: 118.72)
    → c0-08                                   (commute: 125.33)
    → c0-02                                   (commute: 135.84)
  c3-04:
    → c3-00                                   (commute: 245.24)
    → c3-07                                   (commute: 247.90)
    → c3-01                                   (commute: 253.12)
    → c3-03                                   (commute: 253.53)
    → c3-02                                   (commute: 256.57)
  c3-05:
    → c3-00                                   (commute: 245.24)
    → c3-07                                   (commute: 247.90)
    → c3-01                                   (commute: 253.12)
    → c3-03                                   (commute: 253.53)
    → c3-02                                   (commute: 256.57)
  c3-06:
    → c3-01                                   (commute: 69.12)
    → c0-03                                   (commute: 97.75)
    → c0-08                                   (commute: 107.33)
    → c0-02                                   (commute: 118.20)
    → c0-12                                   (commute: 133.74)
  c3-07:
    → c3-00                                   (commute: 78.47)
    → c3-02                                   (commute: 91.62)
    → c0-03                                   (commute: 97.14)
    → c0-08                                   (commute: 105.82)
    → c0-07                                   (commute: 132.29)
  c3-08:
    → c3-01                                   (commute: 84.20)
    → c3-03                                   (commute: 99.07)
    → c3-02                                   (commute: 105.96)
    → c0-03                                   (commute: 127.09)
    → c0-08                                   (commute: 134.57)
  c3-09:
    → c3-00                                   (commute: 140.92)
    → c3-07                                   (commute: 142.28)
    → c3-06                                   (commute: 142.28)
    → c3-03                                   (commute: 155.59)
    → c3-02                                   (commute: 159.20)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `c0-03` (0.076087), `c0-08` (0.048913). Lowest π: `c3-04`/`c3-05` (tied, π=0.005435;
`c3-04` shown).

| From  | To    | Commute time | Within / Cross community |
| ----- | ----- | ------------ | --------------------------- |
| c0-03 | c3-01 | 87.80         | Cross                        |
| c0-03 | c0-05 | 92.85         | Within                       |
| c0-03 | c3-07 | 97.14         | Cross                        |
| c0-03 | c3-06 | 97.75         | Cross                        |
| c0-03 | c0-15 | 100.90        | Within                       |
| c0-08 | c0-11 | 77.98         | Within                       |
| c0-08 | c0-02 | 84.40         | Within                       |
| c0-08 | c0-09 | 86.57         | Within                       |
| c0-08 | c0-10 | 86.64         | Within                       |
| c0-08 | c0-04 | 98.43         | Within                       |
| c3-04 | c3-00 | 245.24        | Within                       |
| c3-04 | c3-07 | 247.90        | Within                       |
| c3-04 | c3-01 | 253.12        | Within                       |
| c3-04 | c3-03 | 253.53        | Within                       |
| c3-04 | c3-02 | 256.57        | Within                       |

**Finding**: At μ=0.1, the top-π node `c0-03` splits 3 cross / 2 within (60% cross), while
`c0-08` (rank-2 π) suggests exclusively within-community links (0% cross) — the two top nodes
diverge sharply in suggestion composition. `c3-04` (lowest π, an almost fully isolated leaf)
suggests exclusively within-community links (5 of 5, 0% cross) at very high commute times
(245–257), all targeting other members of its own small community — consistent with being a
near-orphan node whose only realistic connections remain unexplored within its home community.

### μ=0.3

#### Raw analyze output

```
=== Overview ===
Pages:        50
Edges:        202
Entropy rate: 2.2020 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 50 page(s)
  c0-00
  c0-01
  c0-02
  c0-03
  c0-04
  c0-05
  c0-06
  c0-07
  c0-08
  c0-09
  c0-10
  c0-11
  c0-12
  c0-13
  c0-14
  c0-15
  c1-00
  c1-01
  c1-02
  c1-03
  c1-04
  c1-05
  c1-06
  c1-07
  c1-08
  c1-09
  c1-10
  c1-11
  c1-12
  c1-13
  c2-00
  c2-01
  c2-02
  c2-03
  c2-04
  c2-05
  c2-06
  c2-07
  c2-08
  c2-09
  c3-00
  c3-01
  c3-02
  c3-03
  c3-04
  c3-05
  c3-06
  c3-07
  c3-08
  c3-09

=== Orphan pages (bottom 10% by stationary distribution) ===
  c2-04                                     π=0.009901  → add inbound links
  c1-00                                     π=0.009901  → add inbound links
  c0-06                                     π=0.009901  → add inbound links
  c0-12                                     π=0.009901  → add inbound links
  c1-07                                     π=0.009901  → add inbound links
  c1-11                                     π=0.009901  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. c1-03                                     π=0.079208
  2. c2-00                                     π=0.039604
  3. c3-06                                     π=0.039604
  4. c0-07                                     π=0.039604
  5. c3-01                                     π=0.034653

=== Suggested missing links (lowest commute time, not yet linked) ===
  c0-00:
    → c1-03                                   (commute: 95.28)
    → c3-06                                   (commute: 121.86)
    → c3-05                                   (commute: 122.49)
    → c3-01                                   (commute: 124.97)
    → c2-05                                   (commute: 125.33)
  c0-01:
    → c2-00                                   (commute: 119.14)
    → c3-06                                   (commute: 135.20)
    → c0-07                                   (commute: 135.78)
    → c2-05                                   (commute: 137.41)
    → c3-05                                   (commute: 139.86)
  c0-02:
    → c1-03                                   (commute: 145.19)
    → c2-00                                   (commute: 152.73)
    → c3-05                                   (commute: 158.68)
    → c2-05                                   (commute: 159.03)
    → c3-09                                   (commute: 163.35)
  c0-03:
    → c1-03                                   (commute: 80.88)
    → c3-06                                   (commute: 96.79)
    → c3-01                                   (commute: 100.43)
    → c2-00                                   (commute: 104.63)
    → c2-05                                   (commute: 108.37)
  c0-04:
    → c1-03                                   (commute: 104.86)
    → c3-06                                   (commute: 123.66)
    → c3-05                                   (commute: 130.53)
    → c2-00                                   (commute: 134.62)
    → c3-09                                   (commute: 138.72)
  c0-05:
    → c1-03                                   (commute: 149.18)
    → c3-05                                   (commute: 168.34)
    → c2-00                                   (commute: 168.58)
    → c3-06                                   (commute: 173.35)
    → c2-05                                   (commute: 174.23)
  c0-06:
    → c1-03                                   (commute: 150.18)
    → c0-07                                   (commute: 170.70)
    → c2-00                                   (commute: 171.55)
    → c3-06                                   (commute: 173.05)
    → c3-01                                   (commute: 174.31)
  c0-07:
    → c1-03                                   (commute: 58.24)
    → c2-00                                   (commute: 79.35)
    → c3-06                                   (commute: 80.95)
    → c3-01                                   (commute: 84.13)
    → c2-05                                   (commute: 85.11)
  c0-08:
    → c1-03                                   (commute: 86.64)
    → c3-01                                   (commute: 111.10)
    → c2-00                                   (commute: 114.54)
    → c3-05                                   (commute: 114.57)
    → c3-06                                   (commute: 114.67)
  c0-09:
    → c1-03                                   (commute: 87.40)
    → c3-01                                   (commute: 110.62)
    → c3-05                                   (commute: 113.91)
    → c3-06                                   (commute: 114.49)
    → c0-07                                   (commute: 116.20)
  c0-10:
    → c1-03                                   (commute: 93.14)
    → c0-07                                   (commute: 110.07)
    → c2-00                                   (commute: 110.84)
    → c2-05                                   (commute: 116.39)
    → c3-06                                   (commute: 118.79)
  c0-11:
    → c1-03                                   (commute: 99.81)
    → c3-05                                   (commute: 114.98)
    → c3-01                                   (commute: 118.35)
    → c2-00                                   (commute: 118.46)
    → c1-08                                   (commute: 118.47)
  c0-12:
    → c1-03                                   (commute: 162.08)
    → c2-00                                   (commute: 180.82)
    → c3-06                                   (commute: 194.66)
    → c3-09                                   (commute: 196.68)
    → c3-05                                   (commute: 197.00)
  c0-13:
    → c1-03                                   (commute: 147.60)
    → c2-00                                   (commute: 159.17)
    → c3-05                                   (commute: 166.19)
    → c3-06                                   (commute: 171.74)
    → c2-05                                   (commute: 175.35)
  c0-14:
    → c1-03                                   (commute: 147.81)
    → c2-01                                   (commute: 168.62)
    → c3-05                                   (commute: 170.81)
    → c3-06                                   (commute: 173.68)
    → c2-00                                   (commute: 176.62)
  c0-15:
    → c1-03                                   (commute: 108.64)
    → c3-06                                   (commute: 126.95)
    → c0-08                                   (commute: 131.56)
    → c2-00                                   (commute: 133.95)
    → c3-05                                   (commute: 134.63)
  c1-00:
    → c1-03                                   (commute: 145.06)
    → c3-06                                   (commute: 164.47)
    → c2-00                                   (commute: 167.96)
    → c3-05                                   (commute: 168.13)
    → c0-07                                   (commute: 173.12)
  c1-01:
    → c1-03                                   (commute: 86.72)
    → c0-07                                   (commute: 103.53)
    → c3-06                                   (commute: 112.29)
    → c3-05                                   (commute: 114.50)
    → c3-09                                   (commute: 114.85)
  c1-02:
    → c3-09                                   (commute: 136.71)
    → c2-00                                   (commute: 138.22)
    → c3-06                                   (commute: 139.33)
    → c3-05                                   (commute: 143.39)
    → c3-01                                   (commute: 147.39)
  c1-03:
    → c2-00                                   (commute: 55.56)
    → c0-07                                   (commute: 58.24)
    → c3-01                                   (commute: 59.01)
    → c3-09                                   (commute: 61.28)
    → c2-05                                   (commute: 63.29)
  c1-04:
    → c2-00                                   (commute: 109.88)
    → c3-05                                   (commute: 116.92)
    → c3-06                                   (commute: 117.70)
    → c2-05                                   (commute: 122.41)
    → c2-06                                   (commute: 123.58)
  c1-05:
    → c2-00                                   (commute: 104.22)
    → c3-09                                   (commute: 105.68)
    → c0-07                                   (commute: 109.68)
    → c3-06                                   (commute: 109.98)
    → c3-05                                   (commute: 113.63)
  c1-06:
    → c0-07                                   (commute: 119.37)
    → c3-06                                   (commute: 124.09)
    → c2-05                                   (commute: 129.97)
    → c2-00                                   (commute: 131.58)
    → c3-05                                   (commute: 131.66)
  c1-07:
    → c3-06                                   (commute: 161.69)
    → c3-01                                   (commute: 163.97)
    → c3-05                                   (commute: 164.12)
    → c0-07                                   (commute: 166.37)
    → c2-00                                   (commute: 167.44)
  c1-08:
    → c3-01                                   (commute: 108.23)
    → c2-00                                   (commute: 111.87)
    → c3-06                                   (commute: 112.11)
    → c3-05                                   (commute: 116.14)
    → c0-11                                   (commute: 118.47)
  c1-09:
    → c1-03                                   (commute: 111.47)
    → c2-00                                   (commute: 129.72)
    → c3-06                                   (commute: 142.50)
    → c2-07                                   (commute: 142.56)
    → c3-01                                   (commute: 144.56)
  c1-10:
    → c2-00                                   (commute: 153.40)
    → c3-06                                   (commute: 160.45)
    → c3-07                                   (commute: 164.11)
    → c3-01                                   (commute: 169.01)
    → c3-05                                   (commute: 170.80)
  c1-11:
    → c0-07                                   (commute: 155.54)
    → c3-01                                   (commute: 156.00)
    → c3-06                                   (commute: 161.91)
    → c3-05                                   (commute: 168.07)
    → c2-00                                   (commute: 169.88)
  c1-12:
    → c3-06                                   (commute: 113.14)
    → c3-05                                   (commute: 124.01)
    → c0-07                                   (commute: 124.97)
    → c2-00                                   (commute: 125.81)
    → c3-09                                   (commute: 129.26)
  c1-13:
    → c3-06                                   (commute: 120.71)
    → c3-05                                   (commute: 125.11)
    → c2-00                                   (commute: 126.25)
    → c3-01                                   (commute: 130.90)
    → c0-07                                   (commute: 132.06)
  c2-00:
    → c1-03                                   (commute: 55.56)
    → c3-06                                   (commute: 74.30)
    → c0-07                                   (commute: 79.35)
    → c3-05                                   (commute: 80.13)
    → c3-09                                   (commute: 82.54)
  c2-01:
    → c2-00                                   (commute: 86.46)
    → c3-05                                   (commute: 92.03)
    → c0-07                                   (commute: 92.92)
    → c3-01                                   (commute: 94.44)
    → c3-09                                   (commute: 96.80)
  c2-02:
    → c1-03                                   (commute: 107.73)
    → c3-06                                   (commute: 122.17)
    → c3-01                                   (commute: 130.88)
    → c2-05                                   (commute: 133.85)
    → c3-05                                   (commute: 137.42)
  c2-03:
    → c1-03                                   (commute: 85.90)
    → c2-00                                   (commute: 94.55)
    → c3-06                                   (commute: 99.77)
    → c0-07                                   (commute: 109.43)
    → c3-05                                   (commute: 109.60)
  c2-04:
    → c1-03                                   (commute: 149.47)
    → c2-05                                   (commute: 172.31)
    → c3-06                                   (commute: 175.97)
    → c0-07                                   (commute: 178.78)
    → c2-07                                   (commute: 179.56)
  c2-05:
    → c1-03                                   (commute: 63.29)
    → c3-06                                   (commute: 73.48)
    → c0-07                                   (commute: 85.11)
    → c3-09                                   (commute: 86.37)
    → c3-01                                   (commute: 90.47)
  c2-06:
    → c1-03                                   (commute: 74.55)
    → c0-07                                   (commute: 88.82)
    → c3-06                                   (commute: 90.28)
    → c3-05                                   (commute: 92.13)
    → c3-01                                   (commute: 102.10)
  c2-07:
    → c1-03                                   (commute: 72.68)
    → c2-05                                   (commute: 91.91)
    → c0-07                                   (commute: 94.28)
    → c3-06                                   (commute: 94.63)
    → c3-05                                   (commute: 100.14)
  c2-08:
    → c2-00                                   (commute: 115.39)
    → c3-06                                   (commute: 126.94)
    → c0-07                                   (commute: 131.44)
    → c3-05                                   (commute: 131.63)
    → c2-05                                   (commute: 132.67)
  c2-09:
    → c1-03                                   (commute: 86.59)
    → c3-06                                   (commute: 95.60)
    → c2-05                                   (commute: 95.99)
    → c2-06                                   (commute: 107.10)
    → c2-01                                   (commute: 108.75)
  c3-00:
    → c1-03                                   (commute: 106.60)
    → c2-00                                   (commute: 119.23)
    → c3-01                                   (commute: 120.09)
    → c3-07                                   (commute: 122.72)
    → c2-01                                   (commute: 124.48)
  c3-01:
    → c1-03                                   (commute: 59.01)
    → c3-09                                   (commute: 83.88)
    → c0-07                                   (commute: 84.13)
    → c2-00                                   (commute: 84.53)
    → c3-05                                   (commute: 85.77)
  c3-02:
    → c1-03                                   (commute: 107.81)
    → c3-07                                   (commute: 117.02)
    → c2-00                                   (commute: 128.45)
    → c3-09                                   (commute: 130.26)
    → c2-05                                   (commute: 130.65)
  c3-03:
    → c2-00                                   (commute: 102.42)
    → c3-01                                   (commute: 102.81)
    → c3-06                                   (commute: 103.95)
    → c3-09                                   (commute: 105.74)
    → c3-05                                   (commute: 108.87)
  c3-04:
    → c1-03                                   (commute: 99.89)
    → c3-06                                   (commute: 112.09)
    → c3-05                                   (commute: 124.21)
    → c2-00                                   (commute: 124.67)
    → c0-07                                   (commute: 128.35)
  c3-05:
    → c3-06                                   (commute: 77.84)
    → c2-00                                   (commute: 80.13)
    → c3-01                                   (commute: 85.77)
    → c2-01                                   (commute: 92.03)
    → c2-06                                   (commute: 92.13)
  c3-06:
    → c2-05                                   (commute: 73.48)
    → c2-00                                   (commute: 74.30)
    → c3-05                                   (commute: 77.84)
    → c0-07                                   (commute: 80.95)
    → c2-06                                   (commute: 90.28)
  c3-07:
    → c1-03                                   (commute: 73.41)
    → c2-00                                   (commute: 89.62)
    → c3-05                                   (commute: 95.04)
    → c0-07                                   (commute: 95.97)
    → c2-05                                   (commute: 97.83)
  c3-08:
    → c1-03                                   (commute: 80.88)
    → c3-01                                   (commute: 86.76)
    → c2-00                                   (commute: 93.56)
    → c2-05                                   (commute: 95.04)
    → c3-09                                   (commute: 104.15)
  c3-09:
    → c1-03                                   (commute: 61.28)
    → c2-00                                   (commute: 82.54)
    → c3-01                                   (commute: 83.88)
    → c2-05                                   (commute: 86.37)
    → c0-07                                   (commute: 87.39)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `c1-03` (0.079208), `c2-00` (0.039604, tied with `c3-06`/`c0-07`; `c2-00` shown).
Lowest π: `c2-04`/`c1-00`/`c0-06`/`c0-12`/`c1-07`/`c1-11` (all tied, π=0.009901; `c2-04` shown).

| From  | To    | Commute time | Within / Cross community |
| ----- | ----- | ------------ | --------------------------- |
| c1-03 | c2-00 | 55.56         | Cross                        |
| c1-03 | c0-07 | 58.24         | Cross                        |
| c1-03 | c3-01 | 59.01         | Cross                        |
| c1-03 | c3-09 | 61.28         | Cross                        |
| c1-03 | c2-05 | 63.29         | Cross                        |
| c2-00 | c1-03 | 55.56         | Cross                        |
| c2-00 | c3-06 | 74.30         | Cross                        |
| c2-00 | c0-07 | 79.35         | Cross                        |
| c2-00 | c3-05 | 80.13         | Cross                        |
| c2-00 | c3-09 | 82.54         | Cross                        |
| c2-04 | c1-03 | 149.47        | Cross                        |
| c2-04 | c2-05 | 172.31        | Within                       |
| c2-04 | c3-06 | 175.97        | Cross                        |
| c2-04 | c0-07 | 178.78        | Cross                        |
| c2-04 | c2-07 | 179.56        | Within                       |

**Finding**: At μ=0.3, both top-π nodes (`c1-03`, `c2-00`) suggest **exclusively
cross-community** links (5 of 5 each, 100%) — a marked change from μ=0.1's mixed pattern (60%
and 0%). `c2-04` (lowest π) shows 3 cross / 2 within (60% cross) — even the lowest-π node's
suggestions have shifted toward cross-community as μ rose.

### μ=0.5

#### Raw analyze output

```
=== Overview ===
Pages:        50
Edges:        226
Entropy rate: 2.4113 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 50 page(s)
  c0-00
  c0-01
  c0-02
  c0-03
  c0-04
  c0-05
  c0-06
  c0-07
  c0-08
  c0-09
  c0-10
  c0-11
  c0-12
  c0-13
  c0-14
  c0-15
  c1-00
  c1-01
  c1-02
  c1-03
  c1-04
  c1-05
  c1-06
  c1-07
  c1-08
  c1-09
  c1-10
  c1-11
  c1-12
  c1-13
  c2-00
  c2-01
  c2-02
  c2-03
  c2-04
  c2-05
  c2-06
  c2-07
  c2-08
  c2-09
  c3-00
  c3-01
  c3-02
  c3-03
  c3-04
  c3-05
  c3-06
  c3-07
  c3-08
  c3-09

=== Orphan pages (bottom 10% by stationary distribution) ===
  c1-09                                     π=0.008850  → add inbound links
  c1-06                                     π=0.008850  → add inbound links
  c1-03                                     π=0.008850  → add inbound links
  c1-01                                     π=0.008850  → add inbound links
  c2-09                                     π=0.008850  → add inbound links
  c1-07                                     π=0.008850  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. c2-02                                     π=0.079646
  2. c0-07                                     π=0.044248
  3. c3-06                                     π=0.044248
  4. c3-05                                     π=0.035398
  5. c3-09                                     π=0.035398

=== Suggested missing links (lowest commute time, not yet linked) ===
  c0-00:
    → c2-02                                   (commute: 72.58)
    → c3-06                                   (commute: 82.22)
    → c0-07                                   (commute: 87.23)
    → c3-08                                   (commute: 88.39)
    → c0-09                                   (commute: 93.28)
  c0-01:
    → c2-02                                   (commute: 52.70)
    → c0-07                                   (commute: 59.75)
    → c3-06                                   (commute: 64.79)
    → c3-09                                   (commute: 71.69)
    → c2-05                                   (commute: 71.98)
  c0-02:
    → c2-02                                   (commute: 89.31)
    → c0-07                                   (commute: 98.24)
    → c3-09                                   (commute: 100.51)
    → c3-06                                   (commute: 101.93)
    → c0-01                                   (commute: 103.21)
  c0-03:
    → c3-06                                   (commute: 83.31)
    → c0-01                                   (commute: 90.30)
    → c2-05                                   (commute: 90.43)
    → c3-05                                   (commute: 92.04)
    → c0-08                                   (commute: 94.49)
  c0-04:
    → c2-02                                   (commute: 120.24)
    → c0-07                                   (commute: 121.45)
    → c3-06                                   (commute: 128.65)
    → c0-01                                   (commute: 128.93)
    → c0-09                                   (commute: 131.54)
  c0-05:
    → c2-02                                   (commute: 105.49)
    → c3-06                                   (commute: 117.03)
    → c3-05                                   (commute: 120.90)
    → c0-01                                   (commute: 121.17)
    → c2-05                                   (commute: 124.44)
  c0-06:
    → c2-02                                   (commute: 114.78)
    → c3-06                                   (commute: 118.79)
    → c0-07                                   (commute: 122.72)
    → c0-01                                   (commute: 126.59)
    → c0-08                                   (commute: 128.96)
  c0-07:
    → c3-06                                   (commute: 56.28)
    → c0-01                                   (commute: 59.75)
    → c0-08                                   (commute: 71.05)
    → c0-09                                   (commute: 71.30)
    → c1-02                                   (commute: 73.31)
  c0-08:
    → c2-02                                   (commute: 58.15)
    → c3-06                                   (commute: 68.67)
    → c0-07                                   (commute: 71.05)
    → c3-09                                   (commute: 78.95)
    → c3-05                                   (commute: 80.84)
  c0-09:
    → c0-07                                   (commute: 71.30)
    → c3-06                                   (commute: 72.67)
    → c3-05                                   (commute: 82.02)
    → c2-05                                   (commute: 82.33)
    → c0-08                                   (commute: 84.10)
  c0-10:
    → c2-02                                   (commute: 63.48)
    → c0-07                                   (commute: 79.09)
    → c3-05                                   (commute: 86.40)
    → c3-09                                   (commute: 90.26)
    → c0-08                                   (commute: 92.81)
  c0-11:
    → c2-02                                   (commute: 88.54)
    → c0-07                                   (commute: 101.01)
    → c3-05                                   (commute: 102.39)
    → c0-01                                   (commute: 105.20)
    → c2-05                                   (commute: 106.74)
  c0-12:
    → c2-02                                   (commute: 94.52)
    → c3-06                                   (commute: 101.69)
    → c3-09                                   (commute: 105.24)
    → c0-07                                   (commute: 106.25)
    → c0-01                                   (commute: 111.44)
  c0-13:
    → c2-02                                   (commute: 145.68)
    → c3-06                                   (commute: 158.59)
    → c3-09                                   (commute: 162.38)
    → c2-05                                   (commute: 162.54)
    → c3-05                                   (commute: 162.90)
  c0-14:
    → c2-02                                   (commute: 84.72)
    → c3-06                                   (commute: 95.34)
    → c0-01                                   (commute: 103.16)
    → c3-09                                   (commute: 103.90)
    → c0-09                                   (commute: 109.35)
  c0-15:
    → c0-07                                   (commute: 115.82)
    → c3-06                                   (commute: 116.54)
    → c3-05                                   (commute: 123.06)
    → c0-09                                   (commute: 124.55)
    → c2-05                                   (commute: 124.93)
  c1-00:
    → c2-02                                   (commute: 186.73)
    → c3-09                                   (commute: 193.58)
    → c3-08                                   (commute: 193.73)
    → c3-06                                   (commute: 194.67)
    → c0-07                                   (commute: 194.71)
  c1-01:
    → c3-06                                   (commute: 178.59)
    → c0-07                                   (commute: 179.92)
    → c2-05                                   (commute: 185.69)
    → c0-01                                   (commute: 189.05)
    → c3-09                                   (commute: 191.00)
  c1-02:
    → c2-02                                   (commute: 68.45)
    → c0-07                                   (commute: 73.31)
    → c3-06                                   (commute: 81.16)
    → c2-05                                   (commute: 83.92)
    → c0-01                                   (commute: 86.18)
  c1-03:
    → c2-02                                   (commute: 160.97)
    → c3-09                                   (commute: 171.23)
    → c3-06                                   (commute: 173.71)
    → c0-01                                   (commute: 173.99)
    → c0-07                                   (commute: 174.31)
  c1-04:
    → c2-02                                   (commute: 143.14)
    → c0-07                                   (commute: 154.24)
    → c3-05                                   (commute: 163.95)
    → c3-09                                   (commute: 165.50)
    → c0-01                                   (commute: 166.01)
  c1-05:
    → c2-02                                   (commute: 112.70)
    → c3-09                                   (commute: 119.86)
    → c3-06                                   (commute: 124.11)
    → c2-05                                   (commute: 130.86)
    → c3-05                                   (commute: 130.94)
  c1-06:
    → c3-05                                   (commute: 189.06)
    → c2-02                                   (commute: 192.64)
    → c3-06                                   (commute: 198.22)
    → c0-07                                   (commute: 203.22)
    → c2-05                                   (commute: 205.22)
  c1-07:
    → c2-02                                   (commute: 150.16)
    → c0-07                                   (commute: 163.06)
    → c3-06                                   (commute: 164.65)
    → c0-01                                   (commute: 173.15)
    → c2-05                                   (commute: 174.74)
  c1-08:
    → c2-02                                   (commute: 97.93)
    → c0-07                                   (commute: 103.11)
    → c3-06                                   (commute: 105.66)
    → c0-01                                   (commute: 115.07)
    → c0-08                                   (commute: 116.05)
  c1-09:
    → c2-02                                   (commute: 157.31)
    → c3-06                                   (commute: 165.55)
    → c0-07                                   (commute: 169.47)
    → c3-09                                   (commute: 181.83)
    → c0-08                                   (commute: 182.35)
  c1-10:
    → c2-02                                   (commute: 128.02)
    → c2-05                                   (commute: 131.21)
    → c0-07                                   (commute: 135.86)
    → c3-05                                   (commute: 139.60)
    → c3-06                                   (commute: 140.67)
  c1-11:
    → c2-02                                   (commute: 92.34)
    → c0-07                                   (commute: 106.11)
    → c3-06                                   (commute: 107.30)
    → c0-08                                   (commute: 109.08)
    → c0-01                                   (commute: 109.14)
  c1-12:
    → c3-06                                   (commute: 142.22)
    → c0-07                                   (commute: 144.20)
    → c2-05                                   (commute: 149.11)
    → c0-10                                   (commute: 149.76)
    → c0-01                                   (commute: 152.31)
  c1-13:
    → c2-02                                   (commute: 151.24)
    → c0-07                                   (commute: 168.78)
    → c3-08                                   (commute: 170.09)
    → c0-08                                   (commute: 170.71)
    → c2-05                                   (commute: 174.97)
  c2-00:
    → c0-07                                   (commute: 114.08)
    → c3-06                                   (commute: 115.06)
    → c2-05                                   (commute: 121.36)
    → c0-01                                   (commute: 122.97)
    → c3-09                                   (commute: 122.98)
  c2-01:
    → c3-06                                   (commute: 80.08)
    → c0-07                                   (commute: 82.85)
    → c0-01                                   (commute: 84.61)
    → c3-09                                   (commute: 90.26)
    → c2-05                                   (commute: 90.42)
  c2-02:
    → c0-01                                   (commute: 52.70)
    → c3-09                                   (commute: 53.61)
    → c3-05                                   (commute: 56.57)
    → c0-08                                   (commute: 58.15)
    → c0-10                                   (commute: 63.48)
  c2-03:
    → c0-07                                   (commute: 91.83)
    → c3-06                                   (commute: 94.02)
    → c3-09                                   (commute: 97.90)
    → c2-05                                   (commute: 102.77)
    → c0-08                                   (commute: 104.33)
  c2-04:
    → c2-02                                   (commute: 168.80)
    → c0-07                                   (commute: 176.61)
    → c1-02                                   (commute: 181.10)
    → c3-06                                   (commute: 181.28)
    → c3-05                                   (commute: 186.18)
  c2-05:
    → c3-06                                   (commute: 61.51)
    → c0-01                                   (commute: 71.98)
    → c3-05                                   (commute: 72.37)
    → c3-09                                   (commute: 75.48)
    → c0-09                                   (commute: 82.33)
  c2-06:
    → c0-07                                   (commute: 78.39)
    → c2-05                                   (commute: 90.01)
    → c3-09                                   (commute: 90.93)
    → c0-09                                   (commute: 91.26)
    → c3-05                                   (commute: 91.42)
  c2-07:
    → c3-09                                   (commute: 81.84)
    → c3-06                                   (commute: 83.94)
    → c0-07                                   (commute: 84.10)
    → c2-05                                   (commute: 88.92)
    → c0-01                                   (commute: 89.53)
  c2-08:
    → c3-06                                   (commute: 93.10)
    → c0-07                                   (commute: 97.62)
    → c3-09                                   (commute: 105.34)
    → c0-01                                   (commute: 107.21)
    → c3-05                                   (commute: 109.91)
  c2-09:
    → c3-06                                   (commute: 160.49)
    → c0-07                                   (commute: 160.91)
    → c3-09                                   (commute: 165.32)
    → c2-05                                   (commute: 166.64)
    → c0-01                                   (commute: 168.70)
  c3-00:
    → c2-02                                   (commute: 128.54)
    → c3-06                                   (commute: 136.67)
    → c0-07                                   (commute: 139.30)
    → c3-09                                   (commute: 147.17)
    → c3-05                                   (commute: 148.64)
  c3-01:
    → c3-06                                   (commute: 91.11)
    → c0-07                                   (commute: 100.83)
    → c2-05                                   (commute: 107.98)
    → c0-01                                   (commute: 109.31)
    → c3-09                                   (commute: 111.74)
  c3-02:
    → c2-02                                   (commute: 109.19)
    → c0-07                                   (commute: 125.32)
    → c2-05                                   (commute: 126.93)
    → c3-09                                   (commute: 128.88)
    → c0-09                                   (commute: 131.67)
  c3-03:
    → c3-06                                   (commute: 124.57)
    → c0-07                                   (commute: 128.49)
    → c2-05                                   (commute: 130.79)
    → c0-01                                   (commute: 133.51)
    → c3-09                                   (commute: 139.60)
  c3-04:
    → c0-10                                   (commute: 189.06)
    → c2-02                                   (commute: 190.34)
    → c3-06                                   (commute: 194.99)
    → c0-07                                   (commute: 195.88)
    → c2-05                                   (commute: 204.53)
  c3-05:
    → c2-02                                   (commute: 56.57)
    → c2-05                                   (commute: 72.37)
    → c3-09                                   (commute: 72.39)
    → c0-01                                   (commute: 72.61)
    → c0-08                                   (commute: 80.84)
  c3-06:
    → c0-07                                   (commute: 56.28)
    → c2-05                                   (commute: 61.51)
    → c0-01                                   (commute: 64.79)
    → c0-08                                   (commute: 68.67)
    → c0-09                                   (commute: 72.67)
  c3-07:
    → c2-02                                   (commute: 78.81)
    → c0-07                                   (commute: 89.66)
    → c3-06                                   (commute: 91.55)
    → c3-05                                   (commute: 92.81)
    → c2-05                                   (commute: 95.91)
  c3-08:
    → c3-09                                   (commute: 87.09)
    → c2-05                                   (commute: 87.14)
    → c0-00                                   (commute: 88.39)
    → c3-05                                   (commute: 88.39)
    → c0-01                                   (commute: 88.82)
  c3-09:
    → c2-02                                   (commute: 53.61)
    → c0-01                                   (commute: 71.69)
    → c3-05                                   (commute: 72.39)
    → c2-05                                   (commute: 75.48)
    → c0-08                                   (commute: 78.95)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `c2-02` (0.079646), `c0-07` (0.044248, tied with `c3-06`; `c0-07` shown). Lowest π:
`c1-09`/`c1-06`/`c1-03`/`c1-01`/`c2-09`/`c1-07` (all tied, π=0.008850; `c1-09` shown).

| From  | To    | Commute time | Within / Cross community |
| ----- | ----- | ------------ | --------------------------- |
| c2-02 | c0-01 | 52.70         | Cross                        |
| c2-02 | c3-09 | 53.61         | Cross                        |
| c2-02 | c3-05 | 56.57         | Cross                        |
| c2-02 | c0-08 | 58.15         | Cross                        |
| c2-02 | c0-10 | 63.48         | Cross                        |
| c0-07 | c3-06 | 56.28         | Cross                        |
| c0-07 | c0-01 | 59.75         | Within                       |
| c0-07 | c0-08 | 71.05         | Within                       |
| c0-07 | c0-09 | 71.30         | Within                       |
| c0-07 | c1-02 | 73.31         | Cross                        |
| c1-09 | c2-02 | 157.31        | Cross                        |
| c1-09 | c3-06 | 165.55        | Cross                        |
| c1-09 | c0-07 | 169.47        | Cross                        |
| c1-09 | c3-09 | 181.83        | Cross                        |
| c1-09 | c0-08 | 182.35        | Cross                        |

**Finding**: At μ=0.5, the top-π node `c2-02` suggests **exclusively cross-community** links (5
of 5, 100%). `c0-07` (rank-2 π) splits 2 cross / 3 within (40% cross). `c1-09` (lowest π) shows
**exclusively cross-community** links (5 of 5, 100%) — at the detection boundary, even the
lowest-π node's best available non-edges are all outside its own community, since μ=0.5 means
half its existing edges are already inter-community, leaving few good within-community options
remaining.

### Stationary distribution spread — across parameter sets

| Parameter set | Max π    | Min π    | Ratio | Label                     |
| ------------- | -------- | -------- | ----- | ---------------------------- |
| μ=0.1         | 0.076087 | 0.005435 | 14.0  | Hub-dominated                |
| μ=0.3         | 0.079208 | 0.009901 | 8.0   | Hub-dominated                |
| μ=0.5         | 0.079646 | 0.008850 | 9.0   | Hub-dominated                |

**The π distribution does not become flat within this tested range.** All three ratios remain
firmly in the 5–20 hub-dominated tier; none approaches 1.0 or even the mildly-skewed 2–5 band.
This is the graph's power-law degree distribution asserting itself — regardless of μ, a small
number of nodes always draw disproportionately high degree by construction (`τ1=2.5` produces
a heavy-tailed degree distribution), so π stays hub-dominated across the entire μ sweep tested
here.

### Visual observations

> *The following description is inferred from the graph's structural properties and stationary
> distribution. The force-directed layout can be verified by opening `/tmp/nx-lfr-01.html`,
> `/tmp/nx-lfr-03.html`, and `/tmp/nx-lfr-05.html`.*

At μ=0.1, the layout should show 4 visually distinct clusters (community sizes 16/14/10/10)
with only a light scattering of cross-cluster edges (22 of 184, 12.0%) — `c0-03` (the
highest-degree node, degree 14 after excluding its self-loop) should appear as the visually
largest node, positioned centrally within its 16-member community. At μ=0.3, the clusters
should remain identifiable but noticeably more interconnected (86 of 202 edges, 42.6%, cross
boundaries), with `c1-03` (new highest-π node) prominent. At μ=0.5, cross-community edges
dominate the edge count (144 of 226, 63.7%) — the four communities should be substantially
blended in the force-directed layout, though the underlying power-law degree distribution
still produces one visually dominant node (`c2-02`, degree 18) regardless of how blended the
community structure appears.

### Markov questions — answered

- **Does π rank match the expected degree-hub nodes (per the power-law degree distribution)?**
  Yes, exactly, at every μ tested. μ=0.1: `c0-03` leads π (0.076087) and has the highest
  retained degree (14) in the graph. μ=0.3: `c1-03` leads π (0.079208) with retained degree 16,
  the highest in that set. μ=0.5: `c2-02` leads π (0.079646) with retained degree 18, the
  highest in that set. π rank tracks the power-law-distributed degree hub precisely at every μ.

- **What fraction of `suggest` recommendations are inter-community at each μ?**
  For the top-π node specifically: μ=0.1 `c0-03` is 3/5 (60%) cross; μ=0.3 `c1-03` is 5/5
  (100%) cross; μ=0.5 `c2-02` is 5/5 (100%) cross. The fraction rises from 60% to 100% and
  plateaus at 100% for the top node from μ=0.3 onward.

- **Is π hub-dominated or flat, and how does the max/min ratio trend across the μ-sweep?**
  Hub-dominated at every μ tested (ratios 14.0, 8.0, 9.0 — all within the 5–20 tier). The ratio
  does not trend monotonically with μ: it is highest at μ=0.1 and lower (though still
  hub-dominated) at μ=0.3 and μ=0.5, because π skew here is driven by the power-law degree
  distribution's inherent variability rather than by the mixing parameter.

- **At what μ does the π distribution become effectively flat (max/min ratio approaches 1),
  indicating the community signal has washed out?**
  **Never, within the tested range (μ=0.1 to μ=0.5).** The max/min ratio stays in the
  hub-dominated tier (8.0–14.0) throughout the entire sweep and does not approach 1.0 at any
  tested μ. This is a genuine finding, not an omission: unlike the caveman and planted-partition
  graphs (which have uniform or near-uniform expected degree per node, so π skew is driven
  purely by community mixing), LFR's power-law degree distribution guarantees persistent π skew
  regardless of μ — the community-mixing signal (tracked instead by the rising cross-community
  edge/suggestion fractions above) and the π-skew signal are structurally decoupled in this
  graph family, more so than in any other graph in this documentation series.

---

## References

- Lancichinetti, A., Fortunato, S., & Radicchi, F. (2008). Benchmark graphs for testing
  community detection algorithms. *Physical Review E*, 78, 046110.
  https://doi.org/10.1103/PhysRevE.78.046110
- Wikipedia: Lancichinetti–Fortunato–Radicchi benchmark —
  https://en.wikipedia.org/wiki/Lancichinetti%E2%80%93Fortunato%E2%80%93Radicchi_benchmark
- NetworkX: https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.LFR_benchmark_graph.html
- LFR section from graph-models research notes

---

## Definition of Done

- [x] All sections above filled — no placeholder text remaining
- [x] Slug ordering algorithm documented with a node-id mapping table (not just the pattern) —
      including the set-based community canonicalisation specific to this graph
- [x] `wikigraph analyze --suggest-top 5` output pasted verbatim for all three parameter sets
- [x] Cross-community edge count computed and recorded (Python helper, per parameter set)
- [x] All Markov questions answered with actual numbers from analyze output
- [x] Three parameter-set variants shown (μ=0.1, 0.3, 0.5) with a comparison table before the
      per-set detail sections; no generation failures occurred, so no `--average-degree 5`
      fallback was needed (documented explicitly above)
- [x] Every reference link manually verified to resolve (Physical Review E DOI 302→200;
      Wikipedia 200; NetworkX docs 200)
- [x] File committed to branch — path to be updated once #42 resolves
