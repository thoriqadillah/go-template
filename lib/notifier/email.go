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

func (e *emailer) Send(m Message) error {
	// TODO: send email with background job
	msg := m.Text
	mimetype := "text/plain"

	if m.Template != "" {
		mimetype = "text/html"
		bytes, err := template.ReadFile(fmt.Sprintf("template/%s", m.Template))
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
	message.SetHeader("From", m.From)
	message.SetHeader("To", m.To...)
	message.SetHeader("Subject", m.Subject)
	message.SetHeader("Bcc", m.Bcc...)
	message.SetHeader("Cc", m.Cc...)
	message.SetBody(mimetype, msg)

	return e.mailer.DialAndSend(message)
}

func init() {
	register("email", createMailer)
}
