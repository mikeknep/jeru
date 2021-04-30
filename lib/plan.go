package lib

import (
	"fmt"
	"io"
	"regexp"
)

const commentOutBackendText = "Comment out your backend, then enter 'yes' to continue."
const reminderText = "Remember to restore your backend and re-initialize."

func Plan(
	runtime RuntimeEnvironment,
	changes io.Reader,
	localState NamedWriter,
) error {

	// pull terraform state into localState
	spinner, err := runtime.StartSpinner("Copying current state...")
	err = runtime.Execute(localState, "terraform", "state", "pull")
	if err != nil {
		return err
	}

	// wait for user to comment out backend
	spinner.Stop()
	fmt.Fprintln(runtime.Screen, commentOutBackendText)
	approved, err := runtime.GetApproval()
	if err != nil {
		return err
	}
	if !approved {
		return nil
	}

	// delete .terraform directory
	err = runtime.Execute(runtime.Void, "rm", "-rf", ".terraform")
	if err != nil {
		return err
	}

	spinner, err = runtime.StartSpinner("Re-initializing...")
	// reinitialize
	err = runtime.Execute(runtime.Void, "terraform", "init")
	if err != nil {
		return err
	}

	spinner.UpdateText("Applying changes to copied state...")
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
		err = runtime.Execute(runtime.Void, "bash", "-c", line)
		if err != nil {
			return err
		}
	}

	spinner.UpdateText("Planning...")
	// run plan against the modified state
	planArgs := []string{"plan", "-state", localState.Name()}
	planArgs = append(planArgs, runtime.ExtraArgs...)
	err = runtime.Execute(runtime.Screen, "terraform", planArgs...)
	if err != nil {
		return err
	}

	spinner.Stop()
	// delete .terraform directory again and remind user to re-init
	err = runtime.Execute(runtime.Void, "rm", "-rf", ".terraform")
	if err != nil {
		return err
	}
	fmt.Fprintln(runtime.Screen, reminderText)

	return nil
}
