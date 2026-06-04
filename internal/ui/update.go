package ui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/monkmode/ytune/internal/player"
)

func (m *Model) blurInputs() {
	m.searchFocus = false
	m.jumpFocus = false
	m.search.Blur()
	m.jump.Blur()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerH := 5
		footerH := 5
		listH := msg.Height - headerH - footerH
		if listH < 4 {
			listH = 4
		}
		m.results.SetSize(msg.Width-4, listH)
		m.queueList.SetSize(msg.Width-4, listH)
		w := min(msg.Width-6, 60)
		m.search.Width = w
		m.jump.Width = min(msg.Width-6, 50)
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

		if m.jumpFocus {
			switch msg.String() {
			case "enter":
				raw := strings.TrimSpace(m.jump.Value())
				m.jump.SetValue("")
				m.blurInputs()
				if raw != "" {
					cmds = append(cmds, seekCmd(m.svc, raw))
				}
			case "esc":
				m.blurInputs()
			default:
				var cmd tea.Cmd
				m.jump, cmd = m.jump.Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		if m.searchFocus {
			switch msg.String() {
			case "enter":
				q := strings.TrimSpace(m.search.Value())
				if q == "" {
					return m, nil
				}
				if strings.HasPrefix(q, ":") {
					if err := m.runCommand(q); err != nil {
						m.errMsg = err.Error()
					} else {
						m.errMsg = ""
					}
					m.search.SetValue("")
					return m, nil
				}
				m.loading = true
				m.errMsg = ""
				m.statusLine = "Searching…"
				cmds = append(cmds, searchCmd(m.svc, q))
			case "esc":
				m.blurInputs()
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
			m.blurInputs()
			m.searchFocus = true
			cmds = append(cmds, m.search.Focus())

		case "g":
			m.blurInputs()
			m.jumpFocus = true
			cmds = append(cmds, m.jump.Focus())

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
					if err := m.svc.AddToQueue(context.Background(), v); err != nil {
						m.errMsg = err.Error()
					} else {
						m.errMsg = ""
						m.statusLine = "Added to queue"
						m.refreshQueue()
					}
				}
			}

		case " ":
			_ = m.svc.TogglePause()
			m.syncPlayState()

		case "n":
			m.statusLine = "Next track…"
			cmds = append(cmds, nextCmd(m.svc))

		case "p":
			m.statusLine = "Previous track…"
			cmds = append(cmds, prevCmd(m.svc))

		case "h":
			m.shuffleOn = m.svc.ToggleShuffle()
			m.statusLine = "Shuffle " + onOff(m.shuffleOn)

		case "r":
			m.repeatMode = m.svc.CycleRepeat().Label()
			m.statusLine = "Repeat " + m.repeatMode

		case "left":
			_ = m.svc.SeekRelative(-10)

		case "right":
			_ = m.svc.SeekRelative(10)

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

	case SeekDoneMsg:
		if msg.Err != nil {
			m.errMsg = msg.Err.Error()
		} else {
			m.errMsg = ""
			m.statusLine = "Jumped to " + formatMMSS(msg.Seconds)
		}

	case PlayerTickMsg:
		st := m.svc.Status()
		m.pos = st.Position
		m.dur = st.Duration
		m.volume = st.Volume
		m.shuffleOn = st.Shuffle
		m.repeatMode = st.Repeat.Label()
		m.syncPlayStateFrom(st)
		if st.Title != "" && !m.loading {
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

	if m.showSplash {
		return m, tea.Batch(cmds...)
	}
	if !m.searchFocus && !m.jumpFocus {
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

func (m *Model) runCommand(q string) error {
	q = strings.TrimSpace(q)
	if !strings.HasPrefix(q, ":") {
		return nil
	}
	return m.execCommand(strings.TrimPrefix(q, ":"))
}

func (m *Model) execCommand(body string) error {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil
	}
	switch {
	case strings.HasPrefix(body, "jump "), strings.HasPrefix(body, "seek "):
		sec, err := ParseSeek(body)
		if err != nil {
			return err
		}
		return m.svc.Seek(sec)
	case body == "help":
		m.statusLine = "Commands: :jump m:ss · :seek · g · h shuffle · r repeat"
		return nil
	default:
		if _, err := ParseSeek(body); err == nil {
			sec, _ := ParseSeek(body)
			return m.svc.Seek(sec)
		}
		return nil
	}
}

func (m *Model) syncPlayState() {
	m.syncPlayStateFrom(m.svc.Status())
}

func (m *Model) syncPlayStateFrom(st player.Status) {
	switch {
	case !st.Playing:
		m.playState = "⏹"
	case st.Paused:
		m.playState = "⏸"
	default:
		m.playState = "▶"
	}
}

func onOff(v bool) string {
	if v {
		return "on"
	}
	return "off"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
