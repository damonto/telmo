package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Recipient string

type Recipients []Recipient

func (r Recipients) Int64s() ([]int64, error) {
	ids := make([]int64, 0, len(r))
	for i, raw := range r {
		value := strings.TrimSpace(string(raw))
		if value == "" {
			return nil, fmt.Errorf("recipient %d is empty", i)
		}
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("recipient %d: %w", i, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r Recipients) Strings() []string {
	values := make([]string, 0, len(r))
	for _, raw := range r {
		value := strings.TrimSpace(string(raw))
		if value == "" {
			continue
		}
		values = append(values, value)
	}
	return values
}

func (r *Recipients) UnmarshalTOML(data any) error {
	switch value := data.(type) {
	case []any:
		recipients := make([]Recipient, 0, len(value))
		for i, item := range value {
			recipient, err := parseRecipient(item)
			if err != nil {
				return fmt.Errorf("recipients[%d]: %w", i, err)
			}
			recipients = append(recipients, recipient)
		}
		*r = recipients
		return nil
	case []string:
		recipients := make([]Recipient, 0, len(value))
		for i, item := range value {
			recipient, err := parseRecipient(item)
			if err != nil {
				return fmt.Errorf("recipients[%d]: %w", i, err)
			}
			recipients = append(recipients, recipient)
		}
		*r = recipients
		return nil
	case []int64:
		recipients := make([]Recipient, 0, len(value))
		for i, item := range value {
			recipient, err := parseRecipient(item)
			if err != nil {
				return fmt.Errorf("recipients[%d]: %w", i, err)
			}
			recipients = append(recipients, recipient)
		}
		*r = recipients
		return nil
	case string, int64, int:
		recipient, err := parseRecipient(value)
		if err != nil {
			return err
		}
		*r = []Recipient{recipient}
		return nil
	default:
		return fmt.Errorf("recipients must be strings or integers")
	}
}

func parseRecipient(value interface{}) (Recipient, error) {
	switch v := value.(type) {
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return "", errors.New("recipient cannot be empty")
		}
		return Recipient(trimmed), nil
	case int64:
		return Recipient(strconv.FormatInt(v, 10)), nil
	case int:
		return Recipient(strconv.FormatInt(int64(v), 10)), nil
	default:
		return "", errors.New("recipient must be a string or integer")
	}
}
