package notifier

import (
	"app/env"
	"bytes"
	"context"
	"embed"
	"fmt"
	"log"
	"text/template"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"gopkg.in/gomail.v2"
)

//go:embed template
var templateFs embed.FS
var templates *template.Template

type emailArg struct {
	Message
}

func (emailArg) Kind() string {
	return "email"
}

type emailer struct {
	river *river.Client[pgx.Tx]
}

func createMailer(opt *option) Notifier {
	if opt.river == nil {
		panic("Please provde queue when creating emailer")
	}

	return &emailer{river: opt.river}
}

func (e *emailer) Send(m Message) error {
	_, err := e.river.Insert(context.Background(), emailArg{m}, &river.InsertOpts{
		MaxAttempts: 3,
	})

	return err
}

type emailWorker struct {
	river.WorkerDefaults[emailArg]
	mailer *gomail.Dialer
}

func CreateEmailWorker() *emailWorker {
	mailer := gomail.NewDialer(
		env.EMAIL_HOST,
		env.EMAIL_PORT,
		env.EMAIL_USERNAME,
		env.EMAIL_PASSWORD,
	)

	return &emailWorker{
		mailer: mailer,
	}
}

func (w *emailWorker) Work(ctx context.Context, job *river.Job[emailArg]) error {
	arg := job.Args
	return w.send(arg.Message)
}

func (e *emailWorker) send(m Message) error {
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

	if env.DEV {
		logger.Info(msg)
		return nil
	}

	from := m.From
	if from == "" {
		from = fmt.Sprintf("%s <%s>", env.APP_NAME, env.EMAIL_SENDER)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", m.To...)
	message.SetHeader("Subject", m.Subject)
	message.SetHeader("Bcc", m.Bcc...)
	message.SetHeader("Cc", m.Cc...)
	message.SetBody(mimetype, msg)

	return e.mailer.DialAndSend(message)
}

func (e *emailWorker) CreateWorker(workers *river.Workers) {
	river.AddWorker(workers, e)
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
