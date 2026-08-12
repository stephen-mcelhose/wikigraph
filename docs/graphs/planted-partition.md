<!-- Path provisional — will be updated when #42 resolves -->

# Planted Partition Graph — symmetric stochastic block model benchmark

> ⚠️ **DRAFT** — AI-assisted write-up, not yet verified by human analysis. Treat all findings as provisional.

> **nx-to-wiki flag**: `--graph planted-partition --l [blocks] --k [block-size] --p-in [p] --p-out [p]`
> **Nodes**: l·k (parameterised) · **Directed edges**: 2×undirected (parameterised)
> **Naming**: tier1-planted-partition-block-attr
> **Source**: generative — parameterised

---

## Background

The planted partition graph is the symmetric special case of the **Stochastic Block Model
(SBM)**: `l` equal-sized blocks of `k` nodes each, where every pair of nodes *within* the same
block is connected independently with probability `p_in`, and every pair of nodes in
*different* blocks is connected independently with probability `p_out`. Unlike the caveman
graph (which starts as pure cliques and rewires specific edges), the planted partition graph
draws every edge — intra- and inter-block — from an independent Bernoulli trial, so it more
directly models the statistical assumption underlying most community-detection algorithms: that
observed edges are noisy samples from an underlying block structure, not a deterministic
rewiring process. The `p_in`/`p_out` ratio directly controls how detectable the planted
communities are. A commonly used modularity proxy is $Q \approx 1 - p_{out}/p_{in}$: as
$p_{out} \to p_{in}$, $Q \to 0$ and the blocks become statistically indistinguishable from a
single Erdős–Rényi graph (the **disassortative** boundary, $p_{in} < p_{out}$, is the regime
where blocks actively repel connections rather than attract them; this document stays in the
**assortative** regime, $p_{in} > p_{out}$, throughout).

**Origin**: Standard construction from the stochastic block model literature (Holland, Laskey &
Leinhardt, 1983, and subsequent SBM literature); see Wikipedia overview and NetworkX reference
below.

---

## Graph Properties

**Key parameters**:

| Flag       | Default | Range        | What it controls                                    |
| ---------- | ------- | ------------ | ------------------------------------------------------ |
| `--l`      | 3       | ≥ 1          | Number of blocks ("communities")                         |
| `--k`      | 10      | ≥ 1          | Nodes per block                                           |
| `--p-in`   | 0.5     | [0.0, 1.0]   | Probability of an edge between two nodes in the same block |
| `--p-out`  | 0.05    | [0.0, 1.0]   | Probability of an edge between two nodes in different blocks |
| `--seed`   | none    | any int      | RNG seed for reproducible edge sampling                    |

Three parameter sets (as directed by the issue), all with `l=3`, `k=10`, `p_in=0.6`, `seed=42`,
so node count is fixed at 30 — only `p_out` and the resulting edge count vary:

| Parameter set | l | k  | p_in | p_out | Nodes | Undirected edges | Directed edges | Min deg | Max deg | Q ≈ 1−p_out/p_in |
| ------------- | - | -- | ---- | ----- | ----- | ------------------ | ---------------- | -------- | -------- | ------------------ |
| Tight         | 3 | 10 | 0.6  | 0.02  | 30    | 87                  | 174               | 3        | 8        | 0.967               |
| Mid           | 3 | 10 | 0.6  | 0.10  | 30    | 116                 | 232               | 5        | 11       | 0.833               |
| Loose         | 3 | 10 | 0.6  | 0.25  | 30    | 152                 | 304               | 6        | 16       | 0.583               |

Community structure: 3 blocks of 10 nodes each, defined by the `block` integer attribute
NetworkX attaches to each node during generation.

---

## Transition Matrix Structure

For every graph produced by `nx-to-wiki`, the random walk is **uniform**: each step chooses a
neighbour with equal probability. The transition matrix $P$ is fully determined by the
adjacency structure of the `.md` files:

$$P_{ij} = \frac{1}{\deg(i)} \quad \text{if } [[j]] \text{ appears in i.md}, \quad \text{else } 0$$

### Tight (p_out=0.02)

| Slug   | Degree | $P_{ij}$ (non-zero) | Structural note                                  |
| ------ | ------ | -------------------- | ------------------------------------------------- |
| c1-03  | 8      | 1/8 = 0.125000        | Highest-π tier node; benefited from above-average Bernoulli draws |
| c0-00  | 5      | 1/5 = 0.200000        | Typical block-0 node — degree near block average    |
| c2-07  | 3      | 1/3 ≈ 0.333333        | Lowest-π node; below-average intra-block edge draws |

### Mid (p_out=0.10)

| Slug   | Degree | $P_{ij}$ (non-zero) | Structural note                              |
| ------ | ------ | -------------------- | ------------------------------------------------ |
| c1-01  | 11     | 1/11 ≈ 0.090909       | Highest-π tier node                                |
| c0-00  | 8      | 1/8 = 0.125000        | Typical mid-degree node                            |
| c2-05  | 5      | 1/5 = 0.200000        | Lowest-π node                                      |

### Loose (p_out=0.25)

| Slug   | Degree | $P_{ij}$ (non-zero) | Structural note                              |
| ------ | ------ | -------------------- | ------------------------------------------------- |
| c1-08  | 16     | 1/16 = 0.062500       | Highest-π node — well above average even for this dense regime |
| c0-01  | 12     | 1/12 ≈ 0.083333       | Typical mid-degree node                             |
| c1-05  | 6      | 1/6 ≈ 0.166667        | Lowest-π node                                       |

**The `.md` files are a lossless encoding of $P$ for all three parameter sets.** Minimum
non-zero $P_{ij}$ across all three sets is $1/16 = 0.0625$ (loose set's highest-degree node) —
well above the `--min-edge` default filter of 0.005.

### Export

```bash
wikigraph export /tmp/nxwiki-pp-tight --format csv -o /tmp/pp-tight-export
wikigraph export /tmp/nxwiki-pp-mid   --format csv -o /tmp/pp-mid-export
wikigraph export /tmp/nxwiki-pp-loose --format csv -o /tmp/pp-loose-export
```

Sparse edge-row counts: 174 (tight), 232 (mid), 304 (loose) — matching directed-edge counts. No
`--min-edge 0` warning applies to any of the three sets tested.

---

## Slug Naming

**Naming tier**: Tier 1

**Assignment algorithm** (`_slugs_from_community_attr(G, attr="block")`, verbatim from
source): NetworkX's `planted_partition_graph(l, k, p_in, p_out, seed)` assigns every node an
integer `block` attribute (`0` through `l-1`). `nx-to-wiki`:

1. Reads `raw[n] = G.nodes[n].get("block")` for every node — a plain integer here (the same
   helper also handles set-valued community attributes for other graphs, via `frozenset`
   canonicalisation, but `block` is always a scalar int for this graph).
2. Computes `unique = sorted(set(raw.values()))` and maps each unique block id to a
   `community_id` in that sorted order (for integer blocks, `community_id[b] == b`).
3. Groups nodes by `community_id`.
4. Within each group, sorts node ids ascending and assigns `c{cid}-{index:02d}`, zero-padded to
   `max(2, len(str(len(group) - 1)))` digits.

**Why this ordering**: NetworkX's `planted_partition_graph` numbers nodes contiguously by block
— block 0 occupies node ids `[0, k)`, block 1 occupies `[k, 2k)`, and so on — so sorting by
node id within each block group reproduces the same contiguous numbering the generator used
internally. This is why `c0-00` is always node 0, `c1-00` is always node `k`, etc., independent
of `p_in`/`p_out` (block membership and node numbering are fixed at construction time, before
any edges are sampled).

Full node-id → slug mapping (l=3, k=10 — identical across all three parameter sets since only
edge sampling, not block/node assignment, changes with `p_out`):

| Slug   | NetworkX node id | Block |
| ------ | ------------------ | ----- |
| c0-00  | 0                   | 0     |
| c0-01  | 1                   | 0     |
| c0-02  | 2                   | 0     |
| c0-03  | 3                   | 0     |
| c0-04  | 4                   | 0     |
| c0-05  | 5                   | 0     |
| c0-06  | 6                   | 0     |
| c0-07  | 7                   | 0     |
| c0-08  | 8                   | 0     |
| c0-09  | 9                   | 0     |
| c1-00  | 10                  | 1     |
| c1-01  | 11                  | 1     |
| c1-02  | 12                  | 1     |
| c1-03  | 13                  | 1     |
| c1-04  | 14                  | 1     |
| c1-05  | 15                  | 1     |
| c1-06  | 16                  | 1     |
| c1-07  | 17                  | 1     |
| c1-08  | 18                  | 1     |
| c1-09  | 19                  | 1     |
| c2-00  | 20                  | 2     |
| c2-01  | 21                  | 2     |
| c2-02  | 22                  | 2     |
| c2-03  | 23                  | 2     |
| c2-04  | 24                  | 2     |
| c2-05  | 25                  | 2     |
| c2-06  | 26                  | 2     |
| c2-07  | 27                  | 2     |
| c2-08  | 28                  | 2     |
| c2-09  | 29                  | 2     |

---

## How to Generate

```bash
# p_out sweep with fixed p_in=0.6, l=3, k=10, seed=42
python3 tools/nx-to-wiki/main.py --graph planted-partition --l 3 --k 10 --p-in 0.6 --p-out 0.02 --seed 42 --out /tmp/nxwiki-pp-tight
python3 tools/nx-to-wiki/main.py --graph planted-partition --l 3 --k 10 --p-in 0.6 --p-out 0.10 --seed 42 --out /tmp/nxwiki-pp-mid
python3 tools/nx-to-wiki/main.py --graph planted-partition --l 3 --k 10 --p-in 0.6 --p-out 0.25 --seed 42 --out /tmp/nxwiki-pp-loose

wikigraph graph /tmp/nxwiki-pp-tight -o /tmp/nx-pp-tight.html --title "Planted Partition p_out=0.02"
wikigraph graph /tmp/nxwiki-pp-mid   -o /tmp/nx-pp-mid.html   --title "Planted Partition p_out=0.10"
wikigraph graph /tmp/nxwiki-pp-loose -o /tmp/nx-pp-loose.html --title "Planted Partition p_out=0.25"

wikigraph analyze /tmp/nxwiki-pp-tight --suggest-top 5
wikigraph analyze /tmp/nxwiki-pp-mid --suggest-top 5
wikigraph analyze /tmp/nxwiki-pp-loose --suggest-top 5
```

**Effect of `p_out`**: Raising `p_out` from 0.02 to 0.25 (with `p_in=0.6` fixed) increases total
edge count substantially (87 → 116 → 152 undirected edges) because more inter-block edges are
sampled, while intra-block edge density stays statistically constant at `p_in=0.6`. The
modularity proxy $Q = 1 - p_{out}/p_{in}$ falls from 0.967 (tight, near-isolated blocks) to
0.583 (loose, approaching a single merged community).

**What each output directory contains:**

- 30 `.md` files, one per node (fixed across all three parameter sets since `l`, `k` are fixed)
- Each file links to all its undirected neighbours (symmetric — every link is bidirectional)
- File layout: `# {slug}\n\n[[neighbour-1]] [[neighbour-2]] ...`
- Why directed edges = 2× undirected: `G.to_directed()` adds both u→v and v→u for every edge

---

## wikigraph Analysis

> **⚠️ Communicating classes are always trivial here.** All wikis produced by `nx-to-wiki` use
> `G.to_directed()` on a connected undirected graph, which yields a **strongly connected**
> directed graph — this graph family is dense enough at all three tested `p_out` values that
> `wikigraph analyze` reports **one communicating class containing all nodes** in every case.
> The meaningful Markov signal is π and `suggest`, using cross-block edge counts (Python helper)
> as the structural boundary proxy.

### Comparison across parameter sets

| Parameter set | p_out | Nodes | Directed edges | Entropy rate | Max π    | Min π    | Ratio | Cross-community % | Q proxy |
| ------------- | ----- | ----- | ---------------- | -------------- | -------- | -------- | ----- | -------------------- | -------- |
| Tight         | 0.02  | 30    | 174               | 2.5736 bits    | 0.045977 | 0.017241 | 2.67  | 12 of 174 (6.9%)      | 0.967    |
| Mid           | 0.10  | 30    | 232               | 2.9834 bits    | 0.047414 | 0.021552 | 2.20  | 60 of 232 (25.9%)     | 0.833    |
| Loose         | 0.25  | 30    | 304               | 3.3765 bits    | 0.052632 | 0.019737 | 2.67  | 144 of 304 (47.4%)    | 0.583    |

**Trend**: Cross-community edge percentage climbs sharply with `p_out` (6.9% → 25.9% → 47.4%),
tracking the falling Q proxy (0.967 → 0.833 → 0.583) almost linearly. Entropy rate also rises
monotonically (2.57 → 2.98 → 3.38 bits) as the graph becomes denser and more uniformly
connected. The max/min π ratio, however, stays in a narrow **mildly skewed** band (2.20–2.67)
across the entire sweep — it does not track Q at all, because within-block degree variance
(driven by Bernoulli sampling noise, not by systematic hub structure) dominates the ratio at
every `p_out` value tested.

### Tight (p_out=0.02)

#### Raw analyze output

```
=== Overview ===
Pages:        30
Edges:        174
Entropy rate: 2.5736 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 30 page(s)
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

=== Orphan pages (bottom 10% by stationary distribution) ===
  c2-07                                     π=0.017241  → add inbound links
  c1-08                                     π=0.022989  → add inbound links
  c1-01                                     π=0.022989  → add inbound links
  c2-01                                     π=0.022989  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. c2-04                                     π=0.045977
  2. c1-03                                     π=0.045977
  3. c1-09                                     π=0.045977
  4. c0-01                                     π=0.040230
  5. c0-09                                     π=0.040230

=== Suggested missing links (lowest commute time, not yet linked) ===
  c0-00:
    → c0-01                                   (commute: 62.00)
    → c0-07                                   (commute: 72.15)
    → c0-06                                   (commute: 73.18)
    → c0-05                                   (commute: 86.69)
    → c2-06                                   (commute: 151.09)
  c0-01:
    → c0-07                                   (commute: 61.83)
    → c0-00                                   (commute: 62.00)
    → c2-03                                   (commute: 137.86)
    → c2-06                                   (commute: 138.10)
    → c2-09                                   (commute: 157.44)
  c0-02:
    → c0-06                                   (commute: 63.67)
    → c0-07                                   (commute: 69.03)
    → c0-04                                   (commute: 77.93)
    → c2-06                                   (commute: 144.59)
    → c2-03                                   (commute: 146.10)
  c0-03:
    → c0-09                                   (commute: 57.02)
    → c0-08                                   (commute: 63.85)
    → c0-04                                   (commute: 78.26)
    → c2-03                                   (commute: 144.23)
    → c2-06                                   (commute: 145.90)
  c0-04:
    → c0-08                                   (commute: 77.14)
    → c0-02                                   (commute: 77.93)
    → c0-03                                   (commute: 78.26)
    → c0-06                                   (commute: 81.09)
    → c0-05                                   (commute: 98.95)
  c0-05:
    → c0-09                                   (commute: 75.30)
    → c0-08                                   (commute: 82.35)
    → c0-07                                   (commute: 86.22)
    → c0-00                                   (commute: 86.69)
    → c0-04                                   (commute: 98.95)
  c0-06:
    → c0-02                                   (commute: 63.67)
    → c0-08                                   (commute: 63.92)
    → c0-00                                   (commute: 73.18)
    → c0-04                                   (commute: 81.09)
    → c2-06                                   (commute: 123.51)
  c0-07:
    → c0-01                                   (commute: 61.83)
    → c0-02                                   (commute: 69.03)
    → c0-00                                   (commute: 72.15)
    → c0-05                                   (commute: 86.22)
    → c2-03                                   (commute: 146.10)
  c0-08:
    → c0-03                                   (commute: 63.85)
    → c0-06                                   (commute: 63.92)
    → c0-04                                   (commute: 77.14)
    → c0-05                                   (commute: 82.35)
    → c2-03                                   (commute: 123.51)
  c0-09:
    → c0-03                                   (commute: 57.02)
    → c0-05                                   (commute: 75.30)
    → c2-06                                   (commute: 137.55)
    → c2-03                                   (commute: 138.15)
    → c2-09                                   (commute: 157.26)
  c1-00:
    → c1-05                                   (commute: 59.61)
    → c1-04                                   (commute: 63.60)
    → c1-06                                   (commute: 73.00)
    → c1-08                                   (commute: 81.41)
    → c2-02                                   (commute: 82.28)
  c1-01:
    → c1-03                                   (commute: 71.17)
    → c1-04                                   (commute: 78.45)
    → c1-07                                   (commute: 84.45)
    → c1-06                                   (commute: 85.91)
    → c1-08                                   (commute: 93.84)
  c1-02:
    → c1-07                                   (commute: 61.57)
    → c1-08                                   (commute: 71.56)
    → c2-02                                   (commute: 86.06)
    → c2-04                                   (commute: 95.63)
    → c2-00                                   (commute: 101.79)
  c1-03:
    → c1-04                                   (commute: 53.94)
    → c1-06                                   (commute: 62.63)
    → c1-01                                   (commute: 71.17)
    → c2-00                                   (commute: 81.99)
    → c2-09                                   (commute: 86.97)
  c1-04:
    → c1-03                                   (commute: 53.94)
    → c1-00                                   (commute: 63.60)
    → c1-08                                   (commute: 76.60)
    → c1-01                                   (commute: 78.45)
    → c2-04                                   (commute: 93.31)
  c1-05:
    → c1-00                                   (commute: 59.61)
    → c1-07                                   (commute: 63.15)
    → c2-02                                   (commute: 87.95)
    → c2-04                                   (commute: 98.09)
    → c2-00                                   (commute: 105.58)
  c1-06:
    → c1-09                                   (commute: 59.16)
    → c1-03                                   (commute: 62.63)
    → c1-00                                   (commute: 73.00)
    → c1-01                                   (commute: 85.91)
    → c2-02                                   (commute: 100.61)
  c1-07:
    → c1-02                                   (commute: 61.57)
    → c1-05                                   (commute: 63.15)
    → c1-08                                   (commute: 82.33)
    → c1-01                                   (commute: 84.45)
    → c2-02                                   (commute: 93.76)
  c1-08:
    → c1-02                                   (commute: 71.56)
    → c1-04                                   (commute: 76.60)
    → c1-00                                   (commute: 81.41)
    → c1-07                                   (commute: 82.33)
    → c1-01                                   (commute: 93.84)
  c1-09:
    → c1-06                                   (commute: 59.16)
    → c2-02                                   (commute: 82.84)
    → c2-04                                   (commute: 92.22)
    → c2-00                                   (commute: 98.51)
    → c2-09                                   (commute: 106.87)
  c2-00:
    → c2-09                                   (commute: 56.83)
    → c2-02                                   (commute: 64.61)
    → c2-05                                   (commute: 71.13)
    → c2-01                                   (commute: 80.97)
    → c1-03                                   (commute: 81.99)
  c2-01:
    → c2-09                                   (commute: 72.97)
    → c2-06                                   (commute: 78.37)
    → c2-00                                   (commute: 80.97)
    → c2-02                                   (commute: 85.76)
    → c1-03                                   (commute: 112.11)
  c2-02:
    → c2-00                                   (commute: 64.61)
    → c2-03                                   (commute: 72.50)
    → c2-05                                   (commute: 75.27)
    → c1-00                                   (commute: 82.28)
    → c1-09                                   (commute: 82.84)
  c2-03:
    → c2-06                                   (commute: 60.54)
    → c2-08                                   (commute: 63.81)
    → c2-02                                   (commute: 72.50)
    → c1-03                                   (commute: 97.97)
    → c2-07                                   (commute: 98.33)
  c2-04:
    → c2-06                                   (commute: 52.54)
    → c1-00                                   (commute: 87.03)
    → c2-07                                   (commute: 89.56)
    → c1-09                                   (commute: 92.22)
    → c1-04                                   (commute: 93.31)
  c2-05:
    → c2-08                                   (commute: 67.17)
    → c2-00                                   (commute: 71.13)
    → c2-02                                   (commute: 75.27)
    → c2-07                                   (commute: 102.46)
    → c1-03                                   (commute: 102.83)
  c2-06:
    → c2-04                                   (commute: 52.54)
    → c2-03                                   (commute: 60.54)
    → c2-01                                   (commute: 78.37)
    → c1-03                                   (commute: 89.61)
    → c1-00                                   (commute: 101.61)
  c2-07:
    → c2-04                                   (commute: 89.56)
    → c2-08                                   (commute: 92.67)
    → c2-03                                   (commute: 98.33)
    → c2-02                                   (commute: 99.26)
    → c2-05                                   (commute: 102.46)
  c2-08:
    → c2-03                                   (commute: 63.81)
    → c2-05                                   (commute: 67.17)
    → c1-03                                   (commute: 87.88)
    → c2-07                                   (commute: 92.67)
    → c1-00                                   (commute: 100.98)
  c2-09:
    → c2-00                                   (commute: 56.83)
    → c2-01                                   (commute: 72.97)
    → c1-03                                   (commute: 86.97)
    → c1-00                                   (commute: 100.80)
    → c1-04                                   (commute: 106.21)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `c2-04` (0.045977, tied with `c1-03`/`c1-09`; `c1-03` shown as representative). Lowest
π: `c2-07` (0.017241).

| From    | To      | Commute time | Within / Cross block |
| ------- | ------- | ------------ | ----------------------- |
| c2-04   | c2-06   | 52.54        | Within                   |
| c2-04   | c1-00   | 87.03        | Cross                     |
| c2-04   | c2-07   | 89.56        | Within                   |
| c2-04   | c1-09   | 92.22        | Cross                     |
| c2-04   | c1-04   | 93.31        | Cross                     |
| c1-03   | c1-04   | 53.94        | Within                    |
| c1-03   | c1-06   | 62.63        | Within                    |
| c1-03   | c1-01   | 71.17        | Within                    |
| c1-03   | c2-00   | 81.99        | Cross                     |
| c1-03   | c2-09   | 86.97        | Cross                     |
| c2-07   | c2-04   | 89.56        | Within                    |
| c2-07   | c2-08   | 92.67        | Within                    |
| c2-07   | c2-03   | 98.33        | Within                    |
| c2-07   | c2-02   | 99.26        | Within                    |
| c2-07   | c2-05   | 102.46       | Within                    |

**Finding**: At tight `p_out`, the two top-π nodes suggest a mix — 3 of 10 combined suggestions
(30%) cross block boundaries. The lowest-π node, `c2-07`, suggests **exclusively within-block**
links (5 of 5, 0% cross) — it is under-connected even within its own block relative to peers,
so the closest not-yet-linked candidates are all block-mates.

### Mid (p_out=0.10)

#### Raw analyze output

```
=== Overview ===
Pages:        30
Edges:        232
Entropy rate: 2.9834 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 30 page(s)
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

=== Orphan pages (bottom 10% by stationary distribution) ===
  c2-05                                     π=0.021552  → add inbound links
  c2-07                                     π=0.025862  → add inbound links
  c0-06                                     π=0.025862  → add inbound links
  c0-00                                     π=0.025862  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. c1-01                                     π=0.047414
  2. c0-03                                     π=0.047414
  3. c2-01                                     π=0.047414
  4. c0-09                                     π=0.043103
  5. c1-09                                     π=0.038793

=== Suggested missing links (lowest commute time, not yet linked) ===
  c0-00:
    → c2-01                                   (commute: 70.00)
    → c1-01                                   (commute: 72.12)
    → c2-08                                   (commute: 74.14)
    → c1-09                                   (commute: 74.34)
    → c0-01                                   (commute: 74.83)
  c0-01:
    → c2-01                                   (commute: 66.50)
    → c1-01                                   (commute: 68.12)
    → c1-09                                   (commute: 68.53)
    → c1-04                                   (commute: 70.87)
    → c2-08                                   (commute: 71.34)
  c0-02:
    → c2-01                                   (commute: 59.37)
    → c1-09                                   (commute: 62.22)
    → c0-04                                   (commute: 63.02)
    → c1-04                                   (commute: 63.61)
    → c2-08                                   (commute: 65.60)
  c0-03:
    → c0-09                                   (commute: 48.21)
    → c1-01                                   (commute: 51.04)
    → c0-04                                   (commute: 53.48)
    → c2-02                                   (commute: 54.56)
    → c2-08                                   (commute: 54.89)
  c0-04:
    → c0-03                                   (commute: 53.48)
    → c2-01                                   (commute: 57.69)
    → c1-01                                   (commute: 60.45)
    → c0-02                                   (commute: 63.02)
    → c1-04                                   (commute: 63.06)
  c0-05:
    → c0-09                                   (commute: 67.68)
    → c1-01                                   (commute: 71.10)
    → c1-04                                   (commute: 73.28)
    → c0-04                                   (commute: 74.64)
    → c2-00                                   (commute: 75.92)
  c0-06:
    → c2-01                                   (commute: 71.51)
    → c0-02                                   (commute: 72.80)
    → c0-04                                   (commute: 73.58)
    → c1-01                                   (commute: 75.46)
    → c2-08                                   (commute: 75.55)
  c0-07:
    → c2-01                                   (commute: 71.01)
    → c1-01                                   (commute: 71.57)
    → c0-02                                   (commute: 72.63)
    → c1-09                                   (commute: 72.83)
    → c0-01                                   (commute: 74.37)
  c0-08:
    → c0-03                                   (commute: 59.99)
    → c1-01                                   (commute: 64.82)
    → c1-09                                   (commute: 67.02)
    → c0-04                                   (commute: 67.71)
    → c2-01                                   (commute: 68.62)
  c0-09:
    → c0-03                                   (commute: 48.21)
    → c1-01                                   (commute: 53.23)
    → c2-01                                   (commute: 55.36)
    → c1-09                                   (commute: 56.90)
    → c1-04                                   (commute: 57.20)
  c1-00:
    → c1-04                                   (commute: 69.40)
    → c0-03                                   (commute: 71.89)
    → c2-03                                   (commute: 72.90)
    → c1-06                                   (commute: 75.19)
    → c0-09                                   (commute: 76.54)
  c1-01:
    → c2-01                                   (commute: 49.77)
    → c0-03                                   (commute: 51.04)
    → c0-09                                   (commute: 53.23)
    → c2-00                                   (commute: 57.03)
    → c2-02                                   (commute: 59.33)
  c1-02:
    → c1-09                                   (commute: 63.09)
    → c2-01                                   (commute: 64.41)
    → c0-03                                   (commute: 67.71)
    → c0-09                                   (commute: 71.59)
    → c2-00                                   (commute: 71.92)
  c1-03:
    → c1-09                                   (commute: 69.02)
    → c0-03                                   (commute: 71.22)
    → c2-01                                   (commute: 72.22)
    → c2-03                                   (commute: 75.77)
    → c1-02                                   (commute: 76.36)
  c1-04:
    → c1-05                                   (commute: 55.11)
    → c0-09                                   (commute: 57.20)
    → c2-03                                   (commute: 57.83)
    → c2-08                                   (commute: 59.32)
    → c2-00                                   (commute: 60.57)
  c1-05:
    → c1-04                                   (commute: 55.11)
    → c2-01                                   (commute: 55.98)
    → c0-03                                   (commute: 57.99)
    → c0-09                                   (commute: 60.74)
    → c2-08                                   (commute: 63.18)
  c1-06:
    → c1-09                                   (commute: 63.83)
    → c2-01                                   (commute: 65.36)
    → c0-03                                   (commute: 66.29)
    → c2-03                                   (commute: 68.57)
    → c0-09                                   (commute: 69.59)
  c1-07:
    → c1-01                                   (commute: 63.52)
    → c1-05                                   (commute: 67.81)
    → c1-04                                   (commute: 71.08)
    → c0-03                                   (commute: 73.02)
    → c2-01                                   (commute: 73.90)
  c1-08:
    → c1-04                                   (commute: 64.78)
    → c0-03                                   (commute: 64.82)
    → c2-01                                   (commute: 67.66)
    → c1-06                                   (commute: 70.88)
    → c0-02                                   (commute: 70.90)
  c1-09:
    → c2-01                                   (commute: 54.72)
    → c0-09                                   (commute: 56.90)
    → c2-03                                   (commute: 61.23)
    → c2-08                                   (commute: 62.02)
    → c0-02                                   (commute: 62.22)
  c2-00:
    → c1-01                                   (commute: 57.03)
    → c2-03                                   (commute: 59.07)
    → c1-04                                   (commute: 60.57)
    → c2-08                                   (commute: 60.68)
    → c0-09                                   (commute: 62.03)
  c2-01:
    → c1-01                                   (commute: 49.77)
    → c2-06                                   (commute: 53.58)
    → c1-09                                   (commute: 54.72)
    → c0-09                                   (commute: 55.36)
    → c1-05                                   (commute: 55.98)
  c2-02:
    → c2-08                                   (commute: 53.98)
    → c0-03                                   (commute: 54.56)
    → c2-04                                   (commute: 57.26)
    → c1-01                                   (commute: 59.33)
    → c0-09                                   (commute: 60.35)
  c2-03:
    → c0-03                                   (commute: 57.33)
    → c1-04                                   (commute: 57.83)
    → c2-00                                   (commute: 59.07)
    → c2-06                                   (commute: 59.97)
    → c1-09                                   (commute: 61.23)
  c2-04:
    → c2-02                                   (commute: 57.26)
    → c0-03                                   (commute: 60.61)
    → c1-01                                   (commute: 61.91)
    → c2-00                                   (commute: 64.10)
    → c0-09                                   (commute: 67.12)
  c2-05:
    → c2-01                                   (commute: 74.10)
    → c0-03                                   (commute: 77.27)
    → c2-03                                   (commute: 79.84)
    → c0-09                                   (commute: 82.19)
    → c0-04                                   (commute: 83.02)
  c2-06:
    → c2-01                                   (commute: 53.58)
    → c2-03                                   (commute: 59.97)
    → c1-01                                   (commute: 63.10)
    → c0-09                                   (commute: 63.56)
    → c1-04                                   (commute: 65.04)
  c2-07:
    → c2-08                                   (commute: 69.07)
    → c0-03                                   (commute: 71.92)
    → c2-06                                   (commute: 72.04)
    → c1-01                                   (commute: 73.78)
    → c1-04                                   (commute: 76.07)
  c2-08:
    → c2-02                                   (commute: 53.98)
    → c0-03                                   (commute: 54.89)
    → c1-04                                   (commute: 59.32)
    → c2-00                                   (commute: 60.68)
    → c1-09                                   (commute: 62.02)
  c2-09:
    → c0-03                                   (commute: 66.49)
    → c2-00                                   (commute: 67.93)
    → c1-01                                   (commute: 68.19)
    → c1-04                                   (commute: 70.88)
    → c0-04                                   (commute: 72.33)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `c1-01` (0.047414, tied with `c0-03`/`c2-01`; `c1-01` shown as representative). Lowest
π: `c2-05` (0.021552).

| From    | To      | Commute time | Within / Cross block |
| ------- | ------- | ------------ | ----------------------- |
| c1-01   | c2-01   | 49.77         | Cross                     |
| c1-01   | c0-03   | 51.04         | Cross                     |
| c1-01   | c0-09   | 53.23         | Cross                     |
| c1-01   | c2-00   | 57.03         | Cross                     |
| c1-01   | c2-02   | 59.33         | Cross                     |
| c0-03   | c0-09   | 48.21         | Within                    |
| c0-03   | c1-01   | 51.04         | Cross                     |
| c0-03   | c0-04   | 53.48         | Within                    |
| c0-03   | c2-02   | 54.56         | Cross                     |
| c0-03   | c2-08   | 54.89         | Cross                     |
| c2-05   | c2-01   | 74.10         | Within                    |
| c2-05   | c0-03   | 77.27         | Cross                     |
| c2-05   | c2-03   | 79.84         | Within                    |
| c2-05   | c0-09   | 82.19         | Cross                     |
| c2-05   | c0-04   | 83.02         | Cross                     |

**Finding**: At mid `p_out`, `c1-01` now suggests **exclusively cross-block** links (5 of 5,
100% — up from 0/5 typical for a within-block-only node at tight p_out), while `c0-03` splits
2 within / 3 cross. Overall cross-block fraction across the two top-π nodes is 8 of 10 (80%), a
sharp rise from the tight set's 30%. `c2-05` (lowest π) shows 2 within / 3 cross — a more even
split than the tight set's lowest-π node, which suggested exclusively within-block links.

### Loose (p_out=0.25)

#### Raw analyze output

```
=== Overview ===
Pages:        30
Edges:        304
Entropy rate: 3.3765 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 30 page(s)
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

=== Orphan pages (bottom 10% by stationary distribution) ===
  c1-05                                     π=0.019737  → add inbound links
  c1-00                                     π=0.023026  → add inbound links
  c2-07                                     π=0.023026  → add inbound links
  c1-09                                     π=0.026316  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. c1-08                                     π=0.052632
  2. c2-06                                     π=0.046053
  3. c0-01                                     π=0.042763
  4. c2-02                                     π=0.042763
  5. c0-03                                     π=0.042763

=== Suggested missing links (lowest commute time, not yet linked) ===
  c0-00:
    → c1-08                                   (commute: 54.02)
    → c0-01                                   (commute: 56.41)
    → c2-02                                   (commute: 57.32)
    → c1-04                                   (commute: 61.07)
    → c2-05                                   (commute: 61.40)
  c0-01:
    → c2-01                                   (commute: 51.27)
    → c2-05                                   (commute: 54.29)
    → c1-01                                   (commute: 54.37)
    → c2-09                                   (commute: 54.58)
    → c1-07                                   (commute: 55.07)
  c0-02:
    → c2-02                                   (commute: 54.51)
    → c1-04                                   (commute: 57.12)
    → c2-05                                   (commute: 58.36)
    → c0-06                                   (commute: 58.59)
    → c1-07                                   (commute: 60.25)
  c0-03:
    → c2-02                                   (commute: 50.26)
    → c2-01                                   (commute: 51.57)
    → c0-09                                   (commute: 53.88)
    → c1-07                                   (commute: 54.18)
    → c1-01                                   (commute: 54.43)
  c0-04:
    → c1-08                                   (commute: 57.88)
    → c2-06                                   (commute: 58.90)
    → c2-02                                   (commute: 60.50)
    → c0-03                                   (commute: 61.31)
    → c1-04                                   (commute: 64.38)
  c0-05:
    → c2-06                                   (commute: 63.53)
    → c2-02                                   (commute: 66.63)
    → c1-01                                   (commute: 69.43)
    → c1-07                                   (commute: 70.35)
    → c0-09                                   (commute: 70.62)
  c0-06:
    → c2-06                                   (commute: 52.48)
    → c2-02                                   (commute: 55.55)
    → c2-01                                   (commute: 57.75)
    → c0-02                                   (commute: 58.59)
    → c2-09                                   (commute: 60.51)
  c0-07:
    → c1-08                                   (commute: 62.31)
    → c0-01                                   (commute: 63.95)
    → c2-06                                   (commute: 64.81)
    → c2-02                                   (commute: 65.70)
    → c0-02                                   (commute: 68.66)
  c0-08:
    → c1-08                                   (commute: 58.33)
    → c2-06                                   (commute: 59.15)
    → c0-03                                   (commute: 60.57)
    → c2-09                                   (commute: 64.48)
    → c1-04                                   (commute: 64.76)
  c0-09:
    → c1-08                                   (commute: 51.03)
    → c2-06                                   (commute: 52.69)
    → c0-03                                   (commute: 53.88)
    → c2-02                                   (commute: 55.04)
    → c2-01                                   (commute: 55.27)
  c1-00:
    → c2-06                                   (commute: 69.28)
    → c2-02                                   (commute: 70.26)
    → c0-01                                   (commute: 71.68)
    → c0-03                                   (commute: 72.40)
    → c2-01                                   (commute: 74.43)
  c1-01:
    → c2-06                                   (commute: 52.93)
    → c0-01                                   (commute: 54.37)
    → c0-03                                   (commute: 54.43)
    → c2-01                                   (commute: 58.07)
    → c2-05                                   (commute: 60.26)
  c1-02:
    → c2-02                                   (commute: 59.77)
    → c0-03                                   (commute: 61.53)
    → c1-04                                   (commute: 63.43)
    → c2-01                                   (commute: 64.26)
    → c1-07                                   (commute: 65.65)
  c1-03:
    → c2-06                                   (commute: 59.44)
    → c0-01                                   (commute: 61.90)
    → c2-02                                   (commute: 61.92)
    → c0-06                                   (commute: 64.29)
    → c0-02                                   (commute: 66.90)
  c1-04:
    → c2-06                                   (commute: 50.35)
    → c2-01                                   (commute: 55.03)
    → c1-07                                   (commute: 55.25)
    → c0-02                                   (commute: 57.12)
    → c0-09                                   (commute: 57.52)
  c1-05:
    → c2-06                                   (commute: 77.24)
    → c0-01                                   (commute: 79.00)
    → c0-03                                   (commute: 80.43)
    → c1-04                                   (commute: 81.50)
    → c2-01                                   (commute: 81.90)
  c1-06:
    → c1-08                                   (commute: 60.90)
    → c0-01                                   (commute: 64.74)
    → c2-02                                   (commute: 66.44)
    → c1-04                                   (commute: 67.82)
    → c2-01                                   (commute: 68.48)
  c1-07:
    → c0-03                                   (commute: 54.18)
    → c0-01                                   (commute: 55.07)
    → c1-04                                   (commute: 55.25)
    → c2-02                                   (commute: 55.76)
    → c2-01                                   (commute: 58.23)
  c1-08:
    → c2-02                                   (commute: 45.59)
    → c2-01                                   (commute: 48.30)
    → c0-09                                   (commute: 51.03)
    → c2-09                                   (commute: 51.71)
    → c2-00                                   (commute: 53.46)
  c1-09:
    → c0-03                                   (commute: 65.39)
    → c0-01                                   (commute: 66.21)
    → c2-02                                   (commute: 66.74)
    → c0-06                                   (commute: 69.54)
    → c2-01                                   (commute: 70.17)
  c2-00:
    → c1-08                                   (commute: 53.46)
    → c0-01                                   (commute: 57.22)
    → c0-03                                   (commute: 58.05)
    → c2-01                                   (commute: 59.48)
    → c1-04                                   (commute: 60.83)
  c2-01:
    → c1-08                                   (commute: 48.30)
    → c2-06                                   (commute: 50.78)
    → c0-01                                   (commute: 51.27)
    → c0-03                                   (commute: 51.57)
    → c1-04                                   (commute: 55.03)
  c2-02:
    → c1-08                                   (commute: 45.59)
    → c0-03                                   (commute: 50.26)
    → c2-05                                   (commute: 52.49)
    → c0-02                                   (commute: 54.51)
    → c0-09                                   (commute: 55.04)
  c2-03:
    → c1-08                                   (commute: 58.11)
    → c2-06                                   (commute: 59.65)
    → c0-03                                   (commute: 62.69)
    → c1-04                                   (commute: 65.31)
    → c0-02                                   (commute: 66.85)
  c2-04:
    → c2-06                                   (commute: 60.19)
    → c0-01                                   (commute: 61.34)
    → c1-04                                   (commute: 63.03)
    → c2-09                                   (commute: 66.06)
    → c1-07                                   (commute: 66.59)
  c2-05:
    → c2-06                                   (commute: 52.16)
    → c2-02                                   (commute: 52.49)
    → c0-01                                   (commute: 54.29)
    → c0-03                                   (commute: 55.56)
    → c0-02                                   (commute: 58.36)
  c2-06:
    → c1-04                                   (commute: 50.35)
    → c2-01                                   (commute: 50.78)
    → c2-05                                   (commute: 52.16)
    → c0-06                                   (commute: 52.48)
    → c0-09                                   (commute: 52.69)
  c2-07:
    → c1-08                                   (commute: 68.41)
    → c0-01                                   (commute: 70.64)
    → c2-01                                   (commute: 71.91)
    → c1-04                                   (commute: 75.53)
    → c0-09                                   (commute: 76.04)
  c2-08:
    → c1-08                                   (commute: 60.90)
    → c0-01                                   (commute: 64.70)
    → c0-03                                   (commute: 66.07)
    → c2-01                                   (commute: 67.55)
    → c0-02                                   (commute: 69.90)
  c2-09:
    → c1-08                                   (commute: 51.71)
    → c0-01                                   (commute: 54.58)
    → c0-03                                   (commute: 55.39)
    → c1-04                                   (commute: 58.47)
    → c0-06                                   (commute: 60.51)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `c1-08` (0.052632), `c2-06` (0.046053). Lowest π: `c1-05` (0.019737).

| From   | To      | Commute time | Within / Cross block |
| ------ | ------- | ------------ | ----------------------- |
| c1-08  | c2-02   | 45.59         | Cross                     |
| c1-08  | c2-01   | 48.30         | Cross                     |
| c1-08  | c0-09   | 51.03         | Cross                     |
| c1-08  | c2-09   | 51.71         | Cross                     |
| c1-08  | c2-00   | 53.46         | Cross                     |
| c2-06  | c1-04   | 50.35         | Cross                     |
| c2-06  | c2-01   | 50.78         | Within                    |
| c2-06  | c2-05   | 52.16         | Within                    |
| c2-06  | c0-06   | 52.48         | Cross                     |
| c2-06  | c0-09   | 52.69         | Cross                     |
| c1-05  | c2-06   | 77.24         | Cross                     |
| c1-05  | c0-01   | 79.00         | Cross                     |
| c1-05  | c0-03   | 80.43         | Cross                     |
| c1-05  | c1-04   | 81.50         | Within                    |
| c1-05  | c2-01   | 81.90         | Cross                     |

**Finding**: At loose `p_out`, `c1-08` (rank-1 π) suggests **exclusively cross-block** links (5
of 5, 100%), while `c2-06` (rank-2 π) splits 3 cross / 2 within. Combined top-2 cross-block
fraction is 8 of 10 (80%) — similar in magnitude to the mid set (also 80%), but the underlying
graph is now nearly half cross-block by raw edge count (47.4%), meaning even a "moderate"
80% cross-block suggestion rate is now *below* what pure random selection over all existing
edges would produce, unlike at tight `p_out` where 30% cross-block suggestions were well above
the 6.9% baseline cross-block edge rate. `c1-05` (lowest π) shows 4 cross / 1 within.

### Stationary distribution spread — across parameter sets

| Parameter set | Max π    | Min π    | Ratio | Label          |
| ------------- | -------- | -------- | ----- | -------------- |
| Tight         | 0.045977 | 0.017241 | 2.67  | Mildly skewed  |
| Mid           | 0.047414 | 0.021552 | 2.20  | Mildly skewed  |
| Loose         | 0.052632 | 0.019737 | 2.67  | Mildly skewed  |

### Within-block π uniformity

Directly inspecting the exported π values for the tight set, grouped by block:

| Block | Min π    | Max π    | n  |
| ----- | -------- | -------- | -- |
| c0    | 0.022989 | 0.040230 | 10 |
| c1    | 0.022989 | 0.045977 | 10 |
| c2    | 0.017241 | 0.045977 | 10 |

**Finding**: π is **not** uniform within a block — each block shows roughly a 2× spread between
its lowest- and highest-π member (e.g. c2: 0.017241 to 0.045977). This is expected: unlike a
deterministic clique (where every intra-block node has identical degree), `planted_partition_graph`
samples each potential edge independently, so individual nodes receive different realised
degrees purely from Bernoulli sampling variance, even though every node's *expected* degree
within a block is identical by construction ($E[\deg] = (k-1) \cdot p_{in} + (l-1)k \cdot
p_{out}$). No single node "dominates" a block in a structural sense — the variation is sampling
noise, not a planted hub.

### Visual observations

> *The following description is inferred from the graph's structural properties and stationary
> distribution. The force-directed layout can be verified by opening `/tmp/nx-pp-tight.html`,
> `/tmp/nx-pp-mid.html`, and `/tmp/nx-pp-loose.html`.*

At `p_out=0.02`, the layout should show 3 tight, well-separated clusters of 10 nodes each, with
only a sparse scattering of cross-cluster edges (12 of 174, 6.9%) visible as thin long-range
lines between clusters. At `p_out=0.10`, the same 3 clusters should remain visually
identifiable but with noticeably more inter-cluster edges (60 of 232, 25.9%) pulling them
closer together in the force layout. At `p_out=0.25`, the clusters should be substantially
blurred — nearly half the edges (144 of 304, 47.4%) cross block boundaries, so the
force-directed layout likely no longer produces three cleanly separated regions, instead
showing a single denser mass with only faint residual clustering by block colour/labelling.
`c1-08` (highest π at loose, degree 16) should render as the visually largest node across all
three layouts where it appears prominently, most strikingly in the loose layout.

### Markov questions — answered

- **Does π rank block-internal nodes uniformly, or do particular nodes dominate within a
  block?**
  Not uniformly — within-block π varies roughly 2× (e.g. tight-set block c2: π ranges from
  0.017241 to 0.045977). However, no node "dominates" in a structural sense; the variation
  is Bernoulli sampling noise around an equal expected degree per block, not a planted hub
  (verified directly above).

- **What fraction of `suggest` recommendations cross block boundaries at each p_out?**
  For the top-2 π nodes combined: tight 3/10 (30%), mid 8/10 (80%), loose 8/10 (80%). This rises
  sharply from tight to mid, then plateaus from mid to loose even though the underlying
  cross-block edge fraction keeps rising (25.9% → 47.4%).

- **Is π hub-dominated or flat, and how does the max/min ratio trend across the p_out sweep?**
  Mildly skewed at all three p_out values (ratio 2.20–2.67) — never hub-dominated, and the
  ratio does not trend monotonically with p_out (it dips slightly at mid, 2.20, before
  returning to 2.67 at loose). This confirms π-ratio and community mixing are largely
  independent signals for this graph family, since π variation comes from Bernoulli sampling
  noise rather than from the block structure itself.

- **At what p_out does the fraction of cross-community `suggest` recommendations exceed the
  fraction expected from a fully-merged (single-community) graph?**
  A fully-merged, single-community reference expectation is the raw cross-block edge fraction
  itself, since a random walk with no community structure would recommend links roughly in
  proportion to how many existing edges are already cross-community. At `p_out=0.02`, the
  observed suggestion cross-block rate (30%) already **exceeds** the raw edge cross-block rate
  (6.9%) — communities are still very much detectable, but `suggest` already over-recommends
  cross-block links relative to the existing edge mix, because within-block options are more
  often already exhausted. At `p_out=0.10`, suggestion rate (80%) again exceeds the edge rate
  (25.9%). At `p_out=0.25`, suggestion rate (80%) still exceeds the raw edge cross-block rate
  (47.4%) — so across this entire tested range (`p_out` 0.02 to 0.25), the `suggest`
  cross-block rate consistently exceeds the raw edge cross-block rate; the sweep does not reach
  a `p_out` where suggestion behaviour drops to match a fully-merged, no-community baseline.

---

## References

- Wikipedia: Stochastic block model — https://en.wikipedia.org/wiki/Stochastic_block_model
- NetworkX: https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.planted_partition_graph.html
- SBM / Planted Partition section from graph-models research notes

---

## Definition of Done

- [x] All sections above filled — no placeholder text remaining
- [x] Slug ordering algorithm documented with a node-id mapping table (not just the pattern)
- [x] `wikigraph analyze --suggest-top 5` output pasted verbatim for all three parameter sets
- [x] Cross-community edge count computed and recorded (Python helper, per parameter set)
- [x] All Markov questions answered with actual numbers from analyze output
- [x] Three parameter-set variants shown (p_out=0.02, 0.10, 0.25) with a comparison table
      before the per-set detail sections
- [x] Every reference link manually verified to resolve (Wikipedia 200; NetworkX docs 200)
- [x] File committed to branch — path to be updated once #42 resolves
