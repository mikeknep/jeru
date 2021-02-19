package io

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
)

type FileReader struct {
	scanner *bufio.Scanner
}

func ReadFile(file string) (*FileReader, error) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fileReader := create(bufio.NewScanner(strings.NewReader(string(input))))
	return &fileReader, nil
}

func (fr *FileReader) EachLine(f func(string)) error {
	for fr.scanner.Scan() {
		f(fr.scanner.Text())
	}
	return fr.scanner.Err()
}

func create(scanner *bufio.Scanner) FileReader {
	return FileReader{scanner}
}

func DisplayIntent(lines []string, preliminaryText string) {
	fmt.Println(preliminaryText + "\n")
	for _, line := range lines {
		fmt.Println("\t" + line)
	}
}
