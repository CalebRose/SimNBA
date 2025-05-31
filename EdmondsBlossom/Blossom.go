package edmondsblossom

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Graph represents an undirected graph for Edmonds' Blossom algorithm.
type Graph struct {
	n int     // number of vertices
	G [][]int // adjacency matrix: G[u][v] == 1 if edge exists
	M []int   // matching: M[u] = matched vertex or -1
}

// NewGraph allocates a Graph with n vertices, no edges, and all vertices unmatched.
func NewGraph(n int) *Graph {
	g := &Graph{n: n}
	g.G = make([][]int, n)
	for i := range g.G {
		g.G[i] = make([]int, n)
	}
	g.M = make([]int, n)
	for i := 0; i < n; i++ {
		g.M[i] = -1
	}
	return g
}

// FreeVertices returns all vertices not yet matched.
func (g *Graph) FreeVertices() []int {
	free := []int{}
	for u, m := range g.M {
		if m < 0 {
			free = append(free, u)
		}
	}
	return free
}

// AdjVertices returns all neighbors of u.
func (g *Graph) AdjVertices(u int) []int {
	nbrs := []int{}
	for v := 0; v < g.n; v++ {
		if g.G[u][v] == 1 && u != v {
			nbrs = append(nbrs, v)
		}
	}
	return nbrs
}

// UpdateMatching augments the matching along path P.
func (g *Graph) UpdateMatching(path []int) {
	if len(path) < 2 {
		return
	}
	for i := 0; i < len(path)-1; i += 2 {
		u := path[i]
		v := path[i+1]
		g.M[u] = v
		g.M[v] = u
	}
}

// getPrvVerFromMap finds the original vertex whose mapping equals target.
func getPrvVerFromMap(mapOld []int, target int) int {
	for i, m := range mapOld {
		if m == target {
			return i
		}
	}
	return -1
}

// ContractGraph builds a contracted graph by collapsing blossom vertices in B.
func ContractGraph(old *Graph, mapOld []int, B map[int]struct{}) *Graph {
	newN := old.n - len(B) + 1
	newG := NewGraph(newN)

	// copy matching edges outside the blossom
	for i := 0; i < old.n; i++ {
		if old.M[i] >= 0 {
			if _, inB := B[i]; !inB {
				newG.M[mapOld[i]] = mapOld[old.M[i]]
				newG.M[mapOld[old.M[i]]] = mapOld[i]
			}
		}
	}

	// copy adjacency
	for i := 0; i < old.n; i++ {
		for j := 0; j < old.n; j++ {
			if old.G[i][j] == 1 {
				newG.G[mapOld[i]][mapOld[j]] = 1
			}
		}
	}

	newG.n = newN
	return newG
}

// ShortestPath returns the path from descendant up to ancestor in forest F.
func ShortestPath(F []ForestVertex, descendant, ancestor int) []int {
	path := []int{}
	for descendant != ancestor {
		path = append(path, descendant)
		descendant = F[descendant].Parent
	}
	path = append(path, ancestor)
	return path
}

// ForestVertex is a node in the alternating forest used to find augmenting paths.
type ForestVertex struct {
	InForest   bool
	Parent     int
	Root       int
	DistToRoot int
}

// ReverseInts reverses a slice of ints in place.
func ReverseInts(s []int) {
	i, j := 0, len(s)-1
	for i < j {
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
}

// AddToForest adds edge (vertex—adj) and its matched partner into forest F.
func AddToForest(g *Graph, F []ForestVertex, vertex, adj int, nodes *[]int) {
	mate := g.M[adj]
	*nodes = append(*nodes, mate)

	F[adj].InForest = true
	F[mate].InForest = true
	F[adj].Root = F[vertex].Root
	F[mate].Root = F[vertex].Root
	F[adj].Parent = vertex
	F[mate].Parent = adj
	F[adj].DistToRoot = F[vertex].DistToRoot + 1
	F[mate].DistToRoot = F[adj].DistToRoot + 1
}

// ReturnAugPath splices two root-to-vertex paths into one augmenting path.
func ReturnAugPath(F []ForestVertex, v1, v2 int) []int {
	P1 := ShortestPath(F, v1, F[v1].Root)
	P2 := ShortestPath(F, v2, F[v2].Root)
	ReverseInts(P1)
	return append(P1, P2...)
}

// BlossomRecursion handles discovery of an odd cycle (blossom) and contracts it.
func BlossomRecursion(g *Graph, F []ForestVertex, v, adj int) []int {
	// 1) Build and reverse root-to-vertex paths
	P1 := ShortestPath(F, v, F[v].Root)
	P2 := ShortestPath(F, adj, F[adj].Root)
	ReverseInts(P1)
	ReverseInts(P2)

	// 2) Find the split index
	i := 0
	for ; i < len(P1) && i < len(P2) && P1[i] == P2[i]; i++ {
	}
	start := i - 1

	// 3) Collect blossom vertices into B
	B := make(map[int]struct{})
	for j := start; j < len(P1); j++ {
		B[P1[j]] = struct{}{}
	}
	for j := start; j < len(P2); j++ {
		B[P2[j]] = struct{}{}
	}

	// 4) Build contraction mapOld: old index → new index
	mapOld := make([]int, g.n)
	for idx := range mapOld {
		mapOld[idx] = -1
	}
	newIdx := 0
	for idx := 0; idx < g.n; idx++ {
		if _, inB := B[idx]; !inB {
			mapOld[idx] = newIdx
			newIdx++
		}
	}
	// All blossom vertices map to newIdx
	for b := range B {
		mapOld[b] = newIdx
	}

	// 5) Contract graph and recurse
	contracted := ContractGraph(g, mapOld, B)
	PB := FindAugPath(contracted)
	if len(PB) < 2 {
		return PB
	}

	// 6) Check if we used the super-node
	super := newIdx
	usedSuper := false
	for _, x := range PB {
		if x == super {
			usedSuper = true
			break
		}
	}
	if !usedSuper {
		// No expansion needed — map PB back directly
		res := make([]int, len(PB))
		for k, x := range PB {
			res[k] = getPrvVerFromMap(mapOld, x)
		}
		return res
	}

	// 7) Remove super-node from PB
	pos := -1
	for idx, x := range PB {
		if x == super {
			pos = idx
			break
		}
	}
	if pos < 0 {
		// Shouldn't happen, but fallback
		res := make([]int, len(PB))
		for k, x := range PB {
			res[k] = getPrvVerFromMap(mapOld, x)
		}
		return res
	}
	PB = append(PB[:pos], PB[pos+1:]...)

	// 8) Map remaining PB back to original indices
	mappedPB := make([]int, len(PB))
	for idx, x := range PB {
		mappedPB[idx] = getPrvVerFromMap(mapOld, x)
	}

	// 9) If after removal there's only one vertex, nothing more to splice
	if len(mappedPB) < 2 {
		return mappedPB
	}

	// 10) Determine which blossom-entry vertex (adjInB) to splice in
	side, adjInB := 0, -1
	if pos == 0 {
		side = 1
		for _, cand := range g.AdjVertices(mappedPB[1]) {
			if _, inB := B[cand]; inB {
				adjInB = cand
				break
			}
		}
	} else if pos == len(mappedPB)-1 {
		side = 0
		pen := mappedPB[len(mappedPB)-2]
		for _, cand := range g.AdjVertices(pen) {
			if _, inB := B[cand]; inB {
				adjInB = cand
				break
			}
		}
	} else {
		side = 0
		for _, cand := range g.AdjVertices(mappedPB[pos-1]) {
			if _, inB := B[cand]; inB {
				adjInB = cand
				break
			}
		}
	}
	if adjInB < 0 {
		log.Fatalf("BlossomRecursion: failed to find blossom entry, pos=%d, B=%v", pos, B)
	}

	// 11) Reconstruct the circular path for the blossom
	stem1 := P1[start:]
	stem2 := P2[start:]
	ReverseInts(stem2)
	circularPath := append(stem1, stem2...)

	// 12) Locate adjInB in circularPath
	bIdx := -1
	for idx, x := range circularPath {
		if x == adjInB {
			bIdx = idx
			break
		}
	}
	if bIdx < 0 {
		return mappedPB
	}

	// 13) Extract the correct slice through the blossom
	var pathInB []int
	if bIdx == 0 || bIdx == len(circularPath)-1 {
		pathInB = []int{circularPath[bIdx]}
	} else if bIdx%2 == 0 {
		pathInB = append([]int{}, circularPath[:bIdx+1]...)
		if side == 1 {
			ReverseInts(pathInB)
		}
	} else {
		pathInB = append([]int{}, circularPath[bIdx:]...)
		if side == 0 {
			ReverseInts(pathInB)
		}
	}

	// 14) Splice everything back into the final augmenting path
	result := append([]int{}, mappedPB[:pos]...)
	result = append(result, pathInB...)
	result = append(result, mappedPB[pos:]...)
	return result
}

// FindAugPath finds one augmenting path or returns empty.
func FindAugPath(g *Graph) []int {
	// initialize forest
	F := make([]ForestVertex, g.n)
	nodes := g.FreeVertices()
	for _, u := range nodes {
		F[u].InForest = true
		F[u].Root = u
		F[u].Parent = u
		F[u].DistToRoot = 0
	}

	// mark edges: 0=no edge, 1=unmatched, 2=matched
	GMarked := make([][]int, g.n)
	for i := range GMarked {
		GMarked[i] = make([]int, g.n)
		for j := range GMarked[i] {
			if g.G[i][j] == 1 {
				if j == g.M[i] {
					GMarked[i][j] = 2
				} else {
					GMarked[i][j] = 1
				}
			}
		}
	}

	// BFS-like search
	for idx := 0; idx < len(nodes); idx++ {
		u := nodes[idx]
		for _, v := range g.AdjVertices(u) {
			if GMarked[u][v] != 1 {
				continue
			}
			if !F[v].InForest {
				AddToForest(g, F, u, v, &nodes)
			} else if F[u].Root != F[v].Root && F[v].DistToRoot%2 == 0 {
				return ReturnAugPath(F, u, v)
			} else if F[v].DistToRoot%2 == 0 {
				return BlossomRecursion(g, F, u, v)
			}
			// mark edge
			GMarked[u][v] = 2
			GMarked[v][u] = 2
		}
	}
	return nil
}

// FindMaxMatching repeatedly augments until no more paths exist.
func FindMaxMatching(g *Graph) {
	for {
		path := FindAugPath(g)
		if len(path) < 2 {
			break
		}
		fmt.Print("Augmenting Path (")
		for i, u := range path {
			if i > 0 {
				fmt.Print("-")
			}
			fmt.Print(u)
		}
		fmt.Println(")")
		g.UpdateMatching(path)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: blossom <graph-file>")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	n := len(lines)
	if n == 0 {
		fmt.Println("Empty graph file")
		return
	}

	// count columns in first line
	cols := 0
	for _, ch := range lines[0] {
		if ch == '0' || ch == '1' {
			cols++
		}
	}
	if cols != n {
		fmt.Println("Graph must be square (n x n)")
		return
	}

	g := NewGraph(n)
	for i, line := range lines {
		c := 0
		for _, ch := range line {
			if ch == '0' || ch == '1' {
				val := int(ch - '0')
				g.G[i][c] = val
				fmt.Print(val, " ")
				c++
			}
		}
		fmt.Println()
	}
	g.n = n

	FindMaxMatching(g)

	fmt.Print("Maximum Matching = ")
	for u := 0; u < g.n; u++ {
		fmt.Printf("%d[%d] ", u, g.M[u])
	}
	fmt.Println()
}
