// Package scan provides types and functions to perform TCP port scans on a list of hosts
package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrExists    = errors.New("error: Host already in the list")
	ErrNotExists = errors.New("error: Host not in the list")
)

// HostsList represents a list of hosts to run port scan
type HostsList struct {
	Hosts []string // Hostnames on which to run port scan
}

// search is a private method that searches for a host in the list
// it will be used by other methods to ensure there are no duplicates in the list
func (hl *HostsList) search(host string) (bool, int) {
	// Sort the hosts alphabetically before searching
	sort.Strings(hl.Hosts)
	// Search the list for a specific host and save the index
	i := sort.SearchStrings(hl.Hosts, host)
	// If our index is valid and provides an existing host, return it
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}
	// If we didn't find it, return -1 and false
	return false, -1
}

// Add adds a host to the list
func (hl *HostsList) Add(host string) error {
	// Use search method to see if the host is already in the list
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}
	// If it is not in the list, append it and return nil
	hl.Hosts = append(hl.Hosts, host)
	return nil
}

// Remove deletes a host from the list
func (hl *HostsList) Remove(host string) error {
	// Use search method to ensure the host is in the list
	if found, i := hl.search(host); found {
		// Change the list to include everything except the to-be-deleted host
		hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
		return nil
	}
	// If the host wasn't found, we cannot delete it
	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// Load obtains a list of hosts from a hosts file
func (hl *HostsList) Load(hostsFile string) error {
	// Try to open the host file
	f, err := os.Open(hostsFile)
	// If we can't, return an error
	if err != nil {
		// If we can't because the file doesn't exist, do nothing
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		// If it is another error, return an error
		return err
	}
	// Close the file when this function is done
	defer f.Close()
	// Create a scanner to read the object
	scanner := bufio.NewScanner(f)
	// Read the lines of the file
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}
	return nil
}

// Save attempts to load the list of hosts into a given hostsFile
func (hl *HostsList) Save(hostsFile string) error {
	// Create an empty string to add the hosts to
	var output string
	// For each host, write a new line to the hostsFile containing the hosts name
	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}
	// Return the result of the file write operation
	return os.WriteFile(hostsFile, []byte(output), 0644)
}
