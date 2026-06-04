package player

import (
	"context"

	"github.com/monkmode/ytune/internal/config"
	"github.com/monkmode/ytune/internal/yt"
)

type Status struct {
	Title    string
	Position float64
	Duration float64
	Paused   bool
	Volume   float64
	Playing  bool
	EOF      bool
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

func (s *Service) PlayVideo(ctx context.Context, v yt.Video) error {
	url, err := s.yt.AudioURL(ctx, v.ID)
	if err != nil {
		return err
	}
	track := Track{Video: v, StreamURL: url}
	s.queue.Add(track)
	idx := s.queue.Len() - 1
	s.queue.SetCurrent(idx)
	return s.mpv.PlayURL(url)
}

func (s *Service) PlayTrackAt(ctx context.Context, index int) error {
	t, ok := s.queue.At(index)
	if !ok {
		return nil
	}
	if t.StreamURL == "" {
		url, err := s.yt.AudioURL(ctx, t.Video.ID)
		if err != nil {
			return err
		}
		t.StreamURL = url
		s.queue.tracks[index] = t
	}
	s.queue.SetCurrent(index)
	return s.mpv.PlayURL(t.StreamURL)
}

func (s *Service) AddToQueue(ctx context.Context, v yt.Video) error {
	url, err := s.yt.AudioURL(ctx, v.ID)
	if err != nil {
		return err
	}
	s.queue.Add(Track{Video: v, StreamURL: url})
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

func (s *Service) TogglePause() error {
	return s.mpv.TogglePause()
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
	if s.queue.Current() >= s.queue.Len()-1 {
		return nil
	}
	return s.Next(ctx)
}
