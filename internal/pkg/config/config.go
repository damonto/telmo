package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type AdminId []int64

func (a *AdminId) Set(value string) error {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*a = append(*a, id)
	return nil
}

func (a *AdminId) String() string {
	var s string
	for _, id := range *a {
		s += fmt.Sprintf("%d,", id)
	}
	if len(s) == 0 {
		return ""
	}
	return s[:len(s)-1]
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
