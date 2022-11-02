package main

import "net/http"

// newMux maps incoming requests to the proper handlers based on the URL of the request
// this function takes the name of the file to save the to-do list to and returns a type
// that satisfies the http.Handler interface, a type that responds to an HTTP request.
func newMux(todoFile string) http.Handler {
	// Instantiate a new http.ServeMux which provides a multiplexer that satisfies the http.Handler interface
	// and allows us to mape routes to handler functions
	m := http.NewServeMux()
	// Attach the root route to the mutliplexer
	m.HandleFunc("/", rootHandler)
	return m
}

// replyTextContent wraps plain test into a http response
func replyTextContent(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain") // Specify the response content type
	w.WriteHeader(status)                        // Specify the HTTP response code
	w.Write([]byte(content))                     // Wrap the content into the response
}
