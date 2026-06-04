package yt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type searchEntry struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Channel         string  `json:"channel"`
	ChannelID       string  `json:"channel_id"`
	Uploader        string  `json:"uploader"`
	Duration        float64 `json:"duration"`
	DurationString  string  `json:"duration_string"`
}

func (c *Client) Search(ctx context.Context, query string) ([]Video, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("empty search query")
	}

	target := fmt.Sprintf("ytsearch%d:%s", c.limit, query)
	out, err := c.run(ctx,
		target,
		"--flat-playlist",
		"-j",
		"--no-warnings",
		"--no-download",
	)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	return parseSearchOutput(out)
}

func parseSearchOutput(out []byte) ([]Video, error) {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var videos []Video
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var e searchEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("parse search result: %w", err)
		}
		if e.ID == "" || e.Title == "" {
			continue
		}
		channel := e.Channel
		if channel == "" {
			channel = e.Uploader
		}
		videos = append(videos, Video{
			ID:       e.ID,
			Title:    e.Title,
			Channel:  channel,
			Duration: int(e.Duration),
		})
	}
	if len(videos) == 0 {
		return nil, fmt.Errorf("no results found")
	}
	return videos, nil
}
