package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

type Mailer interface {
	SendWelcome(ctx context.Context, to, name string) error
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	UseTLS   bool
}

type SMTPMailer struct {
	logger *log.Logger
	cfg    SMTPConfig
}

func NewSMTPMailer(logger *log.Logger, cfg SMTPConfig) *SMTPMailer {
	return &SMTPMailer{logger: logger, cfg: cfg}
}

func (m *SMTPMailer) SendWelcome(ctx context.Context, to, name string) error {
	if m.cfg.Host == "" {
		return fmt.Errorf("smtp host is not configured")
	}

	msg := buildWelcomeMessage(m.cfg.FromName, m.cfg.From, to, name)
	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)

	var auth smtp.Auth
	if m.cfg.Username != "" {
		auth = smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	}

	if m.cfg.UseTLS {
		if err := m.sendWithTLS(auth, addr, msg, to); err != nil {
			return err
		}
	} else {
		if err := smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg); err != nil {
			return err
		}
	}

	m.logger.Printf("[mailer] dispatched SMTP welcome email to %s (%s)", name, to)
	return nil
}

func (m *SMTPMailer) sendWithTLS(auth smtp.Auth, addr string, msg []byte, to string) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: m.cfg.Host})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.cfg.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err := client.Mail(m.cfg.From); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func buildWelcomeMessage(fromName, fromEmail, toEmail, toName string) []byte {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, fromEmail))
	buf.WriteString(fmt.Sprintf("To: %s <%s>\r\n", toName, toEmail))
	buf.WriteString("Subject: Welcome to the Shop!\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(fmt.Sprintf("Hi %s,\r\n\r\n", toName))
	buf.WriteString("Thanks for creating an account with us! We're excited to have you on board.\r\n")
	buf.WriteString("If you have any questions, just reply to this email and we'll be glad to help.\r\n\r\n")
	buf.WriteString("Cheers,\r\nThe Shop Team\r\n")
	return buf.Bytes()
}

var _ Mailer = (*SMTPMailer)(nil)
