package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// Constants
const (
	// header containers the HTML header required to view
	// html in a browser
	header = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html"; charset=utf-8">
		<title>Markdown Preview Tool</title>
	</head>
	<body>
`
	footer = `
	</body>
</html>
`
)

// Functions

// main parses flags to determine which file to pass to run
func main() {
	// Override the default help/info message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool Adapted by Rohit Singh, based of the Chapter 3 example from the Pragmatic Bookshelf\n",
			os.Args[0],
		)
		fmt.Fprintf(flag.CommandLine.Output(), "Adapted in October 2022\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("skipPreview", false, "Create HTML file without preview in browser")
	// help := flag.Bool("help", false, "Displays this message")
	flag.Parse()

	if *filename == "" {
		// If no flag is provided, pass a help message
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		// Try to run the program without error
		os.Exit(1)
	}
}

// run coordinates the execution of the remaining functions
func run(filename string, out io.Writer, skipPreview bool) error {
	// Parse the input file for any errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	// Convert the input to HTML data
	htmlData := parseContent(input)
	// Create a temporary file to prevent garbage
	temp, err := os.CreateTemp("", "mdp*.html")
	// Check for errors
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()

	// Print the outName to stdout for clarity and testing
	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}
	if skipPreview {
		return nil
	}

	// Once the run function is done, remove the tempFile
	defer os.Remove(outName)

	return preview(outName)
}

// parseContent goes through the MD input and converts to HTML
//
// The function takes in the MD file as an array of bytes and returns
// the html data as an array of bytes
func parseContent(input []byte) []byte {
	// First we pass it through blackfriday to generate HTML
	output := blackfriday.Run(input)
	// Pass blackfriday output to bluemonday to santize output
	body := bluemonday.UGCPolicy().SanitizeBytes(output)
	// Create buffer to store the content
	var buffer bytes.Buffer
	// Write the html to the buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)
	// Return the generated HTML data
	return buffer.Bytes()
}

// saveHTML saves the content created by parseContent into an html file
func saveHTML(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

// preview takes the tempfile name and tries to open it in the browser
func preview(fname string) error {
	cName := ""
	cParams := []string{}
	// Define the executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}
	// Add the file name to cParams so it is handed to the final command
	cParams = append(cParams, fname)
	// Ensure that the path will be a valid input to the preview cmd
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser time to open the file to prevent a race condition
	time.Sleep(2 * time.Second)
	return err
}
