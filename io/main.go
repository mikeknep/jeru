package io

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mikeknep/jeru/lib"
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

func Run(filename string) error {
	fileCommand := exec.Command(filename)
	fileCommand.Stdout = nil
	fileCommand.Stderr = nil
	return fileCommand.Run()
}

func CreateScript(filename string) (*lib.Script, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	file.Chmod(0777)
	return &lib.Script{Name: filename, W: file}, nil
}

func DisplayIntent(intro string, lines []string) {
	fmt.Println(intro + "\n")
	for _, line := range lines {
		fmt.Println("\t" + line)
	}
}
