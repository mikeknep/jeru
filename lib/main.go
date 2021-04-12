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

func (p TfPlan) MvCandidates() []ChangingResource {
	var movable []ChangingResource

	// "replace" actions are represented as either ["delete", "create"] or ["create", "delete"]
	// These changes cannot be avoided with terraform state mv
	for _, cr := range p.ChangingResources {
		if len(cr.Change.Actions) != 2 {
			movable = append(movable, cr)
		}
	}

	return movable
}

type ChangingResource struct {
	Address      string `json:"address"`
	Change       Change `json:"change"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Type         string `json:"type"`
}

type Change struct {
	Actions []string `json:"actions"`
	After   *map[string]interface{}
	Before  *map[string]interface{}
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
func (cr ChangingResource) String() string {
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
