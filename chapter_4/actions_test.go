package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	// Create anonymous slice of struct with test case definitions
	testCases := []struct {
		// Here we define the properties of our tests
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		// Here we define each test
		{"FilterNoExtension", "testdata/dir.log", "", 0, false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterExtensionSizeNoMatch", "testdata/dir.log", ".log", 20, true},
	}
	// Iterate over the testCases object and run the test
	for _, tc := range testCases {
		// Run a test
		t.Run(tc.name, func(t *testing.T) {
			// Get the file's attributes using Stat
			info, err := os.Stat(tc.file)
			if err != nil {
				// Check for errors
				t.Fatal(err)
			}

			// Run filterOut for the specified inputs
			filterResult := filterOut(tc.file, tc.ext, tc.minSize, info)

			// Check that the result was what we expected
			if filterResult != tc.expected {
				t.Errorf("Expected:\n\t%t\nGot:\n\t%t", tc.expected, filterResult)
			}
		})
	}
}
