package notify

import (
	"fmt"
	"strings"
	"time"
)

type TextMessage struct {
	Text string `json:"text"`
}

func (m TextMessage) String() string {
	return m.Text
}

func (m TextMessage) Markdown() string {
	return escapeMarkdownV2(m.Text)
}

type SMSMessage struct {
	Modem    string    `json:"modem"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Time     time.Time `json:"timestamp,omitempty"`
	Text     string    `json:"text"`
	Incoming bool      `json:"incoming"`
}

func (m SMSMessage) String() string {
	return fmt.Sprintf(
		"SMS received\nModem: %s\nFrom: %s\nTo: %s\nTime: %s\n\n%s",
		m.Modem,
		m.From,
		m.To,
		m.displayTimestamp(),
		m.displayText(),
	)
}

func (m SMSMessage) Markdown() string {
	return fmt.Sprintf(
		"*Modem:* %s\n*From:* %s\n*To:* %s\n*Time:* %s\n\n%s",
		escapeMarkdownV2(m.Modem),
		escapeMarkdownV2(m.From),
		escapeMarkdownV2(m.To),
		escapeMarkdownV2(m.displayTimestamp()),
		escapeMarkdownV2(m.displayText()),
	)
}

func (m SMSMessage) displayText() string {
	text := strings.TrimSpace(m.Text)
	if text == "" {
		return "(empty message)"
	}
	return text
}

func (m SMSMessage) displayTimestamp() string {
	if m.Time.IsZero() {
		return "unknown"
	}
	return m.Time.Format(time.RFC3339)
}

var markdownV2Escaper = strings.NewReplacer(
	"\\", "\\\\",
	"_", "\\_",
	"*", "\\*",
	"[", "\\[",
	"]", "\\]",
	"(", "\\(",
	")", "\\)",
	"~", "\\~",
	"`", "\\`",
	">", "\\>",
	"#", "\\#",
	"+", "\\+",
	"-", "\\-",
	"=", "\\=",
	"|", "\\|",
	"{", "\\{",
	"}", "\\}",
	".", "\\.",
	"!", "\\!",
)

func escapeMarkdownV2(text string) string {
	return markdownV2Escaper.Replace(text)
}
