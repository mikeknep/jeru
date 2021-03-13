package lib

import (
	"io/ioutil"
	"strings"
)

var (
	im = "terraform import resource.a id"
	mv = "terraform state mv module.a module.b"
	rm = "terraform state rm resource.a"

	approve   = func() (bool, error) { return true, nil }
	unapprove = func() (bool, error) { return false, nil }

	void = ioutil.Discard
)

type NamedStringbuilder struct {
	name   string
	writer *strings.Builder
}

func CreateNamedStringbuilder(name string) *NamedStringbuilder {
	var builder strings.Builder
	return &NamedStringbuilder{name: name, writer: &builder}
}

func (nsb *NamedStringbuilder) Name() string {
	return nsb.name
}

func (nsb *NamedStringbuilder) Write(x []byte) (int, error) {
	return nsb.writer.Write(x)
}

func (nsb *NamedStringbuilder) String() string {
	return nsb.writer.String()
}
