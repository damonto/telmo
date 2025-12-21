package at

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type CSIM struct {
	at   *AT
	data []byte
}

func NewCSIM(at *AT) ATCommand { return &CSIM{at: at} }

func (c *CSIM) Run(command []byte) ([]byte, error) {
	cmd := fmt.Sprintf("%X", command)
	cmd = fmt.Sprintf("AT+CSIM=%d,%q", len(cmd), cmd)
	slog.Debug("[AT] CSIM Sending", "command", cmd)
	response, err := c.at.Run(cmd)
	slog.Debug("[AT] CSIM Received", "response", response, "error", err)
	if err != nil {
		return nil, err
	}
	sw, err := c.sw(response)
	if err != nil {
		return nil, err
	}
	if sw[len(sw)-2] != 0x61 && sw[len(sw)-2] != 0x90 {
		return sw, fmt.Errorf("unexpected response: %X", sw)
	}
	if sw[len(sw)-2] == 0x61 {
		return c.read(sw[1:])
	}
	return sw, nil
}

func (c *CSIM) read(length []byte) ([]byte, error) {
	for {
		b, err := c.Run(append([]byte{0x00, 0xC0, 0x00, 0x00}, length...))
		if err != nil {
			return nil, err
		}
		c.data = append(c.data, b[:len(b)-2]...)
		if b[len(b)-2] == 0x90 {
			break
		}
	}
	return c.data, nil
}

func (c *CSIM) sw(sw string) ([]byte, error) {
	lastIdx := strings.LastIndex(sw, ",")
	if lastIdx == -1 {
		return nil, errors.New("invalid response")
	}
	return hex.DecodeString(sw[lastIdx+2 : len(sw)-1])
}
