# .

## Variables

Define the error values as variables to be exported to other files.
By convention, the variables should start with `Err`

```golang
var (
    ErrNotNumber        = errors.New("Data is not numeric")
    ErrInvalidColumn    = errors.New("Invalid Column Number")
    ErrNoFiles          = errors.New("No input files")
    ErrInvalidOperation = errors.New("Invalid Operation")
)
```

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
