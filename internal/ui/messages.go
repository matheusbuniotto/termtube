package ui

import "github.com/matheusbuniotto/termtube/internal/yt"

type SearchDoneMsg struct {
	Results []yt.Video
	Err     error
}

type PlayDoneMsg struct {
	Err error
}

type AddQueueDoneMsg struct {
	Err error
}

type SeekDoneMsg struct {
	Seconds float64
	Err     error
}

type ResolveDoneMsg struct {
	Err error
}

type PlayerTickMsg struct {
	StatusTick
}

type StatusTick struct {
	Advance bool
	Err     error
}
