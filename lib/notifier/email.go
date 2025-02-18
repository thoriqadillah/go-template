package notifier

import (
	"app/env"
	"bytes"
	"embed"
	"log"
	"text/template"

	gomail "gopkg.in/gomail.v2"
)

//go:embed template
var templateFs embed.FS
var templates *template.Template

type emailer struct {
	mailer *gomail.Dialer
}

func createMailer() Notifier {
	mailer := gomail.NewDialer(
		env.EmailHost,
		env.EmailPort,
		env.EmailUsername,
		env.EmailPassword,
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

		var buff bytes.Buffer
		if err := templates.ExecuteTemplate(&buff, m.Template, m.Data); err != nil {
			return err
		}

		msg = buff.String()
	}

	if env.Dev {
		logger.Info(msg)
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
	templ, err := template.ParseFS(templateFs, "template/*.html")
	if err != nil {
		log.Fatalf("Could not parse template fs: %v", err)
		return
	}

	templates = templ
	register("email", createMailer)
}
