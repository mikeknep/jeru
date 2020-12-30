package lib

import (
	"fmt"
	"io"
	"os/exec"
)

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

func Terraform(args []string, stdout io.Writer) *exec.Cmd {
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = stdout
	cmd.Stderr = nil
	return cmd
}
