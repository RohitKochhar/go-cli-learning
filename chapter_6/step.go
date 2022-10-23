package main

import (
	"os/exec"
)

// Step represents the step of a pipeline
type Step struct {
	name    string   // step name
	exe     string   // executable name of the tool we want to execute
	args    []string // arguments to be passed to the executable
	message string   // output message if step was a success
	proj    string   // the target project to execute the task on
}

// NewStep() method is a constructor for Step type
func NewStep(name, exe, message, proj string, args []string) Step {
	return Step{
		name:    name,
		exe:     exe,
		message: message,
		args:    args,
		proj:    proj,
	}
}

// Execute() method runs the exe of a step and returns it's success
func (s Step) Execute() (string, error) {
	// Create a command using the Step parameters
	cmd := exec.Command(s.exe, s.args...)
	// Specify the directory we want to run the command in
	cmd.Dir = s.proj
	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		// If we had an error, create a StepErr value to return
		return "", &StepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}
	// Otherwise, return success
	return s.message, nil
}
