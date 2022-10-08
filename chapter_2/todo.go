package todo

import (
	"fmt"
	"time"
)

// item type is represents a todolist item
//
// Attributes
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

// Description:
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
func (l *List) Add(task string) {
	new_task := item{
		Id:          len(*l),
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, new_task)
}

// Description:
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
	if id <= 0 || id > len(ls) {
		return fmt.Errorf("Item %d does not exist", id)
	}
	ls[id].Done = true
	ls[id].CompletedAt = time.Now()

	return nil
}
