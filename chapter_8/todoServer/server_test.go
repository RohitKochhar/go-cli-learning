package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"todo"
)

// setupAPI is a helper function that creates a test server
func setupAPI(t *testing.T) (string, func()) {
	t.Helper() // Marks the function as a test helper
	// Create a temp file and fill it with some placeholders
	tempTodoFile, err := os.CreateTemp("", "todotest")
	ts := httptest.NewServer(newMux(tempTodoFile.Name())) // Creates a new test server with our custom mux function
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < 3; i++ {
		var body bytes.Buffer
		taskName := fmt.Sprintf("Task number %d", i)
		item := struct {
			Task string `json:"task"`
		}{
			Task: taskName,
		}
		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}

		r, err := http.Post(ts.URL+"/todo/", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to add initial items: Status: %d", r.StatusCode)
		}
	}

	// Return the created server's url and a cleanup function that closes the server when executed
	return ts.URL, func() {
		ts.Close()
		os.Remove(tempTodoFile.Name())
	}
}

// TestGet tests the HTTP GET method on the server's root
func TestGet(t *testing.T) {
	// Using table-driven testing to allow for tests on different paths
	testCases := []struct {
		name       string // Name of the test to be executed
		path       string // Server URL path to test
		expCode    int    // Expected return code from the server
		expItems   int    // Expected number of items returned by query
		expContent string // Expected body content of the response
	}{
		{
			// GetRoot checks that requesting the root is a success
			name:       "GetRoot",
			path:       "/",
			expCode:    http.StatusOK,
			expContent: "There's an API here\n",
		},
		{
			// GetAll checks that requesting all tasks is a success
			name:       "GetAll",
			path:       "/todo",
			expCode:    http.StatusOK,
			expItems:   2,
			expContent: "Task number 1",
		},
		{
			// GetOne checks that requesting a single item is a success
			name:       "GetOne",
			path:       "/todo/1",
			expCode:    http.StatusOK,
			expItems:   1,
			expContent: "Task number 1",
		},
		{
			// NotFound checks that requesting a non-existent URL path returns 404
			name:    "NotFound",
			path:    "/bad/path",
			expCode: http.StatusNotFound,
		},
	}
	// Create a new test server using the helper function
	url, cleanup := setupAPI(t)
	defer cleanup()
	// Loop through test cases and execute each test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				resp struct {
					Results      todo.List `json:"results"`
					Date         int64     `json:"date"`
					TotalResults int       `json:"total_results"`
				}
				body []byte
				err  error
			)
			r, err := http.Get(url + tc.path)
			if err != nil {
				t.Fatalf("Unexpected error while getting path: %q", err.Error())
			}
			defer r.Body.Close()
			// Validate the returned status code
			if r.StatusCode != tc.expCode {
				t.Fatalf("Expected %q, got %q", http.StatusText(tc.expCode), http.StatusText(r.StatusCode))
			}
			// Use a switch statement on the content-type
			switch {
			case r.Header.Get("Content-Type") == "application/json":
				if err = json.NewDecoder(r.Body).Decode(&resp); err != nil {
					t.Error(err)
				}
				if resp.TotalResults != tc.expItems {
					fmt.Printf("%+v", len(resp.Results))
					t.Errorf("Expected %d items, got %d.", tc.expItems, resp.TotalResults)
				}
				if resp.Results[0].Task != tc.expContent {
					t.Errorf("Expected %q, got %q", tc.expContent, resp.Results[0].Task)
				}
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(r.Body); err != nil {
					t.Error(err)
				}
				if !strings.Contains(string(body), tc.expContent) {
					t.Errorf("Expected %q, got %q", tc.expContent, string(body))
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", r.Header.Get("Content-Type"))
			}
		})
	}
}
