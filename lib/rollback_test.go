package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRollbackCommand(t *testing.T) {
	// TODO
	require.Equal(t, 1, 1)
}

func present(captured []string) func(intro string, lines []string) {
	return func(intro string, lines []string) {
		captured = append(captured, intro)
		for _, line := range lines {
			captured = append(captured, line)
		}
	}
}
