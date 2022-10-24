package scan_test

import (
	"errors"
	"os"
	"rohitsingh/pScan/scan"
	"testing"
)

// TestAdd tests the Add method
func TestAdd(t *testing.T) {
	// Using Table-Driven testing
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an empty hosts list
			hl := &scan.HostsList{}
			// Try to add a host to initialize
			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}
			// Try to add a host to test
			err := hl.Add(tc.host)
			// Check if we expected an error
			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Expected error, got nil instead\n")
				}
				if !errors.Is(err, tc.expectErr) {
					t.Fatalf("Expected an error of %q, got %q instead\n", tc.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %q", err)
			}
			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Expected list length %d, got %d instead\n", tc.expectLen, len(hl.Hosts))
			}
			if hl.Hosts[1] != tc.host {
				t.Errorf("Expected host name %q as index 1, got %q instead\n", tc.host, hl.Hosts[1])
			}
		})
	}
}

// TestRemove tests the Remove method
func TestRemove(t *testing.T) {
	// Using Table-Driven testing
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"RemoveExisting", "host1", 1, nil},
		{"RemoveNotFound", "host3", 1, scan.ErrNotExists},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an empty hosts list
			hl := &scan.HostsList{}
			// Initialize hostlist
			for _, h := range []string{"host1", "host2"} {
				if err := hl.Add(h); err != nil {
					t.Fatal(err)
				}
			}
			// Try to add a host to test
			err := hl.Remove(tc.host)
			// Check if we expected an error
			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Expected error, got nil instead\n")
				}
				if !errors.Is(err, tc.expectErr) {
					t.Fatalf("Expected an error of %q, got %q instead\n", tc.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %q", err)
			}
			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Expected list length %d, got %d instead\n", tc.expectLen, len(hl.Hosts))
			}
			if hl.Hosts[0] == tc.host {
				t.Errorf("Expected host name %q to not be in the list", tc.host)
			}
		})
	}
}

// TestSaveLoad tests both the save and the load methods
func TestSaveLoad(t *testing.T) {
	// Init hostslists
	hl1 := scan.HostsList{}
	hl2 := scan.HostsList{}
	// Add a host to one of the lists
	hostName := "host1"
	hl1.Add(hostName)
	// Create a tempfile to save to
	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error while creating temp file: %s", err)
	}
	// Remove the file when we are done the test
	defer os.Remove(tf.Name())
	// Check for errors
	if err := hl1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list from file: %s", err)
	}
	if err := hl2.Load(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}
	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Fatalf("Host %q should match %q host", hl1.Hosts[0], hl2.Hosts[0])
	}
}

// TestLoadNoFile checks what happens if we try to load a file that is non existent
func TestLoadNoFile(t *testing.T) {
	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	if err := os.Remove(tf.Name()); err != nil {
		t.Fatalf("Error deleting temp file: %s", err)
	}
	hl := &scan.HostsList{}
	// Check that loading the deleted file doesn't cause an error
	if err := hl.Load(tf.Name()); err != nil {
		t.Errorf("Expected no error: got %q instead\n", err)
	}
}
