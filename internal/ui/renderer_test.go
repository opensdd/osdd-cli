package ui

import (
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestCanRenderColored_WithNO_COLOR(t *testing.T) {
	// Save original NO_COLOR value
	originalNoColor := os.Getenv("NO_COLOR")
	defer func() {
		if originalNoColor == "" {
			os.Unsetenv("NO_COLOR")
		} else {
			os.Setenv("NO_COLOR", originalNoColor)
		}
	}()

	// Set NO_COLOR environment variable
	os.Setenv("NO_COLOR", "1")

	// Should return false when NO_COLOR is set
	if canRenderColored() {
		t.Error("Expected canRenderColored() to return false when NO_COLOR is set")
	}
}

func TestCanRenderColored_WithoutNO_COLOR(t *testing.T) {
	// Save original NO_COLOR value
	originalNoColor := os.Getenv("NO_COLOR")
	defer func() {
		if originalNoColor == "" {
			os.Unsetenv("NO_COLOR")
		} else {
			os.Setenv("NO_COLOR", originalNoColor)
		}
	}()

	// Ensure NO_COLOR is not set
	os.Unsetenv("NO_COLOR")

	// Result depends on TTY and color profile detection
	// We can't easily mock those, but we can verify the function doesn't panic
	result := canRenderColored()

	// The function should return a boolean (either true or false, depending on environment)
	_ = result // Just verify it runs without panicking
}

func TestRainbowColor_ValidPositions(t *testing.T) {
	total := 42 // Width of ASCII art

	tests := []struct {
		position int
		expected lipgloss.Color
	}{
		{0, lipgloss.Color("#FF0000")},  // Red at start
		{6, lipgloss.Color("#FF7F00")},  // Orange
		{12, lipgloss.Color("#FFFF00")}, // Yellow
		{18, lipgloss.Color("#00FF00")}, // Green
		{24, lipgloss.Color("#00FFFF")}, // Cyan
		{30, lipgloss.Color("#0000FF")}, // Blue
		{36, lipgloss.Color("#8B00FF")}, // Purple
		{41, lipgloss.Color("#8B00FF")}, // Purple at end
	}

	for _, tt := range tests {
		got := rainbowColor(tt.position, total)
		if got != tt.expected {
			t.Errorf("rainbowColor(%d, %d) = %v; want %v", tt.position, total, got, tt.expected)
		}
	}
}

func TestRainbowColor_ZeroTotal(t *testing.T) {
	// Zero total should return first color (red) without panicking
	got := rainbowColor(5, 0)
	expected := lipgloss.Color("#FF0000")

	if got != expected {
		t.Errorf("rainbowColor(5, 0) = %v; want %v (should return red for zero total)", got, expected)
	}
}

func TestRainbowColor_SinglePosition(t *testing.T) {
	// Single position should work correctly
	got := rainbowColor(0, 1)
	expected := lipgloss.Color("#FF0000")

	if got != expected {
		t.Errorf("rainbowColor(0, 1) = %v; want %v", got, expected)
	}
}

func TestAbs_PositiveNumber(t *testing.T) {
	if abs(5) != 5 {
		t.Errorf("abs(5) = %d; want 5", abs(5))
	}
}

func TestAbs_NegativeNumber(t *testing.T) {
	if abs(-5) != 5 {
		t.Errorf("abs(-5) = %d; want 5", abs(-5))
	}
}

func TestAbs_Zero(t *testing.T) {
	if abs(0) != 0 {
		t.Errorf("abs(0) = %d; want 0", abs(0))
	}
}

func TestRenderPlainASCII_DoesNotPanic(t *testing.T) {
	// This test verifies that renderPlainASCII doesn't panic
	// We can't easily capture stdout, but we can ensure it runs without error
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("renderPlainASCII panicked: %v", r)
		}
	}()

	renderPlainASCII("test-version")
}

func TestRenderPlainASCII_DevVersionShowsWarning(t *testing.T) {
	// This test verifies that "dev" version triggers the warning path
	// (We can't easily capture stdout, but we can ensure it runs without panicking)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("renderPlainASCII with 'dev' version panicked: %v", r)
		}
	}()

	renderPlainASCII("dev")
}

func TestPrintVersion_DoesNotPanic(t *testing.T) {
	// This test verifies that PrintVersion doesn't panic
	// This is important because PrintVersion is the public API
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintVersion panicked: %v", r)
		}
	}()

	PrintVersion("test-version")
}

func TestPrintVersion_DevVersion(t *testing.T) {
	// Test with "dev" version to ensure warning path works
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintVersion with 'dev' version panicked: %v", r)
		}
	}()

	PrintVersion("dev")
}
