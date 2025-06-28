package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SplashArt holds configuration for rendering ASCII splash art.
type SplashArt struct {
	Art        []string // ASCII art lines
	ColorName  string   // Foreground color name (e.g., "green", "cyan")
	Background string   // Background color hex (e.g., "#101820")
	Border     lipgloss.Border
	Padding    [2]int // [vertical, horizontal]
	Align      lipgloss.Alignment
	Bold       bool
}

// colorNameToHex maps common color names to hex codes.
var colorNameToHex = map[string]string{
	"green":    "#39FF14",
	"cyan":     "#00FFFF",
	"magenta":  "#FF00FF",
	"yellow":   "#FFFF00",
	"red":      "#FF3131",
	"blue":     "#00BFFF",
	"white":    "#FFFFFF",
	"black":    "#000000",
	"orange":   "#FFA500",
	"purple":   "#A020F0",
	"gray":     "#888888",
	"lime":     "#BFFF00",
	"pink":     "#FF69B4",
	"teal":     "#008080",
	"gold":     "#FFD700",
}

// Render generates the styled splash art as a string.
func (s SplashArt) Render() string {
	fg := lipgloss.Color(colorNameToHex[s.ColorName])
	bg := lipgloss.Color(s.Background)
	box := lipgloss.NewStyle().
		Border(s.Border).
		BorderForeground(fg).
		Background(bg).
		Padding(s.Padding[0], s.Padding[1]).
		Align(s.Align)

	style := lipgloss.NewStyle().
		Foreground(fg).
		Bold(s.Bold)

	var styledLines []string
	for _, line := range s.Art {
		styledLines = append(styledLines, style.Render(line))
	}
	return "\n" + box.Render(strings.Join(styledLines, "\n")) + "\n"
}

// Example usage: Generate a splash artifact for "NEXA"
func main() {
	nexaArt := []string{
		"███╗   ██╗███████╗██╗  ██╗ █████╗      █████╗ ██╗   ██╗████████╗ ██████╗ ",
		"████╗  ██║██╔════╝╚██╗██╔╝██╔══██╗    ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗",
		"██╔██╗ ██║█████╗   ╚███╔╝ ███████║    ███████║██║   ██║   ██║   ██║   ██║",
		"██║╚██╗██║██╔══╝   ██╔██╗ ██╔══██║    ██╔══██║██║   ██║   ██║   ██║   ██║",
		"██║ ╚████║███████╗██╔╝ ██╗██║  ██║    ██║  ██║╚██████╔╝   ██║   ╚██████╔╝",
		"╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝    ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝ ",
	}

	splash := SplashArt{
		Art:        nexaArt,
		ColorName:  "green",
		Background: "#101820",
		Border:     lipgloss.ThickBorder(),
		Padding:    [2]int{1, 4},
		Align:      lipgloss.Center,
		Bold:       true,
	}

	fmt.Print(splash.Render())
}
