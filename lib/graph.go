package lib

func createEdge(x, y *ChangingResource) Edge {
	if x.GetAction() == "create" {
		return Edge{a: x, b: y}
	}
	return Edge{a: y, b: x}
}

type Edge struct {
	a *ChangingResource
	b *ChangingResource
}

func (e Edge) isValid() bool {
	if e.a == nil || e.b == nil {
		return true
	}

	sameType := e.a.GetType() == e.b.GetType()
	differentActions := e.a.GetAction() != e.b.GetAction()

	return sameType && differentActions
}

func isValidSet(set []Edge) bool {
	var seenNodes []*ChangingResource

	for _, edge := range set {
		// a set of edges is only valid if all its component edges are valid
		if !edge.isValid() {
			return false
		}

		// a set of edges is only valid if the edges' nodes are unique
		// aka, each node can only have zero or one edge
		for _, seenNode := range seenNodes {
			if seenNode == edge.a || seenNode == edge.b {
				return false
			}
		}

		seenNodes = append(seenNodes, edge.a)
		seenNodes = append(seenNodes, edge.b)
	}

	return true
}

func validEdgeCombinationsFor(nodes []*ChangingResource) [][]Edge {
	if len(nodes)%2 != 0 {
		nodes = append(nodes, nil)
	}
	results := find(nodes, []Edge{})
	for i, setOfEdges := range results {
		for j, edge := range setOfEdges {
			if edge.a == nil || edge.b == nil {
				setWithoutNilFirstPart := make([]Edge, j)
				setWithoutNilSecondPart := make([]Edge, len(setOfEdges)-j-1)
				copy(setWithoutNilFirstPart, setOfEdges[:j])
				copy(setWithoutNilSecondPart, setOfEdges[j+1:])
				setWithoutNilEdge := append(setWithoutNilFirstPart, setWithoutNilSecondPart...)
				results[i] = setWithoutNilEdge
			}
		}
	}

	return results
}

func find(nodes []*ChangingResource, current []Edge) (results [][]Edge) {
	if len(nodes) < 2 {
		if isValidSet(current) {
			results = append(results, current)
		}
		return
	}

	nodeA := nodes[0] // pluck the first node from the set of nodes

	// ensure we don't change the original nodes. slices do not copy!
	remNodes := make([]*ChangingResource, len(nodes)-1)
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
		nextSetFirstPart := make([]*ChangingResource, i)
		nextSetSecondPart := make([]*ChangingResource, len(remNodes)-i-1)
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
