package lib

import (
	"bufio"
	"io"
)

type NamedWriter interface {
	io.Writer
	Name() string
}

func ConsumeByLine(reader io.Reader, f func(string)) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		f(scanner.Text())
	}
	return scanner.Err()
}
