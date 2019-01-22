// The package reversearch provides an interface for exploring the
// Reverse Search algorithm of Avis and Fukuda:
//
// https://www.sciencedirect.com/science/article/pii/0166218X9500026N
//
// The preliminary implementation here is drawn from David Avis'
// tutorial:
//
// cgm.cs.mcgill.ca/~avis/doc/tutorial1/tutorial.pdf
//
// The algorithm was originally invented to investigate convex hull
// calculations, but is amenable to calculation in a variety of
// combinatorial settings.
//
// The first example here enumerates all permutations, given length n;
// others to follow.
//
// The purpose of this implementation is to take advantage of Go's
// flexible interface structure to untangle the original C
// implementation, find the proper level of generalization, and refactor
// to take advantage of the inherent recursive and concurrent nature of
// the algorithm.
package reversesearch

import (
	"fmt"
)

// ReverseSearch is the interface that wraps the necessary methods
// invoked by the reverse search algorithm, implemented in Enumerate.
type ReverseSearch interface {
	// Adjacent defines the adjacency graph of combinatorial objects,
	// for which the reverse search algorithm finds a spanning tree.
	// Adjacent should be nilpotent: (v.Adjacent).Adjacent is the
	// identity function.
	Adjacent(int) error
	// LocalSearch defines a canonical path from a node to the root
	// vertex. (v.LocalSearch).Adjacent(i) is the identity, for some i.
	Localsearch()
	// Equal tests equality of two nodes.
	Equal(ReverseSearch) bool
	// Copy makes a copy of node data.
	Copy() ReverseSearch
	// IsRoot tests whether a node is the root node of the spanning tree
	// to be determined.
	IsRoot() bool
	// MaxDeg returns the maximum degree of the adjacency graph.
	MaxDeg() int
	// Output returns the node data.
	Output(interface{})
}

// Reverse tests an index i to determine whether Adjacent(i) may be
// undone by a call to LocalSearch, hence whether it determines an edge
// in the spanning tree.
func Reverse(v ReverseSearch, i int) bool {
	w := v.Copy()

	err := w.Adjacent(i)

	if err != nil {
		fmt.Println(err)
		return false
	} else {
		w.Localsearch()
		return (v.Equal(w))
	}
}

// Backtrack returns the index i such that LocalSearch is undone by
// Adjacent, hence backtracking along the spanning tree to return to
// root.
func Backtrack(v ReverseSearch) int {
	var i int = -1
	var found bool = false

	child := v.Copy()

	v.Localsearch()

	for found == false {
		i++
		v.Adjacent(i)
		found = v.Equal(child)
		v.Adjacent(i)
	}

	return i
}

// Enumerate walks the spanning tree of the adjacency graph determined
// by Adjacent and LocalSearch, outputting the nodes.
func Enumerate(v ReverseSearch, c chan []int) {

	var i int = -1
	var count int = 1
	var mdeg int = v.MaxDeg()

	defer close(c)

	v.Output(c)

	for i < mdeg {
		for {
			i++
			if (i > mdeg-1) || (Reverse(v, i) == true) {
				break
			}
		}

		if i < mdeg {
			v.Adjacent(i)
			v.Output(c)
			count++
			i = -1
		} else if v.IsRoot() == false {
			i = Backtrack(v)
		}
	}
}
