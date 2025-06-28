package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	splash := CreateNexaSplash()
	fmt.Print(splash)
}

// CreateNexaSplash prints a colored, boxed ASCII splash for "NEXA" in cyberpunk green style
func CreateNexaSplash() string {
	nexaLines := []string{
		"███╗   ██╗███████╗██╗  ██╗ █████╗      █████╗ ██╗   ██╗████████╗ ██████╗ ",
		"████╗  ██║██╔════╝╚██╗██╔╝██╔══██╗    ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗",
		"██╔██╗ ██║█████╗   ╚███╔╝ ███████║    ███████║██║   ██║   ██║   ██║   ██║",
		"██║╚██╗██║██╔══╝   ██╔██╗ ██╔══██║    ██╔══██║██║   ██║   ██║   ██║   ██║",
		"██║ ╚████║███████╗██╔╝ ██╗██║  ██║    ██║  ██║╚██████╔╝   ██║   ╚██████╔╝",
		"╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝    ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝ ",
	}

	// Cyberpunk green style
	green := lipgloss.Color("#39FF14")
	box := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(green).
		Background(lipgloss.Color("#101820")).
		Padding(1, 4).
		Align(lipgloss.Center)

	var styledLines []string
	style := lipgloss.NewStyle().
		Foreground(green).
		Bold(true)
	for _, line := range nexaLines {
		styledLines = append(styledLines, style.Render(line))
	}

	return "\n" + box.Render(strings.Join(styledLines, "\n")) + "\n"
}
