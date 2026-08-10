---
type: runbook
title: wikigraph Manual Test Runbook
description: Manual test plan covering all subcommands and edge cases, including graph, analyze, goal, and export.
tags: [testing, runbook, qa]
timestamp: 2026-08-09T06:54:46Z
---

# wikigraph — Manual Test Runbook

## Overview

Comprehensive manual test plan covering all `wikigraph` subcommands (`graph`, `analyze`, `goal`, `export`) and edge cases against the repository's `docs/` wiki.

**Binary:** built locally from `~/repos/wikigraph/`  
**Wiki under test:** `docs/` in this repo (23 content pages)  
**Run all commands from:** `~/repos/wikigraph/` (the repo root)  
**Last verified:** 2026-08-09 against current wiki state (26 pages, 117 edges, 1 recurrent class)

---

## Prerequisites

Build the binary:

```bash
cd ~/repos/wikigraph
go build -o wikigraph .
export PATH="$PWD:$PATH"
```

All commands use `docs/` as the wiki path. Run from `~/repos/wikigraph`:

```bash
cd ~/repos/wikigraph
```

The 25 content pages (after default exclusions of `index`, `log`, `AGENTS`):
`adr-001-embedding-layer`, `adr-002-slug-resolution`, `adr-003-orphan-threshold`, `adr-004-quantum-go-example-wiki`, `adr-005-page-type-conventions-and-proposal-storage`, `adr-006-recursive-vault-traversal`, `analyze`, `architecture`, `catrace`, `communicating-classes`, `commute-time`, `entropy-rate`, `export`, `goal`, `graph`, `how-to-docs-plan`, `llm-wiki-pattern`, `markov-model`, `mfpt`, `page-type-conventions`, `quickstart`, `random-walk`, `recurrent-class`, `sink-page`, `stationary-distribution`, `testing-runbook`

Verify:

```bash
wikigraph --help         # must list graph, goal, export, analyze
```

---

## Test Cases

## TC-01 · Root help and persistent flag

**Goal:** Root command lists all four subcommands; `--exclude` is shown as a persistent flag.

```bash
wikigraph --help
```

**Pass criteria:**
- Output lists `graph`, `goal`, `export`, `analyze` under *Available Commands*
- `--exclude` appears under *Flags* (not under a subcommand)
- Exit code 0

---

## TC-02 · graph — baseline render

**Goal:** Full wiki renders to a valid HTML file. See [[graph]].

```bash
wikigraph graph docs/ -o /tmp/wg_graph.html
open /tmp/wg_graph.html
```

**Pass criteria:**
- Stderr prints `Pages: 20` and `Written: /tmp/wg_graph.html`
- File exists and opens in browser showing a force-directed graph
- 20 labelled nodes visible; nodes sized differently (stationary dist)
- 1 colour (1 recurrent class)
- Exit code 0

---

## TC-03 · graph — --exclude removes pages

**Goal:** Excluding a content page reduces node count.

```bash
wikigraph graph docs/ -e index -e log -e AGENTS -e testing-runbook -o /tmp/wg_excl.html
```

**Pass criteria:**
- Stderr prints `Pages: 19`
- `testing-runbook` node absent from rendered graph

---

## TC-04 · graph — --min-edge reduces edge clutter

**Goal:** Raising the edge threshold produces a sparser graph.

```bash
wikigraph graph docs/ --min-edge 0.30 -o /tmp/wg_sparse.html
```

**Pass criteria:**
- File renders with visibly fewer edges than TC-02 output
- No error

---

## TC-05 · graph — --sed patches the HTML

**Goal:** sed expressions are applied to the output.

```bash
wikigraph graph docs/ \
  -s 's/docs wiki/TEST TITLE/' \
  -o /tmp/wg_sed.html
grep "TEST TITLE" /tmp/wg_sed.html
```

**Pass criteria:**
- `grep` finds the patched string
- Exit code 0

---

## TC-06 · graph — missing wiki dir

**Goal:** Helpful error on bad path.

```bash
wikigraph graph /tmp/does-not-exist
```

**Pass criteria:**
- Non-zero exit code
- Error message references the bad path

---

## TC-07 · goal — learning path to a goal page

**Goal:** Subgraph centred on a goal page. See [[goal]].

```bash
wikigraph goal docs/ --goal analyze --top 5 -o /tmp/wg_goal.html
```

**Pass criteria:**
- Stderr prints `Pages: 30` and `Written: /tmp/wg_goal.html (5 nodes, strategy: union)`
- Output file `/tmp/wg_goal.html` exists with 5 nodes
- `analyze` node is present
- Exit code 0

---

## TC-07a · goal — intersection strategy

**Goal:** Surface shared prerequisite pages connecting multiple goals ($\max_g \text{MFPT}$).

```bash
wikigraph goal docs/ --goal mfpt --goal analyze --strategy intersection --top 5 -o /tmp/wg_goal_intersection.html
```

**Pass criteria:**
- Stderr prints `Written: /tmp/wg_goal_intersection.html (5 nodes, strategy: intersection)`
- Both goal nodes (`mfpt`, `analyze`) are present
- Exit code 0

---

## TC-07b · goal — path strategy

**Goal:** Sequential Dijkstra transition chain connecting goals in flag order ($w_{ij} = -\log P_{ij}$).

```bash
wikigraph goal docs/ --goal markov-model --goal analyze --strategy path --top 5 -o /tmp/wg_goal_path.html
```

**Pass criteria:**
- Stderr prints `Written: /tmp/wg_goal_path.html (5 nodes, strategy: path)`
- Output contains the sequential path nodes connecting `markov-model` to `analyze`
- Exit code 0

---

## TC-07c · goal — bottleneck strategy

**Goal:** Gatekeeper chokepoints ranked by random walk betweenness centrality across goal pairs.

```bash
wikigraph goal docs/ --goal markov-model --goal analyze --strategy bottleneck --top 5 -o /tmp/wg_goal_bottleneck.html
```

**Pass criteria:**
- Stderr prints `Written: /tmp/wg_goal_bottleneck.html (5 nodes, strategy: bottleneck)`
- Exit code 0

---

## TC-07d · goal — unknown strategy

**Goal:** Helpful error message when invalid strategy name is passed.

```bash
wikigraph goal docs/ --goal analyze --strategy invalid
```

**Pass criteria:**
- Error message: `unknown strategy "invalid" (valid: union, intersection, path, bottleneck)`
- Non-zero exit code

## TC-08 · goal — multiple goals

**Goal:** Two goal pages, both present in output.

```bash
wikigraph goal docs/ \
  --goal analyze --goal graph \
  --top 5 -o /tmp/wg_goal2.html
```

**Pass criteria:**
- Stderr shows `(5 nodes)`
- Both `analyze` and `graph` present in rendered graph

---

## TC-09 · goal — unknown slug prints valid list

**Goal:** Clear error when a slug doesn't exist.

```bash
wikigraph goal docs/ --goal not-a-real-page -o /tmp/x.html
```

**Pass criteria:**
- Stderr prints `unknown --goal slug "not-a-real-page"`
- Stderr lists valid slugs
- Non-zero exit code; no output file written

---

## TC-10 · goal — missing --goal flag

**Goal:** Requires at least one goal.

```bash
wikigraph goal docs/ -o /tmp/x.html
```

**Pass criteria:**
- Error: `at least one --goal slug is required`
- Non-zero exit code

---

## TC-11 · export — JSON

**Goal:** Valid node-link JSON with correct shape. See [[export]].

```bash
wikigraph export docs/ --format json -o /tmp/wg
jq '.' /tmp/wg.json | head -20
```

**Pass criteria:**
- File `/tmp/wg.json` exists
- Top-level keys are `nodes` and `links`
- Each node has `id` (string), `pi` (float), `class` (int)
- Each link has `source`, `target` (strings), `value` (float)
- `jq` exits 0 (valid JSON)

---

## TC-12 · export — CSV

**Goal:** Two well-formed CSV files.

```bash
wikigraph export docs/ --format csv -o /tmp/wg
head /tmp/wg_nodes.csv
head /tmp/wg_edges.csv
```

**Pass criteria:**
- `wg_nodes.csv` header: `slug,pi,class`
- `wg_edges.csv` header: `source,target,probability`
- Both have 20 data rows in nodes (one per page)
- Values are numeric floats; no empty cells

---

## TC-13 · export — DOT

**Goal:** Valid Graphviz DOT output.

```bash
wikigraph export docs/ --format dot -o /tmp/wg
head /tmp/wg.dot
dot -Tsvg /tmp/wg.dot -o /tmp/wg.svg && echo "dot OK"
```

**Pass criteria:**
- First line: `digraph wiki {`
- Node lines: `"slug" [weight=N.NNNNNN];`
- Edge lines: `"slug" -> "slug" [weight=N.NNNNNN];`
- `dot -Tsvg` succeeds (requires Graphviz; skip if not installed)

---

## TC-14 · export — unknown format

```bash
wikigraph export docs/ --format xml -o /tmp/wg
```

**Pass criteria:**
- Error: `unknown format "xml": choose json, csv, or dot`
- Non-zero exit code

---

## TC-15 · export — --min-edge filters edges

```bash
wikigraph export docs/ --format json --min-edge 0.5 -o /tmp/wg_sparse
jq '.links | length' /tmp/wg_sparse.json
```

**Pass criteria:**
- Link count is significantly lower than with default `--min-edge 0.005`

---

## TC-16 · analyze — full report

**Goal:** All six sections printed; known facts match. See [[analyze]].

```bash
wikigraph analyze docs/
```

**Pass criteria (spot-check against known state):**

| Section         | Expected                                                                              |
| --------------- | ------------------------------------------------------------------------------------- |
| Overview        | Pages: 26, Edges: 117, Entropy rate: ~2.27 bits, Classes: 1                          |
| Classes         | 1 recurrent (26 pages)                                                                |
| Orphans (≤10%)  | `llm-wiki-pattern` (π=0.001321), `quickstart` (π=0.005286), `adr-005-page-type-conventions-and-proposal-storage` (π=0.009690) |
| Sinks           | `(none)`                                                                              |
| Most central #1 | `analyze` (π=0.115652)                                                               |
| Suggestions     | At least one page with 3 suggestions listed                                           |

- Exit code 0
- Completes in < 1 second

---

## TC-17 · analyze — --suggest-top 0 skips commute section

```bash
time wikigraph analyze docs/ --suggest-top 0
```

**Pass criteria:**
- Output ends after the *Most central* section; no *Suggested missing links* section
- Completes in < 0.5 seconds (commute computation not run)
- Exit code 0

---

## TC-18 · analyze — --orphan-pct 0 shows only minimum-pi pages

```bash
wikigraph analyze docs/ --orphan-pct 0 --suggest-top 0
```

**Pass criteria:**
- Orphan section shows `adr-003-orphan-threshold` (π=0.006887) — the single lowest-π page
- Section header says `bottom 0%`

---

## TC-19 · analyze — --orphan-pct 1.0 shows all pages

```bash
wikigraph analyze docs/ --orphan-pct 1.0 --suggest-top 0 2>/dev/null | grep -c "→ add inbound"
```

**Pass criteria:**
- Count equals 23 (all pages shown as orphans at 100th percentile)

---

## TC-20 · --exclude propagates to all subcommands

**Goal:** Persistent flag is inherited; excluded pages disappear from every subcommand.

```bash
# Should reduce page count to 6 in each case
wikigraph graph   docs/ -e index -e log -e AGENTS -e how-to-docs-plan -o /dev/null 2>&1 | grep Pages
wikigraph goal    docs/ -e index -e log -e AGENTS -e how-to-docs-plan --goal analyze -o /dev/null 2>&1 | grep Pages
wikigraph export  docs/ -e index -e log -e AGENTS -e how-to-docs-plan -o /tmp/excl 2>&1 | grep Pages
wikigraph analyze docs/ -e index -e log -e AGENTS -e how-to-docs-plan --suggest-top 0 2>/dev/null | grep Pages
```

**Pass criteria:**
- All four lines print `Pages: 20`
- `how-to-docs-plan` absent from exported JSON nodes list:
  ```bash
  jq -e '[.nodes[].id] | index("how-to-docs-plan") | not' /tmp/excl.json
  ```

---

## TC-21 · alias wikilinks are parsed correctly

**Goal:** `[[slug|Display Name]]` links are resolved to the slug, not dropped.

```bash
mkdir /tmp/aliaswiki
printf '# Test\n\n[[some-page|Display Name]]\n' > /tmp/aliaswiki/test.md
printf '# Some Page\n' > /tmp/aliaswiki/some-page.md
wikigraph graph /tmp/aliaswiki -o /tmp/alias.html 2>&1
```

**Pass criteria:**
- Stderr prints `Pages: 2`
- HTML contains 2 nodes and 1 edge (`test` → `some-page`)
- No error

---

## TC-22 · dangling link does not panic

**Goal:** A wikilink to a non-existent page is silently ignored; no crash.

```bash
mkdir /tmp/tinywiki
printf '# Hello\n\n[[does-not-exist]]\n' > /tmp/tinywiki/hello.md
wikigraph graph /tmp/tinywiki -o /tmp/tiny.html
```

**Pass criteria:**
- Exit code 0
- Stderr prints `Pages: 1`
- HTML file written (single-node graph)

---

## TC-23 · recursive vault traversal (-r)

**Goal:** `-r / --recursive` scans nested subdirectories and skips hidden folders. See [[adr-006-recursive-vault-traversal]].

```bash
rm -rf /tmp/recwiki
mkdir -p /tmp/recwiki/folder1 /tmp/recwiki/.hidden
printf '# Note 1\n\n[[note2]]\n' > /tmp/recwiki/folder1/note1.md
printf '# Note 2\n\n[[note1]]\n' > /tmp/recwiki/folder1/note2.md
printf '# Hidden Note\n\n[[note1]]\n' > /tmp/recwiki/.hidden/note3.md
wikigraph analyze /tmp/recwiki -r --suggest-top 0 2>&1
```

**Pass criteria:**
- Output lists `Pages: 2` (finds `note1` and `note2`, skips `.hidden/note3`)
- Exit code 0

---

## TC-24 · duplicate slug collision in recursive mode

**Goal:** Clear error when duplicate file basenames exist in different subdirectories under `-r`.

```bash
rm -rf /tmp/dupwiki
mkdir -p /tmp/dupwiki/dirA /tmp/dupwiki/dirB
printf '# Note A\n' > /tmp/dupwiki/dirA/same.md
printf '# Note B\n' > /tmp/dupwiki/dirB/same.md
wikigraph analyze /tmp/dupwiki -r 2>&1
```

**Pass criteria:**
- Non-zero exit code
- Stderr contains `duplicate slug "same" found`

---

## TC-25 · multiple isolated communicating classes (knowledge silos)

**Goal:** Verify `wikigraph analyze` correctly detects and reports multiple isolated communicating classes (strongly connected components) across disconnected topic clusters. See [[communicating-classes]].

```bash
rm -rf /tmp/isolatedwiki
mkdir -p /tmp/isolatedwiki/clusterA /tmp/isolatedwiki/clusterB

# Cluster A: Physics (quantum-a <-> quantum-b)
printf '# Quantum A\n\n[[quantum-b]]\n' > /tmp/isolatedwiki/clusterA/quantum-a.md
printf '# Quantum B\n\n[[quantum-a]]\n' > /tmp/isolatedwiki/clusterA/quantum-b.md

# Cluster B: Cooking (pasta-recipe <-> sauce-recipe)
printf '# Pasta Recipe\n\n[[sauce-recipe]]\n' > /tmp/isolatedwiki/clusterB/pasta-recipe.md
printf '# Sauce Recipe\n\n[[pasta-recipe]]\n' > /tmp/isolatedwiki/clusterB/sauce-recipe.md

wikigraph analyze /tmp/isolatedwiki -r --suggest-top 0
```

**Pass criteria:**
- Overview reports `Pages: 4` and `Classes: 2`
- Communicating classes section lists two distinct recurrent classes (`Class 1` with 2 pages, `Class 2` with 2 pages)
- Exit code 0

---

## See Also

- [[graph]]
- [[analyze]]
- [[goal]]
- [[export]]

## Sources

- [`cmd_analyze.go`](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_analyze.go)
- [`cmd_graph.go`](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_graph.go)
- [`cmd_goal.go`](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_goal.go)
- [`cmd_export.go`](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_export.go)
