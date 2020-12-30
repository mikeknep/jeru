package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mikeknep/jeru/lib"
	"github.com/spf13/cobra"
)

var recommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Analyze plan output and recommend state mv changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		tempdir, err := ioutil.TempDir(".", ".jeru")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempdir)

		// the name of the binary file to which terraform writes its plan
		planfile := fmt.Sprintf("%s/planfile", tempdir)

		// generate a planfile
		planCommand := exec.Command("terraform", "plan", "-out", planfile)
		planCommand.Stdout = nil
		planCommand.Stderr = nil
		err = planCommand.Run()
		if err != nil {
			return err
		}

		// convert the planfile to json and decode into a Plan
		var plan lib.Plan
		toJsonCommand := exec.Command("terraform", "show", "-json", planfile)
		stdout, err := toJsonCommand.StdoutPipe()
		if err != nil {
			return err
		}
		if err := toJsonCommand.Start(); err != nil {
			return err
		}
		if err := json.NewDecoder(stdout).Decode(&plan); err != nil {
			return err
		}
		if err := toJsonCommand.Wait(); err != nil {
			return err
		}

		for _, pr := range plan.PossibleRefactors() {
			fmt.Println(pr.AsCommand())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(recommendCmd)
}
