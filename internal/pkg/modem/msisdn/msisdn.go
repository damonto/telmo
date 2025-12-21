package msisdn

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/damonto/sigmo/internal/pkg/cond"
	"github.com/damonto/sigmo/internal/pkg/modem/at"
)

var phoneRE = regexp.MustCompile(`^\+?[0-9]{1,15}$`)

type MSISDN struct {
	at     *at.AT
	runner Runner
}

func New(device string) (*MSISDN, error) {
	conn, err := at.Open(device)
	if err != nil {
		return nil, err
	}
	m := &MSISDN{at: conn}
	if err := m.selectRunner(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return m, nil
}

func (m *MSISDN) Close() error {
	return m.at.Close()
}

func (m *MSISDN) Update(name, number string) error {
	if !phoneRE.MatchString(number) {
		return errors.New("invalid phone number")
	}
	return m.update(strings.HasPrefix(number, "+"), name, number)
}

func (m *MSISDN) selectRunner() error {
	if m.at.Support("AT+CRSM=?") {
		m.runner = NewCRSM(m.at)
		return nil
	}
	if m.at.Support("AT+CSIM=?") {
		m.runner = NewCSIM(m.at)
		return nil
	}
	return errors.New("modem does not support updating MSISDN")
}

type Runner interface {
	Run(data []byte) error
	Select() ([]byte, error)
}

func (m *MSISDN) update(hasPrefix bool, name string, number string) error {
	n, err := m.recordLen()
	if err != nil {
		return err
	}
	if len(name) > n-14 {
		return errors.New("name is too long")
	}
	nb, err := m.encodeBCD(strings.TrimPrefix(number, "+"))
	if err != nil {
		return err
	}
	var data []byte
	data = append(data, m.padRight([]byte(name), n-14)...)
	data = append(data, append(
		[]byte{byte(len(nb) + 1), cond.If(hasPrefix, byte(0x91), byte(0x81))},
		m.padRight(nb, 12)...,
	)...)
	return m.runner.Run(data)
}

func (m *MSISDN) recordLen() (int, error) {
	b, err := m.runner.Select()
	if err != nil {
		return 0, err
	}
	data := m.findTag(b, 0x82)
	if data == nil {
		return 0, fmt.Errorf("unexpected response: %X", b)
	}
	return int(data[4])<<8 + int(data[5]), nil
}

func (m *MSISDN) encodeBCD(value string) ([]byte, error) {
	for _, r := range value {
		if (r < '0' || r > '9') && !(r == 'f' || r == 'F') {
			return nil, errors.New("invalid value")
		}
	}
	if len(value)%2 != 0 {
		value += "F"
	}
	id, err := hex.DecodeString(value)
	if err != nil {
		return nil, err
	}
	for index := range id {
		id[index] = id[index]>>4 | id[index]<<4
	}
	return id, nil
}

func (m *MSISDN) padRight(value []byte, length int) []byte {
	if len(value) >= length {
		return value
	}
	return append(value, bytes.Repeat([]byte{0xFF}, length-len(value))...)
}

func (m *MSISDN) findTag(bs []byte, tag byte) []byte {
	bs = bs[2:]
	for len(bs) > 0 {
		n := int(bs[1])
		if bs[0] == tag {
			return bs[:2+n]
		}
		bs = bs[2+n:]
	}
	return nil
}
