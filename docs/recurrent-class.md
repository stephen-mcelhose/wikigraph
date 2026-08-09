---
type: concept
title: Recurrent vs Transient Classes
description: Recurrent classes are the permanent homes of a Markov random walk — pages here are visited infinitely often. Transient classes are visited finitely many times before the walk escapes forever.
resource: cmd_analyze.go
tags: [recurrent, transient, markov, stationary-distribution, communicating-classes]
timestamp: 2026-08-09T07:31:56Z
---

# Recurrent vs Transient Classes

Every [[communicating-classes|communicating class]] in a Markov chain is
either **recurrent** or **transient**. The distinction determines whether a
page gets a share of the [[stationary-distribution|stationary distribution]] π
— and therefore whether it appears in centrality rankings,
[[commute-time|commute-time]] suggestions, and [[mfpt]] calculations.

## Recurrent class

A class is **recurrent** if, once the [[random-walk|random walk]] enters it, the probability
of ever leaving is zero. The walk returns to every page in the class infinitely
often. Formally: for any state i in a recurrent class, the walk returns to i
with probability 1.

In wiki terms: a recurrent class is a cluster of pages that are all mutually
linked (directly or via paths), and no page in the cluster links exclusively
outward to a different cluster. The walk circulates within the cluster forever.

**Stationary distribution π is defined on recurrent classes.** Each page in a
recurrent class gets a π > 0, reflecting how often the random walk lands there
in the long run. Pages with high π are hubs; pages with low π are peripheral
but still visited.

## Transient class

A class is **transient** if the random walk eventually escapes to a recurrent
class and never returns. Pages in transient classes are visited a finite number
of times in total — their long-run visit frequency is 0, so π = 0.

In wiki terms: a transient class is typically a cluster of pages that links
into the main graph but has no path back. The reader can navigate *out* of
the cluster but the random walk never returns.

**A page can have many outgoing links and still be transient** — what matters
is whether there is a return path back to the class, not how many links exist.

## How wikigraph labels them

After Tarjan SCC decomposition, `catrace.Classes` tests each SCC:

```go
// A class is recurrent if no state in it has a transition to another class.
// Any such transition makes the class transient.
```

The `analyze` output labels each class:

```
Class 1 (recurrent): 9 page(s)         ← healthy
Class 2 (transient — add links out of this class): 2 page(s)  ← broken
```

The fix message "add links out of this class" is slightly misleading — the
real fix is to add links *from* the recurrent class *into* the transient class,
so the walk can return. Adding links out of the transient class into the
recurrent class alone does not help: you need the cycle to close.

## Practical example

```
testing-runbook  →  analyze  →  goal  →  testing-runbook   ← recurrent (cycle exists)
how-to-docs-plan →  analyze                                 ← transient (no return path)
```

`how-to-docs-plan` links to `analyze` but nothing links back to `how-to-docs-plan`
from within the main cluster. It is stranded. Fix: add a link from any recurrent
page to `how-to-docs-plan`.

## Why it matters for wiki health

| Scenario | Impact |
| -------- | ------ |
| Page in transient class | π = 0; never appears in centrality rankings; commute time from it to recurrent pages is ∞ |
| Two separate recurrent classes | Each has its own π; the wiki is two disconnected knowledge graphs |
| One recurrent class, all pages | Ideal — full mutual reachability, meaningful π for every page |

The communicating class count in `wikigraph analyze` is the headline health
number. **Classes: 1** with all pages recurrent is the goal.

## Sources

- `cmd_analyze.go` — `transientSet` construction, recurrent/transient labelling
- `catrace` library — `Classes.Recurrent`, `Classes.Transient`
