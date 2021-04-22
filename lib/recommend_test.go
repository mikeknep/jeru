package lib

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var planfileName = "planfile"

func NewMockReadWriter(returning *bytes.Reader) MockReadWriter {
	var b strings.Builder
	return MockReadWriter{&b, returning}
}

type MockReadWriter struct {
	builder *strings.Builder
	reader  *bytes.Reader
}

func (mrw MockReadWriter) Write(p []byte) (int, error) {
	return mrw.builder.Write(p)
}

func (mrw MockReadWriter) Read(b []byte) (int, error) {
	return mrw.reader.Read(b)
}

func TestRecommendRunsTerraformPlanAndShowCommands(t *testing.T) {
	planfile := CreateNamedStringbuilder(planfileName)
	jsonPlan := NewMockReadWriter(bytes.NewReader([]byte("{}")))
	screen := ioutil.Discard
	var void strings.Builder

	Recommend(planfile, jsonPlan, screen, &void, spyPlanExecute, BestEffortRefactorFinder, []string{})

	expectedVoid := "terraform plan -out planfile\n"
	require.Equal(t, expectedVoid, void.String())

	expectedJsonPlan := "terraform show -json planfile\n"
	require.Equal(t, expectedJsonPlan, jsonPlan.builder.String())
}
