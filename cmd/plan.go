package cmd

import (
	"os"

	"github.com/mikeknep/jeru/lib"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Preview the terraform plan following proposed changes to state",
	RunE: func(cmd *cobra.Command, args []string) error {

		statefile, err := os.Create(".jeru-state.tfstate")
		if err != nil {
			return err
		}
		defer os.Remove(statefile.Name())

		changes, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer changes.Close()

		extraArgs := []string{}
		if len(args) > 1 {
			extraArgs = args[1:]
		}

		runtime := lib.CreateLiveRuntimeEnvironment(extraArgs)

		return lib.Plan(
			runtime,
			changes,
			statefile,
		)
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
