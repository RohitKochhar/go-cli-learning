package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// item type is represents a todolist item
//
// # Attributes
//
// - Id (int): Task ID for accessing
//
// - Task (string): name defining the task
//
// - Done (bool): represents whether the task is done
//
// - CreatedAt (time.Time): time at which the task was created
//
// - CompletedAt (time.Time): time at which the task was completed
//
// This is only used internally in this file, so its name is
// defined starting with a lowercase character
type item struct {
	Id          int
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List type represents a list of ToDo items
//
// This class ensures all objects within the list are of type
// item, preventing runtime errors due to unexpected types
type List []item

// CheckItemId Description
//
// - Checks if item in list exists with specified id
//
// Inputs:
//
// - Id (int): id of task to be found
//
// Outputs:
//
// - error (err|nil): err if task not found, nil else
func (l *List) CheckItemId(id int) error {
	ls := *l
	for idx := 0; idx < len(*l); idx++ {
		if ls[idx].Id == id {
			return nil
		}
	}
	return fmt.Errorf("could not find item with Id=%d in list", id)
}

// Add Description:
//
// - Creates a new todo item and appends it to the list
//
// Inputs:
//
// - task (string): name of the new task to be created
//
// Outputs:
//
// - None
func (l *List) Add(task string) item {
	new_task := item{
		Id:          len(*l),
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, new_task)

	return new_task
}

// Complete Description:
//
// - Marks a todo item as completed by setting Done = True
// and CompletedAt as the current time
//
// Inputs:
//
// - id (int): ID of task to be completed
//
// Outputs:
//
// - error (fmt.Errorf | nil): error if ID is OOB, else nil
func (l *List) Complete(id int) error {
	ls := *l
	if ls.CheckItemId(id) != nil {
		return fmt.Errorf("could not find item with Id=%d in list", id)
	}
	ls[id].Done = true
	ls[id].CompletedAt = time.Now()

	return nil
}

// Delete Description
//
// - Deletes a ToDo item from the list
//
// Inputs:
//
// - id (int): ID of task to be deleted from list
//
// # Outputs
//
// - error (fmt.Errorf | nil): Error if ID is OOB, else nil
func (l *List) Delete(id int) error {
	// Dereference pointer to mutate object
	ls := *l
	if ls.CheckItemId(id) != nil {
		return fmt.Errorf("could not find item with Id=%d in list", id)
	}
	// Go through the list and delete the item with matching ID
	var newList []item
	for idx := 0; idx < len(ls); idx++ {
		if ls[idx].Id != id {
			newList = append(newList, ls[idx])
		}
	}
	*l = newList
	return nil
}

// Save Description
//
// - Uses the json.Marshal function to encode l into JSON
//
// - If json encoding is successful, writes to file specified in args
//
// Inputs:
//
// - filename (string): Name of file to be written to
//
// Outputs:
//
// - error (err|nil): Throws error if there is a problem marshalling item
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, js, 0644)
}

// Get Description
//
// - Opens the provided file name, decodes the JSON data and turns into a list
//
// - Performs the inverse function of the Save method
//
// Inputs:
//
// - filename (string): Name of the file to get list from
//
// # Outputs
//
// - result (err|nil|object): Returns error if file is not found, else returns object
func (l *List) Get(filename string) error {
	// Try opening the file for reading
	file, err := os.ReadFile(filename)
	// If the error exists, check what error
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// If the file doesn't exist, just return a blank object (nil)
			return nil
		}
	}
	// If the file is found but has nothing in it, just return nil
	if len(file) == 0 {
		return nil
	}
	// If the file is found and is not empty, unmarshal it
	return json.Unmarshal(file, l)
}

// Print Description outputs list in human-readable form
func (l *List) Print() {
	ls := *l
	fmt.Println("ToDo list:")
	for idx := 0; idx < len(ls); idx++ {
		fmt.Printf("\tTask ID: %d, Task Name: %s, Done: %t\n", ls[idx].Id, ls[idx].Task, ls[idx].Done)
	}
}
