---
type: runbook
title: wikigraph Manual Test Runbook
description: Manual test plan covering all subcommands and edge cases, including graph, analyze, goal, and export.
tags: [testing, runbook, qa]
timestamp: 2026-08-09T06:54:46Z
---

# wikigraph — Manual Test Runbook

**Binary:** `~/go/bin/wikigraph` (installed via `go install github.com/stephen-mcelhose/wikigraph@latest`)  
**Wiki under test:** any wiki directory with 36 content pages (e.g. `quantum-go/wiki`)  
**Last verified:** 2026-08 against commit `3d0d08e`

---

## Prerequisites

Install the binary:

```bash
go install github.com/stephen-mcelhose/wikigraph@latest
```

All commands below assume you are **in the wiki directory**:

```bash
cd path/to/your/wiki
```

Then verify the binary is available:

```bash
which wikigraph          # must print ~/go/bin/wikigraph
wikigraph --help         # must list graph, goal, export, analyze
```

---

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

**Goal:** Full wiki renders to a valid HTML file. See [[How to generate an interactive wiki graph]].

```bash
wikigraph graph . -o /tmp/wg_graph.html
open /tmp/wg_graph.html
```

**Pass criteria:**
- Stderr prints `Pages: 36` and `Written: /tmp/wg_graph.html`
- File exists and opens in browser showing a force-directed graph
- 36 labelled nodes visible; nodes sized differently (stationary dist)
- All nodes coloured the same (one communicating class)
- Exit code 0

---

## TC-03 · graph — --exclude removes pages

**Goal:** Excluding a content page reduces node count.

```bash
wikigraph graph . -e index -e log -e AGENTS -e gate-zoo -o /tmp/wg_excl.html
```

**Pass criteria:**
- Stderr prints `Pages: 35`
- `gate-zoo` node absent from rendered graph

---

## TC-04 · graph — --min-edge reduces edge clutter

**Goal:** Raising the edge threshold produces a sparser graph.

```bash
wikigraph graph . --min-edge 0.10 -o /tmp/wg_sparse.html
```

**Pass criteria:**
- File renders with visibly fewer edges than TC-02 output
- No error

---

## TC-05 · graph — --sed patches the HTML

**Goal:** sed expressions are applied to the output.

```bash
wikigraph graph . \
  -s 's/wiki wiki/TEST TITLE/' \
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

## TC-07 · goal — learning path to shors-algorithm

**Goal:** 12-node subgraph centred on Shor's algorithm. See [[How to find a learning path through your wiki]].

```bash
wikigraph goal . --goal shors-algorithm --top 12 -o /tmp/wg_goal.html
open /tmp/wg_goal.html
```

**Pass criteria:**
- Stderr prints `Pages: 36` and `Written: /tmp/wg_goal.html (12 nodes)`
- Browser shows exactly 12 nodes
- `shors-algorithm` node is present
- Exit code 0

---

## TC-08 · goal — multiple goals

**Goal:** Two goal pages, both present in output.

```bash
wikigraph goal . \
  --goal shors-algorithm --goal grovers-algorithm \
  --top 8 -o /tmp/wg_goal2.html
```

**Pass criteria:**
- Stderr shows `(8 nodes)`
- Both `shors-algorithm` and `grovers-algorithm` present in rendered graph

---

## TC-09 · goal — unknown slug prints valid list

**Goal:** Clear error when a slug doesn't exist.

```bash
wikigraph goal . --goal not-a-real-page -o /tmp/x.html
```

**Pass criteria:**
- Stderr prints `unknown --goal slug "not-a-real-page"`
- Stderr lists valid slugs
- Non-zero exit code; no output file written

---

## TC-10 · goal — missing --goal flag

**Goal:** Requires at least one goal.

```bash
wikigraph goal . -o /tmp/x.html
```

**Pass criteria:**
- Error: `at least one --goal slug is required`
- Non-zero exit code

---

## TC-11 · export — JSON

**Goal:** Valid node-link JSON with correct shape. See [[How to export your wiki graph for external tools]].

```bash
wikigraph export . --format json -o /tmp/wg
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
wikigraph export . --format csv -o /tmp/wg
head /tmp/wg_nodes.csv
head /tmp/wg_edges.csv
```

**Pass criteria:**
- `wg_nodes.csv` header: `slug,pi,class`
- `wg_edges.csv` header: `source,target,probability`
- Both have 36 data rows in nodes (one per page)
- Values are numeric floats; no empty cells

---

## TC-13 · export — DOT

**Goal:** Valid Graphviz DOT output.

```bash
wikigraph export . --format dot -o /tmp/wg
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
wikigraph export . --format xml -o /tmp/wg
```

**Pass criteria:**
- Error: `unknown format "xml": choose json, csv, or dot`
- Non-zero exit code

---

## TC-15 · export — --min-edge filters edges

```bash
wikigraph export . --format json --min-edge 0.5 -o /tmp/wg_sparse
jq '.links | length' /tmp/wg_sparse.json
```

**Pass criteria:**
- Link count is significantly lower than with default `--min-edge 0.005`

---

## TC-16 · analyze — full report

**Goal:** All six sections printed; known facts match. See [[How to analyse your wiki's health]].

```bash
wikigraph analyze .
```

**Pass criteria (spot-check against known state):**

| Section         | Expected                                                     |
| --------------- | ------------------------------------------------------------ |
| Overview        | Pages: 36, Edges: 295, Entropy rate: ~2.20 bits, Classes: 1 |
| Classes         | 1 recurrent class containing all 36 pages                   |
| Orphans (≤10%)  | None shown, or only `how-to-add-a-new-gate` (others fixed in 2026-08-09 lint pass) |
| Sinks           | `algorithm-comparison`, `fuzz-testing`, `qelib1-standard-gates`, `verification-tests` |
| Most central #1 | `composite-gates` (π ≈ 0.096)                               |
| Suggestions     | At least one page with 3 suggestions listed                  |

- Exit code 0
- Completes in < 2 seconds

---

## TC-17 · analyze — --suggest-top 0 skips commute section

```bash
time wikigraph analyze . --suggest-top 0
```

**Pass criteria:**
- Output ends after the *Most central* section; no *Suggested missing links* section
- Completes in < 0.5 seconds (commute computation not run)
- Exit code 0

---

## TC-18 · analyze — --orphan-pct 0 shows only minimum-pi pages

```bash
wikigraph analyze . --orphan-pct 0 --suggest-top 0
```

**Pass criteria:**
- Orphan section shows only the page(s) with the absolute lowest π
- Section header says `bottom 0%`

---

## TC-19 · analyze — --orphan-pct 1.0 shows all pages

```bash
wikigraph analyze . --orphan-pct 1.0 --suggest-top 0 2>/dev/null | grep -c "→ add inbound"
```

**Pass criteria:**
- Count equals 36 (all pages shown as orphans at 100th percentile)

---

## TC-20 · --exclude propagates to all subcommands

**Goal:** Persistent flag is inherited; excluded pages disappear from every subcommand.

```bash
# Should reduce page count to 35 in each case
wikigraph graph   . -e index -e log -e AGENTS -e gate-zoo -o /dev/null 2>&1 | grep Pages
wikigraph goal    . -e index -e log -e AGENTS -e gate-zoo --goal shors-algorithm -o /dev/null 2>&1 | grep Pages
wikigraph export  . -e index -e log -e AGENTS -e gate-zoo -o /tmp/excl 2>&1 | grep Pages
wikigraph analyze . -e index -e log -e AGENTS -e gate-zoo --suggest-top 0 2>/dev/null | grep Pages
```

**Pass criteria:**
- All four lines print `Pages: 35`
- `gate-zoo` absent from exported JSON nodes list:
  ```bash
  jq -e '[.nodes[].id] | index("gate-zoo") | not' /tmp/excl.json
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

## See Also

- [[How to generate an interactive wiki graph]]
- [[How to analyse your wiki's health]]
- [[How to find a learning path through your wiki]]
- [[How to export your wiki graph for external tools]]
