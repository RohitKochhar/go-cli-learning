package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

// TestParseContent checks the bytes are equal between two the result of parse and golden file
func TestParseContent(t *testing.T) {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	result := parseContent(input)

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("Golden:\n%s\n", expected)
		t.Logf("Result:\n%s\n", result)
		t.Errorf("Result content does not match golden file")
	}
}

// TestRun checks the bytes between result and golden
func TestRun(t *testing.T) {
	// Create a mock stdout pipe for testing
	var mockStdOut bytes.Buffer

	if err := run(inputFile, &mockStdOut); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdOut.String())

	result, err := os.ReadFile(resultFile)
	if err != nil {
		fmt.Println(mockStdOut.String())
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("Golden:\n%s\n", expected)
		t.Logf("Result:\n%s\n", result)
		t.Errorf("Result content does not match golden file")
	}
	os.Remove(resultFile)
}

// TestShaHash checks the sha hash between result and golden
func TestShaHash(t *testing.T) {
	// Create a mock stdout pipe for testing
	var mockStdOut bytes.Buffer

	if err := run(inputFile, &mockStdOut); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdOut.String())
	// Store the result file into a variable for sha256 hash
	result, err := os.Open(resultFile)
	if err != nil {
		t.Fatal(err)
	}
	// Close the file
	defer result.Close()
	// Get the sha256 hash from the result variable
	resultHash := sha256.New()
	if _, err := io.Copy(resultHash, result); err != nil {
		log.Fatal(err)
	}
	// Store the golden file into a variable fo sha256 hash
	expected, err := os.Open(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	// Close the file
	defer expected.Close()
	// Get the sha256 from the golden variable
	goldenHash := sha256.New()
	if _, err := io.Copy(goldenHash, expected); err != nil {
		log.Fatal(err)
	}
	// Sum the hash components for the two individual vars
	resultHashString := fmt.Sprintf("%x", resultHash.Sum(nil))
	goldenHashString := fmt.Sprintf("%x", goldenHash.Sum(nil))
	// Ensure that they are equal
	if goldenHashString != resultHashString {
		t.Logf("Golden Hash:\n%s\n", goldenHashString)
		t.Logf("Result Hash:\n%s\n", resultHashString)
		t.Errorf("Result content does not match golden file")
	}
	os.Remove(resultFile)
}
