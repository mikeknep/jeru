package lib

import "fmt"

type Plan struct {
	ChangingResources []ChangingResource `json:"resource_changes"`
}

type ChangingResource struct {
	Address      string
	Change       Change
	Name         string
	ProviderName string `json:"provider_name"`
	Type         string
}

type Change struct {
	Actions []string
	After   *map[string]interface{}
	Before  *map[string]interface{}
}

type PossibleRefactor struct {
	NewAddress string
	OldAddress string
}

func (pr PossibleRefactor) AsCommand() string {
	return fmt.Sprintf("terraform state mv %s %s", pr.OldAddress, pr.NewAddress)
}

func (plan Plan) PossibleRefactors() []PossibleRefactor {
	return []PossibleRefactor{}
}
