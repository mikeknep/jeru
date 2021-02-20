package io

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
)

func ConsumeFileByLine(file string, f func(string)) error {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(input)))
	for scanner.Scan() {
		f(scanner.Text())
	}
	return scanner.Err()
}

func DisplayIntent(lines []string, preliminaryText string) {
	fmt.Println(preliminaryText + "\n")
	for _, line := range lines {
		fmt.Println("\t" + line)
	}
}
