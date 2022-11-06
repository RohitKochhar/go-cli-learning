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

## Organizing the Interface's Layout
- Once all the widgets have been defined, they need to be organized and laid out logically to compose the user interface.
- In Termdash, the layout is defined using containers represented by the type `container.Container`. Termdash requires at least one container to start the application, but multiple containers can be used to split the screen and organize the widgets.
- Containers can be created in two different ways:
  - Using the container package to split containers resulting in a binary tree layout
  - Using the grid package to define a grid of rows and columns
- For this application, we use the latter approach as it is easier to organize the code to compose the layout that we desire.

## Building the Interactive Interface
- Once the widgets and layout is ready, we can put everything together to create an app that launches and manages the interface.
- Termdash provides two ways to run dashboard applications:
  - `termdash.Run()`: Starts and manages the application automatically.
    - Using this function, Termdash periodically redraws the screen and handles resizing
  - `termdash.NewController()`: Creates a new instance of `termdash.Controller` that allows for manual management of redrawing and resizing processes.
- Although the first choice is easier, the periodic screen draw continuously consumes system resources, so for this application we will use the second option.
