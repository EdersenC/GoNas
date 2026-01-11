package test

import (
	"goNAS/storage"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestFunk tests the Funk function with a real filesystem device
func TestFunk(t *testing.T) {
	// Try to find a valid mounted filesystem to test with
	// We'll use the root filesystem as it's always available
	device := "/"

	total, avail, err := storage.Funk(device)
	if err != nil {
		t.Fatalf("Funk failed: %v", err)
	}

	// Basic sanity checks
	if total == 0 {
		t.Errorf("Expected non-zero total capacity, got %d", total)
	}

	if avail == 0 {
		t.Errorf("Expected non-zero available capacity, got %d", avail)
	}

	if avail > total {
		t.Errorf("Available capacity (%d) should not exceed total capacity (%d)", avail, total)
	}

	t.Logf("Total capacity: %d bytes, Available capacity: %d bytes", total, avail)
}

// TestFunkInvalidDevice tests the Funk function with an invalid device
func TestFunkInvalidDevice(t *testing.T) {
	device := "/nonexistent/device/path"

	_, _, err := storage.Funk(device)
	if err == nil {
		t.Error("Expected error for invalid device, got nil")
	}

	if !strings.Contains(err.Error(), "df command failed") {
		t.Errorf("Expected 'df command failed' error, got: %v", err)
	}
}

// TestFunkWithDfAvailable tests if df command is available
func TestFunkWithDfAvailable(t *testing.T) {
	// Check if df command is available
	_, err := exec.LookPath("df")
	if err != nil {
		t.Skip("df command not available on this system")
	}

	// Test with /tmp which should exist on most systems
	device := "/tmp"
	if _, err := os.Stat(device); os.IsNotExist(err) {
		t.Skip("/tmp directory not available on this system")
	}

	total, avail, err := storage.Funk(device)
	if err != nil {
		t.Fatalf("Funk failed with /tmp: %v", err)
	}

	if total == 0 || avail == 0 {
		t.Errorf("Expected non-zero values, got total=%d, avail=%d", total, avail)
	}
}
