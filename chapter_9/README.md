# Chapter 9: Developing Interactive Terminal Tools

## Overview
- So far, the tools we have built have run mostly unattended, the user runs a command once and then some process executes, and the program terminates. 
- Now, we will work on an application that is suited for an interactive workflow where the user provides and receives feedback continously. 
- In this chapter, we will develop an interactvie Pomodoro timer application.
- Rather than developing a full GUI, we'll design and implement an interactive CLI application that runs directly in the terminal.
- CLI applications require less resources than GUI applications and are significantly more portable.
- For this project, we will implement the Repository Pattern to abstract the data source, decoupling the business logic from the data, allowing us to implement different data stores according to our requirements.

## Storing Data with the Repository Pattern
- We will implement a data store for the Pomodoro intervals using the Repository pattern.
- With this approach, the data store implementation is decoupled from the business logic, allowing for flexibility in storing data.
- This means that the data store implementation can be changed or switched to a different database entirely at a later time without impacting the business logic.
- We will implement two different data stores with this application:
  - An in-memory data store
  - A SQLite database
- The repository pattern requires two components:
    1. An interface that specifies all the methods a given type must implement to qualify as a repository for this application
    2. A custom type which implements that interface working as the repository.
- Our implemented repositories will be stored in `pomo/pomodoro/repository/`

## Building the Interface Widgets
- We will create the basic interface that has the controls required to run and display the Pomodoro status using the [Termdash](https://github.com/mum4k/termdash) dashboard library
- Termdash is a good option because it's cross-platform, under active development and has a good set of features like graphical widgets, dashboard resizing, customizable layouts and handling of mouse and keyboard events.
- Termdash requires a set of other libraries to run as backend, so some additional libraries must be added to the project.


