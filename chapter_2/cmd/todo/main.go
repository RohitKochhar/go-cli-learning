package main

import (
	"fmt"
	"os"
	"strings"
	"todo"
)

// Constants
const todoFileName = ".todo.json"

func main() {
	l := &todo.List{}

	// Check for non-empty file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide how to handle given args
	switch {
	// If no args, print the todo list
	case len(os.Args) == 1:
		l.Print()
	default:
		// By default, if an arg is given, we assume it is a string
		// 	representing the task name
		newTaskName := strings.Join(os.Args[1:], " ")
		// Add the new task to the list
		l.Add(newTaskName)
		// Save the list to the file
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("New task %s saved successfully\n", newTaskName)
		l.Print()

	}
}
