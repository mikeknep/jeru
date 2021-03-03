package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
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
	plan := TfPlan{
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

func parsePlanFile(path string) TfPlan {
	jsonFile, _ := os.Open(path)
	bytes, _ := ioutil.ReadAll(jsonFile)
	var plan TfPlan
	json.Unmarshal(bytes, &plan)
	return plan
}

// TODO: replace with higher level test (addRollbackLine will be private)
func TestProperlyHandlesNoOpLines(t *testing.T) {
	rollbackLines := []string{}
	srcLines := []string{"#!/bin/bash\n", "\n", ""}

	for _, line := range srcLines {
		addRollbackLine(&rollbackLines, line)
	}

	require.Equal(t, []string{}, rollbackLines)
}

// TODO: replace with higher level test (addRollbackLine will be private)
func TestGeneratesRollbackLinesInReverseOrder(t *testing.T) {
	rollbackLines := []string{}
	srcLines := []string{
		"terraform plan",
		"terraform state rm module.a",
		"terraform import module.a identifier",
		"terraform state mv module.a module.b",
	}

	for _, line := range srcLines {
		addRollbackLine(&rollbackLines, line)
	}

	require.Equal(t, "terraform state mv module.b module.a", rollbackLines[0])
	require.Equal(t, "terraform state rm module.a", rollbackLines[1])

	// can't generate rollback for removals
	require.Regexp(t, "^#", rollbackLines[2])       // is a comment
	require.Regexp(t, "module.a", rollbackLines[2]) // includes the address of the removed resource for reference

	// can't generate rollback for non-state command
	require.Regexp(t, "^#", rollbackLines[3])             // is a comment
	require.Regexp(t, "terraform plan", rollbackLines[3]) // includes the original command for reference
}

func TestActingOnReaderLines(t *testing.T) {
	lines := `one
two
three`

	lengths := []int{}
	ConsumeByLine(strings.NewReader(lines), func(line string) {
		lengths = append(lengths, len(line))
	})

	require.Equal(t, []int{3, 3, 5}, lengths)
}
