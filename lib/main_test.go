package lib

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatsRefactorAsTerraformCommand(t *testing.T) {
	r := Refactor{
		NewAddress: "new",
		OldAddress: "old",
	}

	require.Equal(t, "terraform state mv old new", r.AsCommand())
}

func TestIdentifiesASimpleRefactorMatchingOnType(t *testing.T) {
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
	expectedRefactor := Refactor{
		OldAddress: "some_resource.old",
		NewAddress: "some_resource.new",
	}

	require.Equal(t, expectedRefactor, plan.PossibleRefactors()[0])
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

func TestParsesJsonToTfPlan(t *testing.T) {
	jsonPlanContent, _ := ioutil.ReadFile("../fixtures/plan.json")
	tfPlan, _ := NewTfPlan(bytes.NewReader(jsonPlanContent))

	expectedTfPlan := TfPlan{[]ChangingResource{
		ChangingResource{
			Address: "local_file.main",
			Type:    "local_file",
			Change: Change{
				Actions: []string{"delete"},
				Before: &map[string]interface{}{
					"content":              "example",
					"content_base64":       nil,
					"directory_permission": "0777",
					"file_permission":      "0777",
					"filename":             "./example",
					"id":                   "c3499c2729730a7f807efb8676a92dcb6f8a3f8f",
					"sensitive_content":    nil,
					"source":               nil,
				},
			},
			Name:         "main",
			ProviderName: "registry.terraform.io/hashicorp/local",
		},
		ChangingResource{
			Address: "local_file.test",
			Type:    "local_file",
			Change: Change{
				Actions: []string{"create"},
				After: &map[string]interface{}{
					"content":              "example",
					"content_base64":       nil,
					"directory_permission": "0777",
					"file_permission":      "0777",
					"filename":             "./example",
					"sensitive_content":    nil,
					"source":               nil,
				},
			},
			Name:         "test",
			ProviderName: "registry.terraform.io/hashicorp/local",
		},
	}}

	require.Equal(t, 2, len(tfPlan.ChangingResources))
	require.Equal(t, expectedTfPlan, tfPlan)
}
