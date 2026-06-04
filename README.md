# ytune

Terminal music player for YouTube. Search, queue tracks, and listen with **mpv** — no API keys.

## Requirements

- [Go](https://go.dev/) 1.22+
- [mpv](https://mpv.io/)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)

```bash
brew install mpv yt-dlp
```

## Install

```bash
git clone <repo-url> && cd CLI
make install
# or
go install ./cmd/ytune
```

## Usage

```bash
ytune
ytune --no-splash   # skip intro
```

On launch you'll see an ASCII splash screen — press any key (or wait ~2 seconds) to enter the player.

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

Optional file at `~/.config/ytune/config.yaml`:

```yaml
yt_dlp_path: yt-dlp
mpv_path: mpv
ipc_socket: ~/.cache/ytune/mpv.sock
search_limit: 10
```

## Troubleshooting

- **yt-dlp not found** — install with Homebrew or set `yt_dlp_path` in config.
- **Playback fails** — run `brew upgrade yt-dlp`. Age-restricted videos may need cookies (not yet supported).
- **No sound** — ensure `mpv` is installed and not blocked by another instance using the same IPC socket.

## License

MIT
