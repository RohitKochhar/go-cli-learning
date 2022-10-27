package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"rohitsingh/pScan/scan"
	"strconv"
	"strings"
	"testing"
)

func TestHostActions(t *testing.T) {
	// Using table-driven testing
	// Define hosts for action tests
	hosts := []string{"host1", "host2", "host3"}
	testCases := []struct {
		name           string                                  // name of the test to be run
		args           []string                                // args handed to each test
		expectedOut    string                                  // expected result
		initList       bool                                    // whether we need to initialize a list or not
		actionFunction func(io.Writer, string, []string) error // either list, add or delete
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeleteAction",
			args:           []string{"host1", "host2"},
			expectedOut:    "Deleted host: host1\nDeleted host: host2\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run setup for the test
			tf, cleanup := setup(t, hosts, tc.initList)
			// cleanup temp file when the test is done
			defer cleanup()
			// Create a buffer to store output
			var out bytes.Buffer
			if err := tc.actionFunction(&out, tf, tc.args); err != nil {
				t.Fatalf("Unexpected error while running test: %q\n", err)
			}
			// Check if we got what we wanted
			if out.String() != tc.expectedOut {
				t.Errorf("Expected output %q, got %q", tc.expectedOut, out.String())
			}
		})
	}

}

// TestIntegration executes the commands in sequence, like a user would
func TestIntegration(t *testing.T) {
	// Define some initial hosts
	hosts := []string{"host1", "host2", "host3"}
	// setup the test
	tf, cleanup := setup(t, hosts, false)
	defer cleanup()
	// Define the host that will be deleted and the expected result after
	delHost := "host2"
	endHosts := []string{"host1", "host3"}
	// Create a buffer to capture the output
	var out bytes.Buffer
	// Create an empty string to store the expected result
	var expectedOut string
	// Expected result of add operation
	for _, h := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", h)
	}
	// Expected result of list operation
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	// Expected result of delete operation
	expectedOut += fmt.Sprintf("Deleted host: %s\n", delHost)
	// Expected result of list operation after delete
	expectedOut += strings.Join(endHosts, "\n")
	expectedOut += fmt.Sprintln()
	for _, v := range endHosts {
		expectedOut += fmt.Sprintf("%s: Host not found\n", v)
		expectedOut += fmt.Sprintln()
	}
	// Now we can execute the operations in order
	// add hosts
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("Unexpected error while adding host: %q\n", err)
	}
	// list hosts
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Unexpected error while listing hosts: %q\n", err)
	}
	// delete host2
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("Unexpected error while deleting host: %q\n", err)
	}
	// list hosts
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Unexpected error while listing hosts: %q\n", err)
	}
	// scan hosts
	if err := scanAction(&out, tf, nil); err != nil {
		t.Fatalf("unexpected error while scanning hosts: %q\n", err)
	}
	// Check that the output is what we anticipated
	if out.String() != expectedOut {
		t.Errorf("Expected output:\n%q, got:\n%q", expectedOut, out.String())
	}
}

// setup() configures temporary files and initializes a list if required
func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// Create temp file, save the name and then close it
	tf, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatalf("Unexpected error while creating temp file: %q", err)
	}
	tf.Close()
	// Init a list if it is needed
	if initList {
		hl := &scan.HostsList{}
		for _, h := range hosts {
			hl.Add(h)
		}
		if err := hl.Save(tf.Name()); err != nil {
			t.Fatalf("Unexpected error while saving hostsList to temp file: %q", err)
		}
	}
	// Return the temp file name and a cleanup function
	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestScanAction(t *testing.T) {
	// Define the list of hosts for this test
	hosts := []string{"localhost", "unknownhost"}
	// Setup the tests using this list of hosts
	tf, cleanup := setup(t, hosts, true)
	defer cleanup()
	// Initialize the ports, one is open, one is closed
	ports := []int{}
	for i := 0; i < 2; i++ {
		ln, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()
		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}
		ports = append(ports, port)
		if i == 1 {
			ln.Close()
		}
	}
	// Define expected output
	expectedOut := fmt.Sprintln("localhost:")
	expectedOut += fmt.Sprintf("\t%d: open\n", ports[0])
	expectedOut += fmt.Sprintf("\t%d: closed\n", ports[1])
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintln("unknownhost: Host not found")
	expectedOut += fmt.Sprintln()
	// Create a buffer to capture the scan output
	var out bytes.Buffer
	// Execute scan and capture output
	if err := scanAction(&out, tf, ports); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}
	// Test scan output
	if out.String() != expectedOut {
		t.Errorf("Expected output %q, got %q\n", expectedOut, out.String())
	}
}
