package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"todo"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidData = errors.New("invalid data")
)

// rootHandler handles requests to the server root
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Check that the client explicitly requested the root path
	if r.URL.Path != "/" {
		replyError(w, r, http.StatusNotFound, "")
		return
	}
	// Return a generic message response to the client if they did request this path
	content := "There's an API here\n"
	replyTextContent(w, r, http.StatusOK, content)
}

// todoRouter looks at incoming requests to /todo and dispatches it to the appropriate function
func todoRouter(todoFile string, l sync.Locker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Instantiate an empty todo list and load the contents of the todofile into it
		list := &todo.List{}
		// Lock the list while reading it to avoid concurrent read/writes and unlock when we are done
		l.Lock()
		defer l.Unlock()
		// Attempt to load the file and return an error response if unsuccessful
		if err := list.Get(todoFile); err != nil {
			replyError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		// Route the request depending on the path of the request and the HTTP method
		// Check if the request was made at the root path, if it was switch based on method
		if r.URL.Path == "" {
			switch r.Method {
			// GET requests to the root asks for all items
			case http.MethodGet:
				getAllHandler(w, r, list)
			// POST requests to root adds an item to the list
			case http.MethodPost:
				addHandler(w, r, list, todoFile)
			// Any other methods are invalid and unsupported
			default:
				replyError(w, r, http.StatusMethodNotAllowed, "Method not supported")
			}
			return
		}
		// Now we can assume that the path has been given a value, meaning the requester
		// is trying to access a specific task id

		// Validate the URL and return an error if an associated item can't be found
		id, err := validateID(r.URL.Path, list)
		if err != nil {
			// Give the user a clearer sense of the error if the item wasn't found
			if errors.Is(err, ErrNotFound) {
				replyError(w, r, http.StatusNotFound, err.Error())
				return
			}
			replyError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		// Handle the request depending on the method
		switch r.Method {
		// GET requests with an ID returns a single task
		case http.MethodGet:
			getOneHandler(w, r, list, id)
		// DELETE requests with an ID removes a single task from the list
		case http.MethodDelete:
			deleteHandler(w, r, list, id, todoFile)
		// PATCH requests update the information for a specific task
		case http.MethodPatch:
			patchHandler(w, r, list, id, todoFile)
		// By default return a method not allowed error
		default:
			replyError(w, r, http.StatusMethodNotAllowed, "Method not supported")
		}
	}
}

// getAllHandler obtains all to-do items from a list
func getAllHandler(w http.ResponseWriter, r *http.Request, list *todo.List) {
	resp := &todoResponse{
		Results: *list,
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}

// getOneHandler obtains a single to-do item from a list
func getOneHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int) {
	resp := &todoResponse{
		Results: (*list)[id-1 : id],
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}

// deleteHandler deletes an item from a list
func deleteHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int, todoFile string) {
	list.Delete(id)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusNoContent, "")
}

// patchHandler completes a specific item
func patchHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int, todoFile string) {
	// Parse the request for queries
	q := r.URL.Query()
	if _, ok := q["complete"]; !ok {
		message := "Missing query parameter `complete`"
		replyError(w, r, http.StatusBadRequest, message)
		return
	}
	// Attempt to mark the task complete and reply accordingly
	list.Complete(id)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusNoContent, "")
}

// addHandler adds a specified task to the todo-list
func addHandler(w http.ResponseWriter, r *http.Request, list *todo.List, todoFile string) {
	item := struct {
		Task string `json:"task"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("Invalid JSON: %s", err)
		replyError(w, r, http.StatusBadRequest, message)
	}
	list.Add(item.Task)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusCreated, "")
}

// validateID ensures the ID provided by the user is valid
func validateID(path string, list *todo.List) (int, error) {
	id, err := strconv.Atoi(path)
	if err != nil {
		return 0, fmt.Errorf("%w: Invalid ID: %s", ErrInvalidData, err)
	}
	if id < 1 {
		return 0, fmt.Errorf("%w, Invalid ID: Less than one", ErrInvalidData)
	}
	if id > len(*list) {
		return id, fmt.Errorf("%w: ID %d not found", ErrNotFound, id)
	}

	return id, nil
}
