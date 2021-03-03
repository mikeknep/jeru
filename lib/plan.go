package lib

import (
	"fmt"
	"io"
	"regexp"
)

const commentOutBackendText = "Comment out your backend, then enter 'yes' to continue."
const reminderText = "Remember to restore your backend and re-initialize."

type localState interface {
	io.Writer
	Name() string
}

func PlanC(
	changes io.Reader,
	localState localState,
	screen io.Writer,
	void io.Writer,
	getApproval func() (bool, error),
	execute func(io.Writer, string, ...string) error,
	additionalPlanArgs []string,
) error {

	// pull terraform state into localState
	err := execute(localState, "terraform", "state", "pull")
	if err != nil {
		return err
	}

	// wait for user to comment out backend
	fmt.Fprintln(screen, commentOutBackendText)
	approved, err := getApproval()
	if err != nil {
		return err
	}
	if !approved {
		return nil
	}

	// delete .terraform directory
	err = execute(void, "rm", "-rf", ".terraform")
	if err != nil {
		return err
	}

	// reinitialize
	err = execute(screen, "terraform", "init")
	if err != nil {
		return err
	}

	// read changes, making new versions of each state command that target state
	re := regexp.MustCompile(`(terraform state (?:mv|rm))`)
	alteredLines := []string{}
	err = ConsumeByLine(changes, func(line string) {
		alt := re.ReplaceAllString(line, fmt.Sprintf("$1 -state=%s", localState.Name()))
		alteredLines = append(alteredLines, alt)
	})
	if err != nil {
		return err
	}

	// execute the changes
	for _, line := range alteredLines {
		err = execute(screen, "bash", "-c", line)
		if err != nil {
			return err
		}
	}

	// run plan against the modified state
	planArgs := []string{"plan", "-state", localState.Name()}
	planArgs = append(planArgs, additionalPlanArgs...)
	err = execute(screen, "terraform", planArgs...)
	if err != nil {
		return err
	}

	// delete .terraform directory again and remind user to re-init
	err = execute(void, "rm", "-rf", ".terraform")
	if err != nil {
		return err
	}
	fmt.Fprintln(screen, reminderText)

	return nil
}
