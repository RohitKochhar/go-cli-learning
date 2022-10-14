package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// config type is used to package arguments in a custom type
// this helps prevent too many positional arguments,
// which would become hard to read
type config struct {
	// File extensions to filter out
	ext string
	// Min file size
	size int64
	// List files
	list bool
	// delete files
	del bool
}

func main() {
	// Override the default help/info message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool Adapted by Rohit Singh, based of the Chapter 4 example from the Pragmatic Bookshelf\n",
			os.Args[0],
		)
		fmt.Fprintf(flag.CommandLine.Output(), "Adapted in October 2022\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	// Parse flags
	root := flag.String("root", ".", "Root directory to being search in")
	// Action flags
	list := flag.Bool("list", false, "List files only")
	del := flag.Bool("del", false, "Delete files")
	// Filter flags
	ext := flag.String("ext", "", "File extensions to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	flag.Parse()
	// Create an instance of the config struct to store flag info
	c := config{
		ext:  *ext,
		size: *size,
		list: *list,
		del:  *del,
	}
	// Call run
	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run defines the logic to descend into the directory and find all
// sub-directories and files within it
func run(root string, out io.Writer, conf config) error {
	return filepath.Walk(root,
		// filepath.Walk requires a function to know what to do once
		// files are found. We use the first-class property of go
		// functions to hand an anonymoust function to filepath.Walk
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filterOut(path, conf.ext, conf.size, info) {
				return nil
			}
			// If list was set, just return the listed files
			if conf.list {
				return listFile(path, out)
			}
			if conf.del {
				return delFile(path)
			}
			// By default, just list the files
			return listFile(path, out)
		})

}
