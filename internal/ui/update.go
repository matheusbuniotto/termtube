package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerH := 4
		footerH := 5
		listH := msg.Height - headerH - footerH
		if listH < 4 {
			listH = 4
		}
		m.results.SetSize(msg.Width-4, listH)
		m.queueList.SetSize(msg.Width-4, listH)
		m.search.Width = min(msg.Width-6, 60)
		m.progress.Width = msg.Width - 24
		if m.progress.Width < 10 {
			m.progress.Width = 10
		}

	case SplashDoneMsg:
		if m.showSplash {
			return m, m.dismissSplash()
		}

	case tea.KeyMsg:
		if m.showSplash {
			if msg.String() == "ctrl+c" {
				m.svc.Stop()
				return m, tea.Quit
			}
			return m, m.dismissSplash()
		}

		if m.searchFocus {
			switch msg.String() {
			case "enter":
				q := m.search.Value()
				if q != "" {
					m.loading = true
					m.errMsg = ""
					m.statusLine = "Searching…"
					cmds = append(cmds, searchCmd(m.svc, q))
				}
			case "esc":
				m.searchFocus = false
				m.search.Blur()
			default:
				var cmd tea.Cmd
				m.search, cmd = m.search.Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			m.svc.Stop()
			return m, tea.Quit

		case "/":
			m.searchFocus = true
			cmds = append(cmds, m.search.Focus())

		case "tab":
			if m.panel == panelResults {
				m.panel = panelQueue
				m.refreshQueue()
			} else {
				m.panel = panelResults
			}

		case "enter":
			if v, ok := m.selectedVideo(); ok {
				m.loading = true
				m.statusLine = "Loading stream…"
				m.errMsg = ""
				cmds = append(cmds, playCmd(m.svc, v))
			}

		case "a":
			if m.panel == panelResults {
				if v, ok := m.selectedVideo(); ok {
					m.statusLine = "Adding to queue…"
					cmds = append(cmds, addQueueCmd(m.svc, v))
				}
			}

		case " ":
			_ = m.svc.TogglePause()
			st := m.svc.Status()
			if st.Paused {
				m.playState = "⏸"
			} else if st.Playing {
				m.playState = "▶"
			}

		case "n":
			m.statusLine = "Next track…"
			cmds = append(cmds, nextCmd(m.svc))

		case "p":
			m.statusLine = "Previous track…"
			cmds = append(cmds, prevCmd(m.svc))

		case "+", "=":
			_ = m.svc.VolumeUp()
			m.volume = m.svc.Status().Volume

		case "-", "_":
			_ = m.svc.VolumeDown()
			m.volume = m.svc.Status().Volume

		case "d":
			if m.panel == panelQueue {
				idx := m.queueList.Index()
				m.svc.RemoveFromQueue(idx)
				m.refreshQueue()
			}
		}

	case SearchDoneMsg:
		m.loading = false
		if msg.Err != nil {
			m.errMsg = msg.Err.Error()
			m.statusLine = ""
		} else {
			m.setResults(msg.Results)
			m.statusLine = "Found results — Enter to play"
			m.errMsg = ""
			m.panel = panelResults
		}
		cmds = append(cmds, tickCmd(m.svc))

	case PlayDoneMsg:
		m.loading = false
		if msg.Err != nil {
			m.errMsg = msg.Err.Error()
		} else {
			m.errMsg = ""
			m.refreshQueue()
			m.statusLine = "Playing"
		}

	case AddQueueDoneMsg:
		if msg.Err != nil {
			m.errMsg = msg.Err.Error()
		} else {
			m.errMsg = ""
			m.statusLine = "Added to queue"
			m.refreshQueue()
		}

	case PlayerTickMsg:
		st := m.svc.Status()
		m.pos = st.Position
		m.dur = st.Duration
		m.volume = st.Volume
		switch {
		case !st.Playing:
			m.playState = "⏹"
		case st.Paused:
			m.playState = "⏸"
		default:
			m.playState = "▶"
		}
		if st.Title != "" {
			m.statusLine = st.Title
		}
		if msg.Err != nil {
			m.errMsg = msg.Err.Error()
		}
		cmds = append(cmds, tickCmd(m.svc))

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update active list when not in search focus
	if m.showSplash {
		return m, tea.Batch(cmds...)
	}
	if !m.searchFocus {
		var cmd tea.Cmd
		if m.panel == panelQueue {
			m.queueList, cmd = m.queueList.Update(msg)
		} else {
			m.results, cmd = m.results.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
