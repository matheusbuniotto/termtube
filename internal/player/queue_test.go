package player

import (
	"testing"

	"github.com/matheusbuniotto/termtube/internal/yt"
)

func TestQueueAddNextPrev(t *testing.T) {
	q := NewQueue()
	v1 := yt.Video{ID: "a", Title: "A"}
	v2 := yt.Video{ID: "b", Title: "B"}
	q.Add(Track{Video: v1})
	q.Add(Track{Video: v2})

	if q.Len() != 2 {
		t.Fatalf("len=%d", q.Len())
	}
	if q.Current() != 0 {
		t.Fatalf("current=%d", q.Current())
	}

	q.Next()
	if q.Current() != 1 {
		t.Fatalf("after next current=%d", q.Current())
	}
	q.Prev()
	if q.Current() != 0 {
		t.Fatalf("after prev current=%d", q.Current())
	}

	q.Remove(0)
	if q.Len() != 1 || q.Current() != 0 {
		t.Fatalf("after remove len=%d current=%d", q.Len(), q.Current())
	}
}

func TestQueueRepeatAll(t *testing.T) {
	q := NewQueue()
	q.Add(Track{Video: yt.Video{ID: "a"}})
	q.Add(Track{Video: yt.Video{ID: "b"}})
	q.SetCurrent(1)
	q.repeat = RepeatAll

	idx, ok := q.NextIndex()
	if !ok || idx != 0 {
		t.Fatalf("next index = %d ok=%v", idx, ok)
	}
}

func TestQueueIndexByVideoID(t *testing.T) {
	q := NewQueue()
	q.Add(Track{Video: yt.Video{ID: "xyz", Title: "T"}})
	if q.IndexByVideoID("xyz") != 0 {
		t.Fatal("expected index 0")
	}
	if q.IndexByVideoID("missing") != -1 {
		t.Fatal("expected -1")
	}
}
