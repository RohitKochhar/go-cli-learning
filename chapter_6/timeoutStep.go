package main

import (
	"context"
	"os/exec"
	"time"
)

// TimeoutStep extends the Step struct and adds a timeout limit
type TimeoutStep struct {
	Step
	timeout time.Duration
}

// NewTimeoutStep is a constructor for the TimeoutStep struct
func NewTimeoutStep(name, exe, message, proj string, args []string, timeout time.Duration) TimeoutStep {
	// Create new timeout step
	s := TimeoutStep{}
	// Create new step from given parameters and assign to timeout step
	s.Step = NewStep(name, exe, message, proj, args)
	// Add the timeout
	s.timeout = timeout
	// If no timeout was provided, add a default
	if s.timeout == 0 {
		s.timeout = 30 * time.Second
	}
	return s
}

// Execute() is the extended version of the original func that returns error if it takes too long
func (s TimeoutStep) Execute() (string, error) {
	// use a context to carry the timeout value
	// context.WithTimeout returns two values
	// 	the context, and the cancellation function
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	// Run the cancellation function when the execute step is complete
	defer cancel()
	// Use CommandContext to create a command that uses a created context to kill the executing
	// process in the case the context becomes done before the command completes
	cmd := exec.CommandContext(ctx, s.exe, s.args...)
	// Specify the directory to run the command in
	cmd.Dir = s.proj
	// Execute the command and check for errors
	if err := cmd.Run(); err != nil {
		// If the context deadline is exceeeded, return a specific error
		if ctx.Err() == context.DeadlineExceeded {
			return "", &StepErr{
				step:  s.name,
				msg:   "failed timeout",
				cause: context.DeadlineExceeded,
			}
		}
		// If it is a different error, return a more generic message
		return "", &StepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}
	// If no errors, return success message
	return s.message, nil
}
