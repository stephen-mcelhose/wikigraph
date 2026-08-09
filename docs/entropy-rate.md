---
type: concept
title: Entropy Rate
description: The entropy rate H measures randomness per step of the Markov random walk — reported by every wikigraph analyze run as a wiki health signal.
resource: cmd_analyze.go
tags: [entropy-rate, information-theory, markov, health]
timestamp: 2026-08-09T08:05:55Z
---

# Entropy Rate

The **entropy rate** H of a [[random-walk]] on a wiki measures the average
amount of information (in bits) produced by each step of the walk:

```
H = −Σᵢ π(i) Σⱼ Pᵢⱼ log₂ Pᵢⱼ
```

where π is the [[stationary-distribution]] and P is the transition matrix
from [[markov-model]]. It is reported on every `wikigraph analyze` run in
the opening overview section.

## What the range means

Entropy rate is bounded by the link structure of the wiki:

| Value                   | Interpretation                                                            |
| ----------------------- | ------------------------------------------------------------------------- |
| H = 0                   | Every page has exactly one outgoing link — pure linear chain, no choice   |
| 0 < H < log₂(N)        | Varied out-degrees; most real wikis live here                             |
| H = log₂(N)            | Every page links to every other page with equal probability — maximum entropy |

N is the number of pages in the recurrent class. log₂(N) is the theoretical
maximum for N states.

## Why mid-range is healthy

- **Too low**: the wiki funnels readers through a small number of bottleneck
  pages. Navigation is predictable but rigid — readers have little choice.
- **Too high** (rare in practice): links are effectively random noise with no
  topic structure. Almost no real wiki reaches this limit.
- **Healthy target**: H > 50% of log₂(N). The wikigraph docs themselves
  score ~1.65 bits on 11 pages (log₂(11) ≈ 3.46 bits), a ratio of ~48% —
  reasonable for a small technical wiki.

The health table in [[analyze]] uses 70% of log₂(N) as the "healthy"
threshold and 50% as the "needs work" threshold.

## How catrace computes it

`catrace.Kernel.EntropyRate(base)` takes `base` as the logarithm base:

- `base = 2` → result in bits (wikigraph default)
- `base = math.E` → result in nats

The method uses the already-computed [[stationary-distribution|π]] and the
transition matrix P stored in the Kernel, so it is cheap to call after
`Stationary` has converged. See [[catrace]] for the full API.

## Reading the analyze output

```
Pages: 11   Edges: 36   Entropy rate: 1.6509 bits   (log₂(11) = 3.46 bits, 47.7%)
```

The ratio H / log₂(N) is the key figure. A wiki adding pages without adding
cross-links will see this ratio fall over time — a sign that new content is
not being integrated into the knowledge graph.

## Relationship to other concepts

- **[[communicating-classes]]**: entropy rate is computed over the recurrent
  class only. Transient pages contribute π = 0 and do not affect H.
- **[[stationary-distribution]]**: π weights each page's contribution to H.
  Hub pages (high π) with many out-links push H up.
- **[[random-walk]]**: H is a property of the walk's long-run dynamics, not
  of any individual page.

## Sources

- [`cmd_analyze.go` — `kern.EntropyRate(2)`](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_analyze.go)
- [Entropy rate — Wikipedia](https://en.wikipedia.org/wiki/Entropy_rate)
- [Shannon entropy — Wikipedia](https://en.wikipedia.org/wiki/Entropy_(information_theory))
