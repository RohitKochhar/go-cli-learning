package main

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	// Using table-driven testing
	testCases := []struct {
		name   string
		conf   config
		out    string
		expErr error
	}{
		{
			name: "Success",
			conf: config{
				proj:   "./testdata/tool",
				test:   true,
				format: true,
				readme: true,
			},
			out:    "go build: SUCCESS\ngo test: SUCCESS\ngofmt: SUCCESS\ngoreadme: SUCCESS\n",
			expErr: nil,
		},
		{
			name: "Fail",
			conf: config{
				proj:   "./testdata/toolErr",
				test:   true,
				format: false,
				readme: true,
			},
			out:    "",
			expErr: &StepErr{step: "go build"},
		},
		{
			name: "Chapter1",
			conf: config{
				proj:   "../chapter_1",
				test:   true,
				format: true,
				readme: true,
			},
			out:    "",
			expErr: &StepErr{step: "go fmt"},
		},
		{
			name: "Chapter2",
			conf: config{
				proj:   "../chapter_2",
				test:   false,
				format: false,
				readme: true,
			},
			out:    "go build: SUCCESS\ngoreadme: SUCCESS\n",
			expErr: nil,
		},
		{
			name: "Chapter3",
			conf: config{
				proj:   "../chapter_3",
				test:   false,
				format: true,
				readme: false,
			},
			out:    "go build: SUCCESS\ngofmt: SUCCESS\n",
			expErr: nil,
		},
		{
			name: "Chapter4",
			conf: config{
				proj:   "../chapter_4",
				test:   false,
				format: false,
				readme: true,
			},
			out:    "go build: SUCCESS\ngoreadme: SUCCESS\n",
			expErr: nil,
		},
		{
			name: "Chapter5",
			conf: config{
				proj:   "../chapter_5",
				test:   true,
				format: true,
				readme: true,
			},
			out:    "go build: SUCCESS\ngo test: SUCCESS\ngofmt: SUCCESS\ngoreadme: SUCCESS\n",
			expErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TestRun%s", tc.name), func(t *testing.T) {
			// Define a buffer to capture the output
			var out bytes.Buffer
			err := run(tc.conf, &out)
			// Check and handle if we wanted an error
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected an error of %q, got nil instead", tc.expErr)
				}
				if !errors.Is(err, tc.expErr) {
					t.Fatalf("Expected an error of %q, got %q instead", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %q", err)
			}
			if out.String() != tc.out {
				t.Fatalf("Expected %q, instead got %q", tc.out, out.String())
			}
		})
	}
}
