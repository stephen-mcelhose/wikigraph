---
type: decision
title: Make vs Buy for Synthetic Wiki Graph Generation
description: Decision on whether to hand-build wiki-gen topologies or convert graphs from NetworkX for wikigraph topology experiments.
resource: tools/wiki-gen/main.go
tags: [adr, wiki-gen, topology, networkx, graph-generation, experiments]
timestamp: 2026-08-10T19:00:00Z
status: accepted
---

# ADR-009 — Make vs Buy for Synthetic Wiki Graph Generation

## Context

`wiki-gen` is a Go CLI tool that generates synthetic Markdown wikis for studying
Markov chain properties. It currently produces one hardcoded topology (hierarchical
barbell: bridge → hubs → sections → leaves) parameterized by `--sides`, `--sections`,
`--leaves`. The tool is young and its topology is too regular to produce interesting
analysis signals.

During research we discovered that **NetworkX** (Python) already provides over 100
named graph generators across 29 categories — including the exact structures we want
to study:

- `barbell_graph(m1, m2)` — our current topology
- `relaxed_caveman_graph(l, k, p)` — our intended mixing-parameter topology
- `planted_partition_graph(l, k, p_in, p_out)` — the SBM variant most relevant to us
- `watts_strogatz_graph(n, k, p)` — small-world / rewiring experiments
- `karate_club_graph()` — canonical two-community benchmark (34 nodes, real data)
- `krackhardt_kite_graph()` — 10-node graph with explicit hub/broker/periphery roles
- `LFR_benchmark_graph(n, tau1, tau2, mu, ...)` — gold-standard community benchmark

The question is: for new topologies beyond the current barbell, should we hand-build
them in `wiki-gen` (make) or write a thin conversion layer from NetworkX (buy)?

---

## Conversion Analysis — What It Takes to Use NetworkX

NetworkX produces Python graph objects (`nx.Graph`, `nx.DiGraph`). wikigraph expects
a flat directory of `.md` files where wikilinks (`[[slug]]`) encode directed edges.

### Step 1 — Directed graph

Most NetworkX generators produce **undirected** graphs. wikigraph's Markov model
requires directed edges (the transition matrix is row-stochastic). Conversion:

```python
G = nx.karate_club_graph()      # undirected
D = G.to_directed()             # symmetric directed: every edge u–v becomes u→v and v→u
```

This gives a symmetric transition matrix (equivalent to the lazy random walk on
the undirected graph). Acceptable for experiments; asymmetric topologies would
need hand-crafted directed generators.

### Step 2 — Node slug naming

NetworkX nodes are integers by default. wikigraph slugs must be filesystem-safe
and globally unique within a vault. Two strategies:

| Strategy     | Example output  | Pros                        | Cons                            |
| ------------ | --------------- | --------------------------- | ------------------------------- |
| `page-{n}`   | `page-0.md`     | Zero effort                 | No semantic meaning             |
| Label map    | `alice.md`      | Human-readable              | Only available for social nets  |
| Community prefix | `c0-page-3.md` | Encodes block membership | Leaks ground truth into filenames |

For graphs with node attributes (karate club has `club` labels; SBM/LFR have
community assignments) a label map is straightforward. For anonymous graphs,
`page-{n}` is fine for experiments.

### Step 3 — Edge → wikilink

For each node `u` in the directed graph, write `[[slug(v)]]` for every
out-neighbour `v`. The entire converter is ~50 lines of Python:

```python
import networkx as nx, os, sys

def slug(n): return f"page-{n}"

def convert(G, out_dir):
    D = G.to_directed()
    os.makedirs(out_dir, exist_ok=True)
    for u in D.nodes():
        links = " ".join(f"[[{slug(v)}]]" for v in D.successors(u))
        path = os.path.join(out_dir, f"{slug(u)}.md")
        with open(path, "w") as f:
            f.write(f"# {slug(u)}\n\n{links}\n")
```

Total implementation effort: **~1 hour** including CLI argument parsing and
a `--graph` flag that dispatches to named NetworkX generators.

### Step 4 — What we lose vs wiki-gen

| Capability                          | wiki-gen (Go)  | NetworkX converter (Python) |
| ----------------------------------- | -------------- | --------------------------- |
| Human-readable slug hierarchy       | ✓ (`a-s01-p03`) | ✗ (`page-42`)              |
| Hierarchical barbell with knobs     | ✓              | Requires custom DiGraph     |
| Folder depth / `--dir-depth`        | ✓ (planned)    | Trivial to add              |
| Named community presets (SBM, LFR)  | ✗              | ✓ (built in)               |
| Probabilistic edge generation       | ✗              | ✓                           |
| Social network benchmarks           | ✗              | ✓ (karate, florentine, etc) |
| Single binary, no runtime deps      | ✓              | ✗ (requires Python + NetworkX) |
| Go ecosystem (same module)          | ✓              | ✗                           |

### Step 5 — Edge cases to handle

- **Self-loops**: some generators produce them (`scale_free_graph`). Skip or add
  a dampening weight — wikigraph's Markov model ignores self-links.
- **Multi-edges** (`MultiGraph`): deduplicate to single wikilinks.
- **Disconnected graphs**: wikigraph handles multiple communicating classes fine;
  caveman without connection produces disconnected components — expected.
- **Large graphs**: LFR at n=10,000 would be slow for `analyze` with suggestions.
  Keep n ≤ 100 for interactive experiments.

---

## Options Considered

### Option A — Extend wiki-gen in Go for all topologies

Re-implement SBM, relaxed caveman, Watts-Strogatz, etc. in Go inside `wiki-gen`.

*Rejected*: NetworkX generators are well-tested, widely cited, and would take
weeks to reimplement correctly (especially probabilistic generators like LFR).
Reimplementing for no functional gain violates the spirit of "buy vs build".

### Option B — Replace wiki-gen entirely with a Python script

Single Python tool covers all topologies. Drop wiki-gen.

*Rejected*: wiki-gen's hierarchical barbell with human-readable slugs
(`a-s01-p03`) is the primary tool for parameterised depth/branching experiments.
The readability of the output matters for debugging — `page-42` is opaque.
Python also adds a runtime dependency to what is otherwise a pure Go project.

### Option C (Selected) — Hybrid: keep wiki-gen for hierarchical, add nx-to-wiki for named graphs

Keep `wiki-gen` in Go for hierarchical barbell experiments (the parameterised
topology family we control: depth, branching, sides, mixing). Add a separate
thin Python script `tools/nx-to-wiki/main.py` for named NetworkX graphs.

The two tools are complementary:
- `wiki-gen` → topology experiments with full knob control and human-readable slugs
- `nx-to-wiki` → named benchmark graphs and probabilistic generators for
  comparison and validation

---

## Decision

**Hybrid (Option C).** wiki-gen remains the primary tool for parameterised
hierarchical topology experiments. A thin Python converter (`tools/nx-to-wiki/main.py`)
will be written to access named NetworkX graphs without reimplementing them.

The converter does not need to be production-quality — it is a development
experiment tool. A single file with a `--graph` flag and `--out` flag is sufficient.

---

## Consequences

- **No new Go dependency**: NetworkX stays in Python; wiki-gen stays Go.
- **Python required for nx-to-wiki**: acceptable for experiment tooling, not
  acceptable as a production wikigraph dependency.
- **Two tools to maintain**: low burden since nx-to-wiki is thin (~100 lines).
- **Named topology inventory**: the following NetworkX graphs are immediately
  available for wikigraph experiments once nx-to-wiki is written:

  | Graph                          | Nodes  | Why interesting                         |
  | ------------------------------ | ------ | --------------------------------------- |
  | `karate_club_graph()`          | 34     | Canonical 2-community benchmark         |
  | `krackhardt_kite_graph()`      | 10     | Explicit hub/broker/periphery roles     |
  | `florentine_families_graph()`  | 15     | Dense community + broker structure      |
  | `barbell_graph(m1, m2)`        | varies | Our topology in pure clique form        |
  | `relaxed_caveman_graph(l,k,p)` | varies | Mixing parameter experiments            |
  | `planted_partition_graph(...)`  | varies | p_in / p_out sweep                     |
  | `watts_strogatz_graph(n,k,p)`  | varies | Small-world / β rewiring sweep          |
  | `LFR_benchmark_graph(...)`     | varies | μ mixing parameter sweep (n ≤ 100)     |

## Sources

- [[wiki-gen]] — current hierarchical barbell generator
- NetworkX generators: https://networkx.org/documentation/stable/reference/generators.html
- Research notes: `research/graph-topologies.md`, `research/graph-models.md`
