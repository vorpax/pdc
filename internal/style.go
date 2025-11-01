package internal

import "github.com/charmbracelet/lipgloss"

var (
	// SuccessStyle is for success messages.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Padding(1, 2)

	// ErrorStyle is for error messages.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Padding(1, 2)

	// InfoStyle is for informational messages.
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF")).
			Padding(1, 2)

	// VerboseStyle is for verbose output.
	VerboseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")).
			Padding(0, 4)
)
