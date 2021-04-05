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
