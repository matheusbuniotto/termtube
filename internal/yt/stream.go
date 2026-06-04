package yt

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) AudioURL(ctx context.Context, videoID string) (string, error) {
	videoID = strings.TrimSpace(videoID)
	if videoID == "" {
		return "", fmt.Errorf("empty video id")
	}
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	out, err := c.run(ctx,
		"-f", "bestaudio[ext=m4a]/bestaudio/best",
		"-g",
		"--no-warnings",
		"--no-download",
		url,
	)
	if err != nil {
		return "", fmt.Errorf("audio url: %w", err)
	}
	streamURL := strings.TrimSpace(string(out))
	if streamURL == "" {
		return "", fmt.Errorf("no stream url returned")
	}
	// yt-dlp may return multiple lines; use first URL
	if idx := strings.Index(streamURL, "\n"); idx >= 0 {
		streamURL = streamURL[:idx]
	}
	return streamURL, nil
}
