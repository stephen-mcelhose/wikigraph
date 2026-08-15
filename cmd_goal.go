package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

var (
	flagGoalGoals    []string
	flagGoalTop      int
	flagGoalOut      string
	flagGoalMinEdge  float64
	flagGoalSed      []string
	flagGoalStrategy string
)

var goalCmd = &cobra.Command{
	Use:   "goal <wiki-dir>",
	Short: "Compute a learning-path subgraph toward goal pages [PROTOTYPE]",
	Long: `goal builds a subgraph focused on navigating toward one or more goal pages.

PROTOTYPE: goals must currently be exact page slugs. Semantic natural-language
goals ("understand quantum error correction") are planned — see
https://github.com/stephen-mcelhose/wikigraph/issues/1

Available strategies (--strategy):
  - union        (default) Ranks pages by minimum MFPT to any goal (OR-neighborhood)
  - intersection Ranks pages by maximum MFPT across all goals (prerequisites shared by ALL goals)
  - path         [PROTOTYPE] Finds the most probable sequential link chain connecting goals in flag order
                 Implements Dijkstra directly on the raw transition matrix rather than via catrace.
  - bottleneck   [PROTOTYPE] Ranks pages by random-walk betweenness centrality across goal pairs
                 Computes the fundamental matrix N=(I-Q)^-1 directly via gonum/mat rather than via catrace.`,

	Args: cobra.ExactArgs(1),
	RunE: runGoal,
}

func init() {
	goalCmd.Flags().StringArrayVar(&flagGoalGoals, "goal", nil, "goal page slug (repeatable, at least one required)")
	goalCmd.Flags().IntVar(&flagGoalTop, "top", 10, "number of pages in the subgraph")
	goalCmd.Flags().StringVarP(&flagGoalOut, "out", "o", "goal_graph.html", "output HTML file")
	goalCmd.Flags().Float64VarP(&flagGoalMinEdge, "min-edge", "m", 0.02, "omit edges below this transition probability")
	goalCmd.Flags().StringArrayVarP(&flagGoalSed, "sed", "s", nil, "sed expression(s) to apply to the HTML output (repeatable)")
	goalCmd.Flags().StringVar(&flagGoalStrategy, "strategy", "union", "selection strategy: union, intersection, path [PROTOTYPE], bottleneck [PROTOTYPE]")
}

func runGoal(cmd *cobra.Command, args []string) error {
	if len(flagGoalGoals) == 0 {
		return fmt.Errorf("at least one --goal slug is required")
	}

	wikiDir := args[0]
	exclude := makeExcludeMap(flagExclude)

	kern, _, pages, _, err := buildKernelWithOpts(wikiDir, flagRecursive, exclude, flagRelativeLinks, flagAlpha, flagSeed)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Pages: %d\n", len(pages))

	// slug → index map.
	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}

	// Validate --goal slugs.
	goalIdxs := make([]int, 0, len(flagGoalGoals))
	for _, g := range flagGoalGoals {
		i, ok := idx[g]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown --goal slug %q\nValid slugs:\n", g)
			for _, p := range pages {
				fmt.Fprintf(os.Stderr, "  %s\n", p)
			}
			return fmt.Errorf("unknown --goal slug: %s", g)
		}
		goalIdxs = append(goalIdxs, i)
	}

	n := len(pages)
	top := flagGoalTop
	if top > n {
		top = n
	}

	// Math transition matrix for prototype path/bottleneck strategies.
	P := kern.P

	var selectedIndices []int

	switch flagGoalStrategy {
	case "union":
		selectedIndices = selectUnion(kern, n, goalIdxs, top)
	case "intersection":
		selectedIndices = selectIntersection(kern, n, goalIdxs, top)
	case "path":
		fmt.Fprintln(os.Stderr, "Warning: path strategy is a prototype — implements Dijkstra directly rather than via catrace (see ADR-008)")
		selectedIndices = selectPath(P, n, goalIdxs, top)
	case "bottleneck":
		fmt.Fprintln(os.Stderr, "Warning: bottleneck strategy is a prototype — computes fundamental matrix directly rather than via catrace (see ADR-008)")
		var err error
		selectedIndices, err = selectBottleneck(P, n, goalIdxs, top)
		if err != nil {
			return fmt.Errorf("bottleneck strategy failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown strategy %q (valid: union, intersection, path, bottleneck)", flagGoalStrategy)
	}

	subset := selectedIndices
	sort.Ints(subset)

	// Compute effective kernel on subset via trace (teleporting math).
	traceKern, err := kern.Trace(subset, 1e-9)
	if err != nil {
		return fmt.Errorf("trace failed (try increasing --top): %w", err)
	}

	pi, err := traceKern.Stationary(1e-12, 5000)
	if err != nil {
		pi = nil // ToHTML falls back to its own Stationary when NodeMass is empty
	}

	title := filepath.Base(wikiDir) + " → " + flagGoalGoals[0] + " [" + flagGoalStrategy + "]"
	html, err := traceKern.ToHTML(&catrace.VisualiseOptions{
		Title:    title,
		MinEdge:  flagGoalMinEdge,
		Width:    1400,
		Height:   900,
		NodeMass: pi,
	})
	if err != nil {
		return fmt.Errorf("generating HTML: %w", err)
	}

	page := string(html)
	if len(flagGoalSed) > 0 {
		page, err = applySed(page, flagGoalSed)
		if err != nil {
			return err
		}
	}

	if err := os.WriteFile(flagGoalOut, []byte(page), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", flagGoalOut, err)
	}
	fmt.Fprintf(os.Stderr, "Written: %s (%d nodes, strategy: %s)\n", flagGoalOut, len(subset), flagGoalStrategy)
	return nil
}

// selectUnion scores pages by min MFPT to any goal (OR-neighborhood).
func selectUnion(kern *catrace.Kernel, n int, goalIdxs []int, top int) []int {
	const inf = 1e18
	scores := make([]float64, n)
	for i := range scores {
		scores[i] = inf
	}
	for _, gIdx := range goalIdxs {
		scores[gIdx] = 0
		for i := 0; i < n; i++ {
			if i == gIdx {
				continue
			}
			mfpt, err := kern.MeanFirstPassage(i, gIdx)
			if err != nil {
				continue
			}
			if mfpt < scores[i] {
				scores[i] = mfpt
			}
		}
	}

	return rankAndSelect(scores, n, goalIdxs, top)
}

// selectIntersection scores pages by max MFPT to all goals (AND-prerequisites).
func selectIntersection(kern *catrace.Kernel, n int, goalIdxs []int, top int) []int {
	goalSet := indexSet(goalIdxs)
	maxMFPT := make([]float64, n)
	for i := 0; i < n; i++ {
		if goalSet[i] {
			maxMFPT[i] = 0
			continue
		}
		maxMFPT[i] = maxMFPTToGoals(kern, i, goalIdxs)
	}
	return rankAndSelect(maxMFPT, n, goalIdxs, top)
}

func indexSet(idxs []int) map[int]bool {
	s := make(map[int]bool, len(idxs))
	for _, i := range idxs {
		s[i] = true
	}
	return s
}

// maxMFPTToGoals returns the max MeanFirstPassage(from, g) over goals, or +inf
// if any goal is unreachable from from.
func maxMFPTToGoals(kern *catrace.Kernel, from int, goalIdxs []int) float64 {
	const inf = 1e18
	maxVal := 0.0
	for _, g := range goalIdxs {
		mfpt, err := kern.MeanFirstPassage(from, g)
		if err != nil {
			return inf
		}
		if mfpt > maxVal {
			maxVal = mfpt
		}
	}
	return maxVal
}

// Item for priority queue in Dijkstra
type pathItem struct {
	node int
	dist float64
}

type priorityQueue []*pathItem

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*pathItem)
	*pq = append(*pq, item)
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// selectPath computes sequential Dijkstra on -log(P_ij) edge weights.
func selectPath(P *mat.Dense, n int, goalIdxs []int, top int) []int {
	pathNodes := make(map[int]bool, len(goalIdxs))
	for _, g := range goalIdxs {
		pathNodes[g] = true
	}
	for k := 0; k+1 < len(goalIdxs); k++ {
		for _, v := range dijkstraNegLogPath(P, n, goalIdxs[k], goalIdxs[k+1]) {
			pathNodes[v] = true
		}
	}
	return expandPathNeighbors(P, n, pathNodes, top)
}

// dijkstraNegLogPath returns nodes on the shortest path from src to dst using
// edge weights -log(P_ij). If dst is unreachable, returns nil (caller already
// seeded pathNodes with the goals themselves).
func dijkstraNegLogPath(P *mat.Dense, n, src, dst int) []int {
	dist := make([]float64, n)
	prev := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}
	dist[src] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &pathItem{node: src, dist: 0})

	for pq.Len() > 0 {
		curr := heap.Pop(pq).(*pathItem)
		u := curr.node
		if curr.dist > dist[u] {
			continue
		}
		if u == dst {
			break
		}
		relaxNegLogNeighbors(P, n, u, dist, prev, pq)
	}

	if prev[dst] == -1 && src != dst {
		return nil
	}
	var path []int
	for curr := dst; curr != -1; curr = prev[curr] {
		path = append(path, curr)
	}
	return path
}

func relaxNegLogNeighbors(P *mat.Dense, n, u int, dist []float64, prev []int, pq *priorityQueue) {
	for v := 0; v < n; v++ {
		p := P.At(u, v)
		if p <= 0 {
			continue
		}
		w := -math.Log(p)
		if dist[u]+w < dist[v] {
			dist[v] = dist[u] + w
			prev[v] = u
			heap.Push(pq, &pathItem{node: v, dist: dist[v]})
		}
	}
}

type pathNeighbor struct {
	idx  int
	prob float64
}

// expandPathNeighbors returns path node indices, padding up to top with the
// highest bidirectional transition probabilities off the path.
func expandPathNeighbors(P *mat.Dense, n int, pathNodes map[int]bool, top int) []int {
	selected := make(map[int]bool, len(pathNodes))
	for idx := range pathNodes {
		selected[idx] = true
	}
	if len(selected) < top {
		fillPathNeighbors(P, n, pathNodes, selected, top)
	}
	res := make([]int, 0, len(selected))
	for idx := range selected {
		res = append(res, idx)
	}
	return res
}

func fillPathNeighbors(P *mat.Dense, n int, pathNodes, selected map[int]bool, top int) {
	candidates := pathNeighborCandidates(P, n, pathNodes, selected)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].prob > candidates[j].prob
	})
	for _, cand := range candidates {
		if len(selected) >= top {
			return
		}
		selected[cand.idx] = true
	}
}

func pathNeighborCandidates(P *mat.Dense, n int, pathNodes, selected map[int]bool) []pathNeighbor {
	var candidates []pathNeighbor
	for i := 0; i < n; i++ {
		if selected[i] {
			continue
		}
		maxProb := maxBidirectionalProb(P, i, pathNodes)
		if maxProb > 0 {
			candidates = append(candidates, pathNeighbor{idx: i, prob: maxProb})
		}
	}
	return candidates
}

func maxBidirectionalProb(P *mat.Dense, i int, pathNodes map[int]bool) float64 {
	maxProb := 0.0
	for pNode := range pathNodes {
		if pOut := P.At(pNode, i); pOut > maxProb {
			maxProb = pOut
		}
		if pIn := P.At(i, pNode); pIn > maxProb {
			maxProb = pIn
		}
	}
	return maxProb
}

type nodePair struct {
	src int
	dst int
}

// selectBottleneck scores nodes by random walk betweenness centrality across goal pairs.
func selectBottleneck(P *mat.Dense, n int, goalIdxs []int, top int) ([]int, error) {
	scores := make([]float64, n)
	for _, pr := range bottleneckPairs(n, goalIdxs) {
		addFundamentalVisits(P, n, pr.src, pr.dst, scores)
	}
	// rankAndSelect picks lowest scores first; negate so highest betweenness is selected.
	for i := range scores {
		scores[i] = -scores[i]
	}
	return rankAndSelect(scores, n, goalIdxs, top), nil
}

func bottleneckPairs(n int, goalIdxs []int) []nodePair {
	if len(goalIdxs) == 1 {
		g := goalIdxs[0]
		pairs := make([]nodePair, 0, n-1)
		for i := 0; i < n; i++ {
			if i != g {
				pairs = append(pairs, nodePair{src: i, dst: g})
			}
		}
		return pairs
	}
	var pairs []nodePair
	for i := 0; i < len(goalIdxs); i++ {
		for j := 0; j < len(goalIdxs); j++ {
			if i != j {
				pairs = append(pairs, nodePair{src: goalIdxs[i], dst: goalIdxs[j]})
			}
		}
	}
	return pairs
}

// addFundamentalVisits accumulates expected visits on the transient set
// (all nodes except absorbing target t) for walks from s, via (I-Q)^{-1}.
func addFundamentalVisits(P *mat.Dense, n, s, t int, scores []float64) {
	transientMap := make([]int, 0, n-1)
	origToTrans := make(map[int]int, n-1)
	for i := 0; i < n; i++ {
		if i != t {
			origToTrans[i] = len(transientMap)
			transientMap = append(transientMap, i)
		}
	}
	tLen := len(transientMap)
	if tLen == 0 {
		return
	}

	N, ok := fundamentalMatrix(P, transientMap)
	if !ok {
		return
	}
	sTrans, ok := origToTrans[s]
	if !ok {
		return
	}
	for kTrans := 0; kTrans < tLen; kTrans++ {
		scores[transientMap[kTrans]] += N.At(sTrans, kTrans)
	}
}

func fundamentalMatrix(P *mat.Dense, transientMap []int) (*mat.Dense, bool) {
	tLen := len(transientMap)
	iqData := make([]float64, tLen*tLen)
	for i := 0; i < tLen; i++ {
		origI := transientMap[i]
		for j := 0; j < tLen; j++ {
			origJ := transientMap[j]
			val := 0.0
			if i == j {
				val = 1.0
			}
			val -= P.At(origI, origJ)
			iqData[i*tLen+j] = val
		}
	}
	iq := mat.NewDense(tLen, tLen, iqData)
	var N mat.Dense
	diagOnes := make([]float64, tLen)
	for i := range diagOnes {
		diagOnes[i] = 1.0
	}
	if err := N.Solve(iq, mat.NewDiagDense(tLen, diagOnes)); err != nil {
		return nil, false
	}
	return &N, true
}

func rankAndSelect(scores []float64, n int, goalIdxs []int, top int) []int {
	type ranked struct {
		idx   int
		score float64
	}
	all := make([]ranked, n)
	for i := range all {
		all[i] = ranked{idx: i, score: scores[i]}
	}
	sort.Slice(all, func(a, b int) bool { return all[a].score < all[b].score })

	selected := make(map[int]bool)
	for _, g := range goalIdxs {
		selected[g] = true
	}
	for _, r := range all {
		if len(selected) >= top {
			break
		}
		selected[r.idx] = true
	}

	res := make([]int, 0, len(selected))
	for idx := range selected {
		res = append(res, idx)
	}
	return res
}
