package main

import "gonum.org/v1/gonum/mat"

// digraphSCC describes one strongly connected component of the raw wikilink
// digraph. closed means no edge leaves the component (authoring "terminal"
// cluster — includes trivial sink singletons).
type digraphSCC struct {
	Members []int
	Closed  bool
}

// rawSCCs returns Tarjan SCCs of the raw adjacency digraph (edge iff adj>0).
// Order is discovery order; useful as an authoring signal. The teleporting
// math kernel is ergodic by construction and is not used here.
func rawSCCs(adj *mat.Dense) []digraphSCC {
	n, _ := adj.Dims()
	if n == 0 {
		return nil
	}

	index := make([]int, n)
	lowlink := make([]int, n)
	onStack := make([]bool, n)
	for i := range index {
		index[i] = -1
	}
	var stack []int
	var sccs [][]int
	dfsNum := 0

	var strongConnect func(v int)
	strongConnect = func(v int) {
		index[v] = dfsNum
		lowlink[v] = dfsNum
		dfsNum++
		stack = append(stack, v)
		onStack[v] = true

		for w := 0; w < n; w++ {
			if adj.At(v, w) <= 0 {
				continue
			}
			if index[w] == -1 {
				strongConnect(w)
				if lowlink[w] < lowlink[v] {
					lowlink[v] = lowlink[w]
				}
			} else if onStack[w] {
				if index[w] < lowlink[v] {
					lowlink[v] = index[w]
				}
			}
		}

		if lowlink[v] == index[v] {
			var comp []int
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				comp = append(comp, w)
				if w == v {
					break
				}
			}
			sccs = append(sccs, comp)
		}
	}

	for v := 0; v < n; v++ {
		if index[v] == -1 {
			strongConnect(v)
		}
	}

	compOf := make([]int, n)
	for ci, comp := range sccs {
		for _, v := range comp {
			compOf[v] = ci
		}
	}

	out := make([]digraphSCC, len(sccs))
	for ci, comp := range sccs {
		closed := true
		for _, v := range comp {
			for w := 0; w < n; w++ {
				if adj.At(v, w) > 0 && compOf[w] != ci {
					closed = false
					break
				}
			}
			if !closed {
				break
			}
		}
		out[ci] = digraphSCC{Members: comp, Closed: closed}
	}
	return out
}
