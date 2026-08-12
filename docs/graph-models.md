# Random Graph Models — Research Notes

---

## Erdős–Rényi G(n, p)

- **Parameters**:
  - `n` — number of nodes.
  - `p` — probability of each possible edge existing, independently.
- **Controls**:
  - `p` directly controls density. Expected edges = p·n(n−1)/2.
  - Sharp connectivity threshold at p = ln(n)/n. Below → almost surely disconnected; above → almost surely connected.
  - Degree distribution: Binomial (≈ Poisson for large n). No hubs — all nodes similar degree.
- **Markov walk**: Mixes in O(log n / log(np)) steps. No community structure — walk is essentially memoryless.
- **Relevance**: Useful as a baseline "noise" topology or to add random cross-links to existing structure. `p` = "randomness" or "noise density" knob.
- **Links**:
  - https://en.wikipedia.org/wiki/Erd%C5%91s%E2%80%93R%C3%A9nyi_model
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.erdos_renyi_graph.html

---

## Barabási–Albert (BA) / Preferential Attachment

- **Parameters**:
  - `n` — total number of nodes in the final graph.
  - `m` — number of edges each new node attaches when it joins.
  - `m₀` — initial seed network size (must satisfy m₀ ≥ m).
- **Controls**:
  - `m` controls hub richness — higher m means denser attachment, more prominent hubs.
  - Growth + preferential attachment ("rich get richer") produces power-law degree distribution P(k) ~ k⁻³.
  - Produces strong core-periphery structure: a few massive hubs + many low-degree periphery nodes.
  - Clustering coefficient C ~ ln(N)²/N → shrinks as graph grows (unlike WS).
- **Variants**:
  - **Non-linear PA**: attachment probability ∝ k^α. Sub-linear (α<1) → stretched exponential; super-linear (α>1) → winner-take-all.
  - **Bianconi–Barabási**: adds per-node "fitness" so later nodes can outcompete older hubs.
  - **powerlaw_cluster_graph(n, m, p)**: Holme-Kim extension that adds triangle-closing step with probability `p` after each preferential attachment — preserves power-law distribution but raises clustering coefficient.
- **Markov walk**: Very fast mixing due to hub shortcuts. π(hub) ∝ k → hubs dominate stationary distribution heavily.
- **Relevance**: The `m` parameter and fitness variant give us a principled way to create "important" pages that genuinely dominate the random walk.
- **Links**:
  - https://en.wikipedia.org/wiki/Barab%C3%A1si%E2%80%93Albert_model
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.barabasi_albert_graph.html
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.powerlaw_cluster_graph.html

---

## Watts-Strogatz (WS) Small-World

- **Parameters**:
  - `n` — number of nodes.
  - `k` — each node initially connected to its k nearest ring neighbors (must be even).
  - `β` (or `p`) — rewiring probability ∈ [0, 1].
- **Controls**:
  - β=0: regular ring lattice. High clustering C(0)≈3(k−2)/4(k−1). Long paths L(0)≈N/2k.
  - β=1: near-random graph. C≈k/N≈0. Short paths L≈ln(N)/ln(k).
  - Small-world regime (0.001 < β < 0.1): L drops to near-minimum while C stays near maximum. The "sweet spot."
  - C(β) ~ C(0)·(1−β)³ — clustering degrades slowly.
  - L(β) drops precipitously even at tiny β — a few shortcuts collapse the diameter.
- **Markov walk**: At β=0, walk is slow-mixing (diameter O(N)). At any positive β, mixing time drops dramatically as shortcuts appear.
- **Relevance**: β is the canonical "shortcut / rewire" knob. Adding random cross-links to our hierarchical wiki is equivalent to raising β. Even β=0.01 (1 in 100 edges rewired) produces dramatic structural changes.
- **Links**:
  - https://en.wikipedia.org/wiki/Watts%E2%80%93Strogatz_model
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.watts_strogatz_graph.html

---

## Stochastic Block Model (SBM) — Planted Partition

- **Parameters**:
  - `l` — number of communities (blocks).
  - `k` — nodes per community.
  - `p_in` — probability of edge within the same community.
  - `p_out` — probability of edge between different communities.
  - Full SBM: `sizes` (list of block sizes) + `P` (l×l probability matrix).
- **Controls**:
  - p_in > p_out → assortative: dense within communities, sparse between. Classic community structure.
  - p_in < p_out → disassortative: bipartite-like, walk prefers crossing.
  - p_in = p_out → degenerates to Erdős–Rényi G(n,p).
  - Modularity Q ≈ 1 − (p_out/p_in) (approximately, for equal-sized communities).
- **Named variants**:
  - **Planted partition model**: equal-sized communities, same p_in and p_out for all pairs (NetworkX `planted_partition_graph`).
  - **Degree-corrected SBM (DCSBM)**: allows heterogeneous degree distributions within blocks.
  - **Hierarchical SBM (HSBM)**: nested communities — blocks within blocks.
  - **Gaussian partition**: community sizes drawn from Gaussian (NetworkX `gaussian_random_partition_graph`).
- **Markov walk**: Conductance Φ ≈ p_out · (l−1) / (p_in·(k−1) + p_out·(l−1)). Low p_out/p_in → low conductance → slow cross-community mixing.
- **Relevance**: This is the most direct model for what we're building. p_in = intra-hub density, p_out = cross-hub density.
- **Links**:
  - https://en.wikipedia.org/wiki/Stochastic_block_model
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.planted_partition_graph.html
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.stochastic_block_model.html
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.gaussian_random_partition_graph.html

---

## LFR Benchmark (Lancichinetti–Fortunato–Radicchi)

- **Parameters**:
  - `n` — total nodes (typical: 1,000–10,000).
  - `tau1` (γ) — power-law exponent for degree distribution. Range: 2.0–3.0. Lower → more hubs.
  - `tau2` (β) — power-law exponent for community size distribution. Range: 1.0–2.0. Lower → more size variance.
  - `mu` (μ) — **mixing parameter**: fraction of each node's edges that go to other communities. Range: 0–1.
  - `min_degree` (k_min) — minimum node degree (typical: 10–20).
  - `max_degree` (k_max) — maximum node degree (typical: √N or 0.1·N).
  - `min_community` (s_min) — minimum community size. Constraint: s_min > k_min.
  - `max_community` (s_max) — maximum community size. Constraint: s_max > k_max.
- **μ in detail**:
  - μ=0: perfectly isolated communities (no cross-community edges).
  - μ=0.5: community boundary — each node has equal internal and external links.
  - μ>0.5: communities indistinguishable from random.
  - Used as the primary x-axis variable in community detection benchmarks.
- **Relevance**: μ is the gold-standard name for the "cross-community mixing" knob. We should call our inter-hub cross-link parameter `--mu` or `--mixing`.
- **Links**:
  - https://en.wikipedia.org/wiki/Lancichinetti%E2%80%93Fortunato%E2%80%93Radicchi_benchmark
  - https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.LFR_benchmark_graph.html
  - Original paper: Lancichinetti et al. (2008), DOI: 10.1103/PhysRevE.78.046110

---

## Summary — Recommended Knob Names for wiki-gen

| Knob concept                        | Best name(s)            | Source / authority       |
| ----------------------------------- | ----------------------- | ------------------------ |
| Number of communities / sides       | `--communities`         | SBM (`l`), caveman (`l`) |
| Nodes per community                 | n/a (derived from depth + branching) | —           |
| Tree branching factor               | `--branching` / `--arity` | NetworkX `balanced_tree(r, h)` |
| Tree depth (link hierarchy levels)  | `--depth`               | NetworkX `balanced_tree` (`h`) |
| Directory nesting depth             | `--dir-depth`           | filesystem convention    |
| Intra-community link density        | `--density` / `--p-in`  | SBM (`p_in`)             |
| Inter-community link density        | `--mixing` / `--mu`     | LFR (μ), SBM (`p_out`)  |
| Random rewiring / shortcut prob     | `--rewire` / `--beta`   | Watts-Strogatz (β)       |
| Hub preferential attachment         | `--attachment`          | Barabási–Albert (`m`)    |

## Links

- https://en.wikipedia.org/wiki/Erd%C5%91s%E2%80%93R%C3%A9nyi_model
- https://en.wikipedia.org/wiki/Barab%C3%A1si%E2%80%93Albert_model
- https://en.wikipedia.org/wiki/Watts%E2%80%93Strogatz_model
- https://en.wikipedia.org/wiki/Stochastic_block_model
- https://en.wikipedia.org/wiki/Lancichinetti%E2%80%93Fortunato%E2%80%93Radicchi_benchmark
- https://en.wikipedia.org/wiki/Modularity_(networks)
- https://networkx.org/documentation/stable/reference/generators.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.erdos_renyi_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.barabasi_albert_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.powerlaw_cluster_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.random_graphs.watts_strogatz_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.planted_partition_graph.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.stochastic_block_model.html
- https://networkx.org/documentation/stable/reference/generated/networkx.generators.community.LFR_benchmark_graph.html
- https://arxiv.org/abs/0805.4770 (LFR original paper)
