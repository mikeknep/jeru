package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	createBucket1 = Node{Action: "create", Type: "bucket", Name: "bucketOne"}
	deleteBucket1 = Node{Action: "delete", Type: "bucket", Name: "bucketOne"}

	createBucket2 = Node{Action: "create", Type: "bucket", Name: "bucketTwo"}
	deleteBucket2 = Node{Action: "delete", Type: "bucket", Name: "bucketTwo"}

	createDatabase1 = Node{Action: "create", Type: "database", Name: "dbOne"}
	deleteDatabase1 = Node{Action: "delete", Type: "database", Name: "dbOne"}

	createDatabase2 = Node{Action: "create", Type: "database", Name: "dbTwo"}
	deleteDatabase2 = Node{Action: "delete", Type: "database", Name: "dbTwo"}

	createInstance = Node{Action: "create", Type: "instance", Name: "instanceOne"}
	deleteInstance = Node{Action: "delete", Type: "instance", Name: "instanceOne"}
)

func TestGraphs(t *testing.T) {
	tests := map[string]struct {
		nodes        []Node
		expectedSets [][]Edge
	}{
		"Simplest case: a single pair": {
			nodes: []Node{
				createBucket1,
				deleteBucket1,
			},
			expectedSets: [][]Edge{
				[]Edge{
					Edge{a: &createBucket1, b: &deleteBucket1},
				},
			},
		},

		"Two pairs of different types": {
			nodes: []Node{
				createBucket1,
				deleteBucket1,
				createDatabase1,
				deleteDatabase1,
			},
			expectedSets: [][]Edge{
				[]Edge{
					Edge{a: &createBucket1, b: &deleteBucket1},
					Edge{a: &createDatabase1, b: &deleteDatabase1},
				},
			},
		},

		"Two pairs of same type": {
			nodes: []Node{
				createBucket1,
				deleteBucket1,
				createBucket2,
				deleteBucket2,
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

		"Complex": {
			nodes: []Node{
				createBucket1,
				deleteBucket1,
				createBucket2,
				deleteBucket2,
				createDatabase1,
				deleteDatabase1,
				createDatabase2,
				deleteDatabase2,
				createInstance,
				deleteInstance,
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
				require.True(t, isValidSet(set))
			}

			// then execute the test
			computedSets := validEdgeCombinationsFor(tc.nodes)
			require.Equal(t, tc.expectedSets, computedSets)
		})
	}
}
