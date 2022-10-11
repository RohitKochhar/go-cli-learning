package todo_test

import (
	"fmt"
	"testing"
	"todo"
)

// Tests:

// TestAdd tests the Add method for the List type
func TestAdd(t *testing.T) {
	l := todo.List{}
	taskName := "Test #1: Test Add method"
	addedTask := l.Add(taskName)
	if addedTask.Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, addedTask.Task)
	}
	// Check that the item can be found in the list
	if l.CheckItemId(addedTask.Id) != nil {
		t.Errorf("Expected item to be in list, but it was not found")
	}
}

// TestComplete checks that the Complete method updates the Done bool
func TestComplete(t *testing.T) {
	// Create a test list
	l := todo.List{}
	// Add a sample task
	taskName := "Test#2: Test Complete Method"
	addedTask := l.Add(taskName)
	if addedTask.Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Task)
	}
	// Check that the ID was set correctly
	if addedTask.Id != 0 {
		t.Errorf("Expected task.Id to be 0, got %d instead", addedTask.Id)
	}
	// Check that the item can be found in the list
	if l.CheckItemId(addedTask.Id) != nil {
		t.Errorf("Expected item to be in list, but it was not found")
	}
	// Ensure task.Done is false
	if addedTask.Done {
		t.Errorf("Task should not be done by default")
	}
	// Complete the task
	l.Complete(addedTask.Id)
	// Ensure task.Done was updated correctly
	if !l[addedTask.Id].Done {
		t.Errorf("Expected task.Done to be true, instead got false")
	}
}

// TestDelete checks that the Delete method removes the Task from the list
func TestDelete(t *testing.T) {
	// Create a test list
	l := todo.List{}
	// Add a sample item
	taskName := "Test #3: Test Delete Method"
	addedTask := l.Add(taskName)
	// Check the name was added correctly
	if addedTask.Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, addedTask.Task)
	}
	// Check that the ID was set correctly
	if addedTask.Id != 0 {
		t.Errorf("Expected task.Id to be 0, got %d instead", addedTask.Id)
	}
	// Ensure task.Done is false
	if addedTask.Done {
		t.Errorf("Task should not be done by default")
	}
	// Save the original length of the list
	preDeletionLength := len(l)
	// Delete the task
	l.Delete(addedTask.Id)
	// Save the new length of the list
	postDeletionLength := len(l)
	if preDeletionLength == postDeletionLength {
		t.Errorf("Expected len(l) to be %d, instead got %d", preDeletionLength-1, postDeletionLength)
	}
	if l.CheckItemId(addedTask.Id) == nil {
		t.Errorf("Expected item to not be in list, but it was found")
	}
}

// Examples

func ExampleList_Add() {
	exampleList := todo.List{}
	addedTask := exampleList.Add("Example Task Name")
	fmt.Println(addedTask.Task)
	// Output: Example Task Name
}

func ExampleList_Complete() {
	exampleList := todo.List{}
	addedTask := exampleList.Add("Example Task Name")
	exampleList.Complete(addedTask.Id)
	updatedTask := exampleList[addedTask.Id]
	fmt.Println(addedTask.Done, updatedTask.Done)
	// Output: false true
}

func ExampleList_Delete() {
	exampleList := todo.List{}
	addedTask := exampleList.Add("Example Task Name")
	preDeletionLength := len(exampleList)
	exampleList.Delete(addedTask.Id)
	postDeletionLength := len(exampleList)
	fmt.Println(preDeletionLength, postDeletionLength)
	// Output: 1 0
}
