package ui

import "github.com/charmbracelet/lipgloss"

var (
	colorAccent  = lipgloss.Color("#89b4fa")
	colorMuted   = lipgloss.Color("#6c7086")
	colorText    = lipgloss.Color("#cdd6f4")
	colorSurface = lipgloss.Color("#313244")
	colorError   = lipgloss.Color("#f38ba8")
	colorSuccess = lipgloss.Color("#a6e3a1")

	styleApp = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	styleHeader = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true)

	styleMode = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(2)

	styleStatus = lipgloss.NewStyle().
			Foreground(colorMuted)

	styleErr = lipgloss.NewStyle().
			Foreground(colorError)

	styleFooter = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSurface).
			Padding(0, 1).
			MarginTop(1)

	styleNowPlaying = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true)

	styleTime = lipgloss.NewStyle().
			Foreground(colorMuted)

	styleHelp = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)
)
