package player

import (
	"math/rand"
)

type Queue struct {
	tracks     []Track
	current    int // -1 when nothing selected
	shuffle    bool
	repeat     RepeatMode
}

func NewQueue() *Queue {
	return &Queue{current: -1}
}

func (q *Queue) Len() int {
	return len(q.tracks)
}

func (q *Queue) Tracks() []Track {
	return q.tracks
}

func (q *Queue) Current() int {
	return q.current
}

func (q *Queue) CurrentTrack() (Track, bool) {
	if q.current < 0 || q.current >= len(q.tracks) {
		return Track{}, false
	}
	return q.tracks[q.current], true
}

func (q *Queue) Shuffle() bool {
	return q.shuffle
}

func (q *Queue) ToggleShuffle() bool {
	q.shuffle = !q.shuffle
	return q.shuffle
}

func (q *Queue) Repeat() RepeatMode {
	return q.repeat
}

func (q *Queue) CycleRepeat() RepeatMode {
	q.repeat = q.repeat.Cycle()
	return q.repeat
}

func (q *Queue) IndexByVideoID(id string) int {
	for i, t := range q.tracks {
		if t.Video.ID == id {
			return i
		}
	}
	return -1
}

func (q *Queue) Add(t Track) {
	q.tracks = append(q.tracks, t)
	if q.current < 0 {
		q.current = 0
	}
}

func (q *Queue) SetTrack(index int, t Track) {
	if index >= 0 && index < len(q.tracks) {
		q.tracks[index] = t
	}
}

func (q *Queue) Remove(index int) {
	if index < 0 || index >= len(q.tracks) {
		return
	}
	q.tracks = append(q.tracks[:index], q.tracks[index+1:]...)
	switch {
	case len(q.tracks) == 0:
		q.current = -1
	case q.current > index:
		q.current--
	case q.current == index:
		if q.current >= len(q.tracks) {
			q.current = len(q.tracks) - 1
		}
	case q.current >= len(q.tracks):
		q.current = len(q.tracks) - 1
	}
}

func (q *Queue) SetCurrent(index int) bool {
	if index < 0 || index >= len(q.tracks) {
		return false
	}
	q.current = index
	return true
}

func (q *Queue) NextIndex() (int, bool) {
	if len(q.tracks) == 0 {
		return -1, false
	}
	if q.repeat == RepeatOne {
		return q.current, true
	}
	if q.shuffle && len(q.tracks) > 1 {
		return q.randomIndexExcept(q.current), true
	}
	if q.current < len(q.tracks)-1 {
		return q.current + 1, true
	}
	if q.repeat == RepeatAll {
		return 0, true
	}
	return q.current, false
}

func (q *Queue) PrevIndex() (int, bool) {
	if len(q.tracks) == 0 {
		return -1, false
	}
	if q.repeat == RepeatOne {
		return q.current, true
	}
	if q.shuffle && len(q.tracks) > 1 {
		return q.randomIndexExcept(q.current), true
	}
	if q.current > 0 {
		return q.current - 1, true
	}
	if q.repeat == RepeatAll {
		return len(q.tracks) - 1, true
	}
	return q.current, false
}

func (q *Queue) Next() (Track, bool) {
	idx, ok := q.NextIndex()
	if !ok {
		if q.repeat == RepeatOne {
			return q.CurrentTrack()
		}
		return Track{}, false
	}
	q.current = idx
	return q.CurrentTrack()
}

func (q *Queue) Prev() (Track, bool) {
	idx, ok := q.PrevIndex()
	if !ok {
		if q.repeat == RepeatOne {
			return q.CurrentTrack()
		}
		return Track{}, false
	}
	q.current = idx
	return q.CurrentTrack()
}

func (q *Queue) At(index int) (Track, bool) {
	if index < 0 || index >= len(q.tracks) {
		return Track{}, false
	}
	return q.tracks[index], true
}

func (q *Queue) randomIndexExcept(except int) int {
	if len(q.tracks) <= 1 {
		return except
	}
	for {
		i := rand.Intn(len(q.tracks))
		if i != except {
			return i
		}
	}
}
