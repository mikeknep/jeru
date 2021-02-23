package lib

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	mv = "terraform state mv module.a module.b"
	im = "terraform import resource.a id"

	noPresent = func(string, []string) { return }
	noRun     = func(string) error { return nil }
)

func TestRollbackWritesGeneratedRollbackLinesToScriptInReverseOrder(t *testing.T) {
	changes := strings.NewReader(strings.Join([]string{mv, im}, "\n"))
	var builder strings.Builder
	script := Script{Name: "", W: &builder}

	Rollback(changes, &script, noPresent, noRun)

	expectedRollback := `#! /bin/bash
terraform state rm resource.a
terraform state mv module.b module.a`
	require.Equal(t, expectedRollback, builder.String())
}

func TestRollbackExecutesScriptUsingRunFunc(t *testing.T) {
	script := Script{Name: "Rollback", W: ioutil.Discard} // this moves from ioutil to io in Go 1.16

	executedScript := ""
	spyRun := func(name string) error {
		executedScript = name
		return nil
	}

	Rollback(strings.NewReader(mv), &script, noPresent, spyRun)

	require.Equal(t, "Rollback", executedScript)
}
