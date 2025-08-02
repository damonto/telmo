package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type AdminId []string

func (a *AdminId) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func (a *AdminId) String() string {
	return strings.Join(*a, ",")
}

func (a *AdminId) IDs() []int64 {
	var ids []int64
	for _, id := range *a {
		id, err := strconv.Atoi(id)
		if err != nil {
			slog.Error("Failed to convert admin id to int64", "id", id, "error", err)
			continue
		}
		ids = append(ids, int64(id))
	}
	return ids
}

type ModemName map[string]string

func (n *ModemName) String() string {
	var names []string
	for imei, name := range *n {
		names = append(names, fmt.Sprintf("%s:%s", imei, name))
	}
	return strings.Join(names, ",")
}

func (n *ModemName) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return errors.New("invalid format")
	}
	(*n)[parts[0]] = parts[1]
	return nil
}

type Config struct {
	BotToken   string
	AdminId    AdminId
	ModemName  ModemName
	Endpoint   string
	ForceAT    bool
	Slowdown   bool
	Compatible bool
	Verbose    bool
}

var (
	ErrBotTokenRequired = errors.New("bot token is required")
	ErrAdminIdRequired  = errors.New("admin id is required")
)

var C *Config

func Init() {
	C = new(Config)
	C.ModemName = make(ModemName)
}

func (c *Config) IsValid() error {
	if c.BotToken == "" {
		return ErrBotTokenRequired
	}
	if len(c.AdminId) == 0 {
		return ErrAdminIdRequired
	}
	return nil
}
