package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// canRenderColored checks if the terminal supports colored rendering
func canRenderColored() bool {
	// Check 1: NO_COLOR environment variable (user preference to disable colors)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check 2: stdout must be a TTY (not piped or redirected)
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}

	// Check 3: Color profile must support at least 256 colors for rainbow gradient
	profile := lipgloss.DefaultRenderer().ColorProfile()
	return profile <= termenv.ANSI256
}

// renderPlainASCII renders plain ASCII art without animation
func renderPlainASCII(version string) {
	asciiArt := `
   ___                   ____  ____  ____
  / _ \ _ __   ___ _ __ / ___||  _ \|  _ \
 | | | | '_ \ / _ \ '_ \\___ \| | | | | | |
 | |_| | |_) |  __/ | | |___) | |_| | |_| |
  \___/| .__/ \___|_| |_|____/|____/|____/
       |_|
`
	fmt.Print(asciiArt)
	fmt.Printf("\nOpenSDD CLI version %s\n", version)

	if version == "dev" {
		fmt.Println("WARNING: Version not set at build time")
	}
}

// renderRainbowASCII renders static ASCII art with rainbow colors
func renderRainbowASCII(version string) {
	asciiArt := []string{
		"   ___                   ____  ____  ____  ",
		"  / _ \\ _ __   ___ _ __ / ___||  _ \\|  _ \\ ",
		" | | | | '_ \\ / _ \\ '_ \\\\___ \\| | | | | | |",
		" | |_| | |_) |  __/ | | |___) | |_| | |_| |",
		"  \\___/| .__/ \\___|_| |_|____/|____/|____/ ",
		"       |_|                                  ",
	}

	// Render each line with rainbow colors
	for _, line := range asciiArt {
		for charIdx, ch := range line {
			// Calculate rainbow color based on character position
			color := rainbowColor(charIdx, len(line))
			style := lipgloss.NewStyle().Foreground(color)
			fmt.Print(style.Render(string(ch)))
		}
		fmt.Println()
	}

	// Print version
	fmt.Printf("\nOpenSDD CLI version %s\n", version)

	if version == "dev" {
		fmt.Println("WARNING: Version not set at build time")
	}
}

// rainbowColor returns a rainbow color based on position
func rainbowColor(position, total int) lipgloss.Color {
	// Rainbow colors: red → orange → yellow → green → cyan → blue → purple
	rainbow := []string{
		"#FF0000", // Red
		"#FF7F00", // Orange
		"#FFFF00", // Yellow
		"#00FF00", // Green
		"#00FFFF", // Cyan
		"#0000FF", // Blue
		"#8B00FF", // Purple
	}

	// Map position to rainbow color index
	if total == 0 {
		return lipgloss.Color(rainbow[0])
	}

	// Calculate which color to use based on position across the width
	colorIndex := (position * len(rainbow)) / total
	if colorIndex >= len(rainbow) {
		colorIndex = len(rainbow) - 1
	}

	return lipgloss.Color(rainbow[colorIndex])
}

// PrintVersion displays the version with ASCII art in rainbow colors when supported.
//
// The function automatically detects terminal capabilities and renders either:
//   - Rainbow-colored ASCII art (when TTY with ANSI256+ colors)
//   - Plain ASCII art (when piped, in CI/CD, or limited terminal support)
//
// Behavior is controlled by:
//   - NO_COLOR environment variable (disables colors when set)
//   - TTY detection (requires stdout to be a terminal)
//   - Color profile detection (requires ANSI256+ color support)
//
// If rendering panics (rare), the function automatically falls back to plain
// rendering and logs a warning to stderr. This function never returns an error
// and always produces output.
func PrintVersion(version string) {
	// Panic recovery for rendering errors
	defer func() {
		if r := recover(); r != nil {
			// Log the panic for debugging visibility
			fmt.Fprintf(os.Stderr, "Warning: Rainbow rendering failed (%v), falling back to plain rendering\n", r)
			renderPlainASCII(version)
		}
	}()

	if canRenderColored() {
		renderRainbowASCII(version)
	} else {
		renderPlainASCII(version)
	}
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Note: max() and min() are builtin functions in Go 1.21+, no custom implementation needed
