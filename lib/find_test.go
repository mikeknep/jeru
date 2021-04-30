package lib

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var planfileName = "planfile"

func TestFindRunsTerraformPlanAndShowCommands(t *testing.T) {
	planfile := CreateNamedStringbuilder(planfileName)
	var void strings.Builder

	runtime := MockRuntimeEnvironment(CaptureVoidTo(&void))
	flags := FindFlags{
		InteractiveMode: false,
	}

	Find(runtime, flags, planfile)

	expectedVoid := fmt.Sprintf("terraform plan -out %s\n", planfileName)
	require.Equal(t, expectedVoid, void.String())
}
