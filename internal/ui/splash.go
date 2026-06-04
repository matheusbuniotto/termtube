package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const splashDuration = 2 * time.Second

const asciiLogo = `
████████╗███████╗██████╗ ███╗   ███╗████████╗██╗   ██╗██████╗ ███████╗
╚══██╔══╝██╔════╝██╔══██╗████╗ ████║╚══██╔══╝██║   ██║██╔══██╗██╔════╝
   ██║   █████╗  ██████╔╝██╔████╔██║   ██║   ██║   ██║██████╔╝█████╗
   ██║   ██╔══╝  ██╔══██╗██║╚██╔╝██║   ██║   ██║   ██║██╔══██╗██╔══╝
   ██║   ███████╗██║  ██║██║ ╚═╝ ██║   ██║   ╚██████╔╝██████╔╝███████╗
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝   ╚═╝    ╚═════╝ ╚═════╝ ╚══════╝
`

var (
	styleLogo = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	styleSplashTag = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)

	styleSplashHint = lipgloss.NewStyle().
			Foreground(colorText)
)

type SplashDoneMsg struct{}

func splashDoneCmd() tea.Cmd {
	return tea.Tick(splashDuration, func(time.Time) tea.Msg {
		return SplashDoneMsg{}
	})
}

func (m *Model) dismissSplash() tea.Cmd {
	m.showSplash = false
	m.blurInputs()
	return m.appInit()
}

func (m *Model) renderSplash() string {
	logo := styleLogo.Render(strings.TrimRight(asciiLogo, "\n"))
	tag := styleSplashTag.Render("terminal youtube music")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent).
		Padding(0, 1).
		Render(m.search.View())

	hint := styleSplashHint.Render("type to search · enter · esc to skip")

	block := lipgloss.JoinVertical(lipgloss.Center, logo, "", tag, "", box, "", hint)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		block,
	)
}
