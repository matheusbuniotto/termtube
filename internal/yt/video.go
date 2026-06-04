package yt

import "fmt"

type Video struct {
	ID       string
	Title    string
	Channel  string
	Duration int // seconds, 0 if unknown
}

func (v Video) URL() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID)
}
