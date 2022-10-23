// package main
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

// config type holds flag information
type config struct {
	proj      string
	test      bool
	format    bool
	gitCommit string
	readme    bool
}

// executer interface specifies which execute function to use for a step
type executer interface {
	Execute() (string, error)
}

// main parses command-line flags and then calls run
func main() {
	// Parse flags
	proj := flag.String("p", "", "Project Directory")
	test := flag.Bool("test", true, "Run go test -v")
	format := flag.Bool("format", true, "Run go fmt")
	gitCommit := flag.String("commit", "", "Run git commit -m {PROVIDED MESSAGE}")
	readme := flag.Bool("readme", false, "Overwrite project with auto-generated README.md")
	conf := config{
		proj:      *proj,
		test:      *test,
		format:    *format,
		gitCommit: *gitCommit,
		readme:    *readme,
	}
	flag.Parse()
	if err := run(conf, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

// run contains the main logic of the program
func run(conf config, out io.Writer) error {
	if conf.proj == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}
	// Create a new pipeline
	pipeline := []executer{}
	// Add build step to pipeline
	pipeline = append(pipeline,
		NewStep(
			"go build",
			"go",
			"go build: SUCCESS",
			conf.proj,
			[]string{"build", ".", "errors"},
		),
	)
	if conf.test {
		pipeline = append(pipeline,
			NewStep(
				"go test",
				"go",
				"go test: SUCCESS",
				conf.proj,
				[]string{"test", "-v"},
			),
		)
	}
	if conf.format {
		pipeline = append(pipeline,
			NewExceptionStep(
				"go fmt",
				"gofmt",
				"gofmt: SUCCESS",
				conf.proj,
				[]string{"-l", "."},
			),
		)
	}
	if conf.readme {
		pipeline = append(pipeline,
			NewStep(
				"goreadme",
				"bash",
				"goreadme: SUCCESS",
				conf.proj,
				[]string{"-c", "goreadme -types -variabless -functions -methods -recursive -factories -constants > ./README.md"},
			),
		)
	}
	if conf.gitCommit != "" {
		pipeline = append(pipeline,
			NewTimeoutStep(
				"git commit",
				"git",
				"git commit: SUCCESS",
				conf.proj,
				[]string{"-m", conf.gitCommit},
				10*time.Second,
			),
		)
	}
	// Iterate through the pipeline and execute each step
	for _, s := range pipeline {
		// Execute the step
		msg, err := s.Execute()
		if err != nil {
			// If we have an error, stop the pipeline
			return err
		}
		// If we don't have an error, output the result to the interface provided to run()
		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}
	// If we made it here, there were no errors, return nil
	return nil
}
