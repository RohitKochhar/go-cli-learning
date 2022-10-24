# Chapter 7: Using the Cobra CLI Framework

## Overview

- So far, all programs we have written have been completely self designed.
- Now, we will make an CLI application using the Cobra framework, which is used for projects like Kubernetes, Hugo and GitHub.

## Cobra Basics
- Cobra provides a library that allows us to design CLI applications supporting POSIX-compliant flags, subcommands, suggestions, autocompletion and automatic help creation
- Cobra integrates with Viper to provide management of configuration and env vars for applications, and provides a generator program that creates boilerplate code to allow developers to focus on designing the business logic of a tool.
- In this chapter, we develop an application called _pScan_, which is a CLI tools that uses sub-commands similar tyo Git or Kubernetes to execute a TCP port scan, similarly to the Nmap command.
- To use Cobra, we need to first install the `cobra-cli` executable using:

```bash
$ go install github.com/spf13/cobra-cli@latest
```

- The install can be verified by running `cobra-cli -h`
- When Cobra generates code, it automatically includes copyright information, such as the author's name and license as specified by either handing the program command line arguments (-a, -l), or by creating a configuration file `.cobra.yaml` in the your user home directory. For example:

```yaml
author: Rohit Singh
license: MIT
```

### Creating a Cobra Application
- To initialize a Cobra application, we have to first be inside an initialized go module, then we run the following command:
``` bash
$ cobra init
```
- This generates the following file structure
```
.
├── LICENSE
├── cmd
│   └── root.go
├── go.mod
├── go.sum
└── main.go
```
- Once created, we can run our program by first downloading any missing dependancies and then running the application
```bash
$ go get
$ go run main.go
```

### Navigating a New Cobra Application
- Cobra structures applications by creating a simple `main.go` file that only imports the package `cmd` and executes the application, which can be seen in the initial `main.go` file.
- The core functionality of the application is within the `cmd` package. When the command is run, the `main()` function calls `cmd.Execute()` to execute the root command of the application. This can be seen in the initial `cmd/root.go` file.
- The `cobra.Command` type is the main type in the Cobra library. It represents a command or subcommand that the tool executes, which can be combined in a parent-child relationship to form a tree structure of subcommands. When Cobra initializes an application, it starts this structure by defining a variable called `rootCmd` as an isntance of the type `cobra.Command` in the `cmd/root.go` file.
- By default, the root command doesn't execute anything, and instead just serves as the parent for other commands (like with `git` or `kubectl`)
- The short and long description that are automatically generated when a Cobra application is initialized should be updated, they are in `cmd/root.go`
- We can also add a `Version: "X"` field to the `rootCmd` so that a version is printed when the user asks for one, for example:
```bash
$ pScan -v (--version)
pScan version 0.1
```
- Additionally, Cobra defines two additional functions in the `cmd/root.go` file:
  - `init()`: Runs before `main()`, used to include additional functionality in the command that can't be defined as properties, like command line flags
  - `initConfig()`: Used by the `cobra.OnInitialize()` function when the application runs. Uses Viper to include configuration management.

### Adding Subcommands to a Cobra Application
- Once the application is initialized, we can use the Cobra generator to add subcommands to it.
- The generator includes a file in the `cmd` directory for each subcommand, each one containing boilerplate code for it's specific subcommand.
- The generator also adds the subcommand to the parent, forming the tree-like structure
- To add a new subcommand, we use the following command:
```bash
$ cobra-cli add hosts
Using config file: /Users/rohitsingh/.cobra.yaml
hosts created at /Users/rohitsingh/Development/F22/go-cli-learning/chapter_7/pScan
```
- Similarly to the parent command, the descriptions in `cmd/hosts.go` should be updated to reflect what the application does.
- We can see the implementation of our new subcommand by running:
```bash
$ go run main.go help hosts
Description:
        Manages the hosts lists for pScan

        Add hosts with the add command
        Delete hosts with the delete command
        List hosts with the list command

Usage:
  pScan hosts [flags]

Flags:
  -h, --help   help for hosts
```
- Also, very cool, Cobra has autosuggestions, so if we run `go run main.go host`, we are met with:
```bash
Error: unknown command "host" for "pScan"

Did you mean this?
        hosts
```
- When we want to add subcommands to a command (like adding list functionality to the hosts subcommand), we can run the following command to generate boilerplate code
```bash
$ cobra-cli add list -p hostsCmd
Using config file: /Users/rohitsingh/.cobra.yaml
list created at /Users/rohitsingh/Development/F22/go-cli-learning/chapter_7/pScan
```
- Even though the command name is `hosts`, we need to specify the instance variable `hostsCmd` as the value for the parent command so Cobra makes the correct association
- The files generated by Cobra only contain boilerplate code, and have to be updated to point to the functionality that you want.
