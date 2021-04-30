package lib

import (
	"bytes"
	"fmt"
)

type RefactorFinder interface {
	Find(TfPlan) ([]Refactor, error)
}

type FindFlags struct {
	InteractiveMode bool
}

func Find(
	runtime RuntimeEnvironment,
	findFlags FindFlags,
	planfile NamedWriter,
) error {
	planArgs := []string{"plan", "-out", planfile.Name()}
	planArgs = append(planArgs, runtime.ExtraArgs...)
	spinner, err := runtime.StartSpinner("Running latest plan...")
	if err != nil {
		return err
	}
	err = runtime.Execute(runtime.Void, "terraform", planArgs...)
	if err != nil {
		return err
	}

	var jsonPlan bytes.Buffer

	err = runtime.Execute(&jsonPlan, "terraform", "show", "-json", planfile.Name())
	if err != nil {
		return err
	}

	tfPlan, err := NewTfPlan(&jsonPlan)
	if err != nil {
		return err
	}

	var finder RefactorFinder
	if findFlags.InteractiveMode {
		finder = GuidedRefactorFinder{Prompt: runtime.Prompt}
		spinner.Stop()
	} else {
		finder = BestEffortRefactorFinder{}
		spinner.UpdateText("Finding best refactor commands...")
	}

	refactors, err := finder.Find(tfPlan)
	if err != nil {
		return err
	}

	spinner.Stop()
	for _, refactor := range refactors {
		fmt.Fprintln(runtime.Screen, refactor.AsCommand())
	}

	return nil
}
