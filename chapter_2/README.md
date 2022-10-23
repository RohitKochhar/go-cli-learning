# todo

## Types

### type [List](/todo.go#L39)

`type List []item`

List type represents a list of ToDo items

This class ensures all objects within the list are of type
item, preventing runtime errors due to unexpected types

#### func (*List) [Add](/todo.go#L73)

`func (l *List) Add(task string) item`

Add Description:

- Creates a new todo item and appends it to the list

Inputs:

- task (string): name of the new task to be created

Outputs:

- None

```golang
package main

import (
	"fmt"
	"todo"
)

func main() {
	exampleList := todo.List{}
	addedTask := exampleList.Add("Example Task Name")
	fmt.Println(addedTask.Task)
}

```

 Output:

```
Example Task Name
```

#### func (*List) [CheckItemId](/todo.go#L52)

`func (l *List) CheckItemId(id int) error`

CheckItemId Description

- Checks if item in list exists with specified id

Inputs:

- Id (int): id of task to be found

Outputs:

- error (err|nil): err if task not found, nil else

#### func (*List) [Complete](/todo.go#L98)

`func (l *List) Complete(id int) error`

Complete Description:

- Marks a todo item as completed by setting Done = True
and CompletedAt as the current time

Inputs:

- id (int): ID of task to be completed

Outputs:

- error (fmt.Errorf | nil): error if ID is OOB, else nil

```golang
package main

import (
	"fmt"
	"todo"
)

func main() {
	exampleList := todo.List{}
	addedTask := exampleList.Add("Example Task Name")
	exampleList.Complete(addedTask.Id)
	updatedTask := exampleList[addedTask.Id]
	fmt.Println(addedTask.Done, updatedTask.Done)
}

```

 Output:

```
false true
```

#### func (*List) [Delete](/todo.go#L120)

`func (l *List) Delete(id int) error`

Delete Description

- Deletes a ToDo item from the list

Inputs:

- id (int): ID of task to be deleted from list

# Outputs

- error (fmt.Errorf | nil): Error if ID is OOB, else nil

```golang
package main

import (
	"fmt"
	"todo"
)

func main() {
	exampleList := todo.List{}
	addedTask := exampleList.Add("Example Task Name")
	preDeletionLength := len(exampleList)
	exampleList.Delete(addedTask.Id)
	postDeletionLength := len(exampleList)
	fmt.Println(preDeletionLength, postDeletionLength)
}

```

 Output:

```
1 0
```

#### func (*List) [Get](/todo.go#L171)

`func (l *List) Get(filename string) error`

Get Description

- Opens the provided file name, decodes the JSON data and turns into a list

- Performs the inverse function of the Save method

Inputs:

- filename (string): Name of the file to get list from

# Outputs

- result (err|nil|object): Returns error if file is not found, else returns object

#### func (*List) [Print](/todo.go#L190)

`func (l *List) Print()`

Print Description outputs list in human-readable form

#### func (*List) [Save](/todo.go#L150)

`func (l *List) Save(filename string) error`

Save Description

- Uses the json.Marshal function to encode l into JSON

- If json encoding is successful, writes to file specified in args

Inputs:

- filename (string): Name of file to be written to

Outputs:

- error (err|nil): Throws error if there is a problem marshalling item

## Sub Packages

* [cmd/todo](./cmd/todo)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
