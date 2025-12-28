package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const defaultTelegramEndpoint = "https://api.telegram.org"

type Telegram struct {
	client   *http.Client
	baseURL  url.URL
	botToken string
}

type telegramMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func NewTelegram(endpoint string, botToken string, client *http.Client) (*Telegram, error) {
	if strings.TrimSpace(botToken) == "" {
		return nil, errors.New("telegram bot token is required")
	}
	if strings.TrimSpace(endpoint) == "" {
		endpoint = defaultTelegramEndpoint
	}
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parsing telegram endpoint: %w", err)
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, errors.New("telegram endpoint must include scheme and host")
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Telegram{
		client:   client,
		baseURL:  *baseURL,
		botToken: botToken,
	}, nil
}

func (t *Telegram) Send(to []int64, text string) error {
	if len(to) == 0 {
		return errors.New("telegram recipients are required")
	}
	var combined error
	for _, recipient := range to {
		if err := t.sendOne(recipient, text); err != nil {
			combined = errors.Join(combined, err)
		}
	}
	return combined
}

func (t *Telegram) sendOne(to int64, text string) error {
	message := telegramMessage{
		ChatID: to,
		Text:   text,
	}
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("encoding telegram message: %w", err)
	}
	endpoint := t.baseURL
	endpoint.Path = path.Join(endpoint.Path, "bot"+t.botToken, "sendMessage")
	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending telegram message to %d: %w", to, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("telegram response status %s: %s", resp.Status, strings.TrimSpace(string(payload)))
	}
	return nil
}
