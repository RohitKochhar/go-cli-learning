package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ExceptionStep extends the step type and implements another verison of the execute method
type ExceptionStep struct {
	Step
}

// NewExceptionStep() is a constructor for the ExceptionStep struct
func NewExceptionStep(name, exe, message, proj string, args []string) ExceptionStep {
	// Create a new ExceptionStep object
	s := ExceptionStep{}
	// Create a new Step object and assign it to the created ExceptionStep
	s.Step = NewStep(name, exe, message, proj, args)
	return s
}

// Execute() is the extended function that handles steps that don't return explicit errors
func (s ExceptionStep) Execute() (string, error) {
	// Craft the command to be executed
	cmd := exec.Command(s.exe, s.args...)
	// Create a buffer that will store the output of the command
	var out bytes.Buffer
	cmd.Stdout = &out
	// Specify the directory the command will run in
	cmd.Dir = s.proj
	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		return "", &StepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}
	// If the command executes without error, we have to check that the size of the output
	// buffer is non-zero, which would inficate that at least one file in the project doesn't match the format
	if out.Len() > 0 {
		return "", &StepErr{
			step:  s.name,
			msg:   fmt.Sprintf("invalid format: %s", out.String()),
			cause: nil,
		}
	}
	// If we got here, there was no problems with this step
	return s.message, nil
}
