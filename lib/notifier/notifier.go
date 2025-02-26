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
	Subject string
	// Body can be a message string or template name.
	// Example:
	// 1. Hello world -> this will set the body type to be text/plain and pass it as is
	// 2. verify.html -> this will look into template folder and set the body type to be text/html and pass the html file as a template
	Body string
	Data Data
	From string
	To   []string
	Bcc  []string
	Cc   []string
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

// register will register the implementation of the notifier as the provider
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
