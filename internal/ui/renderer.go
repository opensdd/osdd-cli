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
func renderPlainASCII() {
	asciiArt := `
   ___                   ____  ____  ____
  / _ \ _ __   ___ _ __ / ___||  _ \|  _ \
 | | | | '_ \ / _ \ '_ \\___ \| | | | | | |
 | |_| | |_) |  __/ | | |___) | |_| | |_| |
  \___/| .__/ \___|_| |_|____/|____/|____/
       |_|
`
	fmt.Print(asciiArt)
}

// renderRainbowASCII renders static ASCII art with rainbow colors
func renderRainbowASCII() {
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

func PrintLogo() {
	if canRenderColored() {
		renderRainbowASCII()
	} else {
		renderPlainASCII()
	}
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
