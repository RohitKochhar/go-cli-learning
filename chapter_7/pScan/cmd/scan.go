/*
Copyright © 2022 Rohit Singh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"rohitsingh/pScan/scan"

	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the hosts to run port scans
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}
		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func scanAction(out io.Writer, hostsFile string, ports []int) error {
	// Create a host list and try to load data from the file
	hl := &scan.HostsList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	// If load was successful, run the port scan
	results := scan.Run(hl, ports)
	return printResults(out, results)
}

func printResults(out io.Writer, results []scan.Results) error {
	// Create empty var to store result message
	var msg string
	// For each result, add the host name and the list of ports/status
	for _, r := range results {
		msg += fmt.Sprintf("%s:", r.Host)
		if r.NotFound {
			msg += " Host not found\n\n"
			continue
		}
		msg += fmt.Sprintln()
		for _, p := range r.PortStates {
			msg += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}
		msg += fmt.Sprintln()
	}
	_, err := fmt.Fprint(out, msg)
	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443}, "ports to scan")
}
