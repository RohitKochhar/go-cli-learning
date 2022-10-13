# Chapter 4: Navigating the File System

## Overview

- In this chapter, we create a file system crawler that descends into a directory tree to look for files that match a specific criteria and perform some action based on it.
- For this program we use a design pattern which breaks the code into different files, allowing for easier maintenance and version control.
- This tool uses the following packages:
  - `flag`: To handle command-line flags
  - `fmt`: To print formatted output
  - `io`: To access the io.Writer object
  - `path/filepath`: To handle file paths appropriately across OSes
- For testing, we use _table-driven testing_, by which we define test cases as a slice of anonymous struct, containing the data required to run tests and the expected results

```
.
├── main.go         // contains main() and run() functions
├── main_test.go    // test suite for main_test.go
├── actions.go      // contains functions for file actions
├── actions_test.go // test suite for actions.go
└── testdata        // contains test files for test suites
    └── ...
```
