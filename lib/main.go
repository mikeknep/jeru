package lib

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
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
	var possibleRefactors []PossibleRefactor

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

	for _, deleting := range beingDeleted {
		for _, creating := range beingCreated {
			if isPossiblySameResource(deleting, creating) {
				possibleRefactors = append(possibleRefactors, PossibleRefactor{
					NewAddress: creating.Address,
					OldAddress: deleting.Address,
				})
			}
		}
	}

	return possibleRefactors
}

// this is not sophisticated enough yet
func isPossiblySameResource(a, b ChangingResource) bool {
	return a.Type == b.Type && a.ProviderName == b.ProviderName
}

func Terraform(args []string, stdout io.Writer) *exec.Cmd {
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = stdout
	cmd.Stderr = nil
	return cmd
}

func generateRollbackLine(line string) string {
	im := regexp.MustCompile(`terraform import (\S+) \S+`)
	mv := regexp.MustCompile(`terraform state mv (\S+) (\S+)`)
	rm := regexp.MustCompile(`terraform state rm (\S+)`)

	switch {
	case im.FindStringIndex(line) != nil:
		return im.ReplaceAllString(line, fmt.Sprintf("terraform state rm $1"))
	case mv.FindStringIndex(line) != nil:
		return mv.ReplaceAllString(line, fmt.Sprintf("terraform state mv $2 $1"))
	case rm.FindStringIndex(line) != nil:
		return rm.ReplaceAllString(line, fmt.Sprintf("# terraform import $1 ___"))
	default:
		return fmt.Sprintf("# Could not generate rollback command for: %s", line)
	}
}

func AddRollbackLine(rollbackLines *[]string, srcLine string) {
	if isNoopLine(srcLine) {
		return
	}
	rollbackLine := generateRollbackLine(srcLine)
	*rollbackLines = append([]string{rollbackLine}, *rollbackLines...)
}

func isNoopLine(line string) bool {
	return isEmpty(line) || isShebang(line)
}

func isEmpty(line string) bool {
	matched, _ := regexp.MatchString(`^\s*$`, line)
	return matched
}

func isShebang(line string) bool {
	matched, _ := regexp.MatchString(`^#!`, line)
	return matched
}

func OrDefault(priority, defaultVal string) string {
	if priority != "" {
		return priority
	} else {
		return defaultVal
	}
}
