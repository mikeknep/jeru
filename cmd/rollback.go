package cmd

import (
	"os"

	"github.com/mikeknep/jeru/io"
	"github.com/mikeknep/jeru/lib"
	"github.com/spf13/cobra"
)

var dryRun bool
var outfile string

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

		rollbackFilename := lib.OrDefault(outfile, "./.jeru-rollback.sh")
		rollbackScript, file, err := io.CreateScript(rollbackFilename)
		if err != nil {
			return err
		}
		defer file.Close()
		if outfile == "" {
			defer os.Remove(rollbackScript.Name)
		}

		var run func(_ string) error
		if dryRun {
			run = lib.DryRun
		} else {
			run = io.Run
		}

		return lib.Rollback(changes, rollbackScript, io.DisplayIntent, run)
	},
}

func init() {
	rollbackCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Generate rollback script but do not write or execute it.")
	rollbackCmd.Flags().StringVar(&outfile, "out", "", "Write the rollback commands to the given path. For current directory, prefix with './' (e.g. './rollback.sh').")

	rootCmd.AddCommand(rollbackCmd)
}
