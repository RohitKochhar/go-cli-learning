package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		conf     config
		expected string
	}{
		{
			name: "NoFilter",
			root: "testdata",
			conf: config{
				ext:  "",
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n",
		},
		{
			name: "FilterExtensionMatch",
			root: "testdata",
			conf: config{
				ext:  ".log",
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeMatch",
			root: "testdata",
			conf: config{
				ext:  ".log",
				size: 10,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeNoMatch",
			root: "testdata",
			conf: config{
				ext:  ".log",
				size: 20,
				list: true,
			},
			expected: "",
		},
		{
			name: "FilterExtensionNoMatch",
			root: "testdata",
			conf: config{
				ext:  ".gz",
				size: 0,
				list: true},
			expected: "",
		},
	}
	// Iterate through test cases and run em
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer to store the output of run
			var buffer bytes.Buffer
			// Try to run without any errors
			if err := run(tc.root, &buffer, tc.conf); err != nil {
				t.Fatal(err)
			}
			// Convert the buffered value to a string for compare
			res := buffer.String()
			// Check that we got what we expected
			if tc.expected != res {
				t.Errorf("Expected:\n\t%q\nGot:\n\t%q", tc.expected, res)
			}
		})
	}
}
