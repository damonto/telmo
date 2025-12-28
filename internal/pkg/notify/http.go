package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HTTP struct {
	client   *http.Client
	endpoint string
	headers  map[string]string
}

type httpMessage struct {
	To   []string `json:"to,omitempty"`
	Text string   `json:"text"`
}

func NewHTTP(endpoint string, headers map[string]string, client *http.Client) (*HTTP, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, errors.New("http endpoint is required")
	}
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parsing http endpoint: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, errors.New("http endpoint must include scheme and host")
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &HTTP{
		client:   client,
		endpoint: endpoint,
		headers:  headers,
	}, nil
}

func (h *HTTP) Send(to []string, text string) error {
	body, err := json.Marshal(httpMessage{
		To:   to,
		Text: text,
	})
	if err != nil {
		return fmt.Errorf("encoding http message: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range h.headers {
		req.Header.Set(key, value)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending http message: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("http response status %s: %s", resp.Status, strings.TrimSpace(string(payload)))
	}
	return nil
}
