package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatsPossibleRefactorAsTerraformCommand(t *testing.T) {
	pr := PossibleRefactor{
		NewAddress: "new",
		OldAddress: "old",
	}

	require.Equal(t, "terraform state mv old new", pr.AsCommand())
}

func TestParsesPlanFile(t *testing.T) {
	plan := parsePlanFile("../fixtures/plan.json")

	require.Equal(t, 2, len(plan.ChangingResources))
}

func parsePlanFile(path string) Plan {
	jsonFile, _ := os.Open(path)
	bytes, _ := ioutil.ReadAll(jsonFile)
	var plan Plan
	json.Unmarshal(bytes, &plan)
	return plan
}
