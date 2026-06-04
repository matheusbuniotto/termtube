package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/monkmode/ytune/internal/config"
	"github.com/monkmode/ytune/internal/player"
	"github.com/monkmode/ytune/internal/yt"
)

type panel int

const (
	panelResults panel = iota
	panelQueue
)

type item struct {
	video yt.Video
	kind  string // "result" or "queue"
}

func (i item) Title() string       { return i.video.Title }
func (i item) Description() string { return formatDesc(i.video) }
func (i item) FilterValue() string { return i.video.Title + i.video.Channel }

type Model struct {
	cfg      config.Config
	svc      *player.Service
	width    int
	height   int

	panel      panel
	search     textinput.Model
	results    list.Model
	queueList  list.Model
	spin       spinner.Model
	progress   progress.Model
	loading    bool
	searchFocus bool

	resultsData []yt.Video
	statusLine  string
	errMsg      string

	playState string
	pos       float64
	dur       float64
	volume    float64

	showSplash bool
}

func NewModel(cfg config.Config) *Model {
	ti := textinput.New()
	ti.Placeholder = "Search YouTube…"
	ti.CharLimit = 120
	ti.Width = 50

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	prog := progress.New(progress.WithGradient("#89b4fa", "#cba6f7"))

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true

	results := list.New([]list.Item{}, delegate, 0, 0)
	results.Title = "Results"
	results.SetShowStatusBar(false)
	results.SetFilteringEnabled(false)

	queue := list.New([]list.Item{}, delegate, 0, 0)
	queue.Title = "Queue"
	queue.SetShowStatusBar(false)
	queue.SetFilteringEnabled(false)

	return &Model{
		cfg:        cfg,
		svc:        player.NewService(cfg),
		search:     ti,
		results:    results,
		queueList:  queue,
		spin:       sp,
		progress:   prog,
		playState:  "⏹",
		showSplash: true,
	}
}

func (m *Model) Init() tea.Cmd {
	if m.showSplash {
		return splashDoneCmd()
	}
	return m.appInit()
}

func (m *Model) appInit() tea.Cmd {
	return tea.Batch(
		m.spin.Tick,
		tickCmd(m.svc),
		textinput.Blink,
	)
}

func tickCmd(svc *player.Service) tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		err := svc.MaybeAdvance(context.Background())
		return PlayerTickMsg{StatusTick: StatusTick{Err: err}}
	})
}

func (m *Model) setResults(videos []yt.Video) {
	m.resultsData = videos
	items := make([]list.Item, len(videos))
	for i, v := range videos {
		items[i] = item{video: v, kind: "result"}
	}
	m.results.SetItems(items)
}

func (m *Model) refreshQueue() {
	tracks := m.svc.Queue().Tracks()
	items := make([]list.Item, len(tracks))
	for i, t := range tracks {
		items[i] = item{video: t.Video, kind: "queue"}
	}
	m.queueList.SetItems(items)
	cur := m.svc.Queue().Current()
	if cur >= 0 && cur < len(items) {
		m.queueList.Select(cur)
	}
}

func (m *Model) activeList() *list.Model {
	if m.panel == panelQueue {
		return &m.queueList
	}
	return &m.results
}

func (m *Model) selectedVideo() (yt.Video, bool) {
	l := m.activeList()
	it, ok := l.SelectedItem().(item)
	if !ok {
		return yt.Video{}, false
	}
	return it.video, true
}

func searchCmd(svc *player.Service, query string) tea.Cmd {
	return func() tea.Msg {
		results, err := svc.YT().Search(context.Background(), query)
		return SearchDoneMsg{Results: results, Err: err}
	}
}

func playCmd(svc *player.Service, v yt.Video) tea.Cmd {
	return func() tea.Msg {
		err := svc.PlayVideo(context.Background(), v)
		return PlayDoneMsg{Err: err}
	}
}

func addQueueCmd(svc *player.Service, v yt.Video) tea.Cmd {
	return func() tea.Msg {
		err := svc.AddToQueue(context.Background(), v)
		return AddQueueDoneMsg{Err: err}
	}
}

func nextCmd(svc *player.Service) tea.Cmd {
	return func() tea.Msg {
		err := svc.Next(context.Background())
		return PlayDoneMsg{Err: err}
	}
}

func prevCmd(svc *player.Service) tea.Cmd {
	return func() tea.Msg {
		err := svc.Prev(context.Background())
		return PlayDoneMsg{Err: err}
	}
}

func formatDesc(v yt.Video) string {
	ch := v.Channel
	if ch == "" {
		ch = "Unknown"
	}
	if v.Duration > 0 {
		return ch + " · " + formatDuration(v.Duration)
	}
	return ch
}

func formatDuration(sec int) string {
	if sec <= 0 {
		return "--:--"
	}
	return formatMMSS(float64(sec))
}

func formatMMSS(seconds float64) string {
	if seconds < 0 || seconds != seconds {
		return "0:00"
	}
	sec := int(seconds + 0.5)
	m := sec / 60
	s := sec % 60
	return fmt.Sprintf("%d:%02d", m, s)
}
