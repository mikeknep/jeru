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
	mv = "terraform state mv module.a module.b"
	im = "terraform import resource.a id"

	rollbackMv = "terraform state mv module.b module.a"
	rollbackIm = "terraform state rm resource.a"

	changes = func() io.Reader { return strings.NewReader(mv + "\n" + im) }

	approve   = func() (bool, error) { return true, nil }
	unapprove = func() (bool, error) { return false, nil }

	execute = func(string, ...string) error { return nil }
)

func spyExecute() (*strings.Builder, func(string, ...string) error) {
	var executedCommands strings.Builder
	return &executedCommands, func(command string, args ...string) error {
		fullCommand := command + " " + strings.Join(args, " ")
		executedCommands.Write([]byte(fullCommand))
		return nil
	}
}

func TestPrintsGeneratedRollbackLinesToTheScreen(t *testing.T) {
	var screen strings.Builder

	Rollback(changes(), &screen, ioutil.Discard, approve, execute)

	expectedScreenContent := fmt.Sprintf(`%s
	%s
	%s
`, introText, rollbackIm, rollbackMv)
	require.Equal(t, expectedScreenContent, screen.String())
}

func TestWritesGeneratedRollbackLinesToTheOutfile(t *testing.T) {
	var outfile strings.Builder

	Rollback(changes(), ioutil.Discard, &outfile, approve, execute)

	expectedOutfileContent := rollbackIm + "\n" + rollbackMv
	require.Equal(t, expectedOutfileContent, outfile.String())
}

func TestExitsWithoutExecutingIfUserDoesNotApprove(t *testing.T) {
	executedCommands, execute := spyExecute()

	Rollback(changes(), ioutil.Discard, ioutil.Discard, unapprove, execute)

	require.Equal(t, "", executedCommands.String())
}

func TestExecutesRollbackLinesIfUserApproves(t *testing.T) {
	executedCommands, execute := spyExecute()

	Rollback(changes(), ioutil.Discard, ioutil.Discard, approve, execute)

	expectedExecution := "bash -c " + rollbackIm + "bash -c " + rollbackMv
	require.Equal(t, expectedExecution, executedCommands.String())
}
