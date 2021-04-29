package lib

import (
	"fmt"
	"io"
)

type RefactorFinder interface {
	Find(TfPlan) ([]Refactor, error)
}

func Find(
	planfile NamedWriter,
	jsonPlan io.ReadWriter,
	screen io.Writer,
	void io.Writer,
	startSpinner StartSpinner,
	execute func(io.Writer, string, ...string) error,
	finder RefactorFinder,
	additionalPlanArgs []string,
) error {
	planArgs := []string{"plan", "-out", planfile.Name()}
	planArgs = append(planArgs, additionalPlanArgs...)
	spinner, err := startSpinner("Running latest plan...")
	if err != nil {
		return err
	}
	err = execute(void, "terraform", planArgs...)
	if err != nil {
		return err
	}

	err = execute(jsonPlan, "terraform", "show", "-json", planfile.Name())
	if err != nil {
		return err
	}

	tfPlan, err := NewTfPlan(jsonPlan)
	if err != nil {
		return err
	}

	spinner.UpdateText("Finding best refactor commands...")
	refactors, err := finder.Find(tfPlan)
	if err != nil {
		return err
	}

	spinner.Success("Complete!")
	for _, refactor := range refactors {
		fmt.Fprintln(screen, refactor.AsCommand())
	}

	return nil
}
