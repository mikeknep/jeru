package lib

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestFormatsRefactorAsTerraformCommand(t *testing.T) {
	r := Refactor{
		NewAddress: "new",
		OldAddress: "old",
	}

	require.Equal(t, `terraform state mv "old" "new"`, r.AsCommand())
}

func TestFormatsRefactorAsTerraformCommandWithQuotes(t *testing.T) {
	r := Refactor{
		NewAddress: `module.repos["1"]`,
		OldAddress: `module.repositories["1"]`,
	}

	require.Equal(t, `terraform state mv "module.repositories[\"1\"]" "module.repos[\"1\"]"`, r.AsCommand())
}
