package main

import "net/http"

// rootHandler handles requests to the server root
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Check that the client explicitly requested the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r) // Responds with an HTTP Not Found Error (404)
		return
	}
	// Return a generic message response to the client if they did request this path
	content := "There's an API here\n"
	replyTextContent(w, r, http.StatusOK, content)
}
