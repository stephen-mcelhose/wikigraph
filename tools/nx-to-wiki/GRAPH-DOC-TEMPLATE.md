# Graph Documentation Template

Use this template for each graph issue (#34–#41). One filled-out document per graph.
Output location is pending #42 — commit to a branch; final path will be confirmed when #42 resolves.
**Do not wait on #42 to start tasks 1–3. Task 4 (wiki placement) is explicitly blocked on it.**

**Done = one `.md` file committed to a branch with all sections filled, `wikigraph analyze`
output pasted verbatim, all Markov questions answered with actual numbers, and all reference
links verified to resolve.**

---

<!-- Copy everything below this line into a new file named after the graph, e.g. karate-club.md -->

# [Graph Name] — [one-line descriptor]

> **nx-to-wiki flag**: `--graph [flag-value]`
> **Nodes**: N · **Directed edges**: E · **Naming**: [tier1-label / tier2-role-walk / tier3-fallback]
> **Source**: [empirical / generative — fixed / generative — parameterised]

---

## Background

<!-- 2–4 sentences. What is this graph? Where does it come from? Why does it matter as a benchmark? -->

[Background text.]

**Origin**: [Author (year) — full citation below]

---

## Graph Properties

<!-- Fill in the structural facts. Some will be known from the literature; some from running nx-to-wiki and inspecting the output. -->

| Property                | Value / Description |
| ----------------------- | ------------------- |
| Nodes                   |                     |
| Undirected edges        |                     |
| Directed edges (×2)     |                     |
| Degree distribution     | [uniform / power-law / bimodal / ...] |
| Community structure     | [number of communities, how defined] |
| Diameter                |                     |
| Average clustering coef |                     |

**Key parameters** (for generative graphs only — omit for fixed empirical graphs):

| Flag        | Default | Range | What it controls |
| ----------- | ------- | ----- | ---------------- |
| `--[flag]`  |         |       |                  |

---

## Slug Naming

<!-- Explain the naming scheme fully — not just the pattern, but the ordering algorithm.
     An implementer should be able to predict the exact slug for any given node without running the tool. -->

**Naming tier**: [Tier 1 / Tier 2 / Tier 3]

**Assignment algorithm**: [Explain exactly how slugs are assigned. E.g.: "Nodes are sorted by
node-id within each club group. Within the `Mr. Hi` group, node IDs are sorted ascending and
assigned `mhi-00`, `mhi-01`, … zero-padded to the width of the group size. Same for `off-*`."]

**Why this ordering**: [E.g.: "Node 0 in the NetworkX karate graph is the known `Mr. Hi`
instructor — it becomes `mhi-00`. Node 33 is the `Officer` instructor — it becomes `off-17`."]

Example slugs and their node-id mappings:

| Slug | NetworkX node id | Notes |
| ---- | ---------------- | ----- |
|      |                  |       |
|      |                  |       |

---

## How to Generate

```bash
# Minimal — default parameters
python3 tools/nx-to-wiki/main.py --graph [flag] --out /tmp/nxwiki-[name]
wikigraph graph /tmp/nxwiki-[name] -o /tmp/nx-[name].html --title "[Title]"
wikigraph analyze /tmp/nxwiki-[name] --suggest-top 5
```

<!-- For parameterised graphs: show at least two additional examples illustrating the effect of
     varying the key knob. For fixed empirical graphs, one command is sufficient. -->

```bash
# [Describe what this variant does — e.g. "longer bridge, weaker community signal"]
python3 tools/nx-to-wiki/main.py --graph [flag] [params] --out /tmp/nxwiki-[name]-variant
```

**What the output directory contains:**

- `[N]` `.md` files, one per node
- Each file links to all its undirected neighbours (symmetric — every link is bidirectional)
- File layout: `# [slug]\n\n[[neighbour-1]] [[neighbour-2]] ...`
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

<!-- Paste the full output of `wikigraph analyze /tmp/nxwiki-[name] --suggest-top 5` verbatim.
     Do not summarise or edit — the raw output is the record. -->

```
[paste here]
```

### Stationary distribution (π)

<!-- Which nodes have the highest π? Do they match the structurally "important" nodes from the
     literature? Use the actual values from the analyze output above. -->

| Rank | Slug | π value | Expected structural role |
| ---- | ---- | ------- | ------------------------ |
| 1    |      |         |                          |
| 2    |      |         |                          |
| 3    |      |         |                          |

**Finding**: [1–2 sentences. Does the top-π node match the known hub/broker/central node?
For karate club, the expected leaders are `mhi-00` (π ≈ 0.103) and `off-16` (π ≈ 0.109) — the
two known instructor hubs. Use a similar concrete expectation for each graph.]

### Suggested links

<!-- What does `--suggest-top 5` recommend? Classify each suggestion: within-community or
     cross-community? Use slug prefixes / community assignments to determine this. -->

| From | To | Commute time | Within / Cross community |
| ---- | -- | ------------ | ------------------------ |
|      |    |              |                          |

**Finding**: [What fraction of suggestions cross community boundaries? What does this tell us —
are the communities well-separated (few cross-community suggestions) or porous (many)?]

### Cross-community edges (structural boundary proxy)

<!-- Because communicating-class colouring is uninformative, use this as the boundary signal instead.
     Count edges in the generated .md files that cross community boundaries. -->

```bash
# Count cross-community links for graphs with known community prefixes (e.g. mhi-* vs off-*)
grep -l "mhi-" /tmp/nxwiki-[name]/*.md | xargs grep -o "\[\[off-[^]]*\]\]" | wc -l
```

- **Cross-community directed edges**: [n] of [total] ([pct]%)
- **Expected**: [high for well-mixed graphs like WS at high β; low for caveman at low p]

### Stationary distribution spread

<!-- Is π flat (everyone visited equally) or spiky (a few hubs dominate)? -->

- **Max π**: [value] ([slug])
- **Min π**: [value] ([slug])
- **Ratio max/min**: [value]

**Finding**: [Ratio ≈ 1 means uniform random walk (no hubs). High ratio means hub-dominated.
For reference: a pure star graph with k leaves gives ratio = k. What does this graph produce?]

### Visual observations

<!-- From the `wikigraph graph` HTML output. Describe what the force-directed layout reveals.
     Be specific: which nodes are large (high π)? Are there visible clusters? Isolated nodes? -->

[Describe what you see: cluster structure, hub nodes, bridge nodes, isolated periphery, etc.]

### Markov questions — answered

<!-- Answer every question with actual numbers from the analyze output. No "probably" or "seems to". -->

- **Does π rank the expected "important" node highest?**
  [Yes/No — slug `[x]` has π = [value]; expected hub from literature is `[y]` = node-id [n]]

- **Does `suggest` propose cross-community links?**
  [Yes/No — [n] of 5 suggestions cross community boundaries]

- **Is the stationary distribution hub-dominated or flat?**
  [Hub-dominated (max/min ratio = [x]) / Roughly flat (ratio = [x])]

- **[Graph-specific question tailored to this graph's known structure]**
  [Answer with actual value. Examples: "Does the Medici node lead in π?" / "Does the bridge
  node in a barbell have lower π than both clique hubs?"]

---

## References

<!-- Verify every link resolves before marking this done. Dead links are a defect. -->

- [Primary paper]: [Author (year). Title. *Journal*. https://doi.org/...]
- NetworkX: [https://networkx.org/documentation/stable/reference/generated/...]

<!-- Note: research/graph-topologies.md and research/graph-models.md are in the project output
     folder (not committed to the repo). Reference the relevant section by name only — do not
     link to a file path that doesn't exist in the repo. -->

- Graph topology notes: [relevant section name from graph-topologies research, e.g. "Planted Partition / SBM"]
- Graph model notes: [relevant section name from graph-models research, e.g. "LFR Benchmark"]

---

## Definition of Done

- [ ] All sections above filled — no `[placeholder]` text remaining
- [ ] Slug ordering algorithm documented with a node-id mapping table (not just the pattern)
- [ ] `wikigraph analyze --suggest-top 5` output pasted verbatim
- [ ] Cross-community edge count computed and recorded
- [ ] All Markov questions answered with actual numbers from analyze output
- [ ] At least one variant command shown (parameterised graphs) or a note explaining why none apply (fixed graphs)
- [ ] Every reference link manually verified to resolve (DOIs, NetworkX docs)
- [ ] File committed to branch — path to be updated once #42 resolves
