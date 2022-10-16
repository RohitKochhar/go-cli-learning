package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// TestRun tests main logic calls
func TestRun(t *testing.T) {
	// Using table driven testing
	testCases := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{
			name:   "TestRunAvgSingleFile",
			col:    3,
			op:     "avg",
			exp:    "227.6\n",
			files:  []string{"./testdata/example.csv"},
			expErr: nil,
		},
		{
			name:   "TestRunAvgMultiFiles",
			col:    3,
			op:     "avg",
			exp:    "233.75\n",
			files:  []string{"./testdata/example.csv", "./testdata/example2.csv"},
			expErr: nil,
		},
		{
			name:   "TestReadFail",
			col:    2,
			op:     "avg",
			exp:    "",
			files:  []string{"./testdata.example.csv", "./testdata/nonexistentfile.csv"},
			expErr: os.ErrNotExist,
		},
		{
			name:   "TestInvalidColumnFail",
			col:    0,
			op:     "avg",
			exp:    "",
			files:  []string{"./testdata/example.csv", "./testdata/example2.csv"},
			expErr: ErrInvalidColumn,
		},
		{
			name:   "TestNoFilesFail",
			col:    1,
			op:     "avg",
			exp:    "",
			files:  []string{},
			expErr: ErrNoFiles,
		},
		{
			name:   "TestInvalidOperationFail",
			col:    1,
			op:     "minus",
			exp:    "",
			files:  []string{"./testdata/example.csv", "./testdata/example2.csv"},
			expErr: ErrInvalidOperation,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res bytes.Buffer
			err := run(tc.files, tc.op, tc.col, &res)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error, got nil instead")
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error %q, got %q isntead", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}
			if res.String() != tc.exp {
				t.Errorf("Expected %q, got %q instead", tc.exp, &res)
			}
		})
	}
}

// Benchmark takes a single input point to a testing.B type
// This provides the methods and fields used for controlling benchmarking
func BenchmarkRun(b *testing.B) {
	// filenames contains all the files used for the benchmark
	filenames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}
	// Before starting benchmarking, reset benchmark clock
	b.ResetTimer()
	// Run benchmark
	for i := 0; i < b.N; i++ {
		if err := run(filenames, "avg", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}
