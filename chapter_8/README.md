# Chapter 8: Talking to REST APIs

## Overview

- Many applications and services use REST APIs as a clear and concise format to expose their data to other applications
- Using Go to interact with these REST APIs opens the dorr to a large number of services that provide many resources for your command line tools
- In this chapter, we will make a command-line tool that connects to a REST API using Go's `net/http` package, exploring advanced concepts usch as the `http.client` and `http.Request` types
- Instead of relying on an external server to send requests to, we will develop our own REST API server so we have a local server to test
- Finally, we will use several testing techniques to test our API server as well as our command-line client application, including local testas, simulated responses, mock servers and integration testas

### Developing a REST API Server

- First, we need to build an API for our command-line tool to talk to.
- To save time, we will use the to-do api that we developed in the first chapter by importing the module.
    - Since this is not a module that is available in a public repository, we have to do some additional actions to import the module.
    - First, we add the dependancy:
    ```bash
    $ go mod edit -require=rohitsingh/todo@v0.0.0
    ```
    - This will fail, since go cannot find the module online, to make sure Go knows to check locally for the file, we have to do the following:
    ```bash
    $ go mod edit -replace=rohitsingh/todo=../../chapter_2
    ```
    - Running `go list -m all` should return no errors
- Once the dependancies are sorted out, we can develop the REST API server.
- For now, we will create the basic structure of the server and the root route, and add the remaining routes to complete the CRUD functionality (Create, Read, Update, Delete) later.
- Our server is created in the todoServer subdirectory. The `main.go` file is the entrypoint to the program, the `server.go` contains server related functions and `handlers.go` contains functions that are executed depening on the request route.

### Testing the REST API Server

- While we could simply run the server and curl to see the response, it would be better to add some structured tests using Go's testing package, specifically the `net/http/httptest` which includes additional types and functions for testing HTTP servers.
-  Because we have implemented our own multiplexer function, it is good to use an approach that allows for integration testing, in which the a test server is create with a URL that simualtes the server, allowing for requests to be made similarly to using curl on the actual server.
