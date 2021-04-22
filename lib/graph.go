package lib

import "fmt"

type Node interface {
	GetAction() string // "create" or "delete"
	GetAddress() string
	GetType() string
}

func createEdge(x, y Node) Edge {
	if x.GetAction() == "create" {
		return Edge{a: x, b: y}
	}
	return Edge{a: y, b: x}
}

type Edge struct {
	a Node
	b Node
}

func (e Edge) String() string {
	return fmt.Sprintf("<\n\ta: %s\n\tb: %s\n>", e.a, e.b)
}

func (e Edge) isValid() bool {
	sameType := e.a.GetType() == e.b.GetType()
	differentActions := e.a.GetAction() != e.b.GetAction()

	return sameType && differentActions
}

func isValidSet(set []Edge) bool {
	var seenNodes []Node

	for _, edge := range set {
		// a set of edges is only valid if all its component edges are valid
		if !edge.isValid() {
			return false
		}

		// a set of edges is only valid if the edges' nodes are unique
		// aka, each node can only have zero or one edge
		for _, seenNode := range seenNodes {
			if seenNode.GetAddress() == edge.a.GetAddress() || seenNode.GetAddress() == edge.b.GetAddress() {
				return false
			}
		}

		seenNodes = append(seenNodes, edge.a)
		seenNodes = append(seenNodes, edge.b)
	}

	return true
}

func validEdgeCombinationsFor(nodes []Node) [][]Edge {
	return find(nodes, []Edge{})
}

func find(nodes []Node, current []Edge) (results [][]Edge) {
	if len(nodes) == 3 {
		currentPlusOne := append(current, createEdge(nodes[0], nodes[1]))
		currentPlusTwo := append(current, createEdge(nodes[1], nodes[2]))
		currentPlusThree := append(current, createEdge(nodes[0], nodes[2]))

		if isValidSet(currentPlusOne) {
			results = append(results, currentPlusOne)
		}
		if isValidSet(currentPlusTwo) {
			results = append(results, currentPlusTwo)
		}
		if isValidSet(currentPlusThree) {
			results = append(results, currentPlusThree)
		}
		return
	}

	if len(nodes) < 2 {
		if isValidSet(current) {
			results = append(results, current)
		}
		return
	}

	nodeA := nodes[0] // pluck the first node from the set of nodes

	// ensure we don't change the original nodes. slices do not copy!
	remNodes := make([]Node, len(nodes)-1)
	copy(remNodes, nodes[1:])

	for i := 0; i < len(remNodes); i++ {
		nodeB := remNodes[i]             // pluck another node from the set of nodes...
		edge := createEdge(nodeA, nodeB) // ...and create an Edge

		// add the Edge to the set we're currently building
		current = append(current, edge)

		// remove the plucked nodeB from remNodes
		// by copying all nodes up to that node,
		// and all nodes after that node,
		// and stitching those two slices together
		nextSetFirstPart := make([]Node, i)
		nextSetSecondPart := make([]Node, len(remNodes)-i-1)
		copy(nextSetFirstPart, remNodes[:i])
		copy(nextSetSecondPart, remNodes[i+1:])
		nextSet := append(nextSetFirstPart, nextSetSecondPart...)

		// recursively find more edges
		results = append(results, find(nextSet, current)...)

		// clear out the current collection as we bubble up out of recursion
		current = current[:len(current)-1]
	}
	return
}
