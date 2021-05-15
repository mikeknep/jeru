package lib

import (
	"io"
	"strings"
)

var (
	im = "terraform import resource.a id"
	mv = "terraform state mv module.a module.b"
	rm = "terraform state rm resource.a"
)

type NamedStringbuilder struct {
	name   string
	writer *strings.Builder
}

func CreateNamedStringbuilder(name string) *NamedStringbuilder {
	var builder strings.Builder
	return &NamedStringbuilder{name: name, writer: &builder}
}

func (nsb *NamedStringbuilder) Name() string {
	return nsb.name
}

func (nsb *NamedStringbuilder) Write(x []byte) (int, error) {
	return nsb.writer.Write(x)
}

func (nsb *NamedStringbuilder) String() string {
	return nsb.writer.String()
}

func MockRuntimeEnvironment(options ...func(*RuntimeEnvironment)) RuntimeEnvironment {
	runtime := RuntimeEnvironment{
		Execute:      MockExecute,
		ExtraArgs:    []string{},
		GetApproval:  AutoApprove,
		Prompt:       MockPrompt{},
		Screen:       io.Discard,
		StartSpinner: StartSilentSpinner,
		Void:         io.Discard,
	}

	for _, option := range options {
		option(&runtime)
	}

	return runtime
}

func CaptureScreenTo(w io.Writer) func(*RuntimeEnvironment) {
	return func(r *RuntimeEnvironment) {
		r.Screen = w
	}
}

func CaptureVoidTo(w io.Writer) func(*RuntimeEnvironment) {
	return func(r *RuntimeEnvironment) {
		r.Void = w
	}
}

func Unapprove(r *RuntimeEnvironment) {
	r.GetApproval = func() (bool, error) {
		return false, nil
	}
}

func WithArgs(args ...string) func(*RuntimeEnvironment) {
	return func(r *RuntimeEnvironment) {
		r.ExtraArgs = args
	}
}

var MockExecute CommandRunner = func(w io.Writer, name string, args ...string) error {
	fullCommand := name + " " + strings.Join(args, " ") + "\n"
	w.Write([]byte(fullCommand))
	return nil
}

type MockPrompt struct{}

func (p MockPrompt) Confirm(_ string) (bool, error) {
	return true, nil
}

func (p MockPrompt) Select(options []string, _ string) (string, error) {
	return options[0], nil
}
