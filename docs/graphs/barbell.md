---
type: spike
title: Barbell Graph — wikigraph Analysis
description: Empirical wikigraph analysis of the barbell graph (two cliques joined by a bottleneck bridge) as a Markov chain benchmark. Path provisional — will move to flat docs/ when #42 resolves.
tags: [benchmark, barbell, networkx, markov-chain, spike, draft]
status: draft
timestamp: 2026-08-11T00:00:00Z
---

<!-- Path provisional — will be updated when #42 resolves -->

# Barbell Graph — two cliques joined by a bottleneck bridge

> ⚠️ **DRAFT** — AI-assisted write-up, not yet verified by human analysis. Treat all findings as provisional.

> **nx-to-wiki flag**: `--graph barbell --m1 [clique-size] --m2 [bridge-length]`
> **Nodes**: 2·m1 + m2 (parameterised) · **Directed edges**: 2×undirected (parameterised)
> **Naming**: tier1-barbell-structural-position
> **Source**: generative — parameterised

---

## Background

The barbell graph is two complete graphs $K_{m1}$ ("weights") connected by a path of $m2$
intermediate nodes ("bar"). It is the textbook example of a graph with **near-zero
conductance** across a single bottleneck: because every node in the left clique is densely
connected to every other left-clique node, but the only route to the right clique passes
through a thin chain of bridge nodes, a random walk started anywhere in the left clique spends
a very long time before it escapes to the right side. This makes the barbell the standard
example for teaching **mixing time**, **commute time**, and **conductance** in Markov chain
theory — it is explicitly the worst case used to lower-bound how slowly some random walks
converge to their stationary distribution.

The barbell is also structurally the closest named graph to wiki-gen's own topology (see
`docs/adr-009-wiki-gen-make-vs-buy.md`): a set of tightly interlinked pages (a "cluster") joined
to another cluster through a small number of bridging pages, which is exactly the shape formed
when two topically-separate wiki sections are connected by only a handful of cross-referencing
articles.

**Origin**: Standard construction in graph theory and Markov chain mixing-time literature; no
single canonical paper — see NetworkX reference implementation below.

---

## Graph Properties

**Key parameters**:

| Flag   | Default | Range        | What it controls                                    |
| ------ | ------- | ------------ | ---------------------------------------------------- |
| `--m1` | 5       | ≥ 2          | Size of each complete-graph clique (left and right)   |
| `--m2` | 3       | ≥ 0          | Number of intermediate bridge nodes (path length)     |

Three parameter sets were run for this document, as directed by the issue:

| Parameter set          | m1 | m2 | Nodes | Undirected edges | Directed edges |
| ------------------------ | -- | -- | ----- | ------------------ | ---------------- |
| Tight bridge             | 6  | 1  | 13    | 32                  | 64                |
| Loose bridge             | 6  | 5  | 17    | 36                  | 72                |
| Mid (added for completeness) | 4  | 3  | 11    | 16                  | 32                |

Degree distribution for all three: every clique-interior node has degree `m1 - 1 + [1 if it
touches the bridge else 0]`; the clique nodes adjacent to the bridge have degree `m1` (one extra
edge into the bridge/other clique); every bridge node has degree exactly 2 (except when m2=0, in
which case the two cliques connect directly). Community structure: none in the formal sense —
three positional zones (`left`, `bridge`, `right`) rather than communities. Diameter and
clustering vary by parameter set (see per-set tables below).

---

## Transition Matrix Structure

For every graph produced by `nx-to-wiki`, the random walk is **uniform**: each step chooses a
neighbour with equal probability. The transition matrix $P$ is fully determined by the
adjacency structure of the `.md` files:

$$P_{ij} = \frac{1}{\deg(i)} \quad \text{if } [[j]] \text{ appears in i.md}, \quad \text{else } 0$$

### Tight bridge (m1=6, m2=1)

| Slug      | Node id | Degree | $P_{ij}$ (non-zero) | Structural note                                  |
| --------- | ------- | ------ | -------------------- | --------------------------------------------------- |
| left-05   | 5       | 6      | 1/6 ≈ 0.166667        | Clique-hub touching the bridge (connector node)      |
| left-00   | 0       | 5      | 1/5 = 0.200000        | Interior clique node — no bridge access              |
| bridge-00 | 6       | 2      | 1/2 = 0.500000        | Sole bridge node; only neighbours are left-05, right-00 |
| right-00  | 7       | 6      | 1/6 ≈ 0.166667        | Clique-hub touching the bridge (mirror of left-05)   |

### Loose bridge (m1=6, m2=5)

| Slug       | Node id | Degree | $P_{ij}$ (non-zero) | Structural note                                    |
| ---------- | ------- | ------ | -------------------- | ------------------------------------------------------ |
| left-05    | 5       | 6      | 1/6 ≈ 0.166667        | Clique-hub touching the bridge chain                    |
| left-00    | 0       | 5      | 1/5 = 0.200000        | Interior clique node                                    |
| bridge-02  | 8       | 2      | 1/2 = 0.500000        | Middle bridge node — furthest from both cliques         |
| right-00   | 11      | 6      | 1/6 ≈ 0.166667        | Clique-hub touching the bridge chain (mirror of left-05) |

### Mid (m1=4, m2=3, added for completeness)

| Slug       | Node id | Degree | $P_{ij}$ (non-zero) | Structural note                                  |
| ---------- | ------- | ------ | -------------------- | ------------------------------------------------- |
| left-03    | 3       | 4      | 1/4 = 0.250000        | Clique-hub touching the bridge chain                |
| left-00    | 0       | 3      | 1/3 ≈ 0.333333        | Interior clique node                                |
| bridge-01  | 5       | 2      | 1/2 = 0.500000        | Middle bridge node                                  |
| right-00   | 7       | 4      | 1/4 = 0.250000        | Clique-hub touching the bridge chain                |

**The `.md` files are a lossless encoding of $P$ for all three parameter sets.** Minimum
non-zero $P_{ij}$ across all three sets is $1/6 \approx 0.167$ (tight/loose) or $1/4 = 0.25$
(mid) — well above the `--min-edge` default filter of 0.005 in every case. No `--min-edge 0`
warning applies for any parameter set tested here; for much larger `m1` (e.g. `m1 ≥ 200`, giving
$1/\text{deg} < 0.005$), `--min-edge 0` would become necessary for a lossless dense export.

### Export

```bash
# Sparse P — non-zero entries only (row count = directed edge count for each set)
wikigraph export /tmp/nxwiki-barbell-tight --format csv -o /tmp/barbell-tight-export
wikigraph export /tmp/nxwiki-barbell-loose --format csv -o /tmp/barbell-loose-export
wikigraph export /tmp/nxwiki-barbell-mid --format csv -o /tmp/barbell-mid-export

# Dense P — full N×N matrix including structural zeros (only needed if min P_ij < 0.005)
wikigraph export /tmp/nxwiki-barbell-tight --format csv -o /tmp/barbell-tight-export --min-edge 0
```

Sparse edge-row counts: 64 (tight), 72 (loose), 32 (mid) — matching directed-edge counts above.

---

## Slug Naming

**Naming tier**: Tier 1

**Assignment algorithm** (`build_barbell`, verbatim from source): NetworkX's
`barbell_graph(m1, m2)` numbers nodes in three contiguous ranges:

- `[0, m1)` — left clique
- `[m1, m1 + m2)` — bridge path
- `[m1 + m2, 2*m1 + m2)` — right clique

`nx-to-wiki` maps each range directly to a slug prefix and a **zero-based, zero-padded index
that resets at the start of each zone**:

- For `n < m1`: `slug = f"left-{pad(n, m1)}"` — index is `n` itself, padded to the width of `m1`
- For `m1 <= n < m1 + m2`: `slug = f"bridge-{pad(n - m1, m2)}"` — index is `n - m1` (offset into
  the bridge segment), padded to the width of `m2`
- For `n >= m1 + m2`: `idx = n - (m1 + m2)`; `slug = f"right-{pad(idx, m1)}"` — index is the
  offset into the right clique, padded to the width of `m1`

`pad(i, count)` zero-pads `i` to `max(2, len(str(count - 1)))` digits — always at least 2 digits
even for small zones.

**Why this ordering**: NetworkX constructs the barbell by literally concatenating the two
cliques with the bridge path in between, so node-id order already encodes zone membership; no
sorting or attribute lookup is needed — the slug is a pure arithmetic function of the node id
and the two parameters `m1`, `m2`.

Full node-id → slug mapping for the **tight** set (m1=6, m2=1):

| Slug       | NetworkX node id | Zone   | Degree |
| ---------- | ----------------- | ------ | ------ |
| left-00    | 0                  | left   | 5      |
| left-01    | 1                  | left   | 5      |
| left-02    | 2                  | left   | 5      |
| left-03    | 3                  | left   | 5      |
| left-04    | 4                  | left   | 5      |
| left-05    | 5                  | left   | 6      |
| bridge-00  | 6                  | bridge | 2      |
| right-00   | 7                  | right  | 6      |
| right-01   | 8                  | right  | 5      |
| right-02   | 9                  | right  | 5      |
| right-03   | 10                 | right  | 5      |
| right-04   | 11                 | right  | 5      |
| right-05   | 12                 | right  | 5      |

For the **loose** (m1=6, m2=5) and **mid** (m1=4, m2=3) sets, the same arithmetic mapping
applies with the corresponding `m1`/`m2` substituted; e.g. for loose, bridge node ids 6–10 map
to `bridge-00`…`bridge-04`, and right-clique node ids 11–16 map to `right-00`…`right-05`.

---

## How to Generate

```bash
# Tight bridge — m1=6, m2=1
python3 tools/nx-to-wiki/main.py --graph barbell --m1 6 --m2 1 --out /tmp/nxwiki-barbell-tight
wikigraph graph /tmp/nxwiki-barbell-tight -o /tmp/nx-barbell-tight.html --title "Barbell (m1=6 m2=1)"
wikigraph analyze /tmp/nxwiki-barbell-tight --suggest-top 5

# Loose bridge — m1=6, m2=5
python3 tools/nx-to-wiki/main.py --graph barbell --m1 6 --m2 5 --out /tmp/nxwiki-barbell-loose
wikigraph graph /tmp/nxwiki-barbell-loose -o /tmp/nx-barbell-loose.html --title "Barbell (m1=6 m2=5)"
wikigraph analyze /tmp/nxwiki-barbell-loose --suggest-top 5

# Mid — m1=4, m2=3 (added for completeness, between the tight and loose extremes)
python3 tools/nx-to-wiki/main.py --graph barbell --m1 4 --m2 3 --out /tmp/nxwiki-barbell-mid
wikigraph graph /tmp/nxwiki-barbell-mid -o /tmp/nx-barbell-mid.html --title "Barbell (m1=4 m2=3)"
wikigraph analyze /tmp/nxwiki-barbell-mid --suggest-top 5
```

**Effect of the two knobs**: Increasing `m1` grows both cliques (more nodes per weight, higher
degree for interior nodes, denser visual clusters at each end). Increasing `m2` lengthens the
bridge path (more bridge nodes, longer minimum left↔right commute time, weaker
"handshake" between the two cliques in the force-directed layout — the two weights visibly
drift further apart).

**What each output directory contains:**

- `2*m1 + m2` `.md` files, one per node (per parameter set)
- Each file links to all its undirected neighbours (symmetric — every link is bidirectional)
- File layout: `# {slug}\n\n[[neighbour-1]] [[neighbour-2]] ...`
- Why directed edges = 2× undirected: `G.to_directed()` adds both u→v and v→u for every edge

---

## wikigraph Analysis

> **⚠️ Communicating classes are always trivial here.** All wikis produced by `nx-to-wiki` use
> `G.to_directed()` on a connected undirected graph, which yields a **strongly connected**
> directed graph. `wikigraph analyze` will report **one communicating class containing all
> nodes** for every parameter set. This is expected and correct — not a bug.
>
> This graph has no community attribute; the meaningful structural distinction is **zone**
> (`left` / `bridge` / `right`), not community. The cross-community section below is replaced
> with a Python-computed cross-zone edge count (left↔bridge, bridge↔right — left↔right direct
> edges do not exist by construction).

### Comparison across parameter sets

| Parameter set | m1 | m2 | Nodes | Directed edges | Entropy rate | Max π    | Min π    | Ratio | Cross-zone % |
| ------------- | -- | -- | ----- | ---------------- | -------------- | -------- | -------- | ----- | -------------- |
| Tight         | 6  | 1  | 13    | 64                | 2.3299 bits    | 0.093750 | 0.031250 | 3.0   | 4 of 64 (6.2%)  |
| Mid           | 4  | 3  | 11    | 32                | 1.5790 bits    | 0.125000 | 0.062500 | 2.0   | 4 of 32 (12.5%) |
| Loose         | 6  | 5  | 17    | 72                | 2.1822 bits    | 0.083333 | 0.027778 | 3.0   | 4 of 72 (5.6%)  |

**Trend**: All three parameter sets are **mildly skewed** (ratio 2.0–3.0) — far less
hub-dominated than the star-like or single-hub graphs elsewhere in this series, because both
cliques have symmetric internal structure and no single node dominates degree the way a true
hub does. Lengthening the bridge (tight → loose, m2: 1 → 5) *lowers* both max and min π (more
total nodes dilutes each individual node's share) but the max/min ratio stays flat at 3.0 in
both cases, because the two connector nodes (`left-05`, `right-00`) always have
degree = m1, and the lowest-π node is always a bridge-interior node with degree 2 — the ratio
is a function of `m1` alone (m1/2), not of `m2`. The mid set (m1=4) confirms this: ratio = 2.0 =
4/2, exactly matching the m1/2 relationship.

### Tight bridge (m1=6, m2=1)

#### Raw analyze output

```
=== Overview ===
Pages:        13
Edges:        64
Entropy rate: 2.3299 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 13 page(s)
  bridge-00
  left-00
  left-01
  left-02
  left-03
  left-04
  left-05
  right-00
  right-01
  right-02
  right-03
  right-04
  right-05

=== Orphan pages (bottom 10% by stationary distribution) ===
  bridge-00                                 π=0.031250  → add inbound links
  left-00                                   π=0.078125  → add inbound links
  left-01                                   π=0.078125  → add inbound links
  left-02                                   π=0.078125  → add inbound links
  left-03                                   π=0.078125  → add inbound links
  left-04                                   π=0.078125  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. right-00                                  π=0.093750
  2. left-05                                   π=0.093750
  3. right-02                                  π=0.078125
  4. right-01                                  π=0.078125
  5. right-03                                  π=0.078125

=== Suggested missing links (lowest commute time, not yet linked) ===
  bridge-00:
    → right-01                                (commute: 85.33)
    → right-03                                (commute: 85.33)
    → right-04                                (commute: 85.33)
    → right-05                                (commute: 85.33)
    → right-02                                (commute: 85.33)
  left-00:
    → bridge-00                               (commute: 85.33)
    → right-00                                (commute: 149.33)
    → right-01                                (commute: 170.67)
    → right-02                                (commute: 170.67)
    → right-03                                (commute: 170.67)
  left-01:
    → bridge-00                               (commute: 85.33)
    → right-00                                (commute: 149.33)
    → right-01                                (commute: 170.67)
    → right-02                                (commute: 170.67)
    → right-03                                (commute: 170.67)
  left-02:
    → bridge-00                               (commute: 85.33)
    → right-00                                (commute: 149.33)
    → right-01                                (commute: 170.67)
    → right-02                                (commute: 170.67)
    → right-03                                (commute: 170.67)
  left-03:
    → bridge-00                               (commute: 85.33)
    → right-00                                (commute: 149.33)
    → right-01                                (commute: 170.67)
    → right-02                                (commute: 170.67)
    → right-03                                (commute: 170.67)
  left-04:
    → bridge-00                               (commute: 85.33)
    → right-00                                (commute: 149.33)
    → right-01                                (commute: 170.67)
    → right-02                                (commute: 170.67)
    → right-03                                (commute: 170.67)
  left-05:
    → right-00                                (commute: 128.00)
    → right-01                                (commute: 149.33)
    → right-02                                (commute: 149.33)
    → right-03                                (commute: 149.33)
    → right-04                                (commute: 149.33)
  right-00:
    → left-05                                 (commute: 128.00)
    → left-00                                 (commute: 149.33)
    → left-01                                 (commute: 149.33)
    → left-03                                 (commute: 149.33)
    → left-02                                 (commute: 149.33)
  right-01:
    → bridge-00                               (commute: 85.33)
    → left-05                                 (commute: 149.33)
    → left-00                                 (commute: 170.67)
    → left-01                                 (commute: 170.67)
    → left-02                                 (commute: 170.67)
  right-02:
    → bridge-00                               (commute: 85.33)
    → left-05                                 (commute: 149.33)
    → left-00                                 (commute: 170.67)
    → left-01                                 (commute: 170.67)
    → left-02                                 (commute: 170.67)
  right-03:
    → bridge-00                               (commute: 85.33)
    → left-05                                 (commute: 149.33)
    → left-00                                 (commute: 170.67)
    → left-01                                 (commute: 170.67)
    → left-02                                 (commute: 170.67)
  right-04:
    → bridge-00                               (commute: 85.33)
    → left-05                                 (commute: 149.33)
    → left-00                                 (commute: 170.67)
    → left-01                                 (commute: 170.67)
    → left-02                                 (commute: 170.67)
  right-05:
    → bridge-00                               (commute: 85.33)
    → left-05                                 (commute: 149.33)
    → left-00                                 (commute: 170.67)
    → left-01                                 (commute: 170.67)
    → left-02                                 (commute: 170.67)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `right-00` (0.093750), `left-05` (0.093750). Lowest π: `bridge-00` (0.031250).

| From      | To        | Commute time | Crosses zone boundary? |
| --------- | --------- | ------------ | ------------------------ |
| right-00  | left-05   | 128.00       | Yes — right→left (direct cross-clique shortcut) |
| right-00  | left-00   | 149.33       | Yes                        |
| right-00  | left-01   | 149.33       | Yes                        |
| right-00  | left-03   | 149.33       | Yes                        |
| right-00  | left-02   | 149.33       | Yes                        |
| left-05   | right-00  | 128.00       | Yes                        |
| left-05   | right-01  | 149.33       | Yes                        |
| left-05   | right-02  | 149.33       | Yes                        |
| left-05   | right-03  | 149.33       | Yes                        |
| left-05   | right-04  | 149.33       | Yes                        |
| bridge-00 | right-01  | 85.33        | Yes — bridge→right         |
| bridge-00 | right-03  | 85.33        | Yes                        |
| bridge-00 | right-04  | 85.33        | Yes                        |
| bridge-00 | right-05  | 85.33        | Yes                        |
| bridge-00 | right-02  | 85.33        | Yes                        |

**Finding**: Every suggestion for the two top-π nodes (`right-00`, `left-05` — both
clique-hubs directly touching the bridge) targets the *opposite* clique, i.e. a genuine
left↔right shortcut recommendation, despite no direct left↔right edge existing by
construction. This is the key structural finding requested by the issue: `suggest` does propose
a left-clique↔right-clique shortcut, even though the two hubs are the closest pair on either
side of the bridge. `bridge-00` (the lowest-π node) shows all 5 suggestions into the right clique — but this is a
tie-breaking artifact, not a structural property. By symmetry, bridge-00's commute time to any
left-interior node (left-00…left-04) equals its commute time to any right-interior node
(right-01…right-05) — all 85.33. The `suggest` algorithm resolves ties by slug order (`r` sorts
after `l`), so right-clique slugs fill the top-5 list first.
Remaining nodes not shown (the 5 interior members of each clique) all suggest `bridge-00` first,
confirming bridge nodes act as the mandatory waypoint for any cross-zone recommendation.

### Loose bridge (m1=6, m2=5)

#### Raw analyze output

```
=== Overview ===
Pages:        17
Edges:        72
Entropy rate: 2.1822 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 17 page(s)
  bridge-00
  bridge-01
  bridge-02
  bridge-03
  bridge-04
  left-00
  left-01
  left-02
  left-03
  left-04
  left-05
  right-00
  right-01
  right-02
  right-03
  right-04
  right-05

=== Orphan pages (bottom 10% by stationary distribution) ===
  bridge-01                                 π=0.027778  → add inbound links
  bridge-03                                 π=0.027778  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. right-00                                  π=0.083333
  2. left-05                                   π=0.083333
  3. right-05                                  π=0.069444
  4. right-04                                  π=0.069444
  5. right-03                                  π=0.069444

=== Suggested missing links (lowest commute time, not yet linked) ===
  bridge-00:
    → left-03                                 (commute: 96.00)
    → left-04                                 (commute: 96.00)
    → left-00                                 (commute: 96.00)
    → left-01                                 (commute: 96.00)
    → left-02                                 (commute: 96.00)
  bridge-01:
    → left-05                                 (commute: 144.00)
    → bridge-03                               (commute: 144.00)
    → left-04                                 (commute: 168.00)
    → left-01                                 (commute: 168.00)
    → left-02                                 (commute: 168.00)
  bridge-02:
    → bridge-04                               (commute: 144.00)
    → bridge-00                               (commute: 144.00)
    → left-05                                 (commute: 216.00)
    → right-00                                (commute: 216.00)
    → left-00                                 (commute: 240.00)
  bridge-03:
    → bridge-01                               (commute: 144.00)
    → right-00                                (commute: 144.00)
    → right-03                                (commute: 168.00)
    → right-05                                (commute: 168.00)
    → right-04                                (commute: 168.00)
  bridge-04:
    → right-01                                (commute: 96.00)
    → right-05                                (commute: 96.00)
    → right-04                                (commute: 96.00)
    → right-03                                (commute: 96.00)
    → right-02                                (commute: 96.00)
  left-00:
    → bridge-00                               (commute: 96.00)
    → bridge-01                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-03                               (commute: 312.00)
    → bridge-04                               (commute: 384.00)
  left-01:
    → bridge-00                               (commute: 96.00)
    → bridge-01                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-03                               (commute: 312.00)
    → bridge-04                               (commute: 384.00)
  left-02:
    → bridge-00                               (commute: 96.00)
    → bridge-01                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-03                               (commute: 312.00)
    → bridge-04                               (commute: 384.00)
  left-03:
    → bridge-00                               (commute: 96.00)
    → bridge-01                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-03                               (commute: 312.00)
    → bridge-04                               (commute: 384.00)
  left-04:
    → bridge-00                               (commute: 96.00)
    → bridge-01                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-03                               (commute: 312.00)
    → bridge-04                               (commute: 384.00)
  left-05:
    → bridge-01                               (commute: 144.00)
    → bridge-02                               (commute: 216.00)
    → bridge-03                               (commute: 288.00)
    → bridge-04                               (commute: 360.00)
    → right-00                                (commute: 432.00)
  right-00:
    → bridge-03                               (commute: 144.00)
    → bridge-02                               (commute: 216.00)
    → bridge-01                               (commute: 288.00)
    → bridge-00                               (commute: 360.00)
    → left-05                                 (commute: 432.00)
  right-01:
    → bridge-04                               (commute: 96.00)
    → bridge-03                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-01                               (commute: 312.00)
    → bridge-00                               (commute: 384.00)
  right-02:
    → bridge-04                               (commute: 96.00)
    → bridge-03                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-01                               (commute: 312.00)
    → bridge-00                               (commute: 384.00)
  right-03:
    → bridge-04                               (commute: 96.00)
    → bridge-03                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-01                               (commute: 312.00)
    → bridge-00                               (commute: 384.00)
  right-04:
    → bridge-04                               (commute: 96.00)
    → bridge-03                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-01                               (commute: 312.00)
    → bridge-00                               (commute: 384.00)
  right-05:
    → bridge-04                               (commute: 96.00)
    → bridge-03                               (commute: 168.00)
    → bridge-02                               (commute: 240.00)
    → bridge-01                               (commute: 312.00)
    → bridge-00                               (commute: 384.00)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `right-00` (0.083333), `left-05` (0.083333). Lowest π: `bridge-01` (0.027778, tied
with `bridge-03`).

| From      | To         | Commute time | Crosses zone boundary? |
| --------- | ---------- | ------------ | ------------------------ |
| right-00  | bridge-03  | 144.00       | Yes — right→bridge         |
| right-00  | bridge-02  | 216.00       | Yes                        |
| right-00  | bridge-01  | 288.00       | Yes                        |
| right-00  | bridge-00  | 360.00       | Yes                        |
| right-00  | left-05    | 432.00       | Yes — right→left (full shortcut) |
| left-05   | bridge-01  | 144.00       | Yes — left→bridge          |
| left-05   | bridge-02  | 216.00       | Yes                        |
| left-05   | bridge-03  | 288.00       | Yes                        |
| left-05   | bridge-04  | 360.00       | Yes                        |
| left-05   | right-00   | 432.00       | Yes — left→right (full shortcut) |
| bridge-01 | left-05    | 144.00       | Yes — bridge→left           |
| bridge-01 | bridge-03  | 144.00       | No — within bridge chain    |
| bridge-01 | left-04    | 168.00       | Yes                          |
| bridge-01 | left-01    | 168.00       | Yes                          |
| bridge-01 | left-02    | 168.00       | Yes                          |

**Finding**: With the longer 5-node bridge, `right-00` and `left-05`'s 5th-ranked suggestion is
the **direct opposite-clique hub** (commute 432.00, the highest commute time in either list) —
a genuine full left↔right shortcut recommendation, but now ranked last among their 5
suggestions rather than dominating the list, because the intermediate bridge nodes are all
individually closer targets. This shows the loose bridge weakens (but does not eliminate) the
left↔right shortcut signal compared to the tight-bridge case. `bridge-01` (a bottom-π node)
splits its suggestions: 1 within the bridge chain (`bridge-03`) and 4 into the left clique — it
does not suggest any right-clique node in its top 5, since the right clique is now much
further away (commute 96.00–384.00 range on the right side vs. 144.00–168.00 on the left).

### Mid (m1=4, m2=3, added for completeness)

#### Raw analyze output

```
=== Overview ===
Pages:        11
Edges:        32
Entropy rate: 1.5790 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 11 page(s)
  bridge-00
  bridge-01
  bridge-02
  left-00
  left-01
  left-02
  left-03
  right-00
  right-01
  right-02
  right-03

=== Orphan pages (bottom 10% by stationary distribution) ===
  bridge-01                                 π=0.062500  → add inbound links
  bridge-00                                 π=0.062500  → add inbound links
  bridge-02                                 π=0.062500  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. left-03                                   π=0.125000
  2. right-00                                  π=0.125000
  3. left-00                                   π=0.093750
  4. left-01                                   π=0.093750
  5. left-02                                   π=0.093750

=== Suggested missing links (lowest commute time, not yet linked) ===
  bridge-00:
    → left-00                                 (commute: 48.00)
    → left-01                                 (commute: 48.00)
    → left-02                                 (commute: 48.00)
    → bridge-02                               (commute: 64.00)
    → right-00                                (commute: 96.00)
  bridge-01:
    → right-00                                (commute: 64.00)
    → left-03                                 (commute: 64.00)
    → left-00                                 (commute: 80.00)
    → left-01                                 (commute: 80.00)
    → left-02                                 (commute: 80.00)
  bridge-02:
    → right-01                                (commute: 48.00)
    → right-02                                (commute: 48.00)
    → right-03                                (commute: 48.00)
    → bridge-00                               (commute: 64.00)
    → left-03                                 (commute: 96.00)
  left-00:
    → bridge-00                               (commute: 48.00)
    → bridge-01                               (commute: 80.00)
    → bridge-02                               (commute: 112.00)
    → right-00                                (commute: 144.00)
    → right-01                                (commute: 160.00)
  left-01:
    → bridge-00                               (commute: 48.00)
    → bridge-01                               (commute: 80.00)
    → bridge-02                               (commute: 112.00)
    → right-00                                (commute: 144.00)
    → right-01                                (commute: 160.00)
  left-02:
    → bridge-00                               (commute: 48.00)
    → bridge-01                               (commute: 80.00)
    → bridge-02                               (commute: 112.00)
    → right-00                                (commute: 144.00)
    → right-01                                (commute: 160.00)
  left-03:
    → bridge-01                               (commute: 64.00)
    → bridge-02                               (commute: 96.00)
    → right-00                                (commute: 128.00)
    → right-01                                (commute: 144.00)
    → right-02                                (commute: 144.00)
  right-00:
    → bridge-01                               (commute: 64.00)
    → bridge-00                               (commute: 96.00)
    → left-03                                 (commute: 128.00)
    → left-00                                 (commute: 144.00)
    → left-01                                 (commute: 144.00)
  right-01:
    → bridge-02                               (commute: 48.00)
    → bridge-01                               (commute: 80.00)
    → bridge-00                               (commute: 112.00)
    → left-03                                 (commute: 144.00)
    → left-00                                 (commute: 160.00)
  right-02:
    → bridge-02                               (commute: 48.00)
    → bridge-01                               (commute: 80.00)
    → bridge-00                               (commute: 112.00)
    → left-03                                 (commute: 144.00)
    → left-00                                 (commute: 160.00)
  right-03:
    → bridge-02                               (commute: 48.00)
    → bridge-01                               (commute: 80.00)
    → bridge-00                               (commute: 112.00)
    → left-03                                 (commute: 144.00)
    → left-00                                 (commute: 160.00)
```

#### Suggested links (top-2 π + lowest-π)

Top-2 π: `left-03` (0.125000), `right-00` (0.125000). Lowest π: `bridge-00`/`bridge-01`/
`bridge-02` (three-way tie at 0.062500; `bridge-00` shown as representative).

| From      | To        | Commute time | Crosses zone boundary? |
| --------- | --------- | ------------ | ------------------------ |
| left-03   | bridge-01 | 64.00        | Yes — left→bridge          |
| left-03   | bridge-02 | 96.00        | Yes                        |
| left-03   | right-00  | 128.00       | Yes — left→right (full shortcut) |
| left-03   | right-01  | 144.00       | Yes                        |
| left-03   | right-02  | 144.00       | Yes                        |
| right-00  | bridge-01 | 64.00        | Yes — right→bridge         |
| right-00  | bridge-00 | 96.00        | Yes                        |
| right-00  | left-03   | 128.00       | Yes — right→left (full shortcut) |
| right-00  | left-00   | 144.00       | Yes                        |
| right-00  | left-01   | 144.00       | Yes                        |
| bridge-00 | left-00   | 48.00        | No — within left clique     |
| bridge-00 | left-01   | 48.00        | No — within left clique     |
| bridge-00 | left-02   | 48.00        | No — within left clique     |
| bridge-00 | bridge-02 | 64.00        | No — within bridge chain    |
| bridge-00 | right-00  | 96.00        | Yes — bridge→right          |

**Finding**: With the shortest cliques (m1=4) and a mid-length bridge (m2=3), the direct
left↔right shortcut (`left-03`↔`right-00`, commute 128.00) again appears in both hubs'
suggestion lists, this time ranked 3rd of 5 — between the tight-bridge case (ranked 1st via `left-05`↔`right-00` at 128.00) and the loose-bridge case (ranked 5th at 432.00). This
confirms a monotonic trend: longer bridges push the direct-shortcut suggestion further down the
ranked list because more, closer intermediate bridge targets become available.

### Stationary distribution spread — across parameter sets

| Parameter set | Max π    | Min π    | Ratio | Label          |
| ------------- | -------- | -------- | ----- | -------------- |
| Tight         | 0.093750 | 0.031250 | 3.0   | Mildly skewed  |
| Mid           | 0.125000 | 0.062500 | 2.0   | Mildly skewed  |
| Loose         | 0.083333 | 0.027778 | 3.0   | Mildly skewed  |

### Visual observations

> *The following description is inferred from the graph's structural properties and stationary
> distribution. The force-directed layout can be verified by opening `/tmp/nx-barbell-tight.html`,
> `/tmp/nx-barbell-loose.html`, and `/tmp/nx-barbell-mid.html`.*

In the **tight** layout (m1=6, m2=1), the two cliques should render as two dense, near-circular
clusters pulled close together by the single bridge node, which sits directly between them as a
visibly thinner connective link — the overall shape resembles a dumbbell with almost no visible
handle. In the **loose** layout (m1=6, m2=5), the same two clusters should be visibly pulled
apart by the 5-node bridge chain, which renders as a distinct line of small, evenly-spaced nodes
between the two clique masses — the "handle" of the dumbbell becomes long and clearly visible.
In both layouts, the connector nodes (`left-05`/`right-00` in tight; same slugs in loose) should
appear marginally larger than their same-clique peers (reflecting their marginally higher π),
while every bridge-interior node should render as the smallest nodes in the graph, forming a
visibly thin, low-π chain between two large, roughly-equal-sized weights — literally the shape
after which the graph is named.

### Markov questions — answered

- **Does π rank the clique hub nodes (`left-{m1-1}`, `right-00` in the general case; here
  `left-05`/`right-00` for tight, `left-05`/`right-00` for loose, `left-03`/`right-00` for mid)
  above bridge path nodes (`bridge-*`)?**
  Yes, in all three parameter sets. Tight: hubs at π=0.093750 vs. bridge at π=0.031250. Loose:
  hubs at π=0.083333 vs. bridge minimum π=0.027778. Mid: hubs at π=0.125000 vs. all three bridge
  nodes tied at π=0.062500.

- **Does `suggest` propose a direct left-clique↔right-clique shortcut?**
  Yes, in all three parameter sets — but its rank within the top-5 list varies with bridge
  length. Tight (m2=1): the shortcut `left-05`↔`right-00` (commute 128.00) is the #1
  suggestion for both hubs. Mid (m2=3): the shortcut `left-03`↔`right-00` (commute 128.00)
  ranks 3rd of 5. Loose (m2=5): the shortcut `left-05`↔`right-00` (commute 432.00) ranks 5th
  (last) of 5.

- **Is π hub-dominated or flat (max/min ratio), and how does that change with m2?**
  Mildly skewed in all three sets (ratio 2.0–3.0) — never hub-dominated. The ratio does **not**
  change with `m2` (tight m2=1 and loose m2=5 both give ratio 3.0); it is determined entirely
  by `m1` (ratio = m1/2: 6/2=3.0 for tight/loose, 4/2=2.0 for mid), since the max-π node's
  degree is always `m1` and the min-π node's degree is always 2 (a bridge-interior node).

- **Do `bridge-*` nodes consistently show lower π than the clique hubs across all three (m1,
  m2) combinations tested?**
  Yes, without exception. Tight: bridge-00 π=0.031250 vs. hub π=0.093750 (3.0× lower). Loose:
  bridge minimum π=0.027778 vs. hub π=0.083333 (3.0× lower). Mid: all three bridge nodes tied at
  π=0.062500 vs. hub π=0.125000 (2.0× lower). In every case the ratio between hub π and bridge
  π equals the same max/min ratio reported above, because the bridge nodes are always the
  global minimum-π nodes in this graph family.

---

## Sources

- NetworkX: https://networkx.org/documentation/stable/reference/generated/networkx.generators.classic.barbell_graph.html
- Barbell section from graph-topologies research notes
- `docs/adr-009-wiki-gen-make-vs-buy.md` — identifies the barbell as the closest named graph to
  wiki-gen's own topology

---

## Definition of Done

- [x] All sections above filled — no placeholder text remaining
- [x] Slug ordering algorithm documented with a node-id mapping table (not just the pattern)
- [x] `wikigraph analyze --suggest-top 5` output pasted verbatim for all three parameter sets
- [x] Cross-community edge count computed and recorded (cross-zone, computed via Python helper,
      per parameter set)
- [x] All Markov questions answered with actual numbers from analyze output
- [x] Three parameter-set variants shown (tight, mid, loose) as required for parameterised
      graphs, with a comparison table before the per-set detail sections
- [x] Every reference link manually verified to resolve (NetworkX docs 200; ADR file exists in
      repo)
- [x] File committed to branch — path to be updated once #42 resolves
