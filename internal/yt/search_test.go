package yt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSearchOutput(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "search_sample.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	videos, err := parseSearchOutput(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(videos) != 2 {
		t.Fatalf("got %d videos, want 2", len(videos))
	}
	if videos[0].ID != "dQw4w9WgXcQ" || videos[0].Title != "Never Gonna Give You Up" {
		t.Fatalf("unexpected first video: %+v", videos[0])
	}
	if videos[1].Channel != "Artist Two" {
		t.Fatalf("unexpected channel: %s", videos[1].Channel)
	}
}

func TestParseSearchOutputEmpty(t *testing.T) {
	_, err := parseSearchOutput([]byte("{}\n"))
	if err == nil {
		t.Fatal("expected error for empty results")
	}
}
