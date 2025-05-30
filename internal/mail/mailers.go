package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/gelozr/go-dash/internal/config"
)

type SMTPMailer struct {
	Host          string // "smtp.example.com"
	Port          int    // 587
	Username      string
	Password      string
	SkipTLSVerify bool
}

func NewSMTPMailer(cfg *config.Config) *SMTPMailer {
	return &SMTPMailer{
		Host:          cfg.MailHost,
		Port:          cfg.MailPort,
		Username:      cfg.MailUser,
		Password:      cfg.MailPass,
		SkipTLSVerify: cfg.MailSkipTLSVerify,
	}
}

func (s *SMTPMailer) Send(ctx context.Context, m *Message) error {
	// Build RFC-822 body
	var sb strings.Builder

	from := mail.Address{Name: m.From.Name, Address: m.From.Address}
	sb.WriteString(fmt.Sprintf("From: %s\r\n", from.String()))

	if len(m.To) > 0 {
		toList := make([]string, len(m.To))
		for i, a := range m.To {
			to := mail.Address{Name: a.Name, Address: a.Address}
			toList[i] = to.String()
		}
		sb.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(toList, ", ")))
	}

	sb.WriteString("Subject: " + m.Subject + "\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString(`Content-Type: text/html; charset="UTF-8"` + "\r\n\r\n")
	sb.WriteString(m.HTML)

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	tlsCfg := &tls.Config{
		ServerName:         s.Host,
		InsecureSkipVerify: s.SkipTLSVerify,
	}

	// honour context cancellation
	done := make(chan error, 1)

	go func() {
		var client *smtp.Client
		var err error

		if s.Port == 465 {
			conn, e := tls.Dial("tcp", addr, tlsCfg)
			if e != nil {
				done <- fmt.Errorf("dial tls: %w", e)
				return
			}

			client, err = smtp.NewClient(conn, s.Host)
			if err != nil {
				done <- fmt.Errorf("smtp new client: %w", err)
				return
			}
		} else {
			client, err = smtp.Dial(addr)
			if err != nil {
				done <- fmt.Errorf("smtp dial: %w", err)
				return
			}

			if ok, _ := client.Extension("STARTTLS"); ok {
				if err = client.StartTLS(tlsCfg); err != nil {
					done <- fmt.Errorf("start tls: %w", err)
				}
			}
		}

		defer func() {
			if err = client.Quit(); err != nil {
				done <- fmt.Errorf("quit: %w", err)
			}
		}()

		if err = client.Auth(auth); err != nil {
			done <- fmt.Errorf("client auth: %w", err)
			return
		}
		if err = client.Mail(m.From.Address); err != nil {
			done <- fmt.Errorf("client mail: %w", err)
			return
		}
		for _, a := range m.To {
			if err = client.Rcpt(a.Address); err != nil {
				done <- fmt.Errorf("client rcpt: %s: %w", a.Address, err)
				return
			}
		}

		w, err := client.Data()
		if err != nil {
			done <- fmt.Errorf("client data: %w", err)
			return
		}

		_, err = w.Write([]byte(sb.String()))
		if err != nil {
			done <- fmt.Errorf("write: %w", err)
			return
		}

		if err = w.Close(); err != nil {
			done <- fmt.Errorf("writer close: %w", err)
			return
		}

		done <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
