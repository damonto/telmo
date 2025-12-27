package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration
type Config struct {
	App      App                `toml:"app"`
	Database Database           `toml:"database"`
	Channels map[string]Channel `toml:"channels"`
	Modems   map[string]Modem   `toml:"modems"`
	Path     string             `toml:"-"`
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
	config.Path = path
	return &config, nil
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *Config) FindModem(id string) Modem {
	if modem, ok := c.Modems[id]; ok {
		return modem
	}
	return Modem{
		Compatible: false,
		MSS:        240,
	}
}

func (c *Config) Save() error {
	if c.Path == "" {
		return errors.New("config path is required")
	}
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(c); err != nil {
		return fmt.Errorf("encoding config file: %w", err)
	}
	if err := os.WriteFile(c.Path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}
	return nil
}
