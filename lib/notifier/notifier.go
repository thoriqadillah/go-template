package notifier

import (
	"app/env"
	"log"
)

type option struct {
	name        string
	host        string
	username    string
	password    string
	port        int
	secure      bool
	templateDir string
}

type Option func(*option)

func WithHost(host string) Option {
	return func(o *option) {
		o.host = host
	}
}

func WithUsername(username string) Option {
	return func(o *option) {
		o.username = username
	}
}

func WithPassword(password string) Option {
	return func(o *option) {
		o.password = password
	}
}

func WithPort(port int) Option {
	return func(o *option) {
		o.port = port
	}
}

func WithSecure(secure bool) Option {
	return func(o *option) {
		o.secure = secure
	}
}

type Message struct {
	Subject  string
	Text     string
	Template string
	From     string
	To       []string
	Bcc      []string
	Cc       []string
}

type Factory func(option *option) Notifier

var providers = map[string]Factory{}

func register(name string, impl Factory) {
	providers[name] = impl
}

type Notifier interface {
	Send(s Message) error
}

func New(name string, opts ...Option) Notifier {
	opt := &option{
		name:        env.AppName,
		host:        env.EmailHost,
		port:        env.EmailPort,
		secure:      env.EmailSecure,
		username:    env.EmailUsername,
		password:    env.EmailPassword,
		templateDir: "", // TODO: embed the templates
	}

	for _, option := range opts {
		option(opt)
	}

	provider, ok := providers[name]
	if !ok {
		log.Fatalf("Notifier provider %s not found", name)
		return nil
	}

	return provider(opt)
}
