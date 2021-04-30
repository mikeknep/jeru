package cmd

import (
	"os"

	"github.com/mikeknep/jeru/lib"
	"github.com/spf13/cobra"
)

var interactive bool

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Analyze plan output and find state mv changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		planfile, err := os.Create(".jeru-find.tfplan")
		if err != nil {
			return err
		}
		defer os.Remove(planfile.Name())

		extraArgs := []string{}
		if len(args) > 0 {
			extraArgs = args[0:]
		}

		runtime := lib.CreateLiveRuntimeEnvironment(extraArgs)

		findFlags := lib.FindFlags{
			InteractiveMode: interactive,
		}

		return lib.Find(
			runtime,
			findFlags,
			planfile,
		)
	},
}

func init() {
	findCmd.Flags().BoolVar(&interactive, "i", false, "Interactive mode")
	rootCmd.AddCommand(findCmd)
}
