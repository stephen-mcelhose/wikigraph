---
type: concept
title: Mean First Passage Time
description: MFPT(i,j) — expected steps for a random walk from page i to first reach page j — used by the goal subcommand to rank pages and by analyze to suggest missing links.
resource: cmd_goal.go
tags: [mfpt, commute-time, goal, random-walk, fundamental-matrix]
timestamp: 2026-08-09T07:31:56Z
---

# Mean First Passage Time

**MFPT(i, j)** is the expected number of steps a [[random-walk|random walker]] starting at
page i needs before it first lands on page j. It is computed from the
**fundamental matrix** of the Markov chain; the implementation lives in
[[catrace|`catrace.Kernel.MeanFirstPassage`]].

Low MFPT means pages are structurally close — the random walk naturally
travels between them quickly. High or infinite MFPT means pages are in
different communicating classes, or rarely visited together.

## Use in `goal` (cmd_goal.go)

The `goal` subcommand finds the N pages closest to one or more target pages:

1. For each goal page g, compute MFPT(i, g) for every other page i
2. Each page scores `min over goals of MFPT(i, goal)` — the nearest goal wins
3. Goal pages score 0 (already at the destination)
4. Pages unreachable from some transient class score ∞ and are ranked last
5. The top N pages (plus all goal pages) form the **subgraph**

After selection, a **trace kernel** (`kern.Trace(subset, tol)`) projects the
full transition matrix onto the subset, producing a self-contained Markov
chain for visualization. Node size in the rendered graph reflects the
[[stationary-distribution|stationary distribution]] of the trace kernel — not the original full chain.

See [[goal]] for the user-facing interface; see [[markov-model]] for how the
underlying kernel is built.

## Use in `analyze` — commute time (cmd_analyze.go)

**[[commute-time|Commute time]]** CT(i, j) = MFPT(i, j) + MFPT(j, i) is symmetric and acts
as a graph distance metric. It is used in the *Suggested missing links*
section:

- For every pair (i, j) not yet directly linked, compute CT(i, j)
- Pages with low commute time are structurally close but not explicitly connected
- A link between them would be informationally useful — a reader would benefit
- Pairs are ranked ascending by commute time; the top `--suggest-top` (default 3) are shown

```
commute time low  → pages already close in graph structure → good link candidate
commute time high → pages far apart or in different classes → link less useful
```

The commute time is skipped entirely when `--suggest-top 0` is passed,
which can save significant computation on large wikis.

## Infinite MFPT

MFPT(i, j) is infinite when page j is not reachable from page i. This
happens in two situations: i is in a recurrent class that does not contain
j (recurrent classes are closed — a walker never leaves), or i is in a
transient class with no directed path toward j's class. In `goal`, such
pages are left with score `1e18` and typically excluded from the top-N
selection. In `analyze`, `CommuteTime` returns an error for such pairs
and the pair is skipped.

## Sources

- `cmd_goal.go`
- `cmd_analyze.go`
- `catrace` library: `MeanFirstPassage`, `CommuteTime`, `Trace`
