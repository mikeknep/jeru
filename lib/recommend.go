package lib

import (
	"fmt"
	"io"
)

type RefactorFinder func(TfPlan) ([]Refactor, error)

func Recommend(
	planfile NamedWriter,
	jsonPlan io.ReadWriter,
	screen io.Writer,
	void io.Writer,
	execute func(io.Writer, string, ...string) error,
	getRefactors RefactorFinder,
	additionalPlanArgs []string,
) error {
	planArgs := []string{"plan", "-out", planfile.Name()}
	planArgs = append(planArgs, additionalPlanArgs...)
	err := execute(void, "terraform", planArgs...)
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

	refactors, err := getRefactors(tfPlan)
	if err != nil {
		return err
	}
	for _, refactor := range refactors {
		fmt.Fprintln(screen, refactor.AsCommand())
	}

	return nil
}
