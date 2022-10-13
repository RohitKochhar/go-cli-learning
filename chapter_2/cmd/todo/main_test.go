package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = os.Getenv("TODO_FILENAME")
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")
	// Depending on the OS, give a file ext
	if runtime.GOOS == "windows" {
		binName += ".exe"
	} else {
		binName += ".out"
	}
	build := exec.Command("go", "build", "-o", binName)
	// Check there was no problem building
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "--> Error while building tool %s: %s", binName, err)
		os.Exit(1)
	} else {
		fmt.Printf("--> Success! Tool %s built\n", binName)
	}
	// If there was no problems, run tests
	fmt.Println("Running tests...")
	result := m.Run()

	// Cleanup
	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)
	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	// Create a new task name
	task := "Test Task Number 1"
	// Try and get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// Create a string representing the command to be run
	cmdPath := filepath.Join(dir, binName)
	// First test to check a new task can be added using CLI
	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	// Second test to esnure the tool can list the tasks
	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("ToDo list:\n\tTask ID: 0, Task Name: %s, Done: false\n", task)
		if expected != string(out) {
			t.Errorf("Expected:\n\t%q\n Got:\n\t%q\n", expected, string(out))
		}
	})
	// Third test to ensure complete cli call works
	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "0")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("Successfully marked task 0 as complete\nSuccessfully saved updated list\nToDo list:\n\tTask ID: 0, Task Name: %s, Done: true\n", task)
		if expected != string(out) {
			t.Errorf("Expected:\n\t%q\n Got:\n\t%q\n", expected, string(out))
		}
	})
	// Fifth test to ensure that we can delete tasks
	t.Run("DeleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-delete", "0")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := "Successfully deleted task 0\nSuccessfully saved updated list\nToDo list:\n"
		if expected != string(out) {
			t.Errorf("Expected:\n\t%q\n Got:\n\t%q\n", expected, string(out))
		}
	})
	// Fourth test to ensure we can add tasks from stdin
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		task2 := "Some input from STDIN"
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

}
