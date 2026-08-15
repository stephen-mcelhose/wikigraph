---
type: concept
title: From Knowledge Graphs to PDA Agents
description: How wikigraph’s PageRank walk on a wiki is the same Markov craft as catrace’s PDA agent model — and why that leap matters for measuring, comparing, and diagnosing agents and networks of agents.
tags: [pda, agents, pagerank, personalized-pagerank, catrace, markov-model, knowledge-graph, multi-agent]
timestamp: 2026-08-15T21:38:39Z
resource: https://github.com/stephen-mcelhose/catrace
---

# From Knowledge Graphs to PDA Agents

## Overview

`wikigraph` and `catrace` share one craft: **row-stochastic kernels, stationary
mass, teleporting / Personalized PageRank, MFPT, and entropy**. What looks like
“wiki health” in wikigraph is the same mathematics used in catrace to model
**agents** (and networks of agents) via the **PDA triplet** — Perception,
Decision, Action.

The intuitive leap:

> Analyzing a knowledge graph as a random walk is the training ground for
> analyzing an agent’s perceive→decide→act loop — and then a mesh of such loops —
> as Markov dynamics you can measure, compare, and diagnose.

[[adr-012-teleporting-pagerank-default]] aligned the wiki CLI with full
α-damping and `--seed`: the same **intent** ($v$) and **intentionality weight**
($\alpha$) dials catrace uses for goal-directed agents. That does not make
ADR-012 an agents ADR — it makes the wiki walk practice the *same operators*
you will use on PDA kernels, without a second math.

## Why this leap matters

Most people meet PageRank as “which pages matter.” That is useful but small.
The deeper move is:

1. **A knowledge graph is a dynamics object.** Links are not just retrieval
   edges; they define a transition law. Fixing sinks and orphans is *repair on
   that law* — the kernel you chose to identify with the wiki (usually a
   stand-in for a closed loop like $Q$ or a world walk), **not** an edit to
   the three PDA factors **P**, **D**, and **A** separately unless you have
   already factored the agent that way.
2. **Intent is a first-class dial.** Uniform teleport is mild global curiosity.
   Seeding recommend / goal pages is already orchestration: top-down preference
   over a world of documents. On a joint agent kernel, the same dial is
   control vs emergence.
3. **Networks of agents are coupled walks.** Once one loop ($Q = DAP$) is
   measurable, joint / shared-world kernels are a further scale — with the
   usual product-space cost. Hubs, sinks, and MFPT remain useful *operator*
   language there; wikigraph itself stops at one editable kernel.

Wikigraph is the place you *see* those quantities on something editable
(markdown). catrace is where the same quantities live on agents. ADR-012’s
PageRank defaults are what make the CLI dials match that craft.

## Key Properties

### Same toolkit, two readings

| Quantity | Knowledge-graph reading (`wikigraph`) | Agent reading (catrace PDA) |
| -------- | ------------------------------------- | --------------------------- |
| States | Wiki pages (slugs) | Experience / world / action states |
| Dynamics kernel | Teleporting walk on raw wikilinks | Qualia kernel $Q = DAP$ (or world $W = PDA$) |
| $\pi$ (uniform restart) | Global PageRank — long-run attention | Long-run experience under **pure dynamics** (no imposed goal) |
| Restart $v$ (`--seed`) | Personalized PageRank seeds | **Goal / intent** distribution |
| $\alpha$ (`--alpha`) | Teleport probability (default 0.15) | **Intentionality weight** — how hard goals override dynamics |
| MFPT / [[goal]] | Learning-path / access cost between pages | Expected steps to a target experience or task-complete state |
| Entropy rate | How “surprising” each step of the walk is | Behavioral uncertainty under dynamics (and under intent, on the teleporting chain) |
| Sinks / orphans / hubs ([[analyze]]) | Structural wiki smells | Dead-end policies, under-linked percepts, attractor states |
| Multi-page clusters | Communicating structure of the vault | Coupled regions of experience — precursor to joint / multi-agent kernels |

Nothing in that table requires a different linear-algebra library. `wikigraph`
already calls `NewTeleportingKernelFromAdj`, `Stationary`, `MeanFirstPassage`,
and `EntropyRate` on [[catrace]].

**Shared operators ≠ shared state meaning.** The left and right columns use the
same *operators* (stationary mass, PPR, MFPT, entropy, SCC-style connectivity).
They do **not** identify wiki pages with experience states $X$, world states
$W$, or actions $G$. The leap is: learn the operators on an editable one-kernel
object (the wiki), then apply them on PDA kernels whose states mean something
else. Analogy of craft; not identity of ontology.

### The PDA triplet (minimal)

catrace models one agent as three kernels between three spaces:

| Kernel | Map | Meaning |
| ------ | --- | ------- |
| **P** | $W \to X$ | Perception: world → experience |
| **D** | $X \to G$ | Decision: experience → action |
| **A** | $G \to W$ | Action: action → next world |

Composing them yields square “same-space” kernels for the closed loop:

- $Q = DAP$ — experience → experience (**qualia** kernel; usual single-agent analysis)
- $S = APD$ — action → action (strategy)
- $W = PDA$ — world → world (natural for **multi-agent** joint dynamics)

A wiki random walk is not a full PDA agent by itself. It is the cleanest
**one-kernel slice** of the story: a dynamics object you can see, edit (add
`[[wikilinks]]`), and re-measure. That slice is exactly where people build
intuition for $\pi$, $\alpha$, $v$, MFPT, and entropy *before* lifting to $Q$
or a joint $W$.

### Knowledge graph as grounding for perception

In the PDA picture the three spaces are distinct: world $W$, experience $X$,
actions $G$. The vault of markdown pages is best read as a **world fragment**
(documents / tickets / ADRs living in $W$) — not as the perception kernel
itself.

**P** is the map $W \to X$: given world state $w$, what does the agent
*experience*? Perception noise is off-diagonal mass in **P** (world is $w$ but
experience is $x \neq$ the faithful percept). A sparse, sink-ridden knowledge
graph does not automatically *equal* a bad **P**; it is a poor **substrate**
that a retrieval-style **P** must read. Only once you specify how the agent
grounds experience in that graph (e.g. “retrieve / attend over linked pages”)
does wiki hygiene become a lever on **P**.

So wiki work is not a metaphor bolted on later — with that model made explicit:

1. **Improve the world graph** (fix sinks, link orphans, strengthen hubs) →
   richer, better-connected substrate in $W$ for whatever retrieval / grounding
   map you use as **P**.
2. **Measure the walk** (`analyze` / `graph`) → diagnostics on that world
   digraph (attractors, dead ends, access cost, entropy) — the same *operators*
   you will later want on $Q$ or a joint world kernel, not a claim that the
   wiki *is* **P**.
3. **Add intent** (`--seed`, $\alpha$) → the same PPR move as “persistent goals”
   on an agent kernel.

### Intent: why `--seed` and `--alpha` matter beyond viz

On an ergodic kernel, plain power iteration forgets its start: every $x_0$
reaches the same $\pi$. That cannot encode an agent that **keeps** preferring
certain states.

Personalized PageRank is the minimal fix (catrace: *Personalized PageRank and
Agent Modeling*):

$$
x \leftarrow \alpha\, v + (1-\alpha)\, x\, P
$$

| Symbol | Wiki CLI | Agent story |
| ------ | -------- | ----------- |
| $v$ | `--seed` slugs | Goal distribution |
| $\alpha$ | `--alpha` | How strongly intent overrides trained dynamics |
| Fixed point | Node mass / “most central” under PPR | Long-run experience under dynamics **and** intent |

Default wikigraph ($v$ uniform, $\alpha = 0.15$) is “mild global curiosity” —
Google-style damping. Seeding recommend pages (as in firehose-graph) is already
an **orchestrator-style pull**: top-down preference over the document world.
The same dial, on a joint agent kernel, is top-down control vs emergent worker
dynamics.

### From one walk to a network of agents

wikigraph models **one** teleporting kernel on pages. Multi-agent work in
catrace lives on joint kernels (often over product or shared world spaces).
Those state spaces grow fast; this page only claims that the *same operators*
($\pi$, PPR, MFPT, entropy, connectivity) remain the diagnostic language —
not that `wikigraph` builds product kernels or that the leap is turnkey.

| Scale | Object | What you typically measure / vary |
| ----- | ------ | ----------------- |
| Wiki | One teleporting kernel on pages | Link structure, hubs, sinks, curricula ([[mfpt]]) |
| One agent | $Q = DAP$ (or world composition $PDA$) | Policy **D**, perception **P**, consequences **A**; $\pi$ vs PPR gap |
| Many agents | Joint kernels (product / shared world) | Coupling, orchestration $\alpha$, shared-world stationary mass |

Metrics transfer (as operators, not as identical state meaning):

- **$\|\pi - \mathrm{ppr}(v,\alpha)\|_1$** — how much intent actually moves long-run mass (wiki: seed vs unseeded graph; agents: goal regime vs free run).
- **MFPT** — access cost (wiki: path to a concept; agents: latency to task-complete).
- **Entropy rate** — unpredictability per step (wiki: hairball vs structured vault; agents: behavioral uncertainty — high entropy is not automatically “bad coordination”).
- **Trace** — coarse-grain to observed states (wiki: `goal` subgraph; agents: hide internals, keep outcomes).

### What [[analyze]] teaches that transfers to agents

| Wiki signal | Agent / network analogue |
| ----------- | ------------------------ |
| Sink pages | Actions or percepts with no onward transition — absorbing failure modes unless restart/intent handles them |
| Low-$\pi$ orphans | States the dynamics rarely visit — under-trained or under-linked regions of experience |
| High-$\pi$ hubs | Attractors; for agents, states that dominate long-run mass (desired or pathological) |
| Suggested links (commute) | “These two states are close in access cost but not coupled” — candidate edges in the kernel you are editing, or in the world graph that grounds a retrieval-style **P** |
| Raw SCC fragmentation | Disconnected regimes the walk (or joint world kernel) cannot mix without new coupling |

PageRank defaults did not invent this mapping — they **aligned** the wiki CLI
with the same *operators* catrace uses on PDA kernels, so dials and diagnostics
transfer as craft. State labels still mean different things on each side.

## Examples

### Example A — docs wiki as pure dynamics

```bash
wikigraph analyze docs/ --suggest-top 0
wikigraph graph docs/ -o /tmp/wiki.html
```

Uniform restart → global PageRank $\pi$: “where does attention settle if the
reader follows links and occasionally teleports?” That is **dynamics-only**
stationary mass on *page* space — the same *operator* as $\pi(Q)$ on an
experience kernel, not the claim that pages *are* experience states.

### Example B — intent over a corpus (firehose pattern)

```bash
wikigraph graph path/to/docs --relative-links --alpha 0.15 \
  --seed some/03-recommend --seed other/03-recommend \
  -o /tmp/ppr.html
```

Seeds are goal documents; $\alpha$ is pull strength. Node sizes are PPR —
structurally weighted importance **relative to intent**. Same construction as
PPR on an agent kernel with goal distribution $v$.

### Example C — access cost vs visit mass

`wikigraph goal --goal <slug>` ranks by MFPT (how hard to *reach*).
`graph --seed <slug>` sizes by PPR (how much long-run mass *near* the seed).
Agents need both: latency to a goal state vs occupancy under persistent pull.
Issue #22 (`goal --rank ppr`) makes that duality explicit in one command.

### Example D — LLM wiki as agent grounding (mental model)

Treat the vault as a **world** fragment (documents, tickets, ADRs in space
$W$). An LLM+human maintainer is an agent that:

- **P** — maps world → experience (retrieval / attention over pages),
- **D** — decides what to edit, link, or ask next,
- **A** — writes markdown that changes the world graph.

Running [[analyze]] after an ingest is then not only “wiki lint”: it is a
cheap check that the *world* graph still mixes, still has hubs where you want
them, and still has low access cost to concepts the agent must reach. That
improves the **substrate in $W$**. It improves **P** only insofar as **P** is
defined as grounding / retrieval over that graph — a modeling choice you must
state, not a free identity. Measuring with `--seed` practices the same intent
discipline you will use on $Q$.

## Related Concepts

- [[catrace]] — Shared Markov engine; PDA lives in the catrace wiki
- [[adr-012-teleporting-pagerank-default]] — α-damping + raw adj that align wiki dials with PPR craft
- [[teleportation-ergodicity]] — Why teleport / $\alpha$ exists
- [[stationary-distribution]] — $\pi$ as long-run mass
- [[markov-model]] — Wiki → kernel pipeline
- [[mfpt]] — Access cost complementary to PPR
- [[analyze]] — Structural health signals that transfer
- [[graph]] — NodeMass / `--seed` as intent visualization
- [[goal]] — MFPT curricula / subgraphs
- [[llm-wiki-pattern]] — Wiki as compounding codebase an agent (or human+LLM) maintains
- [[architecture]] — How wikigraph wires catrace
- [[random-walk]] — One-kernel foundation
- [[pagerank-foundation-rewrite]] — Proposal that motivated ADR-012

## Sources

- catrace wiki: PDA Triplet Model — `github.com/stephen-mcelhose/catrace` (`docs/wiki/pda-triplet-model.md`)
- catrace wiki: Personalized PageRank and Agent Modeling (`docs/wiki/personalized-pagerank-agent-modeling.md`)
- Hoffman, D., Prakash, C. & Chattopadhyay, S. (2024). *Traces of Consciousness* — PDA formalism source
- Jeh, G. & Widom, J. (2003). Scaling Personalized Web Search. *WWW 2003.*
- [[adr-012-teleporting-pagerank-default]]
- `cmd_graph.go`, `wiki.go` — `--alpha` / `--seed` / teleporting kernel wiring
