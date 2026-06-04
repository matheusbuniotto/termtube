# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`termtube` is a terminal YouTube music player (Bubble Tea TUI) that shells out to external binaries rather than using any API: `yt-dlp` for search and stream-URL resolution, `mpv` for audio playback. No API keys. Module path is `github.com/matheusbuniotto/termtube`; the binary is `termtube`.

## Commands

```bash
make build          # go build -o bin/termtube ./cmd/termtube
make run            # build + run ./bin/termtube
make install        # go install ./cmd/termtube
make test           # go test ./...
go test ./internal/player/ -run TestQueue   # single package / single test
```

Runtime requires `mpv` and `yt-dlp` on PATH (`brew install mpv yt-dlp`). `main.go` calls `cfg.CheckDeps()` on startup and exits if either is missing.

## Architecture

Three layers, strictly one-directional (`ui` → `player` → `yt`/`mpv`):

- **`internal/ui`** — Bubble Tea Model/Update/View. `Model` (model.go) holds all TUI state; `update.go` handles input and async message results; `view.go` renders. The UI never touches `yt-dlp` or `mpv` directly — it only calls `player.Service`.
- **`internal/player`** — `Service` (status.go) is the single orchestration facade the UI talks to. It owns the `Queue`, the `MPV` controller, and the `yt.Client`, and composes them (resolve stream URL → play → advance). `Queue` (queue.go) is pure in-memory state (tracks, current index, shuffle, repeat) with no I/O.
- **`internal/yt`** — thin `yt-dlp` wrapper. `Client.run` (client.go) is the single subprocess entry point with a 45s context timeout; `Search` uses `ytsearch%d:` + `--flat-playlist -j`, `AudioURL` uses `-g` to get a direct stream URL.
- **`internal/mpv`** (`player/mpv.go`) — spawns headless `mpv --no-video --input-ipc-server=<socket>` and controls it over a Unix-domain JSON IPC socket (`command`/`getProperty`). Each command opens a fresh connection. Playback state queried via properties: `time-pos`, `duration`, `eof-reached`.

### Two concurrency patterns to respect

1. **Blocking work runs as `tea.Cmd`s, never inline in Update.** Search, play, seek, next/prev each have a `*Cmd` function (bottom of model.go) that runs the slow `Service` call in a goroutine and returns a `*DoneMsg` (messages.go). `Update` only mutates `Model` in response to those messages. Adding a new action that calls `yt-dlp`/`mpv` means: add a Msg, add a Cmd, handle the key in `update.go`, handle the Msg in `update.go`.

2. **Auto-advance is poll-driven.** `tickCmd` fires every 100ms and calls `Service.MaybeAdvance`, which checks `mpv`'s `eof-reached` and advances the queue (honoring repeat-one/repeat-all). There is no event from mpv — end-of-track is detected by polling. The same tick refreshes the status line (position/duration/volume).

### Stream URL caching

`Track.StreamURL` is resolved lazily and cached on the queue entry (`Service.resolveURL`). Queue-add does *not* resolve; resolution happens on first play. yt-dlp stream URLs are time-limited, so a long-queued track re-resolves naturally only if cleared — keep this in mind if playback of an old queued track fails.

## Config

`internal/config` loads `~/.config/termtube/config.yaml` (optional; falls back to `Default()`), and ensures `~/.cache/termtube/` exists for the mpv IPC socket. Fields: `yt_dlp_path`, `mpv_path`, `ipc_socket`, `search_limit`.

## Tests

Table/unit tests live next to code (`queue_test.go`, `seek_test.go`, `search_test.go`). `internal/yt/testdata/search_sample.jsonl` is captured `yt-dlp -j` output used to test `parseSearchOutput` without invoking the network. Add fixtures there rather than calling yt-dlp in tests.
