package mailer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

type Message struct {
	To       string
	Subject  string
	TextBody string
	HTMLBody string
}

type Sender interface {
	Send(ctx context.Context, msg Message) error
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Timeout  time.Duration
}

type SMTPMailer struct {
	cfg SMTPConfig
}

func NewSMTPMailer(cfg SMTPConfig) (*SMTPMailer, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("smtp host is required")
	}
	if cfg.Port == 0 {
		return nil, fmt.Errorf("smtp port is required")
	}
	if cfg.From == "" {
		return nil, fmt.Errorf("smtp from address is required")
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 15 * time.Second
	}

	return &SMTPMailer{cfg: cfg}, nil
}

func (m *SMTPMailer) Send(ctx context.Context, msg Message) error {
	if strings.TrimSpace(msg.To) == "" {
		return fmt.Errorf("recipient is required")
	}
	if strings.TrimSpace(msg.Subject) == "" {
		return fmt.Errorf("subject is required")
	}
	if msg.TextBody == "" && msg.HTMLBody == "" {
		return fmt.Errorf("email body is required")
	}

	rawMessage := buildMIMEMessage(m.cfg.From, msg)

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)

	var auth smtp.Auth
	if m.cfg.Username != "" {
		auth = smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- smtp.SendMail(addr, auth, m.cfg.From, []string{msg.To}, []byte(rawMessage))
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("send smtp mail: %w", err)
		}
	}

	return nil
}

func buildMIMEMessage(from string, msg Message) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("From: %s\r\n", from))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", msg.To))
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	builder.WriteString("MIME-Version: 1.0\r\n")

	if msg.TextBody != "" && msg.HTMLBody != "" {
		boundary := randomBoundary()
		builder.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", boundary))

		builder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		builder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
		builder.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		builder.WriteString(msg.TextBody)
		builder.WriteString("\r\n")

		builder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		builder.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		builder.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		builder.WriteString(msg.HTMLBody)
		builder.WriteString("\r\n")

		builder.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else if msg.HTMLBody != "" {
		builder.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		builder.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		builder.WriteString(msg.HTMLBody)
		builder.WriteString("\r\n")
	} else {
		builder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
		builder.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		builder.WriteString(msg.TextBody)
		builder.WriteString("\r\n")
	}

	return builder.String()
}

func randomBoundary() string {
	var bytes [12]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "mimeBoundary"
	}
	return hex.EncodeToString(bytes[:])
}
