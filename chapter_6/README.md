# .

package main

## Variables

```golang
var (
    ErrValidation = errors.New("error: Validation failed")
)
```

## Types

### type [ExceptionStep](/exceptionStep.go#L10)

`type ExceptionStep struct { ... }`

ExceptionStep extends the step type and implements another verison of the execute method

#### func [NewExceptionStep](/exceptionStep.go#L15)

`func NewExceptionStep(name, exe, message, proj string, args []string) ExceptionStep`

NewExceptionStep() is a constructor for the ExceptionStep struct

#### func (ExceptionStep) [Execute](/exceptionStep.go#L24)

`func (s ExceptionStep) Execute() (string, error)`

Execute() is the extended function that handles steps that don't return explicit errors

### type [Step](/step.go#L9)

`type Step struct { ... }`

Step represents the step of a pipeline

#### func [NewPipedStep](/step.go#L29)

`func NewPipedStep(name, exe1, exe2, message, proj string, args1 []string, args2 []string) Step`

NewPipedStep() is a constructor for the Step struct that takes two serially piped commands

#### func [NewStep](/step.go#L18)

`func NewStep(name, exe, message, proj string, args []string) Step`

NewStep() method is a constructor for Step type

#### func (Step) [Execute](/step.go#L50)

`func (s Step) Execute() (string, error)`

Execute() method runs the exe of a step and returns it's success

### type [StepErr](/errors.go#L13)

`type StepErr struct { ... }`

StepErr represents a class of errors associated with the CI steps

#### func (*StepErr) [Error](/errors.go#L20)

`func (s *StepErr) Error() string`

Error() method creates a string from a StepErr object

#### func (*StepErr) [Is](/errors.go#L25)

`func (s *StepErr) Is(target error) bool`

Is() method checks if two errors are equivalent

#### func (*StepErr) [Unwrap](/errors.go#L34)

`func (s *StepErr) Unwrap() error`

Unwrap() method returns the error stored in the cause field of an error

### type [TimeoutStep](/timeoutStep.go#L10)

`type TimeoutStep struct { ... }`

TimeoutStep extends the Step struct and adds a timeout limit

#### func [NewTimeoutStep](/timeoutStep.go#L16)

`func NewTimeoutStep(name, exe, message, proj string, args []string, timeout time.Duration) TimeoutStep`

NewTimeoutStep is a constructor for the TimeoutStep struct

#### func (TimeoutStep) [Execute](/timeoutStep.go#L31)

`func (s TimeoutStep) Execute() (string, error)`

Execute() is the extended version of the original func that returns error if it takes too long

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
