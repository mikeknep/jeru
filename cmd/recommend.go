package cmd

import (
	"bytes"
	"io"
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
		planfile, err := os.Create(".jeru-recommend.tfplan")
		if err != nil {
			return err
		}
		defer os.Remove(planfile.Name())

		var jsonPlan bytes.Buffer

		execute := func(w io.Writer, name string, args ...string) error {
			cmd := exec.Command(name, args...)
			cmd.Stdout = w
			cmd.Stderr = os.Stderr // do we care about Stderr?
			return cmd.Run()
		}

		return lib.Recommend(planfile, &jsonPlan, os.Stdout, ioutil.Discard, execute)
	},
}

func init() {
	rootCmd.AddCommand(recommendCmd)
}
