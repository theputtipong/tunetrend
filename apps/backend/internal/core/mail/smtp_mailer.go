package mail

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"tunetrend-backend/internal/domain"
)

type SMTPMailer struct {
	host, port, username, password, to string
}

func NewSMTPMailer() domain.MailSender {
	return &SMTPMailer{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_APP_PASSWORD"),
		to:       os.Getenv("CONTACT_EMAIL_TO"),
	}
}

func (m *SMTPMailer) Send(msg domain.MailMessage) error {
	done := make(chan error, 1)
	go func() { done <- m.sendRaw(msg) }()

	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		return errors.New("smtp: send timed out")
	}
}

func (m *SMTPMailer) sendRaw(msg domain.MailMessage) error {
	addr := m.host + ":" + m.port
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	replyToLine := ""
	if msg.ReplyTo != "" {
		replyToLine = fmt.Sprintf("Reply-To: %s\r\n", msg.ReplyTo)
	}

	raw := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n%sMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.username, m.to, msg.Subject, replyToLine, msg.Body,
	)

	return smtp.SendMail(addr, auth, m.username, []string{m.to}, []byte(raw))
}
