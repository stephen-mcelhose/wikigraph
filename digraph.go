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
	return annotateClosed(adj, tarjanSCCs(adj))
}

type tarjanState struct {
	adj            *mat.Dense
	n              int
	index, lowlink []int
	onStack        []bool
	stack          []int
	dfsNum         int
	sccs           [][]int
}

func tarjanSCCs(adj *mat.Dense) [][]int {
	n, _ := adj.Dims()
	s := &tarjanState{
		adj:     adj,
		n:       n,
		index:   make([]int, n),
		lowlink: make([]int, n),
		onStack: make([]bool, n),
	}
	for i := range s.index {
		s.index[i] = -1
	}
	for v := 0; v < n; v++ {
		if s.index[v] == -1 {
			s.strongConnect(v)
		}
	}
	return s.sccs
}

func (s *tarjanState) strongConnect(v int) {
	s.index[v] = s.dfsNum
	s.lowlink[v] = s.dfsNum
	s.dfsNum++
	s.stack = append(s.stack, v)
	s.onStack[v] = true

	s.considerOutNeighbors(v)

	if s.lowlink[v] == s.index[v] {
		s.sccs = append(s.sccs, s.popSCC(v))
	}
}

func (s *tarjanState) considerOutNeighbors(v int) {
	for w := 0; w < s.n; w++ {
		if s.adj.At(v, w) <= 0 {
			continue
		}
		if s.index[w] == -1 {
			s.strongConnect(w)
			if s.lowlink[w] < s.lowlink[v] {
				s.lowlink[v] = s.lowlink[w]
			}
			continue
		}
		if s.onStack[w] && s.index[w] < s.lowlink[v] {
			s.lowlink[v] = s.index[w]
		}
	}
}

func (s *tarjanState) popSCC(root int) []int {
	var comp []int
	for {
		w := s.stack[len(s.stack)-1]
		s.stack = s.stack[:len(s.stack)-1]
		s.onStack[w] = false
		comp = append(comp, w)
		if w == root {
			return comp
		}
	}
}

func annotateClosed(adj *mat.Dense, comps [][]int) []digraphSCC {
	n, _ := adj.Dims()
	compOf := make([]int, n)
	for ci, comp := range comps {
		for _, v := range comp {
			compOf[v] = ci
		}
	}
	out := make([]digraphSCC, len(comps))
	for ci, comp := range comps {
		out[ci] = digraphSCC{
			Members: comp,
			Closed:  componentIsClosed(adj, comp, compOf, ci),
		}
	}
	return out
}

func componentIsClosed(adj *mat.Dense, members []int, compOf []int, ci int) bool {
	n, _ := adj.Dims()
	for _, v := range members {
		for w := 0; w < n; w++ {
			if adj.At(v, w) > 0 && compOf[w] != ci {
				return false
			}
		}
	}
	return true
}
