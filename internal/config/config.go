package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	YtDlpPath   string `yaml:"yt_dlp_path"`
	MpvPath     string `yaml:"mpv_path"`
	IPCSocket   string `yaml:"ipc_socket"`
	SearchLimit int    `yaml:"search_limit"`
}

func Default() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	cache := filepath.Join(home, ".cache", "termtube")
	return Config{
		YtDlpPath:   "yt-dlp",
		MpvPath:     "mpv",
		IPCSocket:   filepath.Join(cache, "mpv.sock"),
		SearchLimit: 10,
	}
}

func Load() (Config, error) {
	cfg := Default()
	path, err := configPath()
	if err != nil {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if cfg.SearchLimit <= 0 {
		cfg.SearchLimit = 10
	}
	return cfg, nil
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "termtube", "config.yaml"), nil
}

func (c Config) EnsureCacheDir() error {
	return os.MkdirAll(filepath.Dir(c.IPCSocket), 0o755)
}

func (c Config) CheckDeps() error {
	if _, err := exec.LookPath(c.resolve(c.YtDlpPath)); err != nil {
		return err
	}
	if _, err := exec.LookPath(c.resolve(c.MpvPath)); err != nil {
		return err
	}
	return nil
}

func (c Config) resolve(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	return name
}
