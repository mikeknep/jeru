package cmd

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/mikeknep/jeru/lib"
	"github.com/spf13/cobra"
)

var dryRun bool

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Revert a series of state changes",
	RunE: func(cmd *cobra.Command, args []string) error {

		input, err := ioutil.ReadFile(changeScript)
		if err != nil {
			return err
		}
		rollbackLines := []string{}
		scanner := bufio.NewScanner(strings.NewReader(string(input)))
		for scanner.Scan() {
			lib.AddRollbackLine(&rollbackLines, scanner.Text())
		}

		rollbackFile, err := os.Create("./rollback.sh")
		if err != nil {
			return err
		}
		defer rollbackFile.Close()
		rollbackLines = append([]string{"#! /bin/bash"}, rollbackLines...)
		out := strings.Join(rollbackLines, "\n")
		if err = ioutil.WriteFile(rollbackFile.Name(), []byte(out), 0777); err != nil {
			return err
		}
		rollbackFile.Chmod(0777)

		if dryRun {
			return nil
		}

		rollbackCommand := exec.Command(rollbackFile.Name())
		rollbackCommand.Stdout = nil
		rollbackCommand.Stderr = nil
		err = rollbackCommand.Run()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rollbackCmd.Flags().StringVar(&changeScript, "changes", "", "A script containing the terraform state mv|rm changes to make")
	rollbackCmd.MarkFlagRequired("changes")

	rollbackCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Generate rollback script but do not execute it")

	rootCmd.AddCommand(rollbackCmd)
}
