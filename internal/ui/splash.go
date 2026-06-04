package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const splashDuration = 2 * time.Second

const asciiLogo = `
██╗   ██╗████████╗██╗   ██╗███╗   ██╗███████╗
╚██╗ ██╔╝╚══██╔══╝██║   ██║████╗  ██║██╔════╝
 ╚████╔╝    ██║   ██║   ██║██╔██╗ ██║█████╗
  ╚██╔╝     ██║   ██║   ██║██║╚██╗██║██╔══╝
   ██║      ██║   ╚██████╔╝██║ ╚████║███████╗
   ╚═╝      ╚═╝    ╚═════╝ ╚═╝  ╚═══╝╚══════╝
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
	return m.appInit()
}

func (m *Model) renderSplash() string {
	logo := styleLogo.Render(strings.TrimRight(asciiLogo, "\n"))
	tag := styleSplashTag.Render("terminal youtube music")
	hint := styleSplashHint.Render("press any key to continue")

	block := lipgloss.JoinVertical(lipgloss.Center, logo, "", tag, "", hint)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		block,
	)
}
