package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
