package notify

import (
	"errors"
	"fmt"
)

type Message interface {
	fmt.Stringer
	Markdown() string
}

type Sender interface {
	Send(message Message) error
}

type SenderFunc func(message Message) error

func (f SenderFunc) Send(message Message) error {
	return f(message)
}

func Send(sender Sender, message Message) error {
	if sender == nil {
		return errors.New("notify sender is required")
	}
	return sender.Send(message)
}
