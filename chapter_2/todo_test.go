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
	taskNames := []string{
		"Test #3: Test Delete Method",
		"Test #3: Important Task",
		"Test #3: Useless Task",
	}
	addedTask0 := l.Add(taskNames[0])
	addedTask1 := l.Add(taskNames[1])
	addedTask2 := l.Add(taskNames[2])
	// Check the names was added correctly
	if addedTask0.Task != taskNames[0] {
		t.Errorf("Expected %q, got %q instead", taskNames[0], addedTask0.Task)
	}
	if addedTask1.Task != taskNames[1] {
		t.Errorf("Expected %q, got %q instead", taskNames[1], addedTask1.Task)
	}
	if addedTask2.Task != taskNames[2] {
		t.Errorf("Expected %q, got %q instead", taskNames[2], addedTask2.Task)
	}
	// Check that the ID was set correctly
	if addedTask0.Id != 0 {
		t.Errorf("Expected task.Id to be 0, got %d instead", addedTask0.Id)
	}
	if addedTask1.Id != 1 {
		t.Errorf("Expected task.Id to be 1, got %d instead", addedTask1.Id)
	}
	if addedTask2.Id != 2 {
		t.Errorf("Expected task.Id to be 2, got %d instead", addedTask2.Id)
	}
	// Ensure task.Done is false
	if addedTask0.Done || addedTask1.Done || addedTask2.Done {
		t.Errorf("Task should not be done by default")
	}
	// Save the original length of the list
	preDeletionLength := len(l)
	// Delete the task
	l.Delete(addedTask0.Id)
	// Save the new length of the list
	postDeletionLength := len(l)
	if preDeletionLength == postDeletionLength {
		t.Errorf("Expected len(l) to be %d, instead got %d", preDeletionLength-1, postDeletionLength)
	}
	// Save the original length of the list
	preDeletionLength = len(l)
	// Delete the task
	l.Delete(addedTask2.Id)
	// Save the new length of the list
	postDeletionLength = len(l)
	if preDeletionLength == postDeletionLength {
		t.Errorf("Expected len(l) to be %d, instead got %d", preDeletionLength-1, postDeletionLength)
	}
	if (l.CheckItemId(addedTask0.Id) == nil) || (l.CheckItemId(addedTask2.Id) == nil) {
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
