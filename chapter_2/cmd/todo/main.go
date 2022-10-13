package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"todo"
)

var todoFileName = ".todo.json"

func main() {
	// Parsing command line flags
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", -1, "ID of task to be completed")
	delete := flag.Int("delete", -1, "ID of task to be deleted")
	flag.Parse()

	// Set the environment var
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

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
	case *add:
		// Get the task from either args or stdin
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Add the task
		l.Add(t)
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			fmt.Printf("Successfully added new task %s\n", t)
			l.Print()
		}
	case *delete >= 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			fmt.Printf("Successfully deleted task %d\n", *delete)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			fmt.Printf("Successfully saved updated list\n")
			l.Print()
		}
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid Option Provided")
		os.Exit(1)
	}
}

// getTask function decides where the description for a new task should be
// retrived from, either args or stdin
func getTask(r io.Reader, args ...string) (string, error) {
	// ... is a variadic function, allowing the function to accept 0 or more
	// values of type string
	if len(args) > 0 {
		// Check if we have enough arguments to form a string
		return strings.Join(args, " "), nil
	}
	// If we have a 0 length for args, we will accept from stdin instead
	s := bufio.NewScanner(r)
	// Try and read from stdin
	s.Scan()
	if err := s.Err(); err != nil {
		// Check for errors
		return "", err
	}
	if len(s.Text()) == 0 {
		// Check that we aren't taking blank text
		return "", fmt.Errorf("error: Task cannot be blank")
	}
	return s.Text(), nil
}
