---
type: how-to
title: Quickstart Guide — Minimal LLM-Wiki Setup & Wikigraph Analysis
description: Step-by-step guide to setting up a minimal llm-wiki skill, visualising and analysing the wiki, and feeding analysis output into maintenance prompts.
tags: [quickstart, llm-wiki, setup, graph, analyze, maintenance]
timestamp: 2026-08-09T17:15:00Z
---

# Quickstart Guide — Minimal LLM-Wiki Setup & Wikigraph Analysis

This quickstart guides you through establishing a minimal `llm-wiki` knowledge base following the [[llm-wiki-pattern]], inspecting and visualising it with `wikigraph`, and feeding health analysis back into your LLM maintenance workflow.

For comprehensive test coverage and edge cases across all scenarios, refer to [[testing-runbook]].

---

## 1. Setting up a minimal `llm-wiki`

An `llm-wiki` follows the Karpathy pattern ("Obsidian as IDE, LLM as programmer, wiki as codebase"). It requires three base components in a dedicated root directory (e.g. `wiki/`):

### Step 1: Create `AGENTS.md` (Schema)
Create `wiki/AGENTS.md` to define conventions and domain:

```markdown
# Wiki Schema

## Domain
General knowledge base maintained by an LLM skill.

## Conventions
- **Page slugs**: kebab-case flat files (`<slug>.md`)
- **Frontmatter**: OKF format (`type`, `title`, `description`, `timestamp`)
- **Cross-references**: `[[slug]]` or `[[slug|Display Name]]` wikilinks
- **Sources**: Every page ends with a `## Sources` section
```

### Step 2: Initialize `index.md` and `log.md`
Create `wiki/index.md` (catalog):

```markdown
# Wiki Index

| Page | Type | Summary |
| ---- | ---- | ------- |
| [[welcome]] | `concept` | Welcome page for the wiki |
```

And `wiki/log.md` (append-only audit log):

```markdown
# Wiki Log

## [2026-08-09] init | wiki initialized
```

### Step 3: Add initial content pages
Create sample pages with `[[slug]]` cross-references:

`wiki/welcome.md`:
```markdown
---
type: concept
title: Welcome Page
description: Starting node for the wiki.
timestamp: 2026-08-09T12:00:00Z
---

# Welcome

Welcome to the minimal wiki. Read more about [[llm-wiki-pattern]].

## Sources
- Internal setup
```

`wiki/llm-wiki-pattern.md`:
```markdown
---
type: concept
title: LLM Wiki Pattern
description: Compounding knowledge base architecture.
timestamp: 2026-08-09T12:00:00Z
---

# LLM Wiki Pattern

The LLM acts as the programmer maintaining interlinked markdown pages. See [[welcome]].

## Sources
- https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f
```

---

## 2. Visualising a wiki

Use `wikigraph graph` to generate a self-contained, interactive HTML force-directed graph:

```bash
wikigraph graph wiki/ -o /tmp/quickstart_graph.html --title "My Wiki Graph"
```

Open `/tmp/quickstart_graph.html` in any browser to inspect and interact with the graph (drag nodes, zoom, pan, hover for degree & $\pi$ metrics).

### Explanation of visual properties

| Visual Element | Graph Property | Meaning in `wikigraph` |
| -------------- | -------------- | ---------------------- |
| **Node Size** | Stationary Distribution ($\pi$) | Represents systemic importance / visit frequency under a random walk. Hub pages appear larger. |
| **Node Colour** | Communicating Class | Pages in the same strongly connected component share a colour. Isolated clusters show distinct colours. |
| **Edge Width** | Transition Probability | Likelihood of a random walker moving from one page to another along that `[[wikilink]]`. |

### Visualisation flags & techniques

- **Filter weak links**: Pass `--min-edge 0.05` to hide low-probability transition edges and clean up dense graphs.
- **Custom styling / sed post-processing**: Pass `-s` / `--sed` expressions to modify the rendered HTML output on macOS/Linux:
  ```bash
  wikigraph graph wiki/ -o /tmp/graph.html -s 's/background: #1a1a2e/background: #0f172a/'
  ```

For full details on graph configuration, see [[graph]].

---

## 3. Analysing a wiki

Run `wikigraph analyze` to audit health signals, orphan pages, isolated clusters, and missing link candidates:

```bash
wikigraph analyze wiki/
```

### Key output sections
1. **Overview**: Total pages, total directed edges, entropy rate, and class count.
2. **Communicating classes**: Detects disconnected clusters (recurrent vs transient).
3. **Orphan pages**: Pages in the bottom percentile by stationary probability ($\pi$) needing inbound links.
4. **Sink pages**: Pages with zero outbound links (dead ends).
5. **Most central**: Top hubs by stationary distribution ($\pi$).
6. **Suggested missing links**: Page pairs with low commute time that are not yet directly linked.

For detailed interpretations of these metrics, see [[analyze]].

---

## 4. Feeding analysis into a maintenance prompt

Integrate `wikigraph analyze` output directly into your LLM skill or maintenance prompt (e.g. during a `lint` or `refactor` operation):

### Workflow
1. Run `wikigraph analyze wiki/` and save the report or pipe it.
2. Provide the report output to the LLM agent alongside `AGENTS.md`.
3. Instruct the LLM to address structural defects identified in the report:
   - Add inbound `[[wikilinks]]` to listed orphan pages.
   - Add outbound links to sink pages or transient communicating classes.
   - Evaluate high-priority "Suggested missing links" and insert cross-references where semantically appropriate.
   - Log edits in `wiki/log.md` and update `wiki/index.md`.

---

## Covered Scenarios & Testing Runbook Reference

For test automation, expected output benchmarks, and verification steps across all primary usage scenarios, refer to [[testing-runbook]]:

| Scenario | Runbook Test Cases |
| -------- | ------------------ |
| **Setting up a wiki** | [[testing-runbook]] (TC-21, TC-22, TC-23) |
| **Visualising a wiki** | [[testing-runbook]] (TC-02, TC-03, TC-04) |
| **Analysing a wiki** | [[testing-runbook]] (TC-16, TC-17, TC-18, TC-25) |
| **Feeding analysis to maintenance prompt** | [[testing-runbook]] (TC-16) |

---

## Sources

- [[testing-runbook]]
- [[graph]]
- [[analyze]]
- [[AGENTS]]
