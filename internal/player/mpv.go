package player

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/matheusbuniotto/termtube/internal/config"
)

type MPV struct {
	cfg     config.Config
	cmd     *exec.Cmd
	mu      sync.Mutex
	paused  bool
	volume  float64
}

func NewMPV(cfg config.Config) *MPV {
	return &MPV{cfg: cfg, volume: 100}
}

func (m *MPV) SocketPath() string {
	return m.cfg.IPCSocket
}

func (m *MPV) PlayURL(streamURL string) error {
	m.Stop()
	_ = os.Remove(m.cfg.IPCSocket)

	args := []string{
		"--no-video",
		"--really-quiet",
		"--keep-open=no",
		fmt.Sprintf("--input-ipc-server=%s", m.cfg.IPCSocket),
		"--volume=100",
		streamURL,
	}

	cmd := exec.Command(m.cfg.MpvPath, args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start mpv: %w", err)
	}

	// Wait for IPC socket
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(m.cfg.IPCSocket); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	m.mu.Lock()
	m.cmd = cmd
	m.paused = false
	m.mu.Unlock()
	return nil
}

func (m *MPV) Stop() {
	m.mu.Lock()
	cmd := m.cmd
	m.cmd = nil
	m.mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Signal(syscall.SIGTERM)
		_ = cmd.Wait()
	}
	_ = os.Remove(m.cfg.IPCSocket)
}

func (m *MPV) Running() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cmd != nil
}

func (m *MPV) command(args ...interface{}) error {
	payload := map[string]interface{}{"command": args}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	conn, err := net.DialTimeout("unix", m.cfg.IPCSocket, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write(data); err != nil {
		return err
	}
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		var resp struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(scanner.Bytes(), &resp)
		if resp.Error != "" && resp.Error != "success" {
			return fmt.Errorf("mpv: %s", resp.Error)
		}
	}
	return nil
}

func (m *MPV) getProperty(name string) (float64, error) {
	payload := map[string]interface{}{
		"command": []interface{}{"get_property", name},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	data = append(data, '\n')

	conn, err := net.DialTimeout("unix", m.cfg.IPCSocket, 2*time.Second)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	if _, err := conn.Write(data); err != nil {
		return 0, err
	}
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return 0, fmt.Errorf("no mpv response")
	}
	var resp struct {
		Error string      `json:"error"`
		Data  interface{} `json:"data"`
	}
	if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
		return 0, err
	}
	if resp.Error != "" && resp.Error != "success" {
		return 0, fmt.Errorf("mpv: %s", resp.Error)
	}
	switch v := resp.Data.(type) {
	case float64:
		return v, nil
	case nil:
		return 0, nil
	default:
		return 0, nil
	}
}

func (m *MPV) TogglePause() error {
	m.mu.Lock()
	m.paused = !m.paused
	paused := m.paused
	m.mu.Unlock()
	return m.command("set_property", "pause", paused)
}

func (m *MPV) SetPause(paused bool) error {
	m.mu.Lock()
	m.paused = paused
	m.mu.Unlock()
	return m.command("set_property", "pause", paused)
}

func (m *MPV) IsPaused() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.paused
}

func (m *MPV) VolumeUp() error {
	m.mu.Lock()
	m.volume = min(100, m.volume+10)
	v := m.volume
	m.mu.Unlock()
	return m.command("set_property", "volume", v)
}

func (m *MPV) VolumeDown() error {
	m.mu.Lock()
	m.volume = max(0, m.volume-10)
	v := m.volume
	m.mu.Unlock()
	return m.command("set_property", "volume", v)
}

func (m *MPV) Volume() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.volume
}

func (m *MPV) Seek(seconds float64) error {
	if seconds < 0 {
		seconds = 0
	}
	return m.command("seek", seconds, "absolute")
}

func (m *MPV) SeekRelative(delta float64) error {
	return m.command("seek", delta, "relative")
}

func (m *MPV) EOF() bool {
	if !m.Running() {
		return true
	}
	v, err := m.getProperty("eof-reached")
	if err != nil {
		return false
	}
	return v > 0
}
