package lib

import (
	"io"
	"os"
	"os/exec"
)

type RuntimeEnvironment struct {
	Execute      CommandRunner
	ExtraArgs    []string
	GetApproval  GetApproval
	Prompt       Prompt
	Screen       io.Writer
	StartSpinner StartSpinner
	Void         io.Writer
}

type CommandRunner func(io.Writer, string, ...string) error

var liveCommandRunner CommandRunner = func(w io.Writer, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = w
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func CreateLiveRuntimeEnvironment(extraArgs []string) RuntimeEnvironment {
	return RuntimeEnvironment{
		Execute:      liveCommandRunner,
		ExtraArgs:    extraArgs,
		GetApproval:  GetApprovalFromPrompt,
		Prompt:       SurveyPrompt{},
		Screen:       os.Stdout,
		StartSpinner: StartPtermSpinner,
		Void:         io.Discard,
	}
}
