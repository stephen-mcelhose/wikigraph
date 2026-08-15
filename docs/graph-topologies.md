---
type: concept
title: Named Graph Topologies
description: Reference catalog of named graph topologies — structural properties, Markov walk behaviour, and NetworkX equivalents — used as benchmarks and analogues for wikigraph analysis.
tags: [graph-theory, topology, networkx, markov-chain, benchmark]
timestamp: 2026-08-11T00:00:00Z
---

# Named Graph Topologies — Research Notes

---

## Star Graph (Hub-and-Spoke)

- **Also known as**: $S_k$, complete bipartite graph $K_{1,k}$, "claw" (for k=3)
- **Key parameters**:
  - `k` — number of leaves (spokes). Total nodes = k+1.
- **Markov walk**:
  - π(center) = 1/2 regardless of k. π(leaf) = 1/(2k).
  - Walk is periodic (period 2) due to bipartiteness — lazy walk needed for convergence.
  - MFPT(center → specific leaf) = 2k−1. MFPT(leaf → leaf) = 2k.
  - Cover time ≈ 2k·ln(k) (Coupon Collector).
  - Mixing time: O(1) — no bottleneck, fastest possible.
- **Relevance**: Our hub is a star. The center dominates π but in our directed case hub links out to 40+ pages, diluting its effective in-degree.
- **Links**:
  - https://en.wikipedia.org/wiki/Star_graph

---

## Caterpillar Tree

- **Also known as**: "caterpillar", Gutman tree (chemistry variant)
- **Key parameters**:
  - **Spine** (backbone / central path) — the path remaining after removing all leaves.
  - **Legs** (hairs / leaves) — degree-1 nodes attached to spine vertices.
  - **Depth** — max distance from any vertex to spine. Always ≤ 1 by definition.
- **Variants**:
  - **Lobster graph** — generalization where depth ≤ 2 (removing leaves yields a caterpillar, not a path).
  - **k-caterpillar** — generalization to k-trees.
- **Relevance**: A sequential article series with side-references is a caterpillar. The spine = tutorial path; legs = reference pages. Our current "section" layer is a set of disconnected mini-stars, not a caterpillar — adding spine links between sections would make it one.
- **Links**:
  - https://en.wikipedia.org/wiki/Caterpillar_tree
  - https://en.wikipedia.org/wiki/Lobster_graph

---

## Balanced Tree (k-ary Tree)

- **Also known as**: r-ary tree, balanced rooted tree, `balanced_tree(r, h)` in NetworkX
- **Key parameters**:
  - `r` — branching factor (children per node). NetworkX calls this `r`.
  - `h` — height (depth from root to leaves). NetworkX calls this `h`.
  - Total nodes = (r^(h+1) − 1) / (r − 1)
- **Markov walk**:
  - Walk is biased toward leaves (they are the majority of nodes).
  - π is proportional to degree: root has highest degree (r), internal nodes degree r+1, leaves degree 1.
  - High effective resistance root-to-leaf — the walk must traverse the full height.
  - Mixing time: O(h²) — diameter dominates.
- **Relevance**: Our current topology is a 2-level tree (hub → sections → leaves). The `--depth` knob extends this to h levels; the `--branching` knob controls r.
- **Links**:
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.classic.balanced_tree.html
  - https://en.wikipedia.org/wiki/K-ary_tree

---

## Caveman Graph

- **Also known as**: clique graph, `caveman_graph(l, k)` in NetworkX
- **Key parameters**:
  - `l` — number of cliques (caves / communities).
  - `k` — size of each clique (all nodes within a clique are fully connected).
- **Variants**:
  - **Connected caveman** — one edge per clique removed and rewired to connect cliques in a ring. Minimum inter-community connectivity.
  - **Relaxed caveman** — each intra-clique edge rewired to another community with probability `p`. Controls inter-community density continuously.
  - **Ring of cliques** — cliques connected in a ring by single bridge edges (similar to connected caveman).
- **Markov walk**:
  - In pure caveman (disconnected cliques): walk is trapped, graph is disconnected.
  - In connected caveman: very high effective resistance between cliques (only 1 bridge edge). Walk rarely crosses. Commute time between communities scales as O(k²·l).
  - Relaxed caveman: raising `p` is the discrete analog of raising μ (mixing parameter).
- **Relevance**: Our topology is a relaxed caveman with l=2, k=~25, and p≈0 (almost no cross-community edges). The `--mixing` knob maps directly to the relaxed caveman `p`.
- **Named graph**: NetworkX `relaxed_caveman_graph(l, k, p)`
- **Links**:
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.caveman_graph.html
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.relaxed_caveman_graph.html
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.ring_of_cliques.html

---

## Windmill Graph

- **Also known as**: friendship graph (for triangles), `windmill_graph(n, k)` in NetworkX
- **Key parameters**:
  - `n` — number of cliques.
  - `k` — size of each clique (all sharing one central hub node).
- **Structure**: n cliques of size k, all sharing a single universal hub. The hub has degree n·(k−1).
- **Markov walk**:
  - π(hub) = n·(k−1) / total_degree → dominates stationary distribution.
  - Walk constantly returns through hub — very high hub centrality.
- **Relevance**: A windmill is our topology taken to an extreme — if the bridge were replaced by a single shared hub connecting all communities. Contrasts with our bridge approach.
- **Links**:
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.windmill_graph.html
  - https://en.wikipedia.org/wiki/Windmill_graph

---

## Planted Partition / Stochastic Block Model

- **Also known as**: SBM, planted l-partition, `planted_partition_graph(l, k, p_in, p_out)` in NetworkX
- **Key parameters**:
  - `l` — number of communities (blocks).
  - `k` — size of each community.
  - `p_in` — probability of edge between two nodes in the same community (intra-community density).
  - `p_out` — probability of edge between nodes in different communities (inter-community density).
- **Derived quantities**:
  - Modularity Q ≈ 1 − (p_out / p_in) when p_out ≪ p_in.
  - Mixing parameter μ ≈ p_out·(l−1)·k / (p_in·(k−1) + p_out·(l−1)·k).
- **Assortative** (p_in > p_out): classic community structure, walk rarely crosses.
- **Disassortative** (p_in < p_out): bipartite-like, walk prefers crossing.
- **Relevance**: The most direct model for our topology. `p_in` = intra-hub link density, `p_out` = cross-hub link density. These are the two key knobs.
- **Links**:
  - https://en.wikipedia.org/wiki/Stochastic_block_model
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.planted_partition_graph.html
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.stochastic_block_model.html

---

## Small-World (Watts-Strogatz)

- **Also known as**: WS model, `watts_strogatz_graph(n, k, p)` in NetworkX
- **Key parameters**:
  - `n` — number of nodes.
  - `k` — each node connected to k nearest ring neighbors (initial lattice degree).
  - `β` (or `p`) — rewiring probability ∈ [0, 1].
- **β=0**: Regular ring lattice. High clustering C(0) ≈ 3/4. Long paths L(0) ≈ N/2k.
- **β=1**: Near-random graph. Low clustering C(1) ≈ k/N ≈ 0. Short paths L(1) ≈ ln(N)/ln(k).
- **Small-world regime** (0.001 < β < 0.1): C stays high, L drops dramatically. Both properties simultaneously.
- **Relevance**: β is the canonical "shortcut probability" or "rewire" knob. Adding random cross-links to our wiki is equivalent to raising β.
- **Links**:
  - https://en.wikipedia.org/wiki/Watts%E2%80%93Strogatz_model
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.watts_strogatz_graph.html

---

---

## NetworkX Classic Graphs (directly relevant)

### Barbell Graph — `barbell_graph(m1, m2)`
- Two complete graphs K_m1 joined by a path of length m2.
- **This is our current topology** — two hub-communities + a bridge path.
- `m1` = community size (clique), `m2` = bridge length.
- Markov walk: extremely low conductance at the bridge. Commute time between communities scales as O(m1·m2).
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.classic.barbell_graph.html

### Lollipop Graph — `lollipop_graph(m, n)`
- Complete graph K_m welded to a path P_n.
- One dense community + a long tail. Walk gets trapped in the clique; tail is rarely visited.
- Useful for modeling a "main topic" cluster with a chain of prerequisite pages.
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.classic.lollipop_graph.html

### Wheel Graph — `wheel_graph(n)`
- Cycle of n−1 nodes all connected to one central hub.
- Hub has degree n−1, π(hub)=(n−1)/(3n−3) ≈ 1/3. All rim nodes have degree 3.
- A richer hub structure than a pure star — rim nodes have lateral connections.
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.classic.wheel_graph.html

### Tadpole Graph — `tadpole_graph(m, n)`
- Cycle C_m attached to a path P_n.
- Walk circulates in the cycle and rarely exits into the tail.
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.classic.tadpole_graph.html

### Dorogovtsev–Goltsev–Mendes Graph — `dorogovtsev_goltsev_mendes_graph(n)`
- Self-similar hierarchical fractal graph built recursively.
- At generation n: 3·(3^(n−1)+1)/2 nodes. Power-law degree distribution emerges.
- Interesting for studying hierarchical depth and self-similarity — each generation adds a new tier.
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.classic.dorogovtsev_goltsev_mendes_graph.html

---

## NetworkX Small Named Graphs (worth knowing)

These are specific famous graphs with their own names — could be used as preset topologies or test fixtures.

| Graph                    | Nodes | Notes |
| ------------------------ | ----- | ----- |
| `petersen_graph()`       | 10    | Cubic, 3-regular, highly symmetric. Famous counterexample in graph theory. |
| `krackhardt_kite_graph()`| 10    | Social network with distinct center, broker, and periphery roles — very relevant to hub/bridge/leaf structure. |
| `karate_club_graph()`    | 34    | Zachary's karate club — real two-community network with a bridge. The canonical community detection example. |
| `florentine_families_graph()` | 15 | Renaissance Florence elite marriages. Rich community + broker structure. |
| `les_miserables_graph()` | 77    | Character co-appearance. Multi-community with natural hubs. |
| `bull_graph()`           | 5     | Triangle + two pendants. Simple bridge + community toy example. |
| `barbell_graph(m1, m2)`  | varies| Our topology. |
| `wheel_graph(n)`         | varies| Hub + lateral rim connections. |

The **Krackhardt Kite** is particularly interesting — it explicitly has nodes in "center", "broker", and "periphery" roles. This mirrors hub, bridge, leaf exactly.

- https://networkx.org/documentation/stable/reference/generated/networkx.generators.social.krackhardt_kite_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.social.karate_club_graph.html

---

## NetworkX Graph Atlas

`graph_atlas_g()` returns all 1,252 non-isomorphic graphs with up to 7 nodes. Could be used as a source of small topology presets for testing edge cases.

- https://networkx.org/documentation/stable/reference/generated/networkx.generators.atlas.graph_atlas.html

---

## Summary — Most Useful Terminology for wiki-gen Knobs

| Structural concept              | Standard name(s)                     | Source model       |
| ------------------------------- | ------------------------------------ | ------------------ |
| Number of communities           | `l`, `num_cliques`, communities      | SBM, caveman       |
| Community size                  | `k`, clique_size                     | SBM, caveman       |
| Tree depth / height             | `h`, height, depth                   | balanced_tree      |
| Children per node               | `r`, branching factor, arity         | balanced_tree      |
| Intra-community edge density    | `p_in`                               | SBM, planted partition |
| Inter-community edge density    | `p_out`                              | SBM, planted partition |
| Fraction of cross-community links | mixing parameter `μ`               | LFR benchmark      |
| Random shortcut probability     | `β` (beta), rewiring probability     | Watts-Strogatz     |
| Spine / backbone                | spine, central path                  | caterpillar        |
| Peripheral nodes                | leaves, legs, periphery              | caterpillar, star  |

## Related Concepts

- [[markov-model]] — how graph edges become a row-stochastic transition kernel
- [[stationary-distribution]] — π is determined by degree structure; varies sharply across these topologies
- [[communicating-classes]] — topology determines whether graph is strongly connected (e.g. pure caveman = multiple classes)
- [[commute-time]] — boundary conductance drives commute time between communities
- [[random-walk]] — the foundational model these topologies are benchmarked against
- [[graph-models]] — random generative models (ER, BA, WS, SBM, LFR) complementing this named-topology catalog
- Benchmark write-ups (provisional `docs/graphs/` pending issue #42): karate-club, krackhardt-kite, barbell, caveman

## Sources

- https://en.wikipedia.org/wiki/Star_graph
- https://en.wikipedia.org/wiki/Caterpillar_tree
- https://en.wikipedia.org/wiki/Lobster_graph
- https://en.wikipedia.org/wiki/K-ary_tree
- https://en.wikipedia.org/wiki/Watts%E2%80%93Strogatz_model
- https://en.wikipedia.org/wiki/Stochastic_block_model
- https://en.wikipedia.org/wiki/Windmill_graph
- https://en.wikipedia.org/wiki/Modularity_(networks)
- https://networkx.org/documentation/stable/reference/generators.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.classic.balanced_tree.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.caveman_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.relaxed_caveman_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.ring_of_cliques.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.planted_partition_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.stochastic_block_model.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.windmill_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.watts_strogatz_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.barabasi_albert_graph.html
