package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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
		planCommand := lib.Terraform([]string{"plan", "-out", planfile}, nil)
		if err = planCommand.Run(); err != nil {
			return err
		}

		// convert the planfile to json and decode into a TfPlan
		reader, writer := io.Pipe()
		showCommand := lib.Terraform([]string{"show", "-json", planfile}, writer)
		if err = showCommand.Start(); err != nil {
			return err
		}
		var plan lib.TfPlan
		if err = json.NewDecoder(reader).Decode(&plan); err != nil {
			return err
		}
		if err = showCommand.Wait(); err != nil {
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
