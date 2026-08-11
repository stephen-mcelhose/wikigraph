# Zachary's Karate Club — two-community social network benchmark

> **nx-to-wiki flag**: `--graph karate`
> **Nodes**: 34 · **Directed edges**: 156 · **Naming**: tier1-karate-club-label (mhi-* / off-*)
> **Source**: empirical — fixed

---

## Background

Zachary's Karate Club is a social network of 34 members of a university karate club, observed by Wayne Zachary in the early 1970s. Zachary recorded which pairs of members interacted outside club activities, producing 78 friendships (undirected edges). The dataset is the canonical two-community benchmark because the club split into two factions — led by the instructor "Mr. Hi" (node 0) and the club officer (node 33) — and almost every community-detection algorithm is validated against its known ground truth.

**Origin**: Zachary, W.W. (1977). An information flow model for conflict and fission in small groups. *Journal of Anthropological Research*, 33(4), 452–473. https://doi.org/10.1086/jar.33.4.3629752

---

## Graph Properties

| Property                | Value / Description                                    |
| ----------------------- | ------------------------------------------------------ |
| Nodes                   | 34                                                     |
| Undirected edges        | 78                                                     |
| Directed edges (×2)     | 156                                                    |
| Degree distribution     | Right-skewed; min 1, max 17, most nodes degree 2–5     |
| Community structure     | 2 communities (Mr. Hi: 17 nodes, Officer: 17 nodes)    |
| Diameter                | 5                                                      |
| Average clustering coef | 0.5706                                                 |

---

## Transition Matrix Structure

For every graph produced by `nx-to-wiki`, the random walk is **uniform**: each step chooses a
neighbour with equal probability. The transition matrix $P$ is fully determined by the adjacency
structure of the `.md` files:

$$P_{ij} = \frac{1}{\deg(i)} \quad \text{if } [[j]] \text{ appears in i.md}, \quad \text{else } 0$$

For the karate club this gives exact rational values — no approximation:

| Node   | Degree | $P_{ij}$ (non-zero) | Example cross-community entries                         |
| ------ | ------ | ------------------- | ------------------------------------------------------- |
| off-16 | 17     | 1/17 ≈ 0.058824     | → mhi-08, mhi-12, mhi-15 (sole cross-community targets) |
| mhi-00 | 16     | 1/16 = 0.062500     | → off-14 (sole cross-community target)                  |
| off-15 | 12     | 1/12 ≈ 0.083333     | → mhi-02, mhi-08 (2 of 12 targets cross community)      |
| mhi-10 | 1      | 1/1 = 1.000000      | → mhi-00 only (degree-1 leaf; all weight on one edge)   |

**The `.md` files are a lossless encoding of $P$.** The `wikigraph analyze` text output is
intentionally lossy — it reports $\pi$ and top-5 commute times per node, not the 156 individual
$P_{ij}$ values. Use `wikigraph export` to obtain $P$ explicitly (see [[kernel-identifiability]]
for the full signal hierarchy).

### Export

```bash
# Sparse P — non-zero entries only (156 edge rows; all edges exceed the default --min-edge 0.005)
wikigraph export /tmp/nxwiki-karate --format csv -o /tmp/karate-export

# Dense P — full 34×34 matrix including structural zeros (1,156 edge rows)
wikigraph export /tmp/nxwiki-karate --format csv -o /tmp/karate-export --min-edge 0
```

Both commands write two files:

- `karate-export_nodes.csv` — 34 rows: `slug, π, class`
- `karate-export_edges.csv` — 156 rows (sparse) or 1,156 rows (dense): `source, target, probability`

Both representations are lossless for this graph because the minimum non-zero
$P_{ij} = 1/17 \approx 0.059$, well above the default filter threshold of 0.005.

---

## Slug Naming

**Naming tier**: Tier 1

**Assignment algorithm**: Each node carries a `club` attribute set by NetworkX — either `"Mr. Hi"` or `"Officer"`. Nodes with `club == "Mr. Hi"` are assigned prefix `mhi`; all others get prefix `off`. Within each group, node IDs are sorted ascending (integer order). The zero-padded index within the sorted group is then appended: `{prefix}-{index:02d}`. Both groups contain exactly 17 nodes, so the index runs from `00` to `16` and the zero-pad width is 2.

**Why this ordering**: Node 0 is the known Mr. Hi instructor in the NetworkX dataset — it is the lowest node-id in its group, so it becomes `mhi-00`. Node 33 is the known Officer instructor — highest node-id in the Officer group — so it becomes `off-16`. The algorithm makes no special case for these nodes; the ascending-sort rule produces the correct result automatically.

Example slugs and their node-id mappings (full table):

| Slug   | NetworkX node id | Club      | Degree |
| ------ | ---------------- | --------- | ------ |
| mhi-00 | 0                | Mr. Hi    | 16     |
| mhi-01 | 1                | Mr. Hi    | 9      |
| mhi-02 | 2                | Mr. Hi    | 10     |
| mhi-03 | 3                | Mr. Hi    | 6      |
| mhi-04 | 4                | Mr. Hi    | 3      |
| mhi-05 | 5                | Mr. Hi    | 4      |
| mhi-06 | 6                | Mr. Hi    | 4      |
| mhi-07 | 7                | Mr. Hi    | 4      |
| mhi-08 | 8                | Mr. Hi    | 5      |
| mhi-09 | 10               | Mr. Hi    | 3      |
| mhi-10 | 11               | Mr. Hi    | 1      |
| mhi-11 | 12               | Mr. Hi    | 2      |
| mhi-12 | 13               | Mr. Hi    | 5      |
| mhi-13 | 16               | Mr. Hi    | 2      |
| mhi-14 | 17               | Mr. Hi    | 2      |
| mhi-15 | 19               | Mr. Hi    | 3      |
| mhi-16 | 21               | Mr. Hi    | 2      |
| off-00 | 9                | Officer   | 2      |
| off-01 | 14               | Officer   | 2      |
| off-02 | 15               | Officer   | 2      |
| off-03 | 18               | Officer   | 2      |
| off-04 | 20               | Officer   | 2      |
| off-05 | 22               | Officer   | 2      |
| off-06 | 23               | Officer   | 5      |
| off-07 | 24               | Officer   | 3      |
| off-08 | 25               | Officer   | 3      |
| off-09 | 26               | Officer   | 2      |
| off-10 | 27               | Officer   | 4      |
| off-11 | 28               | Officer   | 3      |
| off-12 | 29               | Officer   | 4      |
| off-13 | 30               | Officer   | 4      |
| off-14 | 31               | Officer   | 6      |
| off-15 | 32               | Officer   | 12     |
| off-16 | 33               | Officer   | 17     |

Note: node IDs 9, 14–22, 26–33 map to `off-*` slugs. Node 9 (`off-00`) is the lowest Officer node-id. Node IDs are not contiguous within each group because NetworkX assigns node IDs globally (0–33), not per-community.

---

## How to Generate

```bash
# Minimal — default parameters (no knobs needed; this is a fixed empirical dataset)
python3 tools/nx-to-wiki/main.py --graph karate --out /tmp/nxwiki-karate
wikigraph graph /tmp/nxwiki-karate -o /tmp/nx-karate.html --title "Karate Club"
wikigraph analyze /tmp/nxwiki-karate --suggest-top 5
```

No parameter variants apply — the karate club graph is a fixed empirical dataset with no generative knobs. The `--graph karate` flag always produces the same 34-node, 78-edge graph.

**What the output directory contains:**

- 34 `.md` files, one per node (17 `mhi-*.md` + 17 `off-*.md`)
- Each file links to all its undirected neighbours (symmetric — every link is bidirectional)
- File layout: `# {slug}\n\n[[neighbour-1]] [[neighbour-2]] ...`
- Why directed edges = 2× undirected: `G.to_directed()` adds both u→v and v→u for every edge

---

## wikigraph Analysis

> **⚠️ Communicating classes are always trivial here.** All wikis produced by `nx-to-wiki` use
> `G.to_directed()` on a connected undirected graph, which yields a **strongly connected** directed
> graph. `wikigraph analyze` will report **one communicating class containing all nodes** for every
> graph in this series. This is expected and correct — not a bug. Do **not** frame analysis around
> communicating-class colouring as a community proxy; it will always be uniform.
>
> The meaningful Markov signal is in **stationary distribution (π)** and **suggested links**.
> Use cross-community edge count from the raw `.md` files as a structural boundary proxy if needed.

### Raw analyze output

```
=== Overview ===
Pages:        34
Edges:        156
Entropy rate: 2.5810 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 34 page(s)
  mhi-00
  mhi-01
  mhi-02
  mhi-03
  mhi-04
  mhi-05
  mhi-06
  mhi-07
  mhi-08
  mhi-09
  mhi-10
  mhi-11
  mhi-12
  mhi-13
  mhi-14
  mhi-15
  mhi-16
  off-00
  off-01
  off-02
  off-03
  off-04
  off-05
  off-06
  off-07
  off-08
  off-09
  off-10
  off-11
  off-12
  off-13
  off-14
  off-15
  off-16

=== Orphan pages (bottom 10% by stationary distribution) ===
  mhi-10                                    π=0.006410  → add inbound links
  off-09                                    π=0.012821  → add inbound links
  off-01                                    π=0.012821  → add inbound links
  off-02                                    π=0.012821  → add inbound links
  off-03                                    π=0.012821  → add inbound links
  off-04                                    π=0.012821  → add inbound links
  off-05                                    π=0.012821  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. off-16                                    π=0.108974
  2. mhi-00                                    π=0.102564
  3. off-15                                    π=0.076923
  4. mhi-02                                    π=0.064103
  5. mhi-01                                    π=0.057692

=== Suggested missing links (lowest commute time, not yet linked) ===
  mhi-00:
    → off-16                                  (commute: 39.59)
    → off-15                                  (commute: 46.62)
    → off-13                                  (commute: 64.13)
    → off-06                                  (commute: 71.57)
    → off-10                                  (commute: 74.53)
  mhi-01:
    → off-16                                  (commute: 44.20)
    → off-15                                  (commute: 51.40)
    → mhi-08                                  (commute: 58.19)
    → off-14                                  (commute: 64.67)
    → off-06                                  (commute: 77.08)
  mhi-02:
    → off-16                                  (commute: 33.01)
    → off-14                                  (commute: 53.06)
    → off-13                                  (commute: 60.99)
    → off-06                                  (commute: 63.26)
    → mhi-15                                  (commute: 73.18)
  mhi-03:
    → off-16                                  (commute: 56.96)
    → off-15                                  (commute: 64.12)
    → mhi-08                                  (commute: 70.64)
    → off-14                                  (commute: 75.96)
    → off-13                                  (commute: 81.24)
  mhi-04:
    → mhi-05                                  (commute: 97.16)
    → mhi-01                                  (commute: 102.64)
    → mhi-02                                  (commute: 104.92)
    → mhi-03                                  (commute: 111.54)
    → off-16                                  (commute: 112.12)
  mhi-05:
    → mhi-04                                  (commute: 97.16)
    → mhi-01                                  (commute: 98.54)
    → mhi-02                                  (commute: 100.81)
    → mhi-03                                  (commute: 107.44)
    → off-16                                  (commute: 108.01)
  mhi-06:
    → mhi-09                                  (commute: 97.16)
    → mhi-01                                  (commute: 98.54)
    → mhi-02                                  (commute: 100.81)
    → mhi-03                                  (commute: 107.44)
    → off-16                                  (commute: 108.01)
  mhi-07:
    → off-16                                  (commute: 68.40)
    → mhi-12                                  (commute: 71.38)
    → off-15                                  (commute: 75.07)
    → mhi-08                                  (commute: 81.36)
    → off-14                                  (commute: 87.00)
  mhi-08:
    → mhi-01                                  (commute: 58.19)
    → off-14                                  (commute: 68.12)
    → mhi-12                                  (commute: 68.89)
    → mhi-03                                  (commute: 70.64)
    → off-06                                  (commute: 75.84)
  mhi-09:
    → mhi-06                                  (commute: 97.16)
    → mhi-01                                  (commute: 102.64)
    → mhi-02                                  (commute: 104.92)
    → mhi-03                                  (commute: 111.54)
    → off-16                                  (commute: 112.12)
  mhi-10:
    → mhi-01                                  (commute: 186.12)
    → mhi-02                                  (commute: 188.39)
    → mhi-03                                  (commute: 195.02)
    → off-16                                  (commute: 195.59)
    → mhi-12                                  (commute: 199.49)
  mhi-11:
    → mhi-01                                  (commute: 104.93)
    → mhi-02                                  (commute: 106.28)
    → mhi-12                                  (commute: 115.92)
    → off-16                                  (commute: 116.52)
    → mhi-07                                  (commute: 121.64)
  mhi-12:
    → off-15                                  (commute: 59.78)
    → mhi-08                                  (commute: 68.89)
    → mhi-07                                  (commute: 71.38)
    → off-14                                  (commute: 73.14)
    → off-13                                  (commute: 78.67)
  mhi-13:
    → mhi-00                                  (commute: 130.00)
    → mhi-04                                  (commute: 150.53)
    → mhi-09                                  (commute: 150.53)
    → mhi-01                                  (commute: 160.12)
    → mhi-02                                  (commute: 162.39)
  mhi-14:
    → mhi-02                                  (commute: 104.72)
    → mhi-03                                  (commute: 111.60)
    → off-16                                  (commute: 112.37)
    → mhi-12                                  (commute: 115.22)
    → off-15                                  (commute: 119.48)
  mhi-15:
    → mhi-02                                  (commute: 73.18)
    → off-15                                  (commute: 79.41)
    → mhi-03                                  (commute: 85.75)
    → mhi-12                                  (commute: 85.85)
    → mhi-08                                  (commute: 89.64)
  mhi-16:
    → mhi-02                                  (commute: 104.72)
    → mhi-03                                  (commute: 111.60)
    → off-16                                  (commute: 112.37)
    → mhi-12                                  (commute: 115.22)
    → off-15                                  (commute: 119.48)
  off-00:
    → off-15                                  (commute: 99.97)
    → mhi-00                                  (commute: 105.74)
    → mhi-01                                  (commute: 109.90)
    → mhi-08                                  (commute: 114.30)
    → off-14                                  (commute: 116.90)
  off-01:
    → mhi-02                                  (commute: 108.09)
    → mhi-00                                  (commute: 115.56)
    → off-14                                  (commute: 116.14)
    → mhi-08                                  (commute: 116.47)
    → off-06                                  (commute: 118.64)
  off-02:
    → mhi-02                                  (commute: 108.09)
    → mhi-00                                  (commute: 115.56)
    → off-14                                  (commute: 116.14)
    → mhi-08                                  (commute: 116.47)
    → off-06                                  (commute: 118.64)
  off-03:
    → mhi-02                                  (commute: 108.09)
    → mhi-00                                  (commute: 115.56)
    → off-14                                  (commute: 116.14)
    → mhi-08                                  (commute: 116.47)
    → off-06                                  (commute: 118.64)
  off-04:
    → mhi-02                                  (commute: 108.09)
    → mhi-00                                  (commute: 115.56)
    → off-14                                  (commute: 116.14)
    → mhi-08                                  (commute: 116.47)
    → off-06                                  (commute: 118.64)
  off-05:
    → mhi-02                                  (commute: 108.09)
    → mhi-00                                  (commute: 115.56)
    → off-14                                  (commute: 116.14)
    → mhi-08                                  (commute: 116.47)
    → off-06                                  (commute: 118.64)
  off-06:
    → mhi-02                                  (commute: 63.26)
    → off-14                                  (commute: 64.80)
    → mhi-00                                  (commute: 71.57)
    → mhi-08                                  (commute: 75.84)
    → mhi-01                                  (commute: 77.08)
  off-07:
    → off-16                                  (commute: 84.37)
    → off-15                                  (commute: 90.80)
    → off-06                                  (commute: 94.38)
    → mhi-02                                  (commute: 95.46)
    → mhi-00                                  (commute: 102.26)
  off-08:
    → off-16                                  (commute: 82.86)
    → off-15                                  (commute: 87.87)
    → mhi-02                                  (commute: 96.81)
    → off-10                                  (commute: 98.30)
    → mhi-00                                  (commute: 102.30)
  off-09:
    → off-15                                  (commute: 104.09)
    → mhi-02                                  (commute: 119.46)
    → off-06                                  (commute: 119.90)
    → mhi-00                                  (commute: 126.54)
    → off-14                                  (commute: 126.69)
  off-10:
    → off-15                                  (commute: 62.90)
    → off-14                                  (commute: 71.08)
    → mhi-00                                  (commute: 74.53)
    → mhi-01                                  (commute: 79.95)
    → mhi-08                                  (commute: 82.46)
  off-11:
    → off-15                                  (commute: 73.38)
    → mhi-00                                  (commute: 80.00)
    → mhi-01                                  (commute: 86.18)
    → mhi-08                                  (commute: 90.26)
    → mhi-12                                  (commute: 93.68)
  off-12:
    → mhi-02                                  (commute: 75.66)
    → off-14                                  (commute: 81.88)
    → mhi-00                                  (commute: 83.24)
    → mhi-08                                  (commute: 85.40)
    → mhi-01                                  (commute: 88.15)
  off-13:
    → mhi-02                                  (commute: 60.99)
    → mhi-00                                  (commute: 64.13)
    → off-14                                  (commute: 77.54)
    → mhi-12                                  (commute: 78.67)
    → mhi-03                                  (commute: 81.24)
  off-14:
    → mhi-02                                  (commute: 53.06)
    → mhi-01                                  (commute: 64.67)
    → off-06                                  (commute: 64.80)
    → mhi-08                                  (commute: 68.12)
    → off-10                                  (commute: 71.08)
  off-15:
    → mhi-00                                  (commute: 46.62)
    → mhi-01                                  (commute: 51.40)
    → mhi-12                                  (commute: 59.78)
    → off-10                                  (commute: 62.90)
    → mhi-03                                  (commute: 64.12)
  off-16:
    → mhi-02                                  (commute: 33.01)
    → mhi-00                                  (commute: 39.59)
    → mhi-01                                  (commute: 44.20)
    → mhi-03                                  (commute: 56.96)
    → mhi-07                                  (commute: 68.40)
```

### Stationary distribution (π)

| Rank | Slug   | π value  | Expected structural role                          |
| ---- | ------ | -------- | ------------------------------------------------- |
| 1    | off-16 | 0.108974 | Officer instructor (node 33); highest degree (17) |
| 2    | mhi-00 | 0.102564 | Mr. Hi instructor (node 0); degree 16             |
| 3    | off-15 | 0.076923 | Officer-side hub (node 32); degree 12             |
| 4    | mhi-02 | 0.064103 | Mr. Hi-side connector (node 2); degree 10         |
| 5    | mhi-01 | 0.057692 | Mr. Hi-side connector (node 1); degree 9          |

**Finding**: Both known instructor hubs lead the stationary distribution exactly as expected. `off-16` (node 33, the Officer) ranks first at π=0.109 and `mhi-00` (node 0, Mr. Hi) ranks second at π=0.103, consistent with the issue ground-truth. The top-5 list also includes `off-15` (node 32), the second-highest-degree node in the Officer faction, confirming that π tracks degree-based importance in this graph.

### Suggested links

The table below shows suggestions for a representative set of nodes. Community is inferred from slug prefix: `mhi-*` = Mr. Hi, `off-*` = Officer.

| From   | To     | Commute time | Within / Cross community |
| ------ | ------ | ------------ | ------------------------ |
| mhi-00 | off-16 | 39.59        | Cross                    |
| mhi-00 | off-15 | 46.62        | Cross                    |
| mhi-00 | off-13 | 64.13        | Cross                    |
| mhi-00 | off-06 | 71.57        | Cross                    |
| mhi-00 | off-10 | 74.53        | Cross                    |
| off-16 | mhi-02 | 33.01        | Cross                    |
| off-16 | mhi-00 | 39.59        | Cross                    |
| off-16 | mhi-01 | 44.20        | Cross                    |
| off-16 | mhi-03 | 56.96        | Cross                    |
| off-16 | mhi-07 | 68.40        | Cross                    |
| off-15 | mhi-00 | 46.62        | Cross                    |
| off-15 | mhi-01 | 51.40        | Cross                    |
| off-15 | mhi-12 | 59.78        | Cross                    |
| off-15 | off-10 | 62.90        | Within                   |
| off-15 | mhi-03 | 64.12        | Cross                    |
| mhi-04 | mhi-05 | 97.16        | Within                   |
| mhi-04 | mhi-01 | 102.64       | Within                   |
| mhi-04 | mhi-02 | 104.92       | Within                   |
| mhi-04 | mhi-03 | 111.54       | Within                   |
| mhi-04 | off-16 | 112.12       | Cross                    |

**Finding**: For the two instructor hubs (`mhi-00` and `off-16`), every one of their 5 suggestions crosses the community boundary — reflecting that both nodes are already densely connected within their own faction and the shortest remaining random-walk path targets the other side. More peripheral nodes (e.g. `mhi-04`, degree 3) suggest mostly within-community links because their own community is harder to reach first. Overall, high-degree boundary nodes suggest cross-community links; peripheral nodes suggest within-community links.

### Cross-community edges (structural boundary proxy)

```bash
# Count mhi->off directed edges
grep -l "^# mhi-" /tmp/nxwiki-karate/*.md | xargs grep -o "\[\[off-[^]]*\]\]" | wc -l

# Count off->mhi directed edges
grep -l "^# off-" /tmp/nxwiki-karate/*.md | xargs grep -o "\[\[mhi-[^]]*\]\]" | wc -l
```

- **Cross-community directed edges**: 22 of 156 (14.1%)
- **Expected**: Low, consistent with two well-separated communities. The 11 undirected cross-community edges (Zachary's identified bridge ties) are reproduced exactly.

Community boundary nodes (nodes with at least one cross-community edge):

| Slug   | Node ID | Side    | Cross-community neighbours          |
| ------ | ------- | ------- | ----------------------------------- |
| mhi-00 | 0       | Mr. Hi  | off-14                              |
| mhi-01 | 1       | Mr. Hi  | off-13                              |
| mhi-02 | 2       | Mr. Hi  | off-00, off-10, off-11, off-15      |
| mhi-08 | 8       | Mr. Hi  | off-13, off-15, off-16              |
| mhi-12 | 13      | Mr. Hi  | off-16                              |
| mhi-15 | 19      | Mr. Hi  | off-16                              |
| off-00 | 9       | Officer | mhi-02                              |
| off-10 | 27      | Officer | mhi-02                              |
| off-11 | 28      | Officer | mhi-02                              |
| off-13 | 30      | Officer | mhi-01, mhi-08                      |
| off-14 | 31      | Officer | mhi-00                              |
| off-15 | 32      | Officer | mhi-02, mhi-08                      |
| off-16 | 33      | Officer | mhi-08, mhi-12, mhi-15              |

### Stationary distribution spread

- **Max π**: 0.108974 (off-16, node 33)
- **Min π**: 0.006410 (mhi-10, node 11 — degree-1 leaf)
- **Ratio max/min**: 17.0

**Finding**: Strongly hub-dominated. The ratio of 17 matches the degree of `off-16` (degree 17), which is not coincidental: for a random walk on a graph, the stationary probability of a node is proportional to its degree (π_i = deg_i / (2|E|)). The degree-1 leaf `mhi-10` (node 11) is visited proportionally rarely, while the degree-17 hub `off-16` is visited 17× more often.

### Visual observations

The force-directed HTML output (`wikigraph graph`) renders the karate club as two distinct clusters separated by a narrow bridge region. The largest nodes (sized by π) are `off-16` and `mhi-00`, which sit near the centres of their respective clusters and dominate the visual. `off-15` and `mhi-02` appear as the second-tier nodes within the Officer and Mr. Hi clusters respectively. The 13 community-boundary nodes (listed above) appear at the cluster peripheries, with `off-16`, `off-15`, `mhi-00`, and `mhi-02` positioned centrally — all four are visibly larger than surrounding nodes due to their higher π. The degree-1 node `mhi-10` (node 11) appears at the Mr. Hi cluster periphery as a very small node with a single link.

### Markov questions — answered

- **Does π rank the expected hub nodes highest?**
  Yes — `off-16` (node 33, Officer instructor) has π=0.108974 (rank 1) and `mhi-00` (node 0, Mr. Hi instructor) has π=0.102564 (rank 2). Both lead, matching the ground-truth stated in the issue.

- **Does `suggest` propose cross-community links?**
  Yes — for both hub nodes (`mhi-00` and `off-16`), all 5 suggestions cross the mhi/off community boundary. Across all 34 nodes, the fraction of cross-community suggestions varies: hubs suggest predominantly cross-community links; peripheral nodes (e.g. `mhi-04`, `mhi-13`) suggest mostly within-community links.

- **Is the stationary distribution hub-dominated or flat?**
  Hub-dominated — max/min ratio = 17.0. This is consistent with a degree-proportional stationary distribution (π_i ∝ deg_i), which holds for any random walk on an undirected graph.

- **Do both `mhi-00` and `off-16` outrank all other nodes in π, matching their known instructor-hub roles?**
  Yes — `off-16` ranks 1st (π=0.109) and `mhi-00` ranks 2nd (π=0.103). No other node comes close; the 3rd-ranked node `off-15` has π=0.077, a gap of 0.026 below `mhi-00`.

---

## References

- Zachary, W.W. (1977). An information flow model for conflict and fission in small groups. *Journal of Anthropological Research*, 33(4), 452–473. https://doi.org/10.1086/jar.33.4.3629752
- NetworkX: https://networkx.org/documentation/stable/reference/generated/networkx.generators.social.karate_club_graph.html
- Graph topology notes: Planted Partition / SBM section from graph-topologies research notes
- Graph model notes: Empirical Social Networks section from graph-models research notes

---

## Definition of Done

- [x] All sections above filled — no `[placeholder]` text remaining
- [x] Slug ordering algorithm documented with a node-id mapping table (not just the pattern)
- [x] `wikigraph analyze --suggest-top 5` output pasted verbatim
- [x] Cross-community edge count computed and recorded
- [x] All Markov questions answered with actual numbers from analyze output
- [x] At least one variant command shown (parameterised graphs) or a note explaining why none apply (fixed graphs)
- [x] Every reference link manually verified to resolve (DOIs, NetworkX docs)
- [x] File committed to branch — path to be updated once #42 resolves
