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
	RunE: func(cmd *cobra.Command, args []string) error {

		changes, err := os.Open(changeScript)
		if err != nil {
			return err
		}
		defer changes.Close()

		err = lib.Rollback(
			changes,
			dryRun,
			outfile,
			io.DisplayIntent,
			io.WriteAndRun,
		)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rollbackCmd.Flags().StringVar(&changeScript, "changes", "", "A script containing the terraform state mv|rm changes to make")
	rollbackCmd.MarkFlagRequired("changes")

	rollbackCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Generate rollback script but do not write or execute it. Supersedes --out.")
	rollbackCmd.Flags().StringVar(&outfile, "out", "", "Write the rollback commands to the given path. For current directory, prefix with './' (e.g. './rollback.sh'). Conflicts (fails) with --dry-run.")

	rootCmd.AddCommand(rollbackCmd)
}
