package notify

type Sender interface {
	Send(text string) error
}

type SenderFunc func(text string) error

func (f SenderFunc) Send(text string) error {
	return f(text)
}
