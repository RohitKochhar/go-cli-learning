package word_counter

import (
	"bytes"
	"fmt"
	"testing"
)

func ExampleCount() {
	fmt.Println(Count(bytes.NewBufferString("word1 word2 word3"), false, false))
	// Output: 3
}

// TestWordCounter is a unit test to ensure the program can count words
func TestWordCounter(t *testing.T) {
	testInput := bytes.NewBufferString("Line1\nLine2\nLine3\nTest test test test test\twow\twow\twow\ncheck it out\n19 120 1321 23939\nsomething@somethingelse.com")
	expected := 19
	result := Count(testInput, false, false)

	if result != expected {
		t.Errorf("Expected %d, got %d instead.\n", expected, result)
	}
}

func TestLineCount(t *testing.T) {
	testInput := bytes.NewBufferString("Line1\nLine2\nLine3\nTest test test test test\twow\twow\twow\ncheck it out\n19 120 1321 23939\nsomething@somethingelse.com")
	expected := 7
	result := Count(testInput, true, false)
	if result != expected {
		t.Errorf("Expected %d, got %d instead.\n", expected, result)
	}
}

func TestByteCount(t *testing.T) {
	testInput := bytes.NewBufferString("Line1\nLine2\nLine3\nTest test test test test\twow\twow\twow\ncheck it out\n19 120 1321 23939\nsomething@somethingelse.com")
	expected := 113
	result := Count(testInput, false, true)
	if result != expected {
		t.Errorf("Expected %d, got %d instead.\n", expected, result)
	}
}
