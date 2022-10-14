package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// filterOut evaluates some metadata about the file or directory
// identified by the path walking process
func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		// Filter out file if its a directory or less than the min size
		return true
	}
	if ext != "" && filepath.Ext(path) != ext {
		// Filter out file if the extension is nul not what we want
		return true
	}
	return false
}

// listFile prints the path of the current file to the specified out pipe
func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string) error {
	splitPath := strings.Split(path, "/")
	if splitPath[0] == "tmp" || (splitPath[0] == "." && splitPath[1] == "tmp") {
		return os.Remove(path)
	} else {
		return fmt.Errorf("cannot delete non-temp files: %s", path)
	}
}
