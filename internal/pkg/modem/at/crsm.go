package at

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
)

type CRSMInstruction uint16

const (
	CRSMReadBinary   CRSMInstruction = 0xB0
	CRSMReadRecord   CRSMInstruction = 0xB2
	CRSMGetResponse  CRSMInstruction = 0xC0
	CRSMUpdateBinary CRSMInstruction = 0xD6
	CRSMUpdateRecord CRSMInstruction = 0xDC
	CRSMStatus       CRSMInstruction = 0xF2
)

type CRSM struct{ at *AT }

func NewCRSM(at *AT) ATCommand { return &CRSM{at: at} }

type CRSMCommand struct {
	Instruction CRSMInstruction
	FileID      uint16
	P1          byte
	P2          byte
	Data        []byte
}

func (c CRSMCommand) Bytes() []byte {
	return fmt.Appendf(nil, "%d,%d,%d,%d,%d,\"%X\"", c.Instruction, c.FileID, c.P1, c.P2, len(c.Data), c.Data)
}

func (c *CRSM) Run(command []byte) ([]byte, error) {
	cmd := fmt.Sprintf("AT+CRSM=%s", command)
	slog.Debug("[AT] CRSM Sending", "command", cmd)
	response, err := c.at.Run(cmd)
	slog.Debug("[AT] CRSM Received", "response", response, "error", err)
	if err != nil {
		return nil, err
	}
	return c.sw(response)
}

func (c *CRSM) sw(sw string) ([]byte, error) {
	if !strings.Contains(sw, "+CRSM: 144") {
		return nil, fmt.Errorf("unexpected response: %s", sw)
	}
	data := strings.Replace(sw, "+CRSM: 144,0,", "", 1)
	return hex.DecodeString(data[1 : len(data)-1])
}
