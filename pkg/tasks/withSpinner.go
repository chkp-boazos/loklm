package tasks

import (
	"time"

	"github.com/briandowns/spinner"
)

func WithSpinner(description string) TaskDecorator {
	return func(task Task) Task {
		return func() error {
			s := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
			s.Suffix = " " + description
			s.Start()
			err := task()
			s.Stop()
			return err
		}
	}
}
