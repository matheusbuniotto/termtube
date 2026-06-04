package yt

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/matheusbuniotto/termtube/internal/config"
)

const defaultTimeout = 45 * time.Second

type Client struct {
	bin    string
	limit  int
	timeout time.Duration
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		bin:     cfg.YtDlpPath,
		limit:   cfg.SearchLimit,
		timeout: defaultTimeout,
	}
}

func (c *Client) run(ctx context.Context, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.bin, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("%w: %s", err, bytes.TrimSpace(stderr.Bytes()))
		}
		return nil, err
	}
	return out, nil
}
