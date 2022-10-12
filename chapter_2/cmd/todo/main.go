package main

import (
	"flag"
	"fmt"
	"os"
	"todo"
)

// Constants
const todoFileName = ".todo.json"

func main() {
	// Parsing command line flags
	task := flag.String("task", "", "Task name to be included in the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", -1, "ID of task to be completed")
	flag.Parse()

	// Create an item list
	l := &todo.List{}

	// Check for non-empty file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide how to handle given args
	switch {
	case *list:
		l.Print()
	case *complete >= 0:
		// Complete the specified item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			fmt.Printf("Successfully marked task %d as complete\n", *complete)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			fmt.Printf("Successfully saved updated list\n")
			l.Print()
		}
	// Check the default task string was fixed
	case *task != "":
		// Add the task
		l.Add(*task)
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			fmt.Printf("Successfully added new task %s\n", *task)
			l.Print()
		}

	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid Option Provided")
		os.Exit(1)
	}
}
