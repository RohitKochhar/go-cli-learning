package scan_test

import (
	"net"
	"rohitsingh/pScan/scan"
	"strconv"
	"testing"
)

// TestStateString tests the String method of the state type
func TestStateString(t *testing.T) {
	// Create PortState object
	ps := scan.PortState{}
	// Check if it is closed by default
	if ps.Open.String() != "closed" {
		t.Errorf("Expected %q, got %q instead", "closed", ps.Open.String())
	}
	// Change to open and check that it has changed
	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("Expected %q, got %q instead", "open", ps.Open.String())
	}
}

// TestRunHostFound tests the Run function assuming the host exists
func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}
	// Test on localhost
	host := "localhost"
	hl := &scan.HostsList{}
	hl.Add(host)
	// Define a slice to hold ports to be checked
	ports := []int{}
	for _, tc := range testCases {
		// 0 is used as a general port that will always be available on the host
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatalf("Unexpected error while listning on port 0: %q\n", err)
		}
		defer ln.Close()
		// Get the port number from the created port
		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatalf("Unexpected error while splitting host port: %q\n", err)
		}
		// Convert the string port number to an int
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatalf("Unexpected error while converting %s to int: %q\n", portStr, err)
		}
		// Add the new port to ports to listen to
		ports = append(ports, port)
		// If we are testing closed port, we close it now to ensure we are using an
		// available port that is closed, not a fake port
		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}
	// execute the Run() method on the ports
	res := scan.Run(hl, ports)
	// Check that we only got 1 result packet from localhost
	if len(res) != 1 {
		t.Fatalf("Expected 1 result, instead got %d\n", len(res))
	}
	// Check that the result returned is from our host
	if res[0].Host != host {
		t.Fatalf("Expected host %q, got %q instead\n", host, res[0].Host)
	}
	// Check that something was returned
	if res[0].NotFound {
		t.Errorf("Expected host %q to be found\n", host)
	}
	// Check that we have two ports in the PortsStates slice
	if len(res[0].PortStates) != 2 {
		t.Fatalf("Expected two states, got %d instead\n", len(res[0].PortStates))
	}
	// Verify each port state
	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port %d, got %d instead\n", ports[0], res[0].PortStates[i].Port)
		}
		if res[0].PortStates[i].Open.String() != tc.expectState {
			t.Errorf("Expected port %d to be %s\n", ports[i], tc.expectState)
		}
	}
}

// TestRunHostNotFound tests the case when the host is not found
func TestRunHostNotFound(t *testing.T) {
	// Create a scan hosts instance and add an invalid DNS to it
	host := "389.389.389.389"
	hl := &scan.HostsList{}
	hl.Add(host)
	// Execute Run using an empty slice for the ports argument, since
	// the host doesn't exist, the ports are irrelevant
	res := scan.Run(hl, []int{})
	// Verify the result
	if len(res) != 1 {
		t.Fatalf("Expected 1 result, got %d instead\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}
	if !res[0].NotFound {
		t.Errorf("Expected host %q NOT to be found\n", host)
	}
	if len(res[0].PortStates) != 0 {
		t.Fatalf("Expected 0 port states, got %d instead\n", len(res[0].PortStates))
	}
}
