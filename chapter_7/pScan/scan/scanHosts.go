// Package scan provides types and functions to perform TCP port scans on a list of hosts
package scan

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state for a single TCP port
type PortState struct {
	Port int
	Open state
}

// Results represents the scan resutls for a single host
type Results struct {
	Host       string      // host name
	NotFound   bool        // whether or not the host can be resolved
	PortStates []PortState // status of each port scanned
}

// state uses true or false to indicate if a port is open or closed
type state bool

// String converts the bool value of state into either "open" or "closed"
func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

// scanPort performs a port scan on a single TCP port
func scanPort(host string, port int) PortState {
	// Define an instance of PortState that will be returned
	p := PortState{
		Port: port,
	}
	// To check if the given port is open or closed, we will use net.DialTimeout
	// This function tries to connect to a network address within a given time,
	// If it cannot connect within a given time, it returns an error, which we assume means it is closed
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)
	// If the function returned an error, the port is closed so we return p as is
	if err != nil {
		p.Open = false
		return p
	}
	// If the connection succeeds, close the connection, indicate the port is open and return p
	scanConn.Close()
	p.Open = true
	return p
}

// Run performs a port scan on the hosts list
func Run(hl *HostsList, ports []int) []Results {
	// Initialize a slice of results that will be returned
	res := make([]Results, 0, len(hl.Hosts))
	// Loop through the hostsList
	for _, h := range hl.Hosts {
		// Define a result for each host
		r := Results{
			Host: h,
		}
		// Check if the host can be resolved
		if _, err := net.LookupHost(h); err != nil {
			// Set the result to not found if we can't resolve it
			r.NotFound = true
			res = append(res, r)
			// Skip the portScan on this host and continue to the next one
			continue
		}
		// If the host was resolved, do a port scan by iterating through ports provided
		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p))
		}
		// Save the result
		res = append(res, r)
	}
	// Once we are done iterating, return the collection of results
	return res
}
