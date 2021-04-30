package lib

import "github.com/pterm/pterm"

type Spinner interface {
	UpdateText(string)
	Success(...interface{})
	Stop() error
}

type StartSpinner func(...interface{}) (Spinner, error)

var StartPtermSpinner StartSpinner = func(x ...interface{}) (Spinner, error) {
	return pterm.DefaultSpinner.WithRemoveWhenDone(true).Start(x)
}

var StartSilentSpinner StartSpinner = func(_ ...interface{}) (Spinner, error) {
	s := SilentSpinner{}
	return &s, nil
}

type SilentSpinner struct{}

func (s *SilentSpinner) UpdateText(_ string)      {}
func (s *SilentSpinner) Success(_ ...interface{}) {}
func (s *SilentSpinner) Stop() error              { return nil }
