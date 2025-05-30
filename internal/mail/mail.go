package mail

import (
	"context"
)

type Address struct {
	Name    string
	Address string
}

type Message struct {
	From    Address
	To      []Address
	Subject string
	HTML    string
	Text    string
	Headers map[string]string
}

type MailerDriver interface {
	Send(context.Context, *Message) error
}

type Sender interface {
	Send(context.Context, *Message) error
}

type Manager struct {
	smtp *SMTPMailer
}

func NewManager(smtp *SMTPMailer) *Manager {
	return &Manager{
		smtp: smtp,
	}
}

func (m *Manager) Send(ctx context.Context, msg *Message) error {
	return m.smtp.Send(ctx, msg)
}
