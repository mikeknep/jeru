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
	rollbackMv = "terraform state mv module.b module.a"
	rollbackIm = "terraform state rm resource.a"

	execute = func(string, ...string) error { return nil }
)

func createSourceChanges() io.Reader {
	return strings.NewReader("#!/bin/bash\n\n" + mv + "\n" + im)
}

func spyExecute() (*strings.Builder, func(string, ...string) error) {
	var executedCommands strings.Builder
	return &executedCommands, func(command string, args ...string) error {
		fullCommand := command + " " + strings.Join(args, " ")
		executedCommands.Write([]byte(fullCommand))
		return nil
	}
}

func TestPrintsGeneratedRollbackLinesToTheScreenAndAsksForApproval(t *testing.T) {
	var screen strings.Builder

	Rollback(createSourceChanges(), &screen, ioutil.Discard, approve, execute)

	expectedScreenContent := fmt.Sprintf(`%s
	%s
	%s
%s
`, introText, rollbackIm, rollbackMv, performTheseActionsText)
	require.Equal(t, expectedScreenContent, screen.String())
}

func TestWritesGeneratedRollbackLinesToTheOutfile(t *testing.T) {
	var outfile strings.Builder

	Rollback(createSourceChanges(), ioutil.Discard, &outfile, approve, execute)

	expectedOutfileContent := rollbackIm + "\n" + rollbackMv
	require.Equal(t, expectedOutfileContent, outfile.String())
}

func TestExitsWithoutExecutingIfUserDoesNotApprove(t *testing.T) {
	executedCommands, execute := spyExecute()

	Rollback(createSourceChanges(), ioutil.Discard, ioutil.Discard, unapprove, execute)

	require.Equal(t, "", executedCommands.String())
}

func TestExecutesRollbackLinesIfUserApproves(t *testing.T) {
	executedCommands, execute := spyExecute()

	Rollback(createSourceChanges(), ioutil.Discard, ioutil.Discard, approve, execute)

	expectedExecution := "bash -c " + rollbackIm + "bash -c " + rollbackMv
	require.Equal(t, expectedExecution, executedCommands.String())
}
