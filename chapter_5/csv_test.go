package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"testing/iotest"
)

func TestOperations(t *testing.T) {
	// Create some sample data
	data := [][]float64{
		{10, 20, 15, 30, 45, 50, 100, 30},
		{5.5, 8, 2.2, 9.75, 8.45, 3, 2.5, 10.25, 4.75, 6.1, 7.67, 12.287, 5.47},
		{-10, -20},
		{102, 37, 44, 57, 67, 129},
	}
	// Using table driven testing
	testCases := []struct {
		name     string
		op       statsFunc
		expected []float64
	}{
		{
			name:     "TestSum",
			op:       sum,
			expected: []float64{300, 85.927, -30, 436},
		},
		{
			name:     "TestAvg",
			op:       avg,
			expected: []float64{37.5, 6.609769230769231, -15, 72.666666666666666},
		},
	}
	// Run the tests!
	for _, tc := range testCases {
		// Since we have a slice of expected values, iterate through them
		for k, expected := range tc.expected {
			// Since we run the test on each expected value, make the name more unique
			name := fmt.Sprintf("%sData%d", tc.name, k)
			// Run each test
			t.Run(name, func(t *testing.T) {
				res := tc.op(data[k])
				if res != expected {
					t.Errorf("Expected:\n\t%g\nGot:\n\t%g", tc.expected, res)
				}
			})
		}
	}
}

// TestCSV2Float runs tests on the csv2Float func
func TestCSV2Float(t *testing.T) {
	csvData := `IP Address,Requests,Response Time
		192.168.0.199,2056,236
		192.168.0.88,899,220
		192.168.0.199,3054,226
		192.168.0.100,4133,218
		192.168.0.199,950,238`
	// Using table driven testing
	testCases := []struct {
		name   string
		col    int
		exp    []float64
		expErr error
		r      io.Reader
	}{
		{
			name:   "TestColumn2Read",
			exp:    []float64{2056, 899, 3054, 4133, 950},
			col:    2,
			expErr: nil,
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "TestColumn3",
			exp:    []float64{236, 220, 226, 218, 238},
			col:    3,
			expErr: nil,
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "TestReadFail",
			col:    1,
			exp:    nil,
			expErr: iotest.ErrTimeout,
			r:      iotest.TimeoutReader(bytes.NewReader([]byte{0})),
		},
		{
			name:   "TestNaNFail",
			col:    1,
			exp:    nil,
			expErr: ErrNotNumber,
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "TestInvalidColumnFail",
			col:    4,
			exp:    nil,
			expErr: ErrInvalidColumn,
			r:      bytes.NewBufferString(csvData),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := csv2float(tc.r, tc.col)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error, got nil instead")
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error:\n\t%q\nInstead got:\n\t%q", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}
			for i, exp := range tc.exp {
				if res[i] != exp {
					t.Errorf("Expected:\n\t%g\nInstead got:\n\t%g", exp, res[i])
				}
			}
		})
	}
}
