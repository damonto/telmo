package cond

func When(condition bool, f func() error) error {
	if condition {
		return f()
	}
	return nil
}
