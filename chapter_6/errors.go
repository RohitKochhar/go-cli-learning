package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("error: Validation failed")
)

// StepErr represents a class of errors associated with the CI steps
type StepErr struct {
	step  string // the name of the step that resulted in an error
	msg   string // the message describing the error condition
	cause error  // the underlying error that caused the step error
}

// Error() method creates a string from a StepErr object
func (s *StepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

// Is() method checks if two errors are equivalent
func (s *StepErr) Is(target error) bool {
	t, ok := target.(*StepErr)
	if !ok {
		return false
	}
	return t.step == s.step
}

// Unwrap() method returns the error stored in the cause field of an error
func (s *StepErr) Unwrap() error {
	return s.cause
}
