// Package to count words, lines and bytes of an input
//
// Available flags:
//
// 	`-l`, instructs program to count lines in input
//
// 	`-b`, instructs program to count bytes in input
//
// When no flags are provided, the program counts words by default
package word_counter

import (
	"bufio" // Used to read text
	"flag"  // Used to create CL flags
	"fmt"   // Used to print text
	"io"    // Used for io.Reader interfact
	"os"    // Used to access OS resources
)

// Main parses the flags from the CLI and calls count accordingly
func Main() {
	// Create a bool flag to specify line count vs words
	lineFlag := flag.Bool("l", false, "Count lines")
	byteFlag := flag.Bool("b", false, "Count bytes")
	// Parse the command for flags
	flag.Parse()
	// Call the count function to get number of words
	//	received through stdin and print it
	fmt.Println(Count(os.Stdin, *lineFlag, *byteFlag))
}

// count parses the input and counts the metric based on the provided flags
func Count(r io.Reader, isLineCount bool, isByteCount bool) int {
	// Create a scanner object to read text from an input
	scanner := bufio.NewScanner(r)
	if isLineCount {
		// Define the scanner to split at words rather than the
		//	default of lines
		scanner.Split(bufio.ScanLines)
	} else if isByteCount {
		scanner.Split(bufio.ScanBytes)
	} else {
		scanner.Split(bufio.ScanWords)
	}
	// Define a counter to store counted words
	counter := 0
	// For every word scanned, inc counter
	for scanner.Scan() {
		counter++
	}
	// Return the total number of counted words
	return counter
}
