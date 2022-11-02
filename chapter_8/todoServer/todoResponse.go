package main

import (
	"encoding/json"
	"time"
	"todo"
)

// todoResponse type wraps the results for the api call
type todoResponse struct {
	Results todo.List `json:"results"`
}

// MarshalJson creates an anonymous struct given the original results,
// defines the date and number of results and marshals the information
// into a new json which is returned
func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      todo.List `json:"results"`
		Date         int64     `json:"date"`
		TotalResults int       `json:"total_results"`
	}{
		Results:      r.Results,
		Date:         time.Now().Unix(),
		TotalResults: len(r.Results),
	}
	return json.Marshal(resp)
}
