package player

type Queue struct {
	tracks  []Track
	current int // -1 when nothing selected
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

func (q *Queue) Add(t Track) {
	q.tracks = append(q.tracks, t)
	if q.current < 0 {
		q.current = 0
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

func (q *Queue) Next() (Track, bool) {
	if len(q.tracks) == 0 {
		return Track{}, false
	}
	if q.current < len(q.tracks)-1 {
		q.current++
	}
	return q.CurrentTrack()
}

func (q *Queue) Prev() (Track, bool) {
	if len(q.tracks) == 0 {
		return Track{}, false
	}
	if q.current > 0 {
		q.current--
	}
	return q.CurrentTrack()
}

func (q *Queue) At(index int) (Track, bool) {
	if index < 0 || index >= len(q.tracks) {
		return Track{}, false
	}
	return q.tracks[index], true
}
