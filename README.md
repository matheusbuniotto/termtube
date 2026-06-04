# termtube

Terminal music player for YTB. Search, queue tracks, and listen with **mpv** — no API keys.

## Requirements

- [Go](https://go.dev/) 1.22+
- [mpv](https://mpv.io/)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)

```bash
brew install mpv yt-dlp
```

## Install

```bash
git clone https://github.com/matheusbuniotto/termtube && cd termtube
make install
# or
go install github.com/matheusbuniotto/termtube/cmd/termtube@latest
```

## Usage

```bash
termtube
termtube --no-splash   # skip intro
```

On launch you'll see an ASCII splash with a search box. Start typing and press `Enter` to search right away, press `Esc` to skip into an empty player, or just wait ~2 seconds. After a search you land directly on the results list — use `j`/`k` or arrows to move, `Enter` to play, `a` to queue.

## Keybindings

| Key | Action |
|-----|--------|
| `/` | Focus search |
| `Enter` | Search (when input focused) or play selected |
| `a` | Add selection to queue |
| `Tab` | Switch between results and queue |
| `j`/`k`, arrows | Move selection |
| `Space` | Pause / resume |
| `n` / `p` | Next / previous track |
| `+` / `-` | Volume up / down |
| `d` | Remove from queue |
| `g` | Jump to time (enter `m:ss` or `m s`) |
| `:jump 1:30` | Jump via search bar command |
| `←` / `→` | Seek −10s / +10s |
| `h` | Toggle shuffle |
| `r` | Cycle repeat (off → all → one) |
| `q` | Quit |

## Config

Optional file at `~/.config/termtube/config.yaml`:

```yaml
yt_dlp_path: yt-dlp
mpv_path: mpv
ipc_socket: ~/.cache/termtube/mpv.sock
search_limit: 10
```

## Troubleshooting

- **yt-dlp not found** — install with Homebrew or set `yt_dlp_path` in config.
- **Playback fails** — run `brew upgrade yt-dlp`. Age-restricted videos may need cookies (not yet supported).
- **No sound** — ensure `mpv` is installed and not blocked by another instance using the same IPC socket.

## License

MIT
