package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/matheusbuniotto/termtube/internal/config"
	"github.com/matheusbuniotto/termtube/internal/player"
	"github.com/matheusbuniotto/termtube/internal/ui"
)

func main() {
	noSplash := flag.Bool("no-splash", false, "skip the ASCII splash screen")
	flag.Parse()

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

	// Kill any mpv left playing by a previous session before we start ours,
	// so two players can't overlap (e.g. after an ungraceful exit).
	player.ReapStale(cfg)

	m := ui.NewModel(cfg, ui.Options{NoSplash: *noSplash})
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
