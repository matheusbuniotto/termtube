package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/monkmode/ytune/internal/config"
	"github.com/monkmode/ytune/internal/ui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.CheckDeps(); err != nil {
		fmt.Fprintf(os.Stderr, "missing dependency: %v\n\nInstall:\n  brew install mpv yt-dlp\n", err)
		os.Exit(1)
	}

	if err := cfg.EnsureCacheDir(); err != nil {
		fmt.Fprintf(os.Stderr, "cache dir: %v\n", err)
		os.Exit(1)
	}

	m := ui.NewModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
