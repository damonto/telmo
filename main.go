package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/damonto/sigmo/internal/pkg/config"
)

var (
	BuildVersion string
	configPath   string
)

func init() {
	flag.StringVar(&configPath, "config", "config.toml", "path to config file")
}

func main() {
	flag.Parse()
	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("unable to load config", "error", err)
		os.Exit(1)
	}
	if !cfg.IsProduction() {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}
