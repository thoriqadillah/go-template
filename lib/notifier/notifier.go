package notifier

import (
	"app/lib/log"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

var logger = log.Logger()

type NotifierName string

type Data map[string]interface{}

type Message struct {
	Subject  string
	Text     string
	Template string
	Data     Data
	From     string
	To       []string
	Bcc      []string
	Cc       []string
}

type option struct {
	river *river.Client[pgx.Tx]
}

type Option func(o *option)

type Factory func(o *option) Notifier

func WithRiverQueue(river *river.Client[pgx.Tx]) Option {
	return func(o *option) {
		o.river = river
	}
}

var providers = map[NotifierName]Factory{}

func register(name NotifierName, impl Factory) {
	providers[name] = impl
}

type Notifier interface {
	Send(s Message) error
}

func New(name NotifierName, options ...Option) Notifier {
	opt := &option{}
	for _, option := range options {
		option(opt)
	}

	provider, ok := providers[name]
	if !ok {
		logger.Fatal(fmt.Sprintf("Notifier with %s provider not found", name))
		return nil
	}

	return provider(opt)
}
