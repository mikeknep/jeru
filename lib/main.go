package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

type NamedWriter interface {
	io.Writer
	Name() string
}

type Refactor struct {
	NewAddress string
	OldAddress string
}

func (r Refactor) AsCommand() string {
	return fmt.Sprintf("terraform state mv %s %s", r.OldAddress, r.NewAddress)
}

type TfPlan struct {
	ChangingResources []ChangingResource `json:"resource_changes"`
}

func CandidatesByResourceType(plan TfPlan) map[string]Candidates {
	candidatesByType := make(map[string]Candidates)

	for _, r := range plan.ChangingResources {
		// "replace" actions are represented as either ["delete", "create"] or ["create", "delete"], and cannot be avoided with terraform state mv
		if len(r.Change.Actions) != 1 {
			continue
		}

		resourceType := r.GetType()
		candidates := candidatesByType[resourceType]
		updatedCandidates := candidates.Add(r)
		candidatesByType[resourceType] = updatedCandidates
	}

	// There must be at least one resource in both "creating" and "deleting"
	for t, candidates := range candidatesByType {
		if len(candidates.Creating) == 0 || len(candidates.Deleting) == 0 {
			delete(candidatesByType, t)
		}
	}

	return candidatesByType
}

type Candidates struct {
	Creating []ChangingResource
	Deleting []ChangingResource
}

func (c Candidates) Add(r ChangingResource) Candidates {
	// "replace" actions are represented as either ["delete", "create"] or ["create", "delete"], and cannot be avoided with terraform state mv
	if len(r.Change.Actions) != 1 {
		return c
	}

	switch r.Change.Actions[0] {
	case "create":
		return Candidates{Creating: append(c.Creating, r), Deleting: c.Deleting}
	case "delete":
		return Candidates{Creating: c.Creating, Deleting: append(c.Deleting, r)}
	default:
		return c
	}
}

func (c Candidates) All() []ChangingResource {
	all := c.Creating
	all = append(all, c.Deleting...)
	return all
}

type ChangingResource struct {
	Address      string `json:"address"`
	Change       Change `json:"change"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Type         string `json:"type"`
}

type Change struct {
	Actions []string                `json:"actions"`
	After   *map[string]interface{} `json:"after"`
	Before  *map[string]interface{} `json:"before"`
}

func (cr ChangingResource) GetAction() string {
	return cr.Change.Actions[0]
}
func (cr ChangingResource) GetType() string {
	return cr.Type
}
func (cr ChangingResource) GetAddress() string {
	return cr.Address
}

func NewTfPlan(r io.Reader) (TfPlan, error) {
	var tfPlan TfPlan
	err := json.NewDecoder(r).Decode(&tfPlan)
	return tfPlan, err
}

func ConsumeByLine(reader io.Reader, f func(string)) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		f(scanner.Text())
	}
	return scanner.Err()
}
