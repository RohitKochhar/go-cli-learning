package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// statsFunc represents a class of functions with the signature of (data []float64)
// We can use this whenever we need a new calc function,
// this will make the code more concise and easier to test
type statsFunc func(data []float64) float64

// sum adds the values of a column
func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

// avg takes the avg of the values in a column
func avg(data []float64) float64 {
	sum := sum(data)
	return sum / float64(len(data))
}

// csv2float parses the contents of a csv file into a slice of fpoint numbers
// By handing the function an io.reader as a parameter, we can call the function
// by passing a buffer that contains test data instead of using a file
func csv2float(r io.Reader, column int) ([]float64, error) {
	// Create a csv reader object using the io.Reader
	cr := csv.NewReader(r)
	// Adjust the column index so that users can start counting at 1 rather than 0
	column--
	// Create a var to hold the results of the csv converstion
	var data []float64
	// Benchmarking improvement 1 -----------------------------
	// --------- Old: (5043294)
	// // Use cr.ReadAll to read in the entire CSV data into a variable
	// allData, err := cr.ReadAll()
	// // Check for errors
	// if err != nil {
	// 	// Wrap the original err message with additional info
	// 	return nil, fmt.Errorf("Cannot read data from file: %w", err)
	// }

	// // Loop through allData variable
	// for i, row := range allData {
	// --------- New: (2529282)
	// Run an infinite loop to read CSV until EOF is reached
	cr.ReuseRecord = true
	for i := 0; ; i++ {
		row, err := cr.Read()
		if err == io.EOF {
			// If EOF, terminate loop
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Cannot read data from file: %w", err)
		}
		if i == 0 {
			// 0th column contains metadata, skip
			continue
		}
		if len(row) <= column {
			// The length of the row variable is the number of cells in it
			// If the number cells in a row is greater than the columns in the table
			// we need to raise an error
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}
		// Try to convert csv data into a float using strConv library
		v, err := strconv.ParseFloat(row[column], 64) // 64 bit float
		if err != nil {
			return nil, fmt.Errorf("%w, %s", ErrNotNumber, err)
		}
		data = append(data, v)
	}
	// If no errors occured, we can return the converted data and a nil error
	return data, nil
}
