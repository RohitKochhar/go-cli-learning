package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

// main will parse the command-line arguments and call the run function
func main() {
	// Override flag defaults
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool Adapted by Rohit Singh, based of the Chapter 5 example from the Pragmatic Bookshelf\n",
			os.Args[0],
		)
		fmt.Fprintf(flag.CommandLine.Output(), "Adapted in October 2022\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	// Verify and parse args
	op := flag.String("op", "sum", "Operation to be executed")
	column := flag.Int("col", 1, "CSV column on which to execute information")
	// Parse flags
	flag.Parse()
	// Check that flags were succesful
	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run handles the main logic of the program
//
// Inputs:
//   - filenames ([]string): Slice of strings representing filenames to process
//   - op (string): Operation to execute
//   - column (int): Column to execute operation on
//   - Out (io.Writer): interface to print the results on
func run(filenames []string, op string, column int, out io.Writer) error {
	// Since we defined a type with the signature for stats func, we can create
	// a variable to store the function in
	var opFunc statsFunc
	// Validate the user provided parameters
	if err := validateFlags(filenames, op, column); err != nil {
		return err
	}
	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	}
	// consolidates the data that is extracted from the given column on each input file
	consolidate := make([]float64, 0)
	// Create channels to receive results or errors of operations
	resCh := make(chan []float64)
	errCh := make(chan error)
	// doneCh uses an empty struct since the channel doesn't need to send any data
	// this is the go equivalent of a None type, and ensures the program doesn't
	// allocate any memory for the channel
	doneCh := make(chan struct{})
	// Create a waitgroup to coordinate goroutine execution
	wg := sync.WaitGroup{}
	// Iterate through each input file
	for _, fname := range filenames {
		// Increment the number of items in the waitGroup
		wg.Add(1)
		// Create an anonymous function that takes fname
		go func(fname string) {
			// Once the anon func is done, trigger the wg.Done method
			defer wg.Done()
			// Open the file for reading
			f, err := os.Open(fname)
			if err != nil {
				errCh <- fmt.Errorf("Cannot open file: %w", err)
				return
			}
			// Parse the CSV into a slice of float64 numbers
			data, err := csv2float(f, column)
			if err != nil {
				errCh <- err
			}
			// Try to close the file
			if err := f.Close(); err != nil {
				errCh <- err
			}
			// By using a channel, we mitigate the risk of a race condition
			resCh <- data
		}(fname)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	// Consolidate the results coming from the resCh into the consolidate var
	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err
		}

	}
}

// validateFlags checks the user provided parameters
func validateFlags(filenames []string, op string, col int) error {
	// Check that a file was provided
	if len(filenames) == 0 {
		return ErrNoFiles
	}
	// Check that a valid column number was given
	if col < 1 {
		return ErrInvalidColumn
	}
	// Check that a valid operation was given
	switch op {
	case "sum":
		return nil
	case "avg":
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}
}
