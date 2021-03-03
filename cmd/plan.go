package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"

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

		execute := func(w io.Writer, name string, args ...string) error {
			cmd := exec.Command(name, args...)
			cmd.Stdout = w
			// cmd.Stderr = w // do we care about Stderr?
			return cmd.Run()
		}

		additionalPlanArgs := []string{}
		if len(args) > 1 {
			additionalPlanArgs = args[1:]
		}

		return lib.PlanC(
			changes,
			statefile,
			os.Stdout,
			ioutil.Discard,
			getApprovalFromPrompt,
			execute,
			additionalPlanArgs,
		)
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
