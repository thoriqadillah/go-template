package notifier

import (
	"app/env"
	"log"
)

type option struct {
	name        string
	authMethod  string
	host        string
	username    string
	password    string
	port        int
	secure      bool
	templateDir string
}

type Send struct {
	Subject string
	Message string
	Html    string
	From    string
	To      []string
	Bcc     []string
	Cc      []string
}

type Factory func(option *option) Notifier

var providers = map[string]Factory{}

func register(name string, impl Factory) {
	providers[name] = impl
}

type Notifier interface {
	Send(s Send) error
}

func New(name string) Notifier {
	option := &option{
		name:        env.AppName,
		host:        env.EmailHost,
		port:        env.EmailPort,
		secure:      env.EmailSecure,
		username:    env.EmailUsername,
		password:    env.EmailPassword,
		templateDir: "", // TODO: embed the templates
	}

	provider, ok := providers[name]
	if !ok {
		log.Fatalf("Notifier provider %s not found", name)
		return nil
	}

	return provider(option)
}
