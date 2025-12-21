package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration
type Config struct {
	App      App                `toml:"app"`
	Database Database           `toml:"database"`
	Channel  map[string]Channel `toml:"channel"`
	Modem    map[string]Modem   `toml:"modem"`
}

type App struct {
	Environment   string `toml:"environment"`
	ListenAddress string `toml:"listen_address"`
}

type Database struct {
	Path string `toml:"path"`
}

type Channel struct {
	BotToken string  `toml:"bot_token"`
	AdminID  []int64 `toml:"admin_id"`
}

type Modem struct {
	Alias      string `toml:"alias"`
	Compatible bool   `toml:"compatible"`
	MSS        int    `toml:"mss"`
}

// Load reads and parses the configuration from the given file path
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	var config Config
	if _, err := toml.Decode(string(data), &config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}
	return &config, nil
}
