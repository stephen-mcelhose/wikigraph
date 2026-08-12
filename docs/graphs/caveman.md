<!-- Path provisional — will be updated when #42 resolves -->

# Relaxed Caveman Graph — clique-rewiring mixing benchmark

> ⚠️ **DRAFT** — AI-assisted write-up, not yet verified by human analysis. Treat all findings as provisional.

> **nx-to-wiki flag**: `--graph caveman --l [cliques] --k [clique-size] --p [rewire-prob]`
> **Nodes**: l·k (parameterised) · **Directed edges**: 2×undirected (parameterised)
> **Naming**: tier1-caveman-clique-membership
> **Source**: generative — parameterised

---

## Background

The relaxed caveman graph starts from the "pure" caveman construction — `l` disjoint complete
graphs (cliques) of size `k`, resembling isolated "caves" of fully-connected individuals — and
then rewires each intra-clique edge to a random inter-clique target with independent
probability `p`. At `p=0` the graph is exactly `l` disconnected cliques (no random walk can
cross between them at all). As `p` rises, cross-clique edges accumulate continuously, and the
clique boundaries progressively blur until, at high enough `p`, the graph is statistically
indistinguishable from an Erdős–Rényi random graph with the same edge count. `p` is the
discrete, edge-rewiring analogue of the LFR benchmark's continuous mixing parameter `μ`: both
control how much of a node's connectivity budget is spent outside its home community, and both
are used to study how quickly community-detection algorithms — and random walks — lose the
ability to distinguish separate clusters.

In Markov-chain terms, a node inside a clique of size `k` has degree `k - 1` from its
intra-clique ties, plus roughly `p · (l - 1)` expected inter-clique ties added by rewiring
(one potential rewire target per node, times `(l-1)` other cliques, times success probability
`p`). The random walk's **conductance** across the clique boundary — the fraction of a random
walker's steps that leave its home clique — scales with `p·(l-1) / (k - 1 + p·(l-1))`,
approaching 0 as `p → 0` (walk trapped inside its clique) and approaching `(l-1)/l` as `p → 1`
(walk indifferent to clique membership).

**Origin**: NetworkX's `relaxed_caveman_graph` implementation, based on Watts's "small-world"
caveman construction lineage; no single canonical citation — see NetworkX reference below.

---

## Graph Properties

**Key parameters**:

| Flag   | Default | Range        | What it controls                                             |
| ------ | ------- | ------------ | --------------------------------------------------------------- |
| `--l`  | 3       | ≥ 1          | Number of cliques ("caves")                                     |
| `--k`  | 4       | ≥ 2          | Nodes per clique                                                 |
| `--p`  | 0.1     | [0.0, 1.0]   | Probability each intra-clique edge is rewired to a random inter-clique target |
| `--seed` | none  | any int      | RNG seed for reproducible rewiring                               |

Three parameter sets (as directed by the issue), all with `l=4`, `k=6`, `seed=42`, so nodes and
edge count are fixed at 24/120 — only rewiring outcome changes:

| Parameter set | l | k | p    | Nodes | Undirected edges | Directed edges | Min deg | Max deg |
| ------------- | - | - | ---- | ----- | ------------------ | ---------------- | -------- | -------- |
| p=0.05        | 4 | 6 | 0.05 | 24    | 60                  | 120               | 4        | 6        |
| p=0.15        | 4 | 6 | 0.15 | 24    | 60                  | 120               | 4        | 6        |
| p=0.30        | 4 | 6 | 0.30 | 24    | 60                  | 120               | 3        | 7        |

Community structure: 4 cliques of 6 nodes each, defined structurally by node-id range (not a
NetworkX attribute — `nx-to-wiki` derives clique membership directly from `n // k`).

---

## Transition Matrix Structure

For every graph produced by `nx-to-wiki`, the random walk is **uniform**: each step chooses a
neighbour with equal probability. The transition matrix $P$ is fully determined by the
adjacency structure of the `.md` files:

$$P_{ij} = \frac{1}{\deg(i)} \quad \text{if } [[j]] \text{ appears in i.md}, \quad \text{else } 0$$

### p=0.05

| Slug      | Degree | $P_{ij}$ (non-zero) | Structural note                                       |
| --------- | ------ | -------------------- | -------------------------------------------------------- |
| cave1-02  | 6      | 1/6 ≈ 0.166667        | Highest-π node; retains full intra-clique degree + 1 inbound rewire |
| cave0-00  | 5      | 1/5 = 0.200000        | Boundary node — 4 intra-clique + 1 cross-clique edge (→cave1-02) |
| cave3-05  | 5      | 1/5 = 0.200000        | Interior node of the still-isolated cave3 clique          |

### p=0.15

| Slug      | Degree | $P_{ij}$ (non-zero) | Structural note                                     |
| --------- | ------ | -------------------- | ------------------------------------------------------ |
| cave1-02  | 6      | 1/6 ≈ 0.166667        | Highest-π node, same structural role as at p=0.05        |
| cave2-00  | 5      | 1/5 = 0.200000        | Boundary node with one cross-clique tie                   |
| cave3-05  | 5      | 1/5 = 0.200000        | Interior node — cave3 remains its own communicating class at this p |

### p=0.30

| Slug      | Degree | $P_{ij}$ (non-zero) | Structural note                                          |
| --------- | ------ | -------------------- | ------------------------------------------------------------ |
| cave2-01  | 7      | 1/7 ≈ 0.142857        | Highest-π node; gained extra degree via rewiring at higher p    |
| cave3-00  | 5      | 1/5 = 0.200000        | cave3 now genuinely connected to the rest of the graph          |
| cave3-05  | 3      | 1/3 ≈ 0.333333        | Lowest-π node; lost intra-clique edges to rewiring, degree dropped to 3 |

**The `.md` files are a lossless encoding of $P$ for all three parameter sets.** Minimum
non-zero $P_{ij}$ is $1/7 \approx 0.143$ (p=0.30's highest-degree node) — far above the
`--min-edge` default filter of 0.005 in every case tested.

### Export

```bash
wikigraph export /tmp/nxwiki-caveman-005 --format csv -o /tmp/caveman-005-export
wikigraph export /tmp/nxwiki-caveman-015 --format csv -o /tmp/caveman-015-export
wikigraph export /tmp/nxwiki-caveman-030 --format csv -o /tmp/caveman-030-export
```

Sparse edge-row counts: 120 for all three sets (matching directed-edge counts). No
`--min-edge 0` warning applies to any of the three sets tested.

---

## Slug Naming

**Naming tier**: Tier 1

**Assignment algorithm** (`build_caveman`, verbatim from source): NetworkX's
`relaxed_caveman_graph(l, k, p, seed)` numbers nodes so that clique `i` (0-indexed) occupies the
contiguous range `[i*k, (i+1)*k)`. `nx-to-wiki` computes, for each node `n`:

- `clique_idx = n // k` — integer division gives the clique index
- `within = n % k` — remainder gives the position within that clique
- `slug = f"cave{clique_idx}-{pad(within, k)}"`

`pad(within, k)` zero-pads `within` to `max(2, len(str(k - 1)))` digits.

**Why this ordering**: NetworkX builds the graph as a disjoint union of `l` complete graphs
before rewiring any edges, so node-id ranges already partition nodes by clique membership;
rewiring changes *edges*, not node ids, so the slug mapping is identical across all three `p`
values even though the resulting adjacency differs.

Full node-id → slug mapping (l=4, k=6 — identical across all three parameter sets since only
edges change with `p`, not node numbering):

| Slug      | NetworkX node id | Clique index |
| --------- | ------------------ | -------------- |
| cave0-00  | 0                   | 0              |
| cave0-01  | 1                   | 0              |
| cave0-02  | 2                   | 0              |
| cave0-03  | 3                   | 0              |
| cave0-04  | 4                   | 0              |
| cave0-05  | 5                   | 0              |
| cave1-00  | 6                   | 1              |
| cave1-01  | 7                   | 1              |
| cave1-02  | 8                   | 1              |
| cave1-03  | 9                   | 1              |
| cave1-04  | 10                  | 1              |
| cave1-05  | 11                  | 1              |
| cave2-00  | 12                  | 2              |
| cave2-01  | 13                  | 2              |
| cave2-02  | 14                  | 2              |
| cave2-03  | 15                  | 2              |
| cave2-04  | 16                  | 2              |
| cave2-05  | 17                  | 2              |
| cave3-00  | 18                  | 3              |
| cave3-01  | 19                  | 3              |
| cave3-02  | 20                  | 3              |
| cave3-03  | 21                  | 3              |
| cave3-04  | 22                  | 3              |
| cave3-05  | 23                  | 3              |

---

## How to Generate

```bash
# p-sweep with a fixed seed so all three sets share the same underlying edge-rewiring RNG state
python3 tools/nx-to-wiki/main.py --graph caveman --l 4 --k 6 --p 0.05 --seed 42 --out /tmp/nxwiki-caveman-005
python3 tools/nx-to-wiki/main.py --graph caveman --l 4 --k 6 --p 0.15 --seed 42 --out /tmp/nxwiki-caveman-015
python3 tools/nx-to-wiki/main.py --graph caveman --l 4 --k 6 --p 0.30 --seed 42 --out /tmp/nxwiki-caveman-030

wikigraph graph /tmp/nxwiki-caveman-005 -o /tmp/nx-caveman-005.html --title "Caveman p=0.05"
wikigraph graph /tmp/nxwiki-caveman-015 -o /tmp/nx-caveman-015.html --title "Caveman p=0.15"
wikigraph graph /tmp/nxwiki-caveman-030 -o /tmp/nx-caveman-030.html --title "Caveman p=0.30"

wikigraph analyze /tmp/nxwiki-caveman-005 --suggest-top 5
wikigraph analyze /tmp/nxwiki-caveman-015 --suggest-top 5
wikigraph analyze /tmp/nxwiki-caveman-030 --suggest-top 5
```

**Effect of `p`**: Raising `p` from 0.05 to 0.30 does not change node or edge counts (both
fixed by `l` and `k`) — it redistributes which edges are intra-clique vs. inter-clique. At low
`p`, most cliques remain their own communicating class in `wikigraph analyze` (see below); at
`p=0.30`, enough rewired edges accumulate that the whole graph merges into a single
communicating class and `wikigraph suggest` begins proposing genuine cross-clique links (empty
at p=0.05/0.15).

**What each output directory contains:**

- 24 `.md` files, one per node (fixed across all three parameter sets since `l`, `k` are fixed)
- Each file links to all its undirected neighbours (symmetric — every link is bidirectional)
- File layout: `# {slug}\n\n[[neighbour-1]] [[neighbour-2]] ...`
- Why directed edges = 2× undirected: `G.to_directed()` adds both u→v and v→u for every edge

---

## wikigraph Analysis

> **⚠️ Communicating classes are informative here — an exception to the usual caveat.** Because
> `relaxed_caveman_graph` at low `p` can produce a genuinely **disconnected** undirected graph
> (isolated cliques with zero rewired edges connecting them), `G.to_directed()` does *not*
> always yield a single strongly connected component for this graph family. At p=0.05 and
> p=0.15, `wikigraph analyze` reports **2 communicating classes** — cave3 remains fully
> isolated from cave0/cave1/cave2 at both low-p values with this seed. Only at p=0.30 does the
> graph become a single class. This is the one graph in this series where communicating-class
> count is itself a meaningful community-detection signal.

### Comparison across parameter sets

| Parameter set | p    | Nodes | Directed edges | Classes | Entropy rate | Max π    | Min π    | Ratio | Cross-community % | Boundary nodes |
| ------------- | ---- | ----- | ---------------- | -------- | -------------- | -------- | -------- | ----- | -------------------- | ---------------- |
| p=0.05        | 0.05 | 24    | 120               | 2        | 2.3316 bits    | 0.050000 | 0.033333 | 1.5   | 8 of 120 (6.7%)       | 7                 |
| p=0.15        | 0.15 | 24    | 120               | 2        | 2.3340 bits    | 0.050000 | 0.033333 | 1.5   | 10 of 120 (8.3%)      | 8                 |
| p=0.30        | 0.30 | 24    | 120               | 1        | 2.3488 bits    | 0.058333 | 0.025000 | 2.33  | 20 of 120 (16.7%)     | 13                |

**Trend**: The cross-community edge fraction rises steadily with `p` (6.7% → 8.3% → 16.7%),
and boundary-node count roughly doubles (7 → 8 → 13). The max/min π ratio, by contrast, does
**not** monotonically fall toward 1.0 as `p` rises within this range — it actually *increases*
slightly at p=0.30 (1.5 → 2.33), because rewiring at p=0.30 concentrates additional edges on a
few nodes (cave2-01 gains degree 7, the highest of any node across all three sets) while
stripping edges from others (cave3-05 drops to degree 3, the lowest). All three ratios remain
**below 5** (mildly skewed at worst, effectively flat at p=0.05/0.15), so this p-sweep does not
reach a regime where π itself is strongly hub-dominated — the key story here is the
communicating-class collapse (2→1) and the rising cross-community edge/boundary-node counts,
not extreme π skew.

### p=0.05

#### Raw analyze output

```
=== Overview ===
Pages:        24
Edges:        120
Entropy rate: 2.3316 bits
Classes:      2

=== Communicating classes ===
Class 1 (recurrent): 6 page(s)
  cave3-00
  cave3-01
  cave3-02
  cave3-03
  cave3-04
  cave3-05
Class 2 (recurrent): 18 page(s)
  cave0-00
  cave0-01
  cave0-02
  cave0-03
  cave0-04
  cave0-05
  cave1-00
  cave1-01
  cave1-02
  cave1-03
  cave1-04
  cave1-05
  cave2-00
  cave2-01
  cave2-02
  cave2-03
  cave2-04
  cave2-05

=== Orphan pages (bottom 10% by stationary distribution) ===
  cave2-03                                  π=0.033333  → add inbound links
  cave1-05                                  π=0.033333  → add inbound links
  cave0-02                                  π=0.033333  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. cave1-02                                  π=0.050000
  2. cave1-00                                  π=0.050000
  3. cave1-01                                  π=0.050000
  4. cave2-02                                  π=0.050000
  5. cave0-01                                  π=0.041667

=== Suggested missing links (lowest commute time, not yet linked) ===
  (none)
```

#### Suggested links (top-2 π + lowest-π)

**No suggestions exist at p=0.05.** `wikigraph suggest` output is `(none)` — with only 2
communicating classes and each clique already internally complete (or near-complete after
rewiring), there are no not-yet-linked node pairs with a finite, useful commute time to
propose within the disconnected structure that the tool surfaces. This itself is a finding: at
low `p`, the graph is close enough to "fully saturated" within each reachable component that no
missing-link recommendation applies.

### p=0.15

#### Raw analyze output

```
=== Overview ===
Pages:        24
Edges:        120
Entropy rate: 2.3340 bits
Classes:      2

=== Communicating classes ===
Class 1 (recurrent): 6 page(s)
  cave3-00
  cave3-01
  cave3-02
  cave3-03
  cave3-04
  cave3-05
Class 2 (recurrent): 18 page(s)
  cave0-00
  cave0-01
  cave0-02
  cave0-03
  cave0-04
  cave0-05
  cave1-00
  cave1-01
  cave1-02
  cave1-03
  cave1-04
  cave1-05
  cave2-00
  cave2-01
  cave2-02
  cave2-03
  cave2-04
  cave2-05

=== Orphan pages (bottom 10% by stationary distribution) ===
  cave2-04                                  π=0.033333  → add inbound links
  cave1-05                                  π=0.033333  → add inbound links
  cave0-02                                  π=0.033333  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. cave1-02                                  π=0.050000
  2. cave1-00                                  π=0.050000
  3. cave1-01                                  π=0.050000
  4. cave2-01                                  π=0.050000
  5. cave2-00                                  π=0.050000

=== Suggested missing links (lowest commute time, not yet linked) ===
  (none)
```

#### Suggested links (top-2 π + lowest-π)

**Still no suggestions at p=0.15**, for the same reason as p=0.05 — cave3 remains its own
isolated communicating class (still 2 classes total) and the reachable subgraph offers no
useful not-yet-linked candidate.

### p=0.30

#### Raw analyze output

```
=== Overview ===
Pages:        24
Edges:        120
Entropy rate: 2.3488 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 24 page(s)
  cave0-00
  cave0-01
  cave0-02
  cave0-03
  cave0-04
  cave0-05
  cave1-00
  cave1-01
  cave1-02
  cave1-03
  cave1-04
  cave1-05
  cave2-00
  cave2-01
  cave2-02
  cave2-03
  cave2-04
  cave2-05
  cave3-00
  cave3-01
  cave3-02
  cave3-03
  cave3-04
  cave3-05

=== Orphan pages (bottom 10% by stationary distribution) ===
  cave3-05                                  π=0.025000  → add inbound links
  cave1-03                                  π=0.033333  → add inbound links
  cave1-01                                  π=0.033333  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. cave2-01                                  π=0.058333
  2. cave3-00                                  π=0.050000
  3. cave2-02                                  π=0.050000
  4. cave0-01                                  π=0.050000
  5. cave2-00                                  π=0.050000

=== Suggested missing links (lowest commute time, not yet linked) ===
  cave0-00:
    → cave0-02                                (commute: 55.90)
    → cave2-01                                (commute: 59.78)
    → cave1-05                                (commute: 67.12)
    → cave2-03                                (commute: 68.53)
    → cave2-00                                (commute: 70.30)
  cave0-01:
    → cave2-02                                (commute: 56.43)
    → cave2-00                                (commute: 59.25)
    → cave1-00                                (commute: 63.55)
    → cave2-03                                (commute: 63.86)
    → cave0-05                                (commute: 64.71)
  cave0-02:
    → cave0-00                                (commute: 55.90)
    → cave2-01                                (commute: 86.48)
    → cave1-00                                (commute: 92.29)
    → cave1-02                                (commute: 93.44)
    → cave2-04                                (commute: 94.87)
  cave0-03:
    → cave0-01                                (commute: 64.71)
    → cave2-01                                (commute: 90.74)
    → cave1-00                                (commute: 91.64)
    → cave1-02                                (commute: 92.64)
    → cave2-04                                (commute: 100.25)
  cave0-04:
    → cave2-01                                (commute: 77.60)
    → cave1-00                                (commute: 81.30)
    → cave1-02                                (commute: 82.39)
    → cave2-04                                (commute: 86.47)
    → cave2-02                                (commute: 88.89)
  cave0-05:
    → cave0-01                                (commute: 64.71)
    → cave2-01                                (commute: 90.74)
    → cave1-00                                (commute: 91.64)
    → cave1-02                                (commute: 92.64)
    → cave2-04                                (commute: 100.25)
  cave1-00:
    → cave2-03                                (commute: 55.00)
    → cave2-00                                (commute: 59.87)
    → cave1-03                                (commute: 60.68)
    → cave1-01                                (commute: 61.57)
    → cave0-01                                (commute: 63.55)
  cave1-01:
    → cave1-02                                (commute: 54.82)
    → cave1-00                                (commute: 61.57)
    → cave2-03                                (commute: 64.33)
    → cave2-01                                (commute: 67.07)
    → cave2-02                                (commute: 75.54)
  cave1-02:
    → cave2-01                                (commute: 53.80)
    → cave1-01                                (commute: 54.82)
    → cave2-00                                (commute: 58.27)
    → cave2-02                                (commute: 64.31)
    → cave0-01                                (commute: 65.06)
  cave1-03:
    → cave1-00                                (commute: 60.68)
    → cave2-03                                (commute: 68.59)
    → cave2-01                                (commute: 74.29)
    → cave2-00                                (commute: 74.48)
    → cave0-00                                (commute: 82.84)
  cave1-04:
    → cave2-03                                (commute: 61.23)
    → cave2-01                                (commute: 65.03)
    → cave2-00                                (commute: 66.95)
    → cave0-00                                (commute: 72.77)
    → cave2-02                                (commute: 75.95)
  cave1-05:
    → cave2-01                                (commute: 56.38)
    → cave2-00                                (commute: 58.27)
    → cave2-02                                (commute: 66.41)
    → cave0-00                                (commute: 67.12)
    → cave2-04                                (commute: 71.91)
  cave2-00:
    → cave1-02                                (commute: 58.27)
    → cave1-05                                (commute: 58.27)
    → cave0-01                                (commute: 59.25)
    → cave1-00                                (commute: 59.87)
    → cave1-04                                (commute: 66.95)
  cave2-01:
    → cave1-02                                (commute: 53.80)
    → cave1-05                                (commute: 56.38)
    → cave0-00                                (commute: 59.78)
    → cave1-04                                (commute: 65.03)
    → cave1-01                                (commute: 67.07)
  cave2-02:
    → cave0-01                                (commute: 56.43)
    → cave1-02                                (commute: 64.31)
    → cave1-00                                (commute: 64.91)
    → cave1-05                                (commute: 66.41)
    → cave0-00                                (commute: 71.60)
  cave2-03:
    → cave1-00                                (commute: 55.00)
    → cave2-04                                (commute: 55.57)
    → cave1-04                                (commute: 61.23)
    → cave2-05                                (commute: 61.46)
    → cave0-01                                (commute: 63.86)
  cave2-04:
    → cave2-03                                (commute: 55.57)
    → cave1-00                                (commute: 68.17)
    → cave1-02                                (commute: 69.09)
    → cave0-00                                (commute: 71.07)
    → cave1-05                                (commute: 71.91)
  cave2-05:
    → cave2-03                                (commute: 61.46)
    → cave0-01                                (commute: 69.33)
    → cave1-00                                (commute: 76.64)
    → cave1-02                                (commute: 77.46)
    → cave1-05                                (commute: 79.33)
  cave3-00:
    → cave0-01                                (commute: 88.48)
    → cave2-01                                (commute: 95.22)
    → cave2-04                                (commute: 100.47)
    → cave2-00                                (commute: 102.17)
    → cave2-03                                (commute: 108.45)
  cave3-01:
    → cave2-02                                (commute: 104.89)
    → cave0-01                                (commute: 105.95)
    → cave2-01                                (commute: 117.41)
    → cave2-04                                (commute: 122.72)
    → cave2-00                                (commute: 125.15)
  cave3-02:
    → cave3-05                                (commute: 69.12)
    → cave2-02                                (commute: 88.48)
    → cave2-01                                (commute: 96.97)
    → cave2-04                                (commute: 102.33)
    → cave2-00                                (commute: 105.69)
  cave3-03:
    → cave2-02                                (commute: 104.89)
    → cave0-01                                (commute: 105.95)
    → cave2-01                                (commute: 117.41)
    → cave2-04                                (commute: 122.72)
    → cave2-00                                (commute: 125.15)
  cave3-04:
    → cave3-05                                (commute: 71.82)
    → cave2-02                                (commute: 109.17)
    → cave0-01                                (commute: 109.70)
    → cave2-01                                (commute: 121.48)
    → cave2-04                                (commute: 126.79)
  cave3-05:
    → cave3-02                                (commute: 69.12)
    → cave3-04                                (commute: 71.82)
    → cave2-02                                (commute: 123.35)
    → cave0-01                                (commute: 127.04)
    → cave2-01                                (commute: 136.93)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `cave2-01` (0.058333), `cave3-00` (0.050000, 2nd-place, tied with 3 others but shown
here as it is also part of the graph's smallest/most-recently-connected clique). Lowest π:
`cave3-05` (0.025000).

| From      | To        | Commute time | Within / Cross clique |
| --------- | --------- | ------------ | ------------------------ |
| cave2-01  | cave1-02  | 53.80        | Cross                     |
| cave2-01  | cave1-05  | 56.38        | Cross                     |
| cave2-01  | cave0-00  | 59.78        | Cross                     |
| cave2-01  | cave1-04  | 65.03        | Cross                     |
| cave2-01  | cave1-01  | 67.07        | Cross                     |
| cave3-00  | cave0-01  | 88.48        | Cross                     |
| cave3-00  | cave2-01  | 95.22        | Cross                     |
| cave3-00  | cave2-04  | 100.47       | Cross                     |
| cave3-00  | cave2-00  | 102.17       | Cross                     |
| cave3-00  | cave2-03  | 108.45       | Cross                     |
| cave3-05  | cave3-02  | 69.12        | Within                    |
| cave3-05  | cave3-04  | 71.82        | Within                    |
| cave3-05  | cave2-02  | 123.35       | Cross                     |
| cave3-05  | cave0-01  | 127.04       | Cross                     |
| cave3-05  | cave2-01  | 136.93       | Cross                     |

**Finding**: At p=0.30, all 5 of `cave2-01`'s suggestions and all 5 of `cave3-00`'s
suggestions cross clique boundaries — every suggestion from these two high-π nodes is
cross-community. `cave3-05` (the lowest-π node) suggests 2 within-clique links first (its own
recently-thinned clique, cave3) before crossing to cave2/cave0 for its remaining 3 — reflecting
that cave3 lost internal edges to rewiring and now has genuine gaps to fill locally as well as
externally. For the remaining 21 nodes not shown, the fraction of cross-clique suggestions
tracks each node's own remaining intra-clique degree: nodes that kept most of their original 5
intra-clique ties suggest mostly cross-clique links (little left to gain locally), while nodes
that lost ties to rewiring suggest more within-clique links first.

### Stationary distribution spread — across parameter sets

| Parameter set | Max π    | Min π    | Ratio | Label              |
| ------------- | -------- | -------- | ----- | -------------------- |
| p=0.05        | 0.050000 | 0.033333 | 1.50  | Effectively flat      |
| p=0.15        | 0.050000 | 0.033333 | 1.50  | Effectively flat      |
| p=0.30        | 0.058333 | 0.025000 | 2.33  | Mildly skewed         |

### Visual observations

> *The following description is inferred from the graph's structural properties and stationary
> distribution. The force-directed layout can be verified by opening `/tmp/nx-caveman-005.html`,
> `/tmp/nx-caveman-015.html`, and `/tmp/nx-caveman-030.html`.*

At p=0.05, the layout should render 4 visually distinct dense hexagonal clusters (one per
6-node clique), with cave3 rendered as a fully separate connected component drifting apart from
the other three — matching its status as its own communicating class. A thin scattering of
cross-clique edges (8 of 120, 6.7%) should be visible connecting cave0/cave1/cave2, but cave3
should appear entirely isolated from them. At p=0.15, the layout should look similar but with
slightly more visible cross-clique edges (10 of 120) — still not enough to visually pull cave3
into the main mass. At p=0.30, cave3 should now visibly connect into the rest of the graph
(single communicating class), though its internal cohesion is weaker — `cave3-05`, having lost
intra-clique edges to rewiring, should appear as a comparatively small, loosely-attached node
even within its own cluster, while `cave2-01` (highest π, degree 7) should appear as the
visually largest node in the entire graph.

### Markov questions — answered

- **Does π rank clique-internal hubs highest at low p?**
  Yes — at p=0.05 and p=0.15, `cave1-02` leads π at 0.050000 in both cases, and it is a
  clique-internal node (degree 6, meaning it retained its full 5 intra-clique ties plus exactly
  1 rewired inbound edge) rather than a node that lost intra-clique ties.

- **How does the number of cross-community `suggest` recommendations change across the
  p-sweep?**
  At p=0.05 and p=0.15, `suggest` produces **zero** recommendations (output is `(none)`) because
  the graph remains split into 2 communicating classes and the tool finds no useful
  not-yet-linked candidates to propose. At p=0.30, once the graph merges into a single
  communicating class, `suggest` produces 5 recommendations per node (120 total across 24
  nodes), and for the two top-π nodes shown above, **all 10 of their combined suggestions
  (100%) cross clique boundaries** — a jump from 0 cross-community suggestions to a fully
  cross-community suggestion set once the classes merge.

- **Is π hub-dominated or flat at each p, and how does the max/min ratio trend as p rises?**
  p=0.05: ratio 1.5, effectively flat. p=0.15: ratio 1.5, effectively flat (unchanged). p=0.30:
  ratio 2.33, mildly skewed — the ratio rises as p increases from 0.15 to 0.30, but not
  monotonically smoothly from p=0.05 to 0.15 (it is flat across that range with this seed) —
  driven by rewiring concentrating extra edges on `cave2-01` (degree 7) while stripping edges
  from `cave3-05` (degree 3) at p=0.30.

- **At what p does the max/min π ratio flatten out, indicating communities have effectively
  merged?**
  The ratio does not flatten toward 1.0 within this sweep — it *rises* from 1.5 to 2.33 between
  p=0.15 and p=0.30. The more direct merge signal in this p-sweep is the **communicating-class
  count**: it collapses from 2 classes (p=0.05, p=0.15) to 1 class (p=0.30), which is the point
  at which the last isolated clique (cave3) first receives inter-clique connectivity in this
  particular seeded run. Within the tested range, p=0.30 is therefore the p-value at which
  communities have structurally merged (single communicating class), even though π's max/min
  ratio actually increases slightly at that same point rather than flattening — the two signals
  (class-count merge vs. π-ratio flattening) diverge here, and class-count is the more reliable
  merge indicator for this graph family.

---

## Sources

- NetworkX: https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.relaxed_caveman_graph.html
- Caveman Graph section from [[graph-topologies]]
- LFR μ section from [[graph-models]]

---

## Definition of Done

- [x] All sections above filled — no placeholder text remaining
- [x] Slug ordering algorithm documented with a node-id mapping table (not just the pattern)
- [x] `wikigraph analyze --suggest-top 5` output pasted verbatim for all three parameter sets
- [x] Cross-community edge count computed and recorded (Python helper, per parameter set)
- [x] All Markov questions answered with actual numbers from analyze output
- [x] Three parameter-set variants shown (p=0.05, 0.15, 0.30) with a comparison table before
      the per-set detail sections
- [x] Every reference link manually verified to resolve (NetworkX docs 200)
- [x] File committed to branch — path to be updated once #42 resolves
