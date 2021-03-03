package basicgraph

// graph.go GraphStructure handles graph queries.
// Every hypergraph may be represented by a bipartite graph.
import (
	"fmt"

	hd "github.com/dgnabasik/acmsearchlib/headers"

	"github.com/yourbasic/graph" // The vertices of all graphs in this package are numbered 0..n-1.
)

// GraphInterface interface functions are not placed into acmsearchlib.
type GraphInterface interface {
	BuildUndirectedGraph(condProbList []hd.ConditionalProbability) *graph.Mutable
}

// GraphNode struct reflects IGraphNode interface in react-app-env.d.ts
type GraphNode struct {
	hd.Vocabulary // {Id, Word, RowCount, Frequency, WordRank, Probability, SpeechPart}
	hd.WordScore  // {Id, Word, Timeinterval, Density, Linkage, Growth, Score}
}

// GraphLink struct reflects IGraphLink interface in react-app-env.d.ts
type GraphLink struct {
	hd.ConditionalProbability // {Id, WordList, Probability, Timeinterval, FirstDate, LastDate, Pmi, DateUpdated}
}

// GraphStructure struct implements GraphInterface; reflects
type GraphStructure struct {
	Nodes []GraphNode
	Links []GraphLink
}

// func (g1 *Virtual) Tensor(g2 *Virtual) *Virtual
/*************************************************************************************/

// BuildUndirectedGraph func
func (bgs *GraphStructure) BuildUndirectedGraph(gs GraphStructure) *graph.Mutable {
	g := graph.New(4)
	g.AddBoth(0, 1) //  0 -- 1
	g.AddBoth(0, 2) //  |    |
	g.AddBoth(2, 3) //  2 -- 3
	g.AddBoth(1, 3)

	// Visit all edges of a graph.
	for v := 0; v < g.Order(); v++ {
		g.Visit(v, func(w int, c int64) (skip bool) {
			// Visiting edge (v, w) of cost c.
			return
		})
	}

	// The immutable data structure created by Sort has an Iterator that returns neighbors in increasing order.
	// Visit the edges in order.
	for v := 0; v < g.Order(); v++ {
		graph.Sort(g).Visit(v, func(w int, c int64) (skip bool) {
			// Visiting edge (v, w) of cost c.
			return
		})
	}

	// The return value of an iterator function is used to break out of the iteration. Visit, in turn, returns a boolean indicating if it was aborted.
	// Skip the iteration at the first edge (v, w) with w equal to 3.
	for v := 0; v < g.Order(); v++ {
		aborted := graph.Sort(g).Visit(v, func(w int, c int64) (skip bool) {
			fmt.Println(v, w)
			if w == 3 {
				skip = true // Aborts the call to Visit.
			}
			return
		})
		if aborted {
			break
		}
	}

	return g
}
