package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var changeScript string

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan some state changes",
	RunE: func(cmd *cobra.Command, args []string) error {

		// create some temp files
		tempdir, err := ioutil.TempDir(".", ".jeru")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempdir)

		statefile, err := ioutil.TempFile(tempdir, "state-")
		if err != nil {
			return err
		}
		defer statefile.Close()
		defer os.Remove(statefile.Name())

		changefile, err := ioutil.TempFile(tempdir, "change-")
		if err != nil {
			return err
		}
		defer changefile.Close()
		defer os.Remove(changefile.Name())

		// make a copy of the current state
		pullCommand := exec.Command("terraform", "state", "pull")
		pullCommand.Stdout = statefile
		pullCommand.Stderr = nil
		err = pullCommand.Run()
		if err != nil {
			return err
		}

		// make a copy of the state changes script that targets the *copy* of the current state
		input, err := ioutil.ReadFile(changeScript)
		if err != nil {
			return err
		}
		lines := strings.Split(string(input), "\n")
		re := regexp.MustCompile(`(terraform state (?:mv|rm))`)
		for i, line := range lines {
			lines[i] = re.ReplaceAllString(line, fmt.Sprintf("$1 -state=%s", statefile.Name()))
		}
		out := strings.Join(lines, "\n")
		err = ioutil.WriteFile(changefile.Name(), []byte(out), 0777)
		if err != nil {
			return err
		}
		os.Chmod(changefile.Name(), 0777)

		// execute that altered script
		changeCommand := exec.Command(changefile.Name())
		changeCommand.Stdout = nil
		changeCommand.Stderr = nil
		err = changeCommand.Run()
		if err != nil {
			return err
		}

		// // run plan against the modified state
		planCommand := exec.Command("terraform", "plan", "-state", statefile.Name())
		planCommand.Stdout = os.Stdout
		planCommand.Stderr = os.Stderr
		return planCommand.Run()
	},
}

func init() {
	planCmd.Flags().StringVar(&changeScript, "changes", "", "A script containing the terraform state mv|rm changes to make")
	planCmd.MarkFlagRequired("changes")

	rootCmd.AddCommand(planCmd)
}
