package notifier

import (
	"app/env"
	"log"

	gomail "gopkg.in/gomail.v2"
)

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
	if env.Dev {
		log.Println("Sending email...")
		if s.Html != "" {
			log.Println(s.Html)
		} else {
			log.Println(s.Message)
		}

		return nil
	}

	message := gomail.NewMessage()
	message.SetHeader("From", s.From)
	message.SetHeader("To", s.To...)
	message.SetHeader("Subject", s.Subject)
	message.SetHeader("Bcc", s.Bcc...)
	message.SetHeader("Cc", s.Cc...)

	if s.Html != "" {
		message.SetBody("text/html", s.Html)
	} else {
		message.SetBody("text/plain", s.Message)
	}

	return e.mailer.DialAndSend(message)
}

func init() {
	register("email", createMailer)
}
