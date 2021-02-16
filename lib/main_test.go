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

func TestIdentifiesASimplePossibleRefactorMatchingOnType(t *testing.T) {
	plan := Plan{
		ChangingResources: []ChangingResource{
			ChangingResource{
				Address:      "some_resource.old",
				Change:       Change{Actions: []string{"delete"}},
				Name:         "old",
				ProviderName: "some_provider",
				Type:         "some_resource",
			},
			ChangingResource{
				Address:      "some_resource.new",
				Change:       Change{Actions: []string{"create"}},
				Name:         "new",
				ProviderName: "some_provider",
				Type:         "some_resource",
			},
			ChangingResource{
				Address:      "completely_different.foo",
				Change:       Change{Actions: []string{"create"}},
				Name:         "foo",
				ProviderName: "some_provider",
				Type:         "completely_different",
			},
		},
	}
	expectedPossibleRefactor := PossibleRefactor{
		OldAddress: "some_resource.old",
		NewAddress: "some_resource.new",
	}

	require.Equal(t, expectedPossibleRefactor, plan.PossibleRefactors()[0])
}

func parsePlanFile(path string) Plan {
	jsonFile, _ := os.Open(path)
	bytes, _ := ioutil.ReadAll(jsonFile)
	var plan Plan
	json.Unmarshal(bytes, &plan)
	return plan
}

func TestRollbackStateMv(t *testing.T) {
	rollback := GenerateRollbackLine("terraform state mv module.a module.b")

	require.Equal(t, "terraform state mv module.b module.a", rollback)
}

func TestRollbackImport(t *testing.T) {
	rollback := GenerateRollbackLine("terraform import module.a identifier")

	require.Equal(t, "terraform state rm module.a", rollback)
}

func TestRollbackStateRm(t *testing.T) {
	rollback := GenerateRollbackLine("terraform state rm module.a")

	require.Regexp(t, "^#", rollback)       // is a comment
	require.Regexp(t, "module.a", rollback) // includes the address of the removed resource
}

func TestRollbackUnrecognizable(t *testing.T) {
	rollback := GenerateRollbackLine("terraform plan")

	require.Regexp(t, "^#", rollback)             // is a comment
	require.Regexp(t, "terraform plan", rollback) // includes the original command
}
