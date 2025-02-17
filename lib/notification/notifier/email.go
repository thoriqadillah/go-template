package notifier

import (
	"app/env"
	"app/lib/logger"
	"embed"
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

var zap = logger.Logger()

//go:embed template
var template embed.FS

type emailer struct {
	mailer *gomail.Dialer
}

func createMailer(option *option) Notifier {
	mailer := gomail.NewDialer(
		option.host,
		option.port,
		option.username,
		option.password,
	)

	return &emailer{
		mailer: mailer,
	}
}

func (e *emailer) Send(s Send) error {
	msg := s.Message
	mimetype := "text/plain"

	if s.Template != "" {
		mimetype = "text/html"
		bytes, err := template.ReadFile(fmt.Sprintf("template/%s", s.Template))
		if err != nil {
			return err
		}

		msg = string(bytes)
	}

	if env.Dev {
		zap.Info(msg)
		return nil
	}

	message := gomail.NewMessage()
	message.SetHeader("From", s.From)
	message.SetHeader("To", s.To...)
	message.SetHeader("Subject", s.Subject)
	message.SetHeader("Bcc", s.Bcc...)
	message.SetHeader("Cc", s.Cc...)
	message.SetBody(mimetype, msg)

	return e.mailer.DialAndSend(message)
}

func init() {
	register("email", createMailer)
}
