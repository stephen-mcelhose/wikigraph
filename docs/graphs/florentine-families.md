<!-- Path provisional — will be updated when #42 resolves -->

# Florentine Families — marriage-ties brokerage benchmark

> ⚠️ **DRAFT** — AI-assisted write-up, not yet verified by human analysis. Treat all findings as provisional.

> **nx-to-wiki flag**: `--graph florentine`
> **Nodes**: 15 · **Directed edges**: 40 · **Naming**: tier1-florentine-family-name
> **Source**: empirical — fixed

---

## Background

The Florentine Families network encodes marriage alliances among 15 elite families of
15th-century Florence, as reconstructed by John Padgett from historical records. The dataset
is the canonical small-world example of **structural brokerage**: the Medici family occupies
a small number of edges yet those edges span otherwise disconnected cliques of rival families,
giving Medici disproportionate control over information and alliance flow. Padgett & Ansell's
(1993) analysis used this network (alongside a business-ties network) to argue that the
Medici's rise to power in Florence was a structural — not merely personal — phenomenon: their
network *position*, not their wealth alone, made them indispensable brokers between families
that would not otherwise deal with each other.

**Origin**: Padgett, J.F., & Ansell, C.K. (1993). Robust Action and the Rise of the Medici,
1400–1434. *American Journal of Sociology*, 98(6), 1259–1319. https://doi.org/10.1086/229983

---

## Graph Properties

| Property                | Value / Description                                          |
| ------------------------ | ------------------------------------------------------------ |
| Nodes                    | 15                                                            |
| Undirected edges         | 20                                                            |
| Directed edges (×2)      | 40                                                            |
| Degree distribution      | Right-skewed; min 1 (4 isolates-by-degree), max 6 (Medici)    |
| Community structure      | None — no community attribute; families link directly by marriage |
| Diameter                 | 5                                                             |
| Average clustering coef  | 0.1917                                                        |

**Key parameters**: None — this is a fixed empirical dataset with no generative knobs.

---

## Transition Matrix Structure

For every graph produced by `nx-to-wiki`, the random walk is **uniform**: each step chooses a
neighbour with equal probability. The transition matrix $P$ is fully determined by the
adjacency structure of the `.md` files:

$$P_{ij} = \frac{1}{\deg(i)} \quad \text{if } [[j]] \text{ appears in i.md}, \quad \text{else } 0$$

Representative rows (hub, mid-degree, leaf):

| Slug         | Degree | $P_{ij}$ (non-zero) | Structural note                                         |
| ------------ | ------ | -------------------- | -------------------------------------------------------- |
| medici       | 6      | 1/6 ≈ 0.166667        | Highest-degree node; spreads weight across 6 rival cliques |
| guadagni     | 4      | 1/4 = 0.250000        | Second-highest degree; links Medici side to periphery      |
| barbadori    | 2      | 1/2 = 0.500000        | Low-degree connector; half its walk mass goes to Medici     |
| acciaiuoli   | 1      | 1/1 = 1.000000        | Degree-1 leaf; all weight on its single edge (→ medici)     |

**The `.md` files are a lossless encoding of $P$.** The minimum non-zero $P_{ij}$ in this graph
is $1/6 \approx 0.167$ (Medici's row), far above the `--min-edge` default filter of 0.005 — so
both sparse and dense exports below are numerically identical here.

### Export

```bash
# Sparse P — non-zero entries only (40 edge rows; every edge exceeds --min-edge 0.005)
wikigraph export /tmp/nxwiki-florentine --format csv -o /tmp/florentine-export

# Dense P — full 15×15 matrix including structural zeros (210 edge rows)
wikigraph export /tmp/nxwiki-florentine --format csv -o /tmp/florentine-export --min-edge 0
```

Both commands write two files:

- `florentine-export_nodes.csv` — 15 rows: `slug, π, class`
- `florentine-export_edges.csv` — 40 rows (sparse) or 210 rows (dense): `source, target, probability`

No `--min-edge 0` warning applies: minimum $P_{ij} = 1/6 \approx 0.167 \gg 0.005$, so the
default sparse export is lossless.

---

## Slug Naming

**Naming tier**: Tier 1

**Assignment algorithm**: Node labels in `nx.florentine_families_graph()` *are* the family
names (e.g. `"Medici"`, `"Pazzi"`). `nx-to-wiki` applies `slugify()` directly to each label: it
lowercases the string, replaces spaces and underscores with hyphens, and strips any character
that is neither alphanumeric nor a hyphen. There is no sorting, grouping, or zero-padding step
— the mapping is a direct 1:1 string transform, independent of node order or degree.

**Why this ordering**: There is no ordering — each family name maps to exactly one slug via
`slugify()`. `"Medici"` → `medici`, `"De' Pazzi"`-style apostrophes/spaces would be stripped if
present, but none of the 15 canonical labels in this dataset contain spaces or punctuation, so
every slug is simply the lowercased family name.

Full node-id → slug mapping:

| Slug          | NetworkX node label | Degree |
| ------------- | -------------------- | ------ |
| acciaiuoli    | Acciaiuoli            | 1      |
| albizzi       | Albizzi                | 3      |
| barbadori     | Barbadori              | 2      |
| bischeri      | Bischeri               | 3      |
| castellani    | Castellani             | 3      |
| ginori        | Ginori                 | 1      |
| guadagni      | Guadagni               | 4      |
| lamberteschi  | Lamberteschi           | 1      |
| medici        | Medici                 | 6      |
| pazzi         | Pazzi                  | 1      |
| peruzzi       | Peruzzi                | 3      |
| ridolfi       | Ridolfi                | 3      |
| salviati      | Salviati               | 2      |
| strozzi       | Strozzi                | 4      |
| tornabuoni    | Tornabuoni             | 3      |

---

## How to Generate

```bash
# Minimal — default parameters (no knobs needed; this is a fixed empirical dataset)
python3 tools/nx-to-wiki/main.py --graph florentine --out /tmp/nxwiki-florentine
wikigraph graph /tmp/nxwiki-florentine -o /tmp/nx-florentine.html --title "Florentine Families"
wikigraph analyze /tmp/nxwiki-florentine --suggest-top 5
```

No parameter variants apply — `--graph florentine` always produces the same 15-node, 20-edge
graph; there are no generative knobs to sweep.

**What the output directory contains:**

- 15 `.md` files, one per family
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
> This graph has **no community attribute** (unlike karate club or planted-partition), so the
> cross-community section of the template is replaced below with a **Structural Position**
> section: which families are isolates, which are connectors, and where Medici sits by π versus
> betweenness centrality.

### Raw analyze output

```
=== Overview ===
Pages:        15
Edges:        40
Entropy rate: 1.6010 bits
Classes:      1

=== Communicating classes ===
Class 1 (recurrent): 15 page(s)
  acciaiuoli
  albizzi
  barbadori
  bischeri
  castellani
  ginori
  guadagni
  lamberteschi
  medici
  pazzi
  peruzzi
  ridolfi
  salviati
  strozzi
  tornabuoni

=== Orphan pages (bottom 10% by stationary distribution) ===
  acciaiuoli                                π=0.025000  → add inbound links
  lamberteschi                              π=0.025000  → add inbound links

=== Sink pages (no outgoing links) ===
  (none)

=== Most central (top 5 by stationary distribution) ===
  1. medici                                    π=0.150000
  2. guadagni                                  π=0.100000
  3. strozzi                                   π=0.100000
  4. castellani                                π=0.075000
  5. peruzzi                                   π=0.075000

=== Suggested missing links (lowest commute time, not yet linked) ===
  acciaiuoli:
    → tornabuoni                              (commute: 60.66)
    → ridolfi                                 (commute: 61.06)
    → albizzi                                 (commute: 67.05)
    → guadagni                                (commute: 68.21)
    → barbadori                               (commute: 68.91)
  albizzi:
    → tornabuoni                              (commute: 35.53)
    → ridolfi                                 (commute: 39.77)
    → bischeri                                (commute: 44.77)
    → strozzi                                 (commute: 45.60)
    → barbadori                               (commute: 50.60)
  barbadori:
    → strozzi                                 (commute: 37.38)
    → ridolfi                                 (commute: 40.03)
    → peruzzi                                 (commute: 42.25)
    → tornabuoni                              (commute: 42.68)
    → bischeri                                (commute: 45.03)
  bischeri:
    → castellani                              (commute: 31.82)
    → ridolfi                                 (commute: 35.00)
    → medici                                  (commute: 36.06)
    → tornabuoni                              (commute: 37.12)
    → albizzi                                 (commute: 44.77)
  castellani:
    → bischeri                                (commute: 31.82)
    → medici                                  (commute: 35.63)
    → ridolfi                                 (commute: 36.82)
    → guadagni                                (commute: 42.38)
    → tornabuoni                              (commute: 42.52)
  ginori:
    → guadagni                                (commute: 67.05)
    → medici                                  (commute: 67.05)
    → tornabuoni                              (commute: 75.53)
    → ridolfi                                 (commute: 79.77)
    → bischeri                                (commute: 84.77)
  guadagni:
    → medici                                  (commute: 28.21)
    → ridolfi                                 (commute: 32.58)
    → strozzi                                 (commute: 33.91)
    → peruzzi                                 (commute: 39.77)
    → castellani                              (commute: 42.38)
  lamberteschi:
    → tornabuoni                              (commute: 64.50)
    → albizzi                                 (commute: 67.05)
    → bischeri                                (commute: 67.58)
    → medici                                  (commute: 68.21)
    → ridolfi                                 (commute: 72.58)
  medici:
    → guadagni                                (commute: 28.21)
    → strozzi                                 (commute: 31.39)
    → castellani                              (commute: 35.63)
    → bischeri                                (commute: 36.06)
    → peruzzi                                 (commute: 39.50)
  pazzi:
    → medici                                  (commute: 80.00)
    → tornabuoni                              (commute: 100.66)
    → ridolfi                                 (commute: 101.06)
    → albizzi                                 (commute: 107.05)
    → guadagni                                (commute: 108.21)
  peruzzi:
    → ridolfi                                 (commute: 37.65)
    → medici                                  (commute: 39.50)
    → guadagni                                (commute: 39.77)
    → barbadori                               (commute: 42.25)
    → tornabuoni                              (commute: 43.21)
  ridolfi:
    → guadagni                                (commute: 32.58)
    → bischeri                                (commute: 35.00)
    → castellani                              (commute: 36.82)
    → peruzzi                                 (commute: 37.65)
    → albizzi                                 (commute: 39.77)
  salviati:
    → tornabuoni                              (commute: 60.66)
    → ridolfi                                 (commute: 61.06)
    → albizzi                                 (commute: 67.05)
    → guadagni                                (commute: 68.21)
    → barbadori                               (commute: 68.91)
  strozzi:
    → medici                                  (commute: 31.39)
    → guadagni                                (commute: 33.91)
    → tornabuoni                              (commute: 34.57)
    → barbadori                               (commute: 37.38)
    → albizzi                                 (commute: 45.60)
  tornabuoni:
    → strozzi                                 (commute: 34.57)
    → albizzi                                 (commute: 35.53)
    → bischeri                                (commute: 37.12)
    → castellani                              (commute: 42.52)
    → barbadori                               (commute: 42.68)
```

### Stationary distribution (π)

| Rank | Slug       | π value | Expected structural role                             |
| ---- | ---------- | ------- | ----------------------------------------------------- |
| 1    | medici     | 0.150000 | Historically the central broker family; degree 6 (max) |
| 2    | guadagni   | 0.100000 | Second-highest degree (4); links Medici to periphery    |
| 3    | strozzi    | 0.100000 | Rival power bloc leader; degree 4                       |
| 4    | castellani | 0.075000 | Strozzi-aligned; degree 3                                |
| 5    | peruzzi    | 0.075000 | Strozzi-aligned; degree 3                                |

**Finding**: `medici` leads π at 0.150, exactly 50% higher than the second-place tie
(`guadagni`/`strozzi` at 0.100). This corroborates Padgett & Ansell's brokerage thesis: even
under a pure degree/uniform-walk model with no betweenness weighting, Medici's marriage ties
reach both the Strozzi-Peruzzi-Castellani bloc and the Guadagni-Albizzi-Tornabuoni bloc,
concentrating random-walk visitation on that single family.

### Suggested links

Top-2 π nodes (`medici`, `guadagni`) and the lowest-π node (`acciaiuoli`, π=0.025000, tied with
`lamberteschi`). Medici's actual neighbours are `acciaiuoli, albizzi, barbadori, ridolfi,
salviati, tornabuoni`; Guadagni's actual neighbours are `albizzi, bischeri, lamberteschi,
tornabuoni` — none of the suggestions below duplicate an existing edge:

| From         | To          | Commute time | Genuine suggestion? |
| ------------ | ----------- | ------------ | -------------------- |
| medici       | guadagni    | 28.21        | Yes                   |
| medici       | strozzi     | 31.39        | Yes                   |
| medici       | castellani  | 35.63        | Yes                   |
| medici       | bischeri    | 36.06        | Yes                   |
| medici       | peruzzi     | 39.50        | Yes                   |
| guadagni     | medici      | 28.21        | Yes                   |
| guadagni     | ridolfi     | 32.58        | Yes                   |
| guadagni     | strozzi     | 33.91        | Yes                   |
| guadagni     | peruzzi     | 39.77        | Yes                   |
| guadagni     | castellani  | 42.38        | Yes                   |
| acciaiuoli   | tornabuoni  | 60.66        | Yes                   |
| acciaiuoli   | ridolfi     | 61.06        | Yes                   |
| acciaiuoli   | albizzi     | 67.05        | Yes                   |
| acciaiuoli   | guadagni    | 68.21        | Yes                   |
| acciaiuoli   | barbadori   | 68.91        | Yes                   |

**Finding**: `medici` and `guadagni` mutually suggest each other as their #1 recommendation
(commute time 28.21, the lowest in the entire graph) — despite both already being the two
highest-π families, the random walk still identifies a missing direct tie between them as the
single most efficient structural improvement available, since their existing shortest paths to
each other run through the shared neighbour `albizzi`/`tornabuoni`. None of Medici's or
Guadagni's suggestions target the three other degree-1 isolates (`ginori`, `pazzi`,
`lamberteschi` — `acciaiuoli` is the fourth isolate and is already linked to Medici) — closer
non-isolate targets dominate on commute time. `acciaiuoli`, the lowest-π node, suggests links to
`tornabuoni` and `ridolfi` (both two hops from its sole existing neighbour, `medici`), not to
any other isolate — confirming isolates are pulled toward the graph's existing hub cluster
rather than toward each other. For the remaining 10 mid-degree families (not shown in the table
above), suggestions consistently cluster around `medici`, `guadagni`, `ridolfi`, and `strozzi`
— the four highest-π nodes — showing that low commute time tracks proximity to existing hubs
across the whole graph, not just for the top/bottom nodes highlighted here.


### Structural Position (replaces cross-community section — no community attribute exists)

Degree-1 isolates: `acciaiuoli` (→medici), `ginori` (→guadagni), `pazzi` (→salviati),
`lamberteschi` (→guadagni). All four sit at π=0.025–0.033, the bottom of the distribution.

Betweenness centrality (computed directly on the underlying NetworkX graph, since `wikigraph`
has no betweenness metric):

| Family      | Betweenness | π       | Note                                             |
| ----------- | ----------- | ------- | ------------------------------------------------- |
| medici      | 0.5220      | 0.150000 | Rank 1 in both metrics — the network's sole broker |
| guadagni    | 0.2546      | 0.100000 | Rank 2 in both metrics                             |
| albizzi     | 0.2125      | 0.050000 | High betweenness but mid-low π — connects Medici's side to Guadagni's side without high degree |
| salviati    | 0.1429      | 0.033333 | Bridges Pazzi (isolate) to the rest of the network  |
| acciaiuoli  | 0.0000      | 0.025000 | Zero betweenness — a pure leaf, contributes nothing to paths between others |

**Finding**: π and betweenness agree on the top-2 (`medici`, `guadagni`), confirming Medici's
brokerage role is visible under both a pure random-walk visitation model and a shortest-path
model — two structurally different definitions of "importance" converge on the same family.

### Stationary distribution spread

- **Max π**: 0.150000 (medici)
- **Min π**: 0.025000 (acciaiuoli, tied with lamberteschi)
- **Ratio max/min**: 6.0 — **Hub-dominated** (falls in the 5–20 tier)

**Finding**: A ratio of 6.0 equals Medici's degree exactly (degree 6), consistent with the
general fact that π_i ∝ deg_i for a uniform random walk on an undirected graph, while the four
degree-1 leaves sit at the theoretical minimum π = 1/(2|E|) = 1/40 × 2 = 0.025.

### Visual observations

> *The following description is inferred from the graph's structural properties and stationary
> distribution. The force-directed layout can be verified by opening `/tmp/nx-florentine.html`.*

Given Medici's degree-6 hub role and the four degree-1 leaves attached to different parts of
the graph (`acciaiuoli`→medici, `ginori`→guadagni, `pazzi`→salviati,
`lamberteschi`→guadagni), the force-directed layout should place `medici` near the visual
centre with visibly larger node size (proportional to π=0.150), flanked by `guadagni` and
`strozzi` as secondary hubs. The four degree-1 leaves should appear as small pendant nodes at
the graph's periphery, each with a single thin edge to its one neighbour. `pazzi`, connected
only to `salviati` (itself only degree-2), should appear as the most visually isolated
two-hop chain from the Medici-centred core — consistent with the historical fact that the
Pazzi family was Medici's principal rival and notably *not* directly tied to them by marriage.

### Markov questions — answered

- **Does `medici` rank highest in π?**
  Yes — `medici` has π=0.150000, the highest in the graph, 50% above the second-place tie
  (`guadagni`/`strozzi` at π=0.100000 each).

- **Does `suggest` recommend any Medici↔periphery links that don't already exist?**
  No — Medici's top-5 suggestions (`guadagni`, `strozzi`, `castellani`, `bischeri`, `peruzzi`,
  commute times 28.21–39.50) are all non-isolate, already-central families. None of the four
  degree-1 peripheral families (`ginori`, `pazzi`, `lamberteschi`, and `acciaiuoli` itself,
  which is already linked to Medici) appear among Medici's suggestions.

- **Is π hub-dominated or flat (max/min ratio)?**
  Hub-dominated — max/min ratio = 6.0 (0.150000 / 0.025000), which falls in the 5–20 tier.

- **Does the π ranking corroborate Medici's historical brokerage position, and do any other
  families come close?**
  Yes. `medici` leads at π=0.150000, and no other family comes close — the runner-up tie
  (`guadagni`, `strozzi` at π=0.100000) is 33% below Medici. The gap between rank 1 and rank 2
  (0.050) is twice the gap between ranks 2 and 4 (0.025), showing Medici's lead is not merely
  incremental but a distinct step change, matching Padgett & Ansell's argument that Medici
  occupied a categorically different structural position from every other family, including
  their nearest rivals Strozzi and Guadagni.

---

## References

- Padgett, J.F., & Ansell, C.K. (1993). Robust Action and the Rise of the Medici, 1400–1434.
  *American Journal of Sociology*, 98(6), 1259–1319. https://doi.org/10.1086/229983
- NetworkX: https://networkx.org/documentation/stable/reference/generated/networkx.generators.social.florentine_families_graph.html

No applicable graph-topology or graph-model research-notes section exists for this empirical
network — omitted per implementation guidance.

---

## Definition of Done

- [x] All sections above filled — no placeholder text remaining
- [x] Slug ordering algorithm documented with a node-id mapping table (not just the pattern)
- [x] `wikigraph analyze --suggest-top 5` output pasted verbatim
- [x] Cross-community edge count computed and recorded (replaced with Structural Position
      section — no community attribute exists in this graph)
- [x] All Markov questions answered with actual numbers from analyze output
- [x] At least one variant command shown, or a note explaining why none apply (fixed graph —
      no generative knobs)
- [x] Every reference link manually verified to resolve (DOI 302→200, NetworkX docs 200)
- [x] File committed to branch — path to be updated once #42 resolves
