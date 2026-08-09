---
type: concept
title: catrace Library
description: catrace is the Go library that implements all Markov chain mathematics for wikigraph — Kernel struct, stationary distribution, SCC decomposition, MFPT, commute time, trace kernels, and HTML rendering.
resource: go.mod
tags: [catrace, library, go, markov, kernel, api]
timestamp: 2026-08-09T08:05:55Z
---

# catrace Library

`github.com/stephen-mcelhose/catrace` is a Go library for
**finite-state Markov models of autonomous-agent networks**. wikigraph uses
it as a pure-math backend: wikigraph handles filesystem I/O, wikilink parsing,
and CLI output; catrace handles all linear algebra and Markov chain
mathematics.

See [[architecture]] for how the two packages fit together.

## Repository and version

- Repository: `https://github.com/stephen-mcelhose/catrace`
- Version pinned in `go.mod`: `v0.0.0-20260805012807-4f4561824de9`
- Dependency: `gonum.org/v1/gonum` (matrix operations)

## The Kernel struct

The central type is `Kernel` — a square row-stochastic matrix with named
states:

```go
type Kernel struct {
    P          *mat.Dense // row-stochastic transition matrix
    StateNames []string
}
```

wikigraph constructs a Kernel via `NewRandomWalkKernel`, which row-normalises
a raw adjacency matrix. An error is returned if any row is all-zero (a
[[sink-page|sink]]) — so [[markov-model|`buildAdjacency`]] applies uniform
teleportation to sinks *before* calling this constructor.

## Constructor

```go
func NewRandomWalkKernel(adj *mat.Dense, names []string) (*Kernel, error)
```

Divides each row of `adj` by its sum, producing `P(i,j) = adj[i][j] / Σₖ adj[i][k]`.
This is the standard [[random-walk]] on a graph — edge weights become
transition probabilities proportional to weight.

## Analysis methods used by wikigraph

| Method                             | Returns                  | Used in                   |
| ---------------------------------- | ------------------------ | ------------------------- |
| `Stationary(tol, maxIter)`         | `[]float64` (π)          | [[analyze]], [[export]], [[graph]] |
| `EntropyRate(base)`                | `float64`                | [[analyze]]               |
| `Classes(tol)`                     | `*ClassDecomposition`    | [[analyze]], [[export]]   |
| `MeanFirstPassage(i, j int)`       | `float64, error`         | [[analyze]], [[goal]]     |
| `CommuteTime(i, j int)`            | `float64, error`         | [[analyze]]               |
| `Trace(subset []int, tol float64)` | `*Kernel, error`         | [[goal]]                  |
| `ToHTML(opts *VisualiseOptions)`   | `[]byte, error`          | [[graph]]                 |

### Stationary(tol, maxIter)

Power iteration for the [[stationary-distribution]]. Starts from a uniform
vector and multiplies by P until L₁ convergence < `tol`.

### EntropyRate(base)

Computes [[entropy-rate]] H = −Σᵢ π(i) Σⱼ Pᵢⱼ log\_base(Pᵢⱼ).
wikigraph passes `base = 2` for bits.

### Classes(tol)

Kosaraju's SCC algorithm over the graph of non-negligible transitions
(`Pᵢⱼ > tol`). Returns a `ClassDecomposition` with `Recurrent` and
`Transient` slices of state-index sets, plus a `Periods` map.
See [[communicating-classes]] and [[recurrent-class]].

### MeanFirstPassage(i, j)

Expected steps from state i to first reach state j via the fundamental
matrix of the chain. Returns `+Inf` when j is unreachable from i.
See [[mfpt]].

### CommuteTime(i, j)

`MeanFirstPassage(i,j) + MeanFirstPassage(j,i)` — symmetric distance metric.
See [[commute-time]].

### Trace(subset, tol)

Projects the full transition matrix onto `subset` states, integrating out
all excursions through the complement:

```
P_A = a + b (I - c)⁻¹ d
```

The stationary distribution of the trace kernel equals the parent π
restricted and renormalized to the subset. Used by [[goal]] to build
focused subgraph visualisations.

### ToHTML(opts)

Renders a D3 force-directed graph. Node size encodes [[stationary-distribution|π]];
edge opacity encodes transition probability. Used by [[graph]].

## Additional Kernel methods

| Method                    | Purpose                                           |
| ------------------------- | ------------------------------------------------- |
| `Clone()`                 | Deep copy of the Kernel                           |
| `NumStates()`             | Number of states (rows/columns of P)              |
| `NormalizeRows(tol)`      | Re-normalize rows in place (mutates P)            |
| `Multiply(other)`         | Matrix product k·other as a new Kernel            |
| `LeftAction(dist)`        | Evolve a distribution one step: π' = π·P          |
| `IsTraceOf(parent, ...)`  | Verify that this Kernel is a valid trace of parent |
| `Sample(rowIdx, rng)`     | Sample one step forward given a starting state    |

## Agent model (not used by wikigraph)

catrace's primary design is for modelling autonomous-agent networks via three
rectangular row-stochastic maps:

```
D : X → G  (decision:   experience → action)
A : G → W  (effect:     action     → world)
P : W → X  (perception: world      → experience)
```

Composing these in cyclic order gives square kernels on each space
(`QualiaKernel`, `StrategyKernel`, `WorldKernel`) that share eigenvalues —
three perspectives on the same closed-loop system. wikigraph does not use
this model but it is the library's primary motivation.

## Estimation and trajectory tools

| Function                             | Purpose                                              |
| ------------------------------------ | ---------------------------------------------------- |
| `EstimateKernelFromSequence`         | Empirical transition counts → kernel (with pseudocount smoothing) |
| `SampleTraceFromSequence`            | Filter a trajectory to observed states               |
| `WindowedTraceEstimates`             | Sliding-window kernel estimates for drift detection  |

## Sources

- [`go.mod` — version pin](https://github.com/stephen-mcelhose/wikigraph/blob/main/go.mod)
- [catrace repository](https://github.com/stephen-mcelhose/catrace)
- [`wiki.go` — `NewRandomWalkKernel` call site](https://github.com/stephen-mcelhose/wikigraph/blob/main/wiki.go)
- [`cmd_analyze.go` — Stationary, Classes, CommuteTime](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_analyze.go)
- [`cmd_goal.go` — MeanFirstPassage, Trace](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_goal.go)
- [`cmd_graph.go` — ToHTML](https://github.com/stephen-mcelhose/wikigraph/blob/main/cmd_graph.go)
