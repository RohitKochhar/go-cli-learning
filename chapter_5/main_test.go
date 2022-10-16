package main

import (
	"bytes"
	"errors"
	"os"
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
