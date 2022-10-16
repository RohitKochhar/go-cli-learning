package main

import "errors"

// Define the error values as variables to be exported to other files.
// By convention, the variables should start with `Err`
var (
	ErrNotNumber        = errors.New("Data is not numeric")
	ErrInvalidColumn    = errors.New("Invalid Column Number")
	ErrNoFiles          = errors.New("No input files")
	ErrInvalidOperation = errors.New("Invalid Operation")
)
