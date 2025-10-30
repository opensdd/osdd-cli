package version

import (
	"testing"
)

func TestGetVersion_WhenVersionIsEmpty(t *testing.T) {
	// Save original version
	originalVersion := Version
	defer func() { Version = originalVersion }()

	// Set Version to empty string
	Version = ""

	got := GetVersion()
	want := "dev"

	if got != want {
		t.Errorf("GetVersion() = %q; want %q", got, want)
	}
}

func TestGetVersion_WhenVersionIsSet(t *testing.T) {
	// Save original version
	originalVersion := Version
	defer func() { Version = originalVersion }()

	// Set Version to a specific value
	Version = "v1.2.3"

	got := GetVersion()
	want := "v1.2.3"

	if got != want {
		t.Errorf("GetVersion() = %q; want %q", got, want)
	}
}

func TestIsSet_WhenVersionIsEmpty(t *testing.T) {
	// Save original version
	originalVersion := Version
	defer func() { Version = originalVersion }()

	// Set Version to empty string
	Version = ""

	got := IsSet()
	want := false

	if got != want {
		t.Errorf("IsSet() = %v; want %v", got, want)
	}
}

func TestIsSet_WhenVersionIsSet(t *testing.T) {
	// Save original version
	originalVersion := Version
	defer func() { Version = originalVersion }()

	// Set Version to a specific value
	Version = "v1.2.3"

	got := IsSet()
	want := true

	if got != want {
		t.Errorf("IsSet() = %v; want %v", got, want)
	}
}
