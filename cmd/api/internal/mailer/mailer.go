package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/jordan-wright/email"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpAddr = "smtp.gmail.com:587"
)

type Mailer struct {
	username string
	password string
	sender   string
}

//go:embed "templates"
var templatesFS embed.FS

func NewMailer(username, password, sender string) *Mailer {
	return &Mailer{
		sender:   sender,
		password: password,
		username: username,
	}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templatesFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainTextBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainTextBody, "plainTextBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	e := email.Email{
		From:    fmt.Sprintf("%s <%s>", m.sender, m.username),
		To:      []string{recipient},
		Subject: subject.String(),
		Text:    plainTextBody.Bytes(),
		HTML:    htmlBody.Bytes(),
	}

	err = e.Send(smtpAddr, smtp.PlainAuth("", m.username, m.password, smtpHost))
	return err
}
