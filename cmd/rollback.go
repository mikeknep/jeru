package cmd

import (
	"github.com/mikeknep/jeru/io"
	"github.com/mikeknep/jeru/lib"
	"github.com/spf13/cobra"
)

var dryRun bool

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Revert a series of state changes",
	RunE: func(cmd *cobra.Command, args []string) error {

		rollbackLines := []string{}
		if err := io.ConsumeFileByLine(changeScript, func(line string) {
			lib.AddRollbackLine(&rollbackLines, line)
		}); err != nil {
			return err
		}

		io.DisplayIntent(rollbackLines, "Jeru has generated the following rollback commands:")

		if dryRun {
			return nil
		}

		if err := io.WriteAndRun("./.jeru-rollback.sh", rollbackLines); err != nil {
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
