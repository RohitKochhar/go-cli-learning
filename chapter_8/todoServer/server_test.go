package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setupAPI is a helper function that creates a test server
func setupAPI(t *testing.T) (string, func()) {
	t.Helper()                           // Marks the function as a test helper
	ts := httptest.NewServer(newMux("")) // Creates a new test server with our custom mux function
	// Return the created server's url and a cleanup function that closes the server when executed
	return ts.URL, func() {
		ts.Close()
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
