package lib

import (
	"encoding/json"
	"fmt"
	"io"
)

type Refactor struct {
	NewAddress string
	OldAddress string
}

func (r Refactor) AsCommand() string {
	return fmt.Sprintf("terraform state mv %q %q", r.OldAddress, r.NewAddress)
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

func (c Candidates) Remove(address string) Candidates {
	var newCreating []ChangingResource
	var newDeleting []ChangingResource

	for _, r := range c.Creating {
		if r.GetAddress() != address {
			newCreating = append(newCreating, r)
		}
	}

	for _, r := range c.Deleting {
		if r.GetAddress() != address {
			newDeleting = append(newDeleting, r)
		}
	}

	return Candidates{Creating: newCreating, Deleting: newDeleting}
}

func (c Candidates) All() []ChangingResource {
	all := c.Creating
	all = append(all, c.Deleting...)
	return all
}

func (c Candidates) NewAddresses() []string {
	var addrs []string
	for _, r := range c.Creating {
		addrs = append(addrs, r.GetAddress())
	}
	return addrs
}

func (c Candidates) OldAddresses() []string {
	var addrs []string
	for _, r := range c.Deleting {
		addrs = append(addrs, r.GetAddress())
	}
	return addrs
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
