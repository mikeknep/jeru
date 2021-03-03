package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	localStateName = "local.tfstate"

	noOpExecute = func(io.Writer, string, ...string) error { return nil }

	void = ioutil.Discard

	noExtraArgs = []string{}
)

func spyPlanExecute(executedCommands io.Writer, command string, args ...string) error {
	fullCommand := command + " " + strings.Join(args, " ") + "\n"
	executedCommands.Write([]byte(fullCommand))
	return nil
}

type StringbuilderLocalState struct {
	name   string
	writer *strings.Builder
}

func NewLocalState(name string) *StringbuilderLocalState {
	var builder strings.Builder
	return &StringbuilderLocalState{name: name, writer: &builder}
}

func (state *StringbuilderLocalState) Name() string {
	return state.name
}

func (state *StringbuilderLocalState) Write(x []byte) (int, error) {
	return state.writer.Write(x)
}

func (state *StringbuilderLocalState) String() string {
	return state.writer.String()
}

func TestPlanCmd(t *testing.T) {
	changes := strings.NewReader(mv + "\n" + rm)
	localState := NewLocalState(localStateName)
	var screen strings.Builder
	var void strings.Builder

	Plan(
		changes,
		localState,
		&screen,
		&void,
		approve,
		spyPlanExecute,
		noExtraArgs,
	)

	expectedScreen := fmt.Sprintf(`%s
terraform init
bash -c terraform state mv -state=local.tfstate module.a module.b
bash -c terraform state rm -state=local.tfstate resource.a
terraform plan -state local.tfstate
%s
`, commentOutBackendText, reminderText)
	require.Equal(t, expectedScreen, screen.String())

	expectedLocalState := "terraform state pull\n"
	require.Equal(t, expectedLocalState, localState.String())

	expectedVoid := "rm -rf .terraform\nrm -rf .terraform\n"
	require.Equal(t, expectedVoid, void.String())
}

func TestEndsIfUserDoesNotConfirmComentingOutBackend(t *testing.T) {
	changes := strings.NewReader(mv + "\n" + rm)
	localState := NewLocalState(localStateName)
	var screen strings.Builder
	var void strings.Builder

	Plan(
		changes,
		localState,
		&screen,
		&void,
		unapprove,
		spyPlanExecute,
		noExtraArgs,
	)

	expectedScreen := fmt.Sprintln(commentOutBackendText)
	require.Equal(t, expectedScreen, screen.String())

	expectedLocalState := "terraform state pull\n"
	require.Equal(t, expectedLocalState, localState.String())

	require.Equal(t, "", void.String())
}

func TestAppendsExtraArgumentsToFinalPlan(t *testing.T) {
	changes := strings.NewReader(mv + "\n" + rm)
	localState := NewLocalState(localStateName)
	var screen strings.Builder
	void := ioutil.Discard

	Plan(
		changes,
		localState,
		&screen,
		void,
		approve,
		spyPlanExecute,
		[]string{"-var-file", "dev.tfvars"},
	)

	planWithExtraArgs := "terraform plan -state local.tfstate -var-file dev.tfvars"
	require.Contains(t, screen.String(), planWithExtraArgs)
}
