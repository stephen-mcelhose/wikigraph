---
title: How-To Page Rubric
description: Quality criteria for evaluating and writing how-to guides in this wiki.
scope: local
---

# How-To Page Rubric

Use this rubric when reviewing or writing `type: how-to` pages.
Adapted from the Divio documentation system and the Google developer docs style guide.

---

## 1. Framing — does it answer a concrete question?

A how-to guides a reader through a specific task. It is not a tutorial
(learning-oriented) or a reference (information-oriented).

| Criterion                                                     | Pass | Fail example                                     |
| ------------------------------------------------------------- | ---- | ------------------------------------------------ |
| Title or Goal section names the task as an action             | ✓    | "The graph subcommand" (topic, not task)         |
| Reader can judge within 5 s whether the page applies to them  | ✓    | No goal statement, jumps straight into steps     |
| Scope is narrow enough to complete in one sitting             | ✓    | "How to document your entire project"            |

---

## 2. Prerequisites — complete and honest

| Criterion                                                     | Pass | Fail example                                     |
| ------------------------------------------------------------- | ---- | ------------------------------------------------ |
| Lists every tool, binary, or file the reader needs            | ✓    | Assumes `wikigraph` is on PATH without saying so |
| States the minimum knowledge assumed                          | ✓    | Uses Markov terminology without a link           |
| Does NOT repeat background theory that belongs in a concept page | ✓ | Two paragraphs explaining MFPT in a how-to       |

---

## 3. Steps — numbered, imperative, independently runnable

| Criterion                                                     | Pass | Fail example                                     |
| ------------------------------------------------------------- | ---- | ------------------------------------------------ |
| Steps are numbered and use imperative mood ("Run", not "You run") | ✓ | "Next we will want to run..."                   |
| Each step has a concrete command or action                    | ✓    | "Configure the flags appropriately"              |
| Each step shows expected output or a way to verify success    | ✓    | Step ends with no indication of what should happen |
| Steps are executable in order with no hidden prerequisites    | ✓    | Step 3 depends on a variable set in Step 1 with no link |
| Steps avoid branching ("if X, do Y, else do Z")               | ✓    | Multiple conditional branches in one step        |

---

## 4. Verification — reader can confirm they succeeded

| Criterion                                                     | Pass | Fail example                                     |
| ------------------------------------------------------------- | ---- | ------------------------------------------------ |
| Dedicated `## Verification` section (or equivalent)           | ✓    | No way to tell if it worked                      |
| Checks are observable (file exists, stderr prints X, UI shows Y) | ✓ | "It should work correctly"                      |
| Expected values match current real output (not stale)         | ✓    | `Pages: 19` when wiki has 20 pages               |

---

## 5. Troubleshooting — anticipates the common failures

| Criterion                                                     | Pass | Fail example                                     |
| ------------------------------------------------------------- | ---- |--------------------------------------------------|
| `## Troubleshooting` section present                          | ✓    | Page ends after Verification                     |
| Each entry has symptom → cause → fix                          | ✓    | "Check your PATH" (no symptom, no cause)         |
| Covers at least the top 2–3 realistic failure modes           | ✓    | Only covers the happy path                       |

---

## 6. Cross-references — links to related pages, not redundant content

| Criterion                                                     | Pass | Fail example                                     |
| ------------------------------------------------------------- | ---- | ------------------------------------------------ |
| `## See also` (or inline `[[slug]]` links) to concept pages   | ✓    | No links out; reader is stuck if they need more  |
| Links to sibling how-to guides where workflow continues       | ✓    | `analyze` guide doesn't link to `goal` guide     |
| Does NOT duplicate content that already lives in a concept page | ✓  | Full MFPT explanation copied from `mfpt.md`      |

---

## 7. OKF Frontmatter — required fields

```yaml
type: how-to
title: <verb phrase, e.g. "Render your wiki as an interactive graph">
description: <one sentence, what the reader will achieve>
tags: [how-to, <subcommand>]
timestamp: <ISO 8601>
```

`resource` is optional but recommended when the page documents a specific CLI flag or API surface.

---

## 8. Sources — traceable to code or specs

| Criterion                                                     | Pass | Fail example                                     |
| ------------------------------------------------------------- | ---- | ------------------------------------------------ |
| `## Sources` section present                                  | ✓    | Page ends without sources                        |
| Each source is a file path, GitHub URL, or issue link         | ✓    | "wikigraph source code" (not traceable)          |
| Sources are specific enough to be verified                    | ✓    | `cmd_graph.go` ✓ vs. "the codebase" ✗            |

---

## Scoring

Count passes. A page needs 14+ / 18 to be considered complete.
Flag any `## Verification` failure or stale expected value as a blocker regardless of total score.
