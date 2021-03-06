package lib

import (
	"sync"
)

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
	sameType := e.a.GetType() == e.b.GetType()
	differentActions := e.a.GetAction() != e.b.GetAction()

	return sameType && differentActions
}

func isValidSet(set []Edge) bool {
	for _, edge := range set {
		if !edge.isValid() {
			return false
		}
	}
	return true
}

func validEdgeCombinationsFor(nodes []*ChangingResource) [][]Edge {
	if len(nodes)%2 != 0 {
		nodes = append(nodes, nil)
	}
	channel := make(chan []Edge)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		waitGroup.Wait()
		close(channel)
	}()

	find(nodes, []Edge{}, channel, waitGroup)

	var results [][]Edge
	for edgeSet := range channel {
		results = append(results, edgeSet)
	}
	return results
}

func find(nodes []*ChangingResource, current []Edge, channel chan []Edge, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	if !isValidSet(current) {
		return
	}

	if len(nodes) < 2 {
		if isValidSet(current) {
			channel <- current
		}
		return
	}

	nodeA := nodes[0]

	remNodes := make([]*ChangingResource, len(nodes)-1)
	copy(remNodes, nodes[1:])

	for i := 0; i < len(remNodes); i++ {
		nodeB := remNodes[i]

		newCurrentCurrent := make([]Edge, len(current))
		copy(newCurrentCurrent, current)

		if nodeA != nil && nodeB != nil {
			newCurrentCurrent = append(newCurrentCurrent, createEdge(nodeA, nodeB))
		}

		nextSetFirstPart := make([]*ChangingResource, i)
		nextSetSecondPart := make([]*ChangingResource, len(remNodes)-i-1)
		copy(nextSetFirstPart, remNodes[:i])
		copy(nextSetSecondPart, remNodes[i+1:])
		nextSet := append(nextSetFirstPart, nextSetSecondPart...)

		waitGroup.Add(1)
		go find(nextSet, newCurrentCurrent, channel, waitGroup)
	}
}
