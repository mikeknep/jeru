package lib

type scoredRefactor struct {
	Refactor *Refactor
	Score    float64
}

func BestEffortRefactorFinder(plan TfPlan) ([]Refactor, error) {
	// Collect all changing resources that could potentially be refactored,
	// organized into those being deleted and those being created
	var beingDeleted []ChangingResource
	var beingCreated []ChangingResource
	for _, cr := range plan.ChangingResources {
		actions := cr.Change.Actions
		switch {
		case len(actions) == 2:
			// "replace" actions are represented as either ["delete", "create"] or ["create", "delete"]
			// These changes cannot be avoided with terraform state mv, and so are ignored here
		case actions[0] == "delete":
			beingDeleted = append(beingDeleted, cr)
		case actions[0] == "create":
			beingCreated = append(beingCreated, cr)
		}
	}

	// Calculate a similarity score for every create/destroy pair of same-type resources
	var scoredRefactors []scoredRefactor
	for _, deleting := range beingDeleted {
		for _, creating := range beingCreated {
			if isSameResourceType(creating, deleting) {
				refactor := Refactor{
					OldAddress: deleting.Address,
					NewAddress: creating.Address,
				}
				scoredRefactors = append(scoredRefactors, scoredRefactor{
					Refactor: &refactor,
					Score:    getScore(creating, deleting),
				})
			}
		}
	}

	// Organize scored refactors into sets where no OldAddress or NewAddress is duplicated
	var setsOfScoredRefactors [][]scoredRefactor

	// Choose the set of scored refactors with the highest cumulative score
	var bestSet []scoredRefactor
	var bestScore float64
	for _, scoredRefactorSet := range setsOfScoredRefactors {
		thisScore := cumulativeScore(scoredRefactorSet)
		if thisScore > bestScore {
			bestSet = scoredRefactorSet
			bestScore = thisScore
		}
	}

	// Return the refactors from that set
	var bestRefactors []Refactor
	for _, scoredRefactor := range bestSet {
		bestRefactors = append(bestRefactors, *scoredRefactor.Refactor)
	}

	return bestRefactors, nil
}

func cumulativeScore(scoredRefactorSet []scoredRefactor) float64 {
	var score float64
	for _, scoredRefactor := range scoredRefactorSet {
		score = score + scoredRefactor.Score
	}
	return score
}

func isSameResourceType(a, b ChangingResource) bool {
	return a.Type == b.Type && a.ProviderName == b.ProviderName
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
