package lib

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var localStateName = "local.tfstate"

func TestPlanCmd(t *testing.T) {
	changes := strings.NewReader(mv + "\n" + rm)
	localState := CreateNamedStringbuilder(localStateName)
	var screen strings.Builder
	var void strings.Builder

	runtime := MockRuntimeEnvironment(CaptureScreenTo(&screen), CaptureVoidTo(&void))

	Plan(runtime, changes, localState)

	expectedScreen := fmt.Sprintf(`%s
terraform plan -state local.tfstate
%s
`, commentOutBackendText, reminderText)
	require.Equal(t, expectedScreen, screen.String())

	expectedLocalState := "terraform state pull\n"
	require.Equal(t, expectedLocalState, localState.String())

	expectedVoid := `rm -rf .terraform
terraform init
bash -c terraform state mv -state=local.tfstate module.a module.b
bash -c terraform state rm -state=local.tfstate resource.a
rm -rf .terraform
`
	require.Equal(t, expectedVoid, void.String())
}

func TestEndsIfUserDoesNotConfirmComentingOutBackend(t *testing.T) {
	changes := strings.NewReader(mv + "\n" + rm)
	localState := CreateNamedStringbuilder(localStateName)
	var screen strings.Builder
	var void strings.Builder

	runtime := MockRuntimeEnvironment(CaptureScreenTo(&screen), CaptureVoidTo(&void), Unapprove)

	Plan(runtime, changes, localState)

	expectedScreen := fmt.Sprintln(commentOutBackendText)
	require.Equal(t, expectedScreen, screen.String())

	expectedLocalState := "terraform state pull\n"
	require.Equal(t, expectedLocalState, localState.String())

	require.Equal(t, "", void.String())
}

func TestAppendsExtraArgumentsToFinalPlan(t *testing.T) {
	changes := strings.NewReader(mv + "\n" + rm)
	localState := CreateNamedStringbuilder(localStateName)
	var screen strings.Builder

	runtime := MockRuntimeEnvironment(CaptureScreenTo(&screen), WithArgs("-var-file", "dev.tfvars"))

	Plan(runtime, changes, localState)

	planWithExtraArgs := "terraform plan -state local.tfstate -var-file dev.tfvars"
	require.Contains(t, screen.String(), planWithExtraArgs)
}
