package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
	"github.com/mikeknep/jeru/lib"
	"github.com/spf13/cobra"
)

var autoApprove bool
var dryRun bool
var out string

const getApprovalText = "\nDo you want to perform these actions? Only 'yes' will be accepted."
const labelText = "\tEnter a value"

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Revert a series of state changes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		changes, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer changes.Close()

		var outfile = ioutil.Discard
		if out != "" {
			file, err := os.Create(out)
			if err != nil {
				return err
			}
			defer file.Close()
			outfile = file
		}

		getApproval := func() (bool, error) {
			if autoApprove {
				return true, nil
			}
			return getApprovalFromPrompt()
		}

		execute := func(name string, arg ...string) error {
			if dryRun {
				return nil
			}
			return exec.Command(name, arg...).Run()
		}

		return lib.Rollback(changes, os.Stdout, outfile, getApproval, execute)
	},
}

func getApprovalFromPrompt() (bool, error) {
	fmt.Println(getApprovalText)
	prompt := promptui.Prompt{
		Label: labelText,
	}
	input, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return input == "yes", nil
}

func init() {
	rollbackCmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "Apply the rollback script without asking for approval.")
	rollbackCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Generate rollback script but do not write or execute it.")
	rollbackCmd.Flags().StringVar(&out, "out", "", "Write the rollback commands to the given path. For current directory, prefix with './' (e.g. './rollback.sh').")

	rootCmd.AddCommand(rollbackCmd)
}
