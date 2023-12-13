package services

import (
	"embed"
	"html/template"

	"github.com/wneessen/go-mail"
)

type MailService struct {
	Client   *mail.Client
	Sender   string
	Template embed.FS
}

func NewMailService(sender string, client *mail.Client, template embed.FS) *MailService {
	return &MailService{
		Client:   client,
		Sender:   sender,
		Template: template,
	}
}

func (s MailService) SendOTP(destination string, code string) error {
	tmpl, err := template.ParseFS(s.Template, "templates/email/otp.html")
	if err != nil {
		return err
	}

	data := map[string]any{
		"code": code,
	}

	msg := mail.NewMsg()
	msg.To(destination)
	msg.From(s.Sender)
	msg.Subject("Foodia One-Time Passcode")
	msg.SetBodyHTMLTemplate(tmpl, data)

	if err := s.Client.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
