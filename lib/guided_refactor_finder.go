package lib

import "fmt"

type GuidedRefactorFinder struct {
	Prompt Prompt
}

func (f GuidedRefactorFinder) Find(plan TfPlan) ([]Refactor, error) {
	var refactors []Refactor

	groupedCandidates := CandidatesByResourceType(plan)
	for t, candidates := range groupedCandidates {
		selectFromThisGroup, err := f.Prompt.Confirm(fmt.Sprintf("Do you want to review %s resources?", t))
		if err != nil {
			return nil, err
		}
		if !selectFromThisGroup {
			continue
		}

		remainingCandidates := &candidates

		for true {
			match, rem, err := getMatch(f.Prompt, remainingCandidates)
			if err != nil {
				return nil, err
			}
			if match != nil {
				refactors = append(refactors, *match)
			}
			remainingCandidates = rem

			if len(remainingCandidates.All()) < 2 {
				break
			}
			continueThisGroup, err := f.Prompt.Confirm(fmt.Sprintf("Are there any more %s resources you want to move?", t))
			if err != nil {
				return nil, err
			}
			if !continueThisGroup {
				break
			}
		}
	}

	return refactors, nil
}

func getMatch(prompt Prompt, candidates *Candidates) (*Refactor, *Candidates, error) {
	cancel := "[ cancel ]"
	remainingCandidates := *candidates

	dSelection, err := prompt.Select(append(remainingCandidates.OldAddresses(), cancel), "Choose a resource currently planned for deletion")
	if err != nil {
		return nil, nil, err
	}
	if dSelection == cancel {
		return nil, candidates, nil
	}
	remainingCandidates = remainingCandidates.Remove(dSelection)
	cMatch, err := prompt.Select(append(remainingCandidates.NewAddresses(), cancel), fmt.Sprintf("Choose the resource planned for creation that matches %s", dSelection))
	if err != nil {
		return nil, nil, err
	}
	if cMatch == cancel {
		return nil, candidates, nil
	}
	remainingCandidates = remainingCandidates.Remove(cMatch)

	return &Refactor{OldAddress: dSelection, NewAddress: cMatch}, &remainingCandidates, nil
}
