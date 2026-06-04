package player

import (
	"testing"

	"github.com/monkmode/ytune/internal/yt"
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
