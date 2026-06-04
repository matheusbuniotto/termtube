package player

import (
	"context"

	"github.com/monkmode/ytune/internal/config"
	"github.com/monkmode/ytune/internal/yt"
)

type Status struct {
	Title      string
	Position   float64
	Duration   float64
	Paused     bool
	Volume     float64
	Playing    bool
	EOF        bool
	Shuffle    bool
	Repeat     RepeatMode
}

type Service struct {
	cfg   config.Config
	yt    *yt.Client
	mpv   *MPV
	queue *Queue
}

func NewService(cfg config.Config) *Service {
	return &Service{
		cfg:   cfg,
		yt:    yt.NewClient(cfg),
		mpv:   NewMPV(cfg),
		queue: NewQueue(),
	}
}

func (s *Service) YT() *yt.Client {
	return s.yt
}

func (s *Service) Queue() *Queue {
	return s.queue
}

func (s *Service) Stop() {
	s.mpv.Stop()
}

func (s *Service) resolveURL(ctx context.Context, index int) (string, error) {
	t, ok := s.queue.At(index)
	if !ok {
		return "", nil
	}
	if t.StreamURL != "" {
		return t.StreamURL, nil
	}
	url, err := s.yt.AudioURL(ctx, t.Video.ID)
	if err != nil {
		return "", err
	}
	t.StreamURL = url
	s.queue.SetTrack(index, t)
	return url, nil
}

func (s *Service) PlayNow(ctx context.Context, v yt.Video) error {
	if idx := s.queue.IndexByVideoID(v.ID); idx >= 0 {
		s.queue.SetCurrent(idx)
		url, err := s.resolveURL(ctx, idx)
		if err != nil {
			return err
		}
		if url == "" {
			return nil
		}
		return s.mpv.PlayURL(url)
	}

	url, err := s.yt.AudioURL(ctx, v.ID)
	if err != nil {
		return err
	}
	track := Track{Video: v, StreamURL: url}
	s.queue.Add(track)
	s.queue.SetCurrent(s.queue.Len() - 1)
	return s.mpv.PlayURL(url)
}

func (s *Service) PlayVideo(ctx context.Context, v yt.Video) error {
	return s.PlayNow(ctx, v)
}

func (s *Service) PlayTrackAt(ctx context.Context, index int) error {
	url, err := s.resolveURL(ctx, index)
	if err != nil {
		return err
	}
	if url == "" {
		return nil
	}
	s.queue.SetCurrent(index)
	return s.mpv.PlayURL(url)
}

func (s *Service) AddToQueue(ctx context.Context, v yt.Video) error {
	if s.queue.IndexByVideoID(v.ID) >= 0 {
		return nil
	}
	s.queue.Add(Track{Video: v})
	return nil
}

func (s *Service) Next(ctx context.Context) error {
	if _, ok := s.queue.Next(); !ok {
		return nil
	}
	return s.PlayTrackAt(ctx, s.queue.Current())
}

func (s *Service) Prev(ctx context.Context) error {
	if _, ok := s.queue.Prev(); !ok {
		return nil
	}
	return s.PlayTrackAt(ctx, s.queue.Current())
}

func (s *Service) ToggleShuffle() bool {
	return s.queue.ToggleShuffle()
}

func (s *Service) CycleRepeat() RepeatMode {
	return s.queue.CycleRepeat()
}

func (s *Service) TogglePause() error {
	return s.mpv.TogglePause()
}

func (s *Service) Seek(seconds float64) error {
	if !s.mpv.Running() {
		return nil
	}
	return s.mpv.Seek(seconds)
}

func (s *Service) SeekRelative(delta float64) error {
	if !s.mpv.Running() {
		return nil
	}
	return s.mpv.SeekRelative(delta)
}

func (s *Service) VolumeUp() error {
	return s.mpv.VolumeUp()
}

func (s *Service) VolumeDown() error {
	return s.mpv.VolumeDown()
}

func (s *Service) RemoveFromQueue(index int) {
	s.queue.Remove(index)
}

func (s *Service) Status() Status {
	st := Status{
		Volume:  s.mpv.Volume(),
		Playing: s.mpv.Running(),
		Paused:  s.mpv.IsPaused(),
		EOF:     s.mpv.EOF(),
		Shuffle: s.queue.Shuffle(),
		Repeat:  s.queue.Repeat(),
	}
	if t, ok := s.queue.CurrentTrack(); ok {
		st.Title = t.Video.Title
	}
	if s.mpv.Running() {
		if pos, err := s.mpv.getProperty("time-pos"); err == nil {
			st.Position = pos
		}
		if dur, err := s.mpv.getProperty("duration"); err == nil {
			st.Duration = dur
		}
	}
	return st
}

func (s *Service) MaybeAdvance(ctx context.Context) error {
	if !s.mpv.Running() {
		return nil
	}
	if !s.mpv.EOF() {
		return nil
	}

	if s.queue.Repeat() == RepeatOne {
		return s.PlayTrackAt(ctx, s.queue.Current())
	}

	if s.queue.Current() >= s.queue.Len()-1 && s.queue.Repeat() != RepeatAll {
		return nil
	}
	return s.Next(ctx)
}
