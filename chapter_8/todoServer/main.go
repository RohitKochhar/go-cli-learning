package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

// main parses the input flags and hands their values to the main program logic
func main() {
	// Parse input flags
	host := flag.String("h", "localhost", "Server Host")              // Server hostname
	port := flag.Int("p", 8080, "Server Port")                        // Server listening port
	todoFile := flag.String("f", "todoServer.json", "todo JSON file") // Filename to save the to-do list
	flag.Parse()
	// Create an instance of the http.Server type to serve HTTP content, instead of using ListenAndServe
	//   that allows us to serve HTTP without creating a custom server interface, giving us more control
	//   over server options such as read/write timeouts
	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", *host, *port), // The HTTP server listening address (hostname:port)
		Handler:      newMux(*todoFile),                  // The handler to dispatch routes
		ReadTimeout:  10 * time.Second,                   // The time limit to read the entire request including body (if available)
		WriteTimeout: 10 * time.Second,                   // The time limit to send the response back to the client
	}
	// Execute ListenAndServe from the custom defined server to listen for incoming requests
	if err := s.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
