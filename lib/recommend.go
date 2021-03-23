package lib

import (
	"fmt"
	"io"
)

func Recommend(
	planfile NamedWriter,
	jsonPlan io.ReadWriter,
	screen io.Writer,
	void io.Writer,
	execute func(io.Writer, string, ...string) error,
) error {
	err := execute(void, "terraform", "plan", "-out", planfile.Name())
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

	for _, refactor := range tfPlan.PossibleRefactors() {
		fmt.Fprintln(screen, refactor.AsCommand())
	}

	return nil
}
