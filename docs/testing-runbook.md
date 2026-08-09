---
type: runbook
title: wikigraph Manual Test Runbook
description: Manual test plan covering all subcommands and edge cases, including graph, analyze, goal, and export.
tags: [testing, runbook, qa]
timestamp: 2026-08-09T06:54:46Z
---

# wikigraph — Manual Test Runbook

**Binary:** built locally from `~/repos/wikigraph/`  
**Wiki under test:** `docs/` in this repo (7 content pages)  
**Run all commands from:** `~/repos/wikigraph/` (the repo root)  
**Last verified:** 2026-08 against commit `0d68bb3`

---

## Prerequisites

### 1. Build the binary

```bash
cd ~/repos/wikigraph
go build -o wikigraph .
export PATH="$PWD:$PATH"
```

Verify:

```bash
wikigraph --help   # must list graph, goal, export, analyze
```

### 2. Stay in the repo root

All commands use `docs/` as the wiki path. Run everything from `~/repos/wikigraph`:

```bash
cd ~/repos/wikigraph
```

`docs/` contains 10 `.md` files across root and subdirectories. Three are excluded by default (`index`, `log`, `AGENTS` — matched by basename), leaving **7 content pages**:

```
adr/001-embedding-layer
how-to/analyze
how-to/export
how-to/goal
how-to/graph
proposal/how-to-docs-plan
testing-runbook
```

> **Slug format:** wikigraph scans subdirectories recursively. Slugs are relative paths without `.md` (e.g. `how-to/analyze`). Exclusions match on full slug OR basename.

> **Re-specify defaults when adding `--exclude`:** passing `-e` replaces the default list. Always include `-e index -e log -e AGENTS` alongside any custom exclusions.

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

**Goal:** Full wiki renders to a valid HTML file. See [[how-to/graph]].

```bash
wikigraph graph docs/ -o /tmp/wg_graph.html
open /tmp/wg_graph.html
```

**Pass criteria:**
- Stderr prints `Pages: 7` and `Written: /tmp/wg_graph.html`
- File exists and opens in browser showing a force-directed graph
- 7 labelled nodes visible; nodes sized differently (stationary dist reflects link structure)
- 3 colours visible (3 communicating classes)
- Exit code 0

---

## TC-03 · graph — --exclude removes pages

**Goal:** Excluding a content page reduces node count.

```bash
wikigraph graph docs/ -e index -e log -e AGENTS -e testing-runbook -o /tmp/wg_excl.html
```

**Pass criteria:**
- Stderr prints `Pages: 6`
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

> Note: the default title is `<dir> wiki` where `<dir>` is the last path component. Running with `docs/` gives title `docs wiki`.

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

**Goal:** Subgraph centred on a goal page. See [[how-to/goal]].

```bash
wikigraph goal docs/ --goal how-to/analyze --top 5 -o /tmp/wg_goal.html
open /tmp/wg_goal.html
```

**Pass criteria:**
- Stderr prints `Pages: 7` and `Written: /tmp/wg_goal.html (5 nodes)`
- Browser shows exactly 5 nodes
- `how-to/analyze` node is present
- Exit code 0

---

## TC-08 · goal — multiple goals

**Goal:** Two goal pages, both present in output.

```bash
wikigraph goal docs/ \
  --goal how-to/analyze --goal how-to/graph \
  --top 5 -o /tmp/wg_goal2.html
```

**Pass criteria:**
- Stderr shows `(5 nodes)`
- Both `how-to/analyze` and `how-to/graph` present in rendered graph

---

## TC-09 · goal — unknown slug prints valid list

**Goal:** Clear error when a slug doesn't exist.

```bash
wikigraph goal docs/ --goal not-a-real-page -o /tmp/x.html
```

**Pass criteria:**
- Prints `unknown --goal slug "not-a-real-page"`
- Lists valid slugs (e.g. `how-to/analyze`, `adr/001-embedding-layer`)
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

**Goal:** Valid node-link JSON with correct shape. See [[how-to/export]].

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
- 7 nodes, 18 links

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
- `wg_nodes.csv` has 7 data rows (one per page)
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
- Node lines: `"slug" [weight=N.NNNNNN];` (slugs may contain `/`, e.g. `"how-to/analyze"`)
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
- Link count is lower than default (default: 18 links; with `--min-edge 0.5`: 4 links)

---

## TC-16 · analyze — full report

**Goal:** All six sections printed; known facts match. See [[how-to/analyze]].

```bash
wikigraph analyze docs/
```

**Pass criteria (spot-check against known state):**

| Section         | Expected                                                                              |
| --------------- | ------------------------------------------------------------------------------------- |
| Overview        | Pages: 7, Edges: 18, Entropy rate: ~1.02 bits, Classes: 3                            |
| Classes         | 1 recurrent (5 pages), 2 transient (1 page each: `testing-runbook`, `proposal/how-to-docs-plan`) |
| Orphans (≤10%)  | `proposal/how-to-docs-plan` and `testing-runbook` (π=0.000000)                       |
| Sinks           | `(none)`                                                                              |
| Most central #1 | `how-to/analyze` (π=0.375000)                                                        |
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
- Orphan section shows only the page(s) with the absolute lowest π (`proposal/how-to-docs-plan` and `testing-runbook`, both π=0.000000)
- Section header says `bottom 0%`

---

## TC-19 · analyze — --orphan-pct 1.0 shows all pages

```bash
wikigraph analyze docs/ --orphan-pct 1.0 --suggest-top 0 2>/dev/null | grep -c "→ add inbound"
```

**Pass criteria:**
- Count equals 7 (all pages shown as orphans at 100th percentile)

---

## TC-20 · --exclude propagates to all subcommands

**Goal:** Persistent flag is inherited; excluded pages disappear from every subcommand.

```bash
# Should reduce page count to 6 in each case
wikigraph graph   docs/ -e index -e log -e AGENTS -e how-to/analyze -o /dev/null 2>&1 | grep Pages
wikigraph goal    docs/ -e index -e log -e AGENTS -e how-to/analyze --goal how-to/graph -o /dev/null 2>&1 | grep Pages
wikigraph export  docs/ -e index -e log -e AGENTS -e how-to/analyze -o /tmp/excl 2>&1 | grep Pages
wikigraph analyze docs/ -e index -e log -e AGENTS -e how-to/analyze --suggest-top 0 2>/dev/null | grep Pages
```

**Pass criteria:**
- All four lines print `Pages: 6`
- `how-to/analyze` absent from exported JSON nodes list:
  ```bash
  jq -e '[.nodes[].id] | index("how-to/analyze") | not' /tmp/excl.json
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

- [[how-to/graph]]
- [[how-to/analyze]]
- [[how-to/goal]]
- [[how-to/export]]
