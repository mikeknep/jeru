package lib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func testNode(action string, resourceType string, name string) ChangingResource {
	return ChangingResource{
		Address: fmt.Sprintf("%s.%s.%s", action, resourceType, name),
		Change: Change{
			Actions: []string{action},
			After:   nil,
			Before:  nil,
		},
		Name:         name,
		ProviderName: "test",
		Type:         resourceType,
	}
}

var (
	createBucket1 = testNode("create", "bucket", "bucketOne")
	deleteBucket1 = testNode("delete", "bucket", "bucketOne")

	createBucket2 = testNode("create", "bucket", "bucketTwo")
	deleteBucket2 = testNode("delete", "bucket", "bucketTwo")

	createDatabase1 = testNode("create", "database", "dbOne")
	deleteDatabase1 = testNode("delete", "database", "dbOne")

	createDatabase2 = testNode("create", "database", "dbTwo")
	deleteDatabase2 = testNode("delete", "database", "dbTwo")

	createInstance = testNode("create", "instance", "instanceOne")
	deleteInstance = testNode("delete", "instance", "instanceOne")
)

func TestGraphs(t *testing.T) {
	tests := map[string]struct {
		nodes        []*ChangingResource
		expectedSets [][]Edge
	}{
		"Simplest case: a single pair": {
			nodes: []*ChangingResource{
				&createBucket1,
				&deleteBucket1,
			},
			expectedSets: [][]Edge{
				[]Edge{
					Edge{a: &createBucket1, b: &deleteBucket1},
				},
			},
		},

		"Simple odd number": {
			nodes: []*ChangingResource{
				&createBucket1,
				&deleteBucket1,
				&createBucket2,
			},
			expectedSets: [][]Edge{
				[]Edge{
					Edge{a: &createBucket1, b: &deleteBucket1},
				},
				[]Edge{
					Edge{a: &createBucket2, b: &deleteBucket1},
				},
			},
		},

		"Two pairs of different types": {
			nodes: []*ChangingResource{
				&createBucket1,
				&deleteBucket1,
				&createDatabase1,
				&deleteDatabase1,
			},
			expectedSets: [][]Edge{
				[]Edge{
					Edge{a: &createBucket1, b: &deleteBucket1},
					Edge{a: &createDatabase1, b: &deleteDatabase1},
				},
			},
		},

		"Two pairs of same type": {
			nodes: []*ChangingResource{
				&createBucket1,
				&deleteBucket1,
				&createBucket2,
				&deleteBucket2,
			},
			expectedSets: [][]Edge{
				[]Edge{
					// This set contains the "correct" pairs
					Edge{a: &createBucket1, b: &deleteBucket1},
					Edge{a: &createBucket2, b: &deleteBucket2},
				},
				[]Edge{
					// This set is also *valid*
					Edge{a: &createBucket1, b: &deleteBucket2},
					Edge{a: &createBucket2, b: &deleteBucket1},
				},
			},
		},

		"Larger odd number": {
			nodes: []*ChangingResource{
				&createBucket1,
				&deleteBucket1,
				&createInstance,
				&createBucket2,
				&deleteBucket2,
			},
			expectedSets: [][]Edge{
				[]Edge{
					Edge{a: &createBucket1, b: &deleteBucket1},
					Edge{a: &createBucket2, b: &deleteBucket2},
				},
				[]Edge{
					Edge{a: &createBucket1, b: &deleteBucket2},
					Edge{a: &createBucket2, b: &deleteBucket1},
				},
			},
		},

		"Larger odd number different order": {
			nodes: []*ChangingResource{
				&createInstance,
				&createBucket1,
				&deleteBucket1,
				&createBucket2,
				&deleteBucket2,
			},
			expectedSets: [][]Edge{
				[]Edge{
					Edge{a: &createBucket1, b: &deleteBucket1},
					Edge{a: &createBucket2, b: &deleteBucket2},
				},
				[]Edge{
					Edge{a: &createBucket1, b: &deleteBucket2},
					Edge{a: &createBucket2, b: &deleteBucket1},
				},
			},
		},

		"Complex": {
			nodes: []*ChangingResource{
				&createBucket1,
				&deleteBucket1,
				&createBucket2,
				&deleteBucket2,
				&createDatabase1,
				&deleteDatabase1,
				&createDatabase2,
				&deleteDatabase2,
				&createInstance,
				&deleteInstance,
			},
			expectedSets: [][]Edge{
				[]Edge{
					// correct buckets, correct databases
					Edge{a: &createBucket1, b: &deleteBucket1},
					Edge{a: &createBucket2, b: &deleteBucket2},
					Edge{a: &createDatabase1, b: &deleteDatabase1},
					Edge{a: &createDatabase2, b: &deleteDatabase2},
					Edge{a: &createInstance, b: &deleteInstance},
				},
				[]Edge{
					// correct buckets, wrong databases
					Edge{a: &createBucket1, b: &deleteBucket1},
					Edge{a: &createBucket2, b: &deleteBucket2},
					Edge{a: &createDatabase1, b: &deleteDatabase2},
					Edge{a: &createDatabase2, b: &deleteDatabase1},
					Edge{a: &createInstance, b: &deleteInstance},
				},
				[]Edge{
					// wrong buckets, correct databases
					Edge{a: &createBucket1, b: &deleteBucket2},
					Edge{a: &createBucket2, b: &deleteBucket1},
					Edge{a: &createDatabase1, b: &deleteDatabase1},
					Edge{a: &createDatabase2, b: &deleteDatabase2},
					Edge{a: &createInstance, b: &deleteInstance},
				},
				[]Edge{
					// wrong buckets, wrong databases
					Edge{a: &createBucket1, b: &deleteBucket2},
					Edge{a: &createBucket2, b: &deleteBucket1},
					Edge{a: &createDatabase1, b: &deleteDatabase2},
					Edge{a: &createDatabase2, b: &deleteDatabase1},
					Edge{a: &createInstance, b: &deleteInstance},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// first, sanity check that test setup is valid
			for _, set := range tc.expectedSets {
				require.True(t, isValidSet(set), "Test setup is incorrect (includes invalid set)")
			}

			// then execute the test
			computedSets := validEdgeCombinationsFor(tc.nodes)
			require.ElementsMatch(t, tc.expectedSets, computedSets)
		})
	}
}
