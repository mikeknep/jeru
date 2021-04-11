package lib

// import (
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestBestEffortRefactorFinder(t *testing.T) {
// 	tests := map[string]struct {
// 		changes  []ChangingResource
// 		expected []Refactor
// 	}{
// 		"OneUpOneDownOneExtra": {
// 			changes: []ChangingResource{
// 				{
// 					Address:      "some_resource.old",
// 					Change:       Change{Actions: []string{"delete"}, Before: &legacyDB},
// 					Name:         "old",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 				{
// 					Address:      "some_resource.new",
// 					Change:       Change{Actions: []string{"create"}, After: &legacyDB},
// 					Name:         "new",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 				{
// 					Address:      "completely_different.foo",
// 					Change:       Change{Actions: []string{"create"}, After: &greenfieldDB},
// 					Name:         "foo",
// 					ProviderName: "some_provider",
// 					Type:         "completely_different",
// 				},
// 			},
// 			expected: []Refactor{
// 				{OldAddress: "some_resource.old", NewAddress: "some_resource.new"},
// 			},
// 		},
// 		"TwoUpTwoDown": {
// 			changes: []ChangingResource{
// 				{
// 					Address:      "some_resource.old_one",
// 					Change:       Change{Actions: []string{"delete"}, Before: &legacyDB},
// 					Name:         "old_one",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 				{
// 					Address:      "some_resource.new_one",
// 					Change:       Change{Actions: []string{"create"}, After: &legacyDB},
// 					Name:         "new_one",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 				{
// 					Address:      "some_resource.old_two",
// 					Change:       Change{Actions: []string{"delete"}, Before: &greenfieldDB},
// 					Name:         "old_two",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 				{
// 					Address:      "some_resource.new_two",
// 					Change:       Change{Actions: []string{"create"}, After: &greenfieldDB},
// 					Name:         "new_two",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 			},
// 			expected: []Refactor{
// 				{OldAddress: "some_resource.old_one", NewAddress: "some_resource.new_one"},
// 				{OldAddress: "some_resource.old_two", NewAddress: "some_resource.new_two"},
// 			},
// 		},
// 		"BestMatchAmongTwoNewOptions": {
// 			// one refactor plus one new of same type
// 			// should choose the "closer" of the two new resource
// 			changes: []ChangingResource{
// 				{
// 					Address:      "some_resource.main",
// 					Change:       Change{Actions: []string{"delete"}, Before: &legacyDB},
// 					Name:         "main",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 				{
// 					Address:      "some_resource.greenfield",
// 					Change:       Change{Actions: []string{"create"}, After: &greenfieldDB},
// 					Name:         "greenfield",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 				{
// 					Address:      "some_resource.legacy",
// 					Change:       Change{Actions: []string{"create"}, After: &legacyDB},
// 					Name:         "legacy",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 			},
// 			expected: []Refactor{
// 				{OldAddress: "some_resource.main", NewAddress: "some_resource.legacy"},
// 			},
// 		},
// 		"BestMatchAmongTwoOldOptions": {
// 			// one refactor plus one destroy of same type
// 			// should not recommend two destroyed objects both be moved to the same new address
// 			changes: []ChangingResource{
// 				{
// 					Address:      "some_resource.main",
// 					Change:       Change{Actions: []string{"create"}, After: &legacyDB},
// 					Name:         "main",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 				{
// 					Address:      "some_resource.greenfield",
// 					Change:       Change{Actions: []string{"delete"}, Before: &greenfieldDB},
// 					Name:         "greenfield",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 				{
// 					Address:      "some_resource.legacy",
// 					Change:       Change{Actions: []string{"delete"}, Before: &legacyDB},
// 					Name:         "legacy",
// 					ProviderName: "some_provider",
// 					Type:         "some_resource",
// 				},
// 			},
// 			expected: []Refactor{
// 				{OldAddress: "some_resource.legacy", NewAddress: "some_resource.main"},
// 			},
// 		},
// 	}
//
// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			tfPlan := TfPlan{tc.changes}
// 			got, _ := BestEffortRefactorFinder(tfPlan)
// 			assert.Equal(t, tc.expected, got)
// 		})
// 	}
// }
//
// var (
// 	legacyDB = map[string]interface{}{
// 		"engine":         "postgres",
// 		"engine_version": "9.6",
// 	}
// 	greenfieldDB = map[string]interface{}{
// 		"engine":         "postgres",
// 		"engine_version": "13.1",
// 	}
// )
