package lib

import (
	"io"
	"regexp"
)

func Rollback(
	changes io.Reader,
	dryRun bool,
	outfile string,
	presenter func(lines []string, msg string),
	writerRunner func(filename string, lines []string, persist bool) error,
) error {

	rollbackLines := []string{}
	err := ConsumeByLine(changes, func(line string) {
		addRollbackLine(&rollbackLines, line)
	})
	if err != nil {
		return err
	}

	presenter(rollbackLines, "Jeru has generated the following rollback commands:")

	if dryRun {
		return nil
	}

	filename := OrDefault(outfile, "./.jeru-rollback.sh")
	persist := outfile != ""
	err = writerRunner(filename, rollbackLines, persist)
	if err != nil {
		return err
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
