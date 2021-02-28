package lib

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

const performTheseActionsText = "\nDo you want to perform these actions? Only 'yes' will be accepted."
const introText = "Jeru has generated the following rollback commands:\n"

func Rollback(
	changes io.Reader,
	screen io.Writer,
	outfile io.Writer,
	getApproval func() (bool, error),
	execute func(string, ...string) error,
) error {

	// Generate rollback lines for the source changes
	rollbackLines := []string{}
	err := ConsumeByLine(changes, func(line string) {
		addRollbackLine(&rollbackLines, line)
	})
	if err != nil {
		return err
	}

	// Show the user what was generated
	fmt.Fprintln(screen, introText)
	for _, line := range rollbackLines {
		fmt.Fprintln(screen, "\t"+line)
	}

	// Write the generated lines to the outfile
	out := strings.Join(rollbackLines, "\n")
	_, err = outfile.Write([]byte(out))
	if err != nil {
		return err
	}

	// Exit if user does not approve changes
	fmt.Fprintln(screen, performTheseActionsText)
	approved, err := getApproval()
	if err != nil {
		return err
	}
	if !approved {
		return nil
	}

	// Execute the changes
	for _, line := range rollbackLines {
		err = execute("bash", "-c", line)
		if err != nil {
			return err
		}
	}
	return nil
}

func addRollbackLine(rollbackLines *[]string, srcLine string) {
	if isNoopLine(srcLine) {
		return
	}
	rollbackLine := generateRollbackLine(srcLine)
	*rollbackLines = append([]string{rollbackLine}, *rollbackLines...)
}

func generateRollbackLine(line string) string {
	im := regexp.MustCompile(`terraform import (\S+) \S+`)
	mv := regexp.MustCompile(`terraform state mv (\S+) (\S+)`)
	rm := regexp.MustCompile(`terraform state rm (\S+)`)

	switch {
	case im.FindStringIndex(line) != nil:
		return im.ReplaceAllString(line, fmt.Sprintf("terraform state rm $1"))
	case mv.FindStringIndex(line) != nil:
		return mv.ReplaceAllString(line, fmt.Sprintf("terraform state mv $2 $1"))
	case rm.FindStringIndex(line) != nil:
		return rm.ReplaceAllString(line, fmt.Sprintf("# terraform import $1 ___"))
	default:
		return fmt.Sprintf("# Could not generate rollback command for: %s", line)
	}
}

func isNoopLine(line string) bool {
	return isEmpty(line) || isShebang(line)
}

func isEmpty(line string) bool {
	matched, _ := regexp.MatchString(`^\s*$`, line)
	return matched
}

func isShebang(line string) bool {
	matched, _ := regexp.MatchString(`^#!`, line)
	return matched
}
