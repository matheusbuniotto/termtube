package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.showSplash {
		w, h := m.width, m.height
		if w == 0 {
			w = 80
		}
		if h == 0 {
			h = 24
		}
		tmp := *m
		tmp.width, tmp.height = w, h
		return tmp.renderSplash()
	}

	var b strings.Builder

	header := styleHeader.Render("ytune") +
		styleMode.Render("  ") +
		modeLabel(m.panel) +
		styleMode.Render(fmt.Sprintf("  · shuffle %s · repeat %s", onOff(m.shuffleOn), m.repeatMode))
	if m.loading {
		header += "  " + m.spin.View()
	}
	b.WriteString(styleApp.Render("♫ ") + header)
	b.WriteString("\n\n")

	if m.jumpFocus {
		line := m.jump.View()
		line = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorAccent).Padding(0, 1).Render(line)
		b.WriteString(styleStatus.Render("Jump") + " " + line)
		b.WriteString("\n\n")
	} else {
		searchLine := m.search.View()
		if m.searchFocus {
			searchLine = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorAccent).Padding(0, 1).Render(searchLine)
		}
		b.WriteString(searchLine)
		b.WriteString("\n\n")
	}

	if m.panel == panelQueue {
		b.WriteString(m.queueList.View())
	} else {
		b.WriteString(m.results.View())
	}

	b.WriteString("\n")
	b.WriteString(m.renderFooter())

	if m.errMsg != "" {
		b.WriteString("\n")
		b.WriteString(styleErr.Render("✗ " + m.errMsg))
	}

	b.WriteString("\n")
	b.WriteString(styleHelp.Render(
		"/ search · :jump 1:30 · g jump · Tab · Enter play · a add · Space · ←→ ±10s · h shuffle · r repeat · n/p · q quit",
	))

	return b.String()
}

func modeLabel(p panel) string {
	if p == panelQueue {
		return styleMode.Render("[Queue]")
	}
	return styleMode.Render("[Results]")
}

func (m *Model) renderFooter() string {
	title := m.statusLine
	if title == "" {
		title = "Nothing playing"
	}
	if len(title) > m.width-10 && m.width > 10 {
		title = title[:m.width-13] + "..."
	}

	prog := 0.0
	if m.dur > 0 {
		prog = m.pos / m.dur
	}
	bar := m.progress.ViewAs(prog)

	timeStr := styleTime.Render(fmt.Sprintf("%s / %s", formatMMSS(m.pos), formatMMSS(m.dur)))
	state := m.playState
	vol := fmt.Sprintf("vol %.0f", m.volume)

	inner := lipgloss.JoinVertical(lipgloss.Left,
		styleNowPlaying.Render(state+" "+title),
		lipgloss.JoinHorizontal(lipgloss.Left, bar, " ", timeStr, "  ", vol),
	)
	return styleFooter.Render(inner)
}
