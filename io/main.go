package io

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func WriteAndRun(filename string, lines []string, persist bool) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if !persist {
		defer os.Remove(file.Name())
	}

	shebang := regexp.MustCompile(`^#!`)
	if shebang.FindStringIndex(lines[0]) == nil {
		lines = append([]string{"#! /bin/bash"}, lines...)
	}

	out := strings.Join(lines, "\n")
	if err = ioutil.WriteFile(file.Name(), []byte(out), 0777); err != nil {
		return err
	}
	file.Chmod(0777)

	fileCommand := exec.Command(file.Name())
	fileCommand.Stdout = nil
	fileCommand.Stderr = nil
	err = fileCommand.Run()
	if err != nil {
		return err
	}

	return nil
}
