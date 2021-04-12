package lib

func BestEffortRefactorFinder(plan TfPlan) ([]Refactor, error) {
	candidates := plan.MvCandidates()
	candidateNodes := make([]Node, len(candidates))
	for i := range candidates {
		candidateNodes[i] = candidates[i]
	}

	validSets := validEdgeCombinationsFor(candidateNodes)

	var bestSet []Refactor
	var bestScore float64
	for _, set := range validSets {
		score := cumulativeScore(set)
		if score > bestScore {
			bestSet = asRefactors(set)
			bestScore = score
		}
	}

	return bestSet, nil
}

func asRefactors(set []Edge) []Refactor {
	var refactors []Refactor
	for _, edge := range set {
		refactors = append(refactors, Refactor{
			NewAddress: edge.a.GetAddress(),
			OldAddress: edge.b.GetAddress(),
		})
	}
	return refactors
}

func cumulativeScore(set []Edge) float64 {
	var score float64
	for _, edge := range set {
		score = score + getScore(edge.a.(ChangingResource), edge.b.(ChangingResource))
	}
	return score
}

func getScore(a, b ChangingResource) float64 {
	var numerator float64
	var denominator float64

	for k, v := range *a.Change.After {
		before := *b.Change.Before
		if before[k] == v {
			numerator = numerator + 1
		}
		denominator = denominator + 1
	}

	return numerator / denominator
}
