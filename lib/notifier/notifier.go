package notifier

import (
	"app/lib/log"
	"fmt"
)

var logger = log.Logger()

type Message struct {
	Subject  string
	Text     string
	Template string
	From     string
	To       []string
	Bcc      []string
	Cc       []string
}

type Factory func() Notifier

var providers = map[string]Factory{}

func register(name string, impl Factory) {
	providers[name] = impl
}

type Notifier interface {
	Send(s Message) error
}

func New(name string) Notifier {
	provider, ok := providers[name]
	if !ok {
		logger.Fatal(fmt.Sprintf("Notifier provider %s not found", name))
		return nil
	}

	return provider()
}
