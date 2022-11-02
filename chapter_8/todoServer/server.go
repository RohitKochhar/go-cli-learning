package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// newMux maps incoming requests to the proper handlers based on the URL of the request
// this function takes the name of the file to save the to-do list to and returns a type
// that satisfies the http.Handler interface, a type that responds to an HTTP request.
func newMux(todoFile string) http.Handler {
	// Instantiate a new http.ServeMux which provides a multiplexer that satisfies the http.Handler interface
	// and allows us to mape routes to handler functions
	m := http.NewServeMux()
	mu := &sync.Mutex{}
	// Attach the root route to the mutliplexer
	m.HandleFunc("/", rootHandler)
	t := todoRouter(todoFile, mu)
	m.Handle("/todo", http.StripPrefix("/todo", t))
	m.Handle("/todo/", http.StripPrefix("/todo/", t))
	return m
}

// replyTextContent wraps plain test into a http response
func replyTextContent(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain") // Specify the response content type
	w.WriteHeader(status)                        // Specify the HTTP response code
	w.Write([]byte(content))                     // Wrap the content into the response
}

// replyJSONContent wraps json object into a http response
func replyJSONContent(w http.ResponseWriter, r *http.Request, status int, resp *todoResponse) {
	// Attempt to unmarshal the response into a json object, if we cannot, return a 500 error
	body, err := json.Marshal(resp)
	if err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

// replyError wraps an error into a http response
func replyError(w http.ResponseWriter, r *http.Request, status int, message string) {
	log.Printf("%s %s: Error: %d %s", r.URL, r.Method, status, message)
	http.Error(w, http.StatusText(status), status)
}
