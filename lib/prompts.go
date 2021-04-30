package lib

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
)

// survey

type Prompt interface {
	Confirm(string) (bool, error)
	Select([]string, string) (string, error)
}

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

// promptui

type GetApproval func() (bool, error)

var GetApprovalFromPrompt GetApproval = func() (bool, error) {
	prompt := promptui.Prompt{
		Label: "\tEnter a value",
	}
	input, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return input == "yes", nil
}

var AutoApprove GetApproval = func() (bool, error) {
	return true, nil
}
