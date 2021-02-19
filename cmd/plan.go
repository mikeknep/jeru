package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/mikeknep/jeru/io"
	"github.com/mikeknep/jeru/lib"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Preview the terraform plan following proposed changes to state",
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
		pullCommand := lib.Terraform([]string{"state", "pull"}, statefile)
		if err = pullCommand.Run(); err != nil {
			return err
		}

		// ensure the user comments out their configured backend, and reinitialize locally
		fmt.Println("Comment out your backend, then enter 'yes' to continue.")
		prompt := promptui.Prompt{
			Label: "  Ready to proceed?",
		}
		response, err := prompt.Run()
		if err != nil {
			return err
		}
		if response != "yes" {
			return nil
		}

		deleteDotTerraformDirCommand := exec.Command("rm", "-rf", ".terraform")
		deleteDotTerraformDirCommand.Stdout = nil
		deleteDotTerraformDirCommand.Stderr = nil
		err = deleteDotTerraformDirCommand.Run()
		if err != nil {
			return err
		}

		initCommand := lib.Terraform([]string{"init"}, os.Stdout)
		initCommand.Run()

		// make a copy of the state changes script that targets the *copy* of the current state
		reader, err := io.ReadFile(changeScript)
		if err != nil {
			return err
		}
		re := regexp.MustCompile(`(terraform state (?:mv|rm))`)
		alteredLines := []string{}
		err = reader.EachLine(func(line string) {
			alt := re.ReplaceAllString(line, fmt.Sprintf("$1 -state=%s", statefile.Name()))
			alteredLines = append(alteredLines, alt)
		})
		out := strings.Join(alteredLines, "\n")
		if err = ioutil.WriteFile(changefile.Name(), []byte(out), 0777); err != nil {
			return err
		}
		changefile.Chmod(0777)

		// execute that altered script
		changeCommand := exec.Command(changefile.Name())
		changeCommand.Stdout = nil
		changeCommand.Stderr = nil
		err = changeCommand.Run()
		if err != nil {
			return err
		}

		// // run plan against the modified state
		planArgs := []string{"plan", "-state", statefile.Name()}
		planArgs = append(planArgs, args...)
		planCommand := lib.Terraform(planArgs, os.Stdout)
		err = planCommand.Run()
		if err != nil {
			return err
		}

		cleanupCommand := exec.Command("rm", "-rf", ".terraform")
		cleanupCommand.Stdout = nil
		cleanupCommand.Stderr = nil
		err = cleanupCommand.Run()
		if err != nil {
			return err
		}

		fmt.Println("Remember to restore your backend and re-initialize")
		return nil
	},
}

func init() {
	planCmd.Flags().StringVar(&changeScript, "changes", "", "A script containing the terraform state mv|rm changes to make")
	planCmd.MarkFlagRequired("changes")

	rootCmd.AddCommand(planCmd)
}
