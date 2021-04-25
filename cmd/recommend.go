package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mikeknep/jeru/lib"
	"github.com/spf13/cobra"
)

type SurveyPrompt struct{}

func (p SurveyPrompt) Confirm(msg string) (bool, error) {
	confirmation := false
	prompt := &survey.Confirm{Message: msg}
	err := survey.AskOne(prompt, &confirmation)
	return confirmation, err
}

func (p SurveyPrompt) Select(options []string, msg string) (string, error) {
	selection := ""
	prompt := &survey.Select{
		Message: msg,
		Options: options,
	}
	err := survey.AskOne(prompt, &selection)
	return selection, err
}

var interactive bool

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

		additionalPlanArgs := []string{}
		if len(args) > 1 {
			additionalPlanArgs = args[1:]
		}

		var finder lib.RefactorFinder
		if interactive {
			finder = lib.GuidedRefactorFinder{Prompt: SurveyPrompt{}}
		} else {
			finder = lib.BestEffortRefactorFinder{}
		}

		return lib.Recommend(planfile, &jsonPlan, os.Stdout, ioutil.Discard, execute, finder, additionalPlanArgs)
	},
}

func init() {
	recommendCmd.Flags().BoolVar(&interactive, "i", false, "Interactive mode")
	rootCmd.AddCommand(recommendCmd)
}
