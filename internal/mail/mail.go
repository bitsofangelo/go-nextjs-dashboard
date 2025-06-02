package mail

import (
	"context"
	"errors"
)

type Mailer string

const (
	PlainSTMPMailer Mailer = "plain_smtp"
	CustomMailer    Mailer = "custom"
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

type Sender interface {
	Send(context.Context, *Message) error
}

type MailerDriver = Sender

type Factory interface {
	Sender
	Mailer(Mailer) (Sender, error)
}

type Manager struct {
	mailers       map[Mailer]MailerDriver
	defaultMailer Mailer
	smtp          *SMTPMailer
}

func NewManager(smtp *SMTPMailer) *Manager {
	mailers := make(map[Mailer]MailerDriver)
	mailers[PlainSTMPMailer] = smtp

	return &Manager{
		mailers:       mailers,
		defaultMailer: PlainSTMPMailer,
		smtp:          smtp,
	}
}

func (m *Manager) Mailer(mailer Mailer) (Sender, error) {
	if ml, ok := m.mailers[mailer]; ok {
		return ml, nil
	}
	return nil, errors.New("mailer not found")
}

func (m *Manager) Send(ctx context.Context, msg *Message) error {
	mailer, err := m.Mailer(m.defaultMailer)
	if err != nil {
		return err
	}

	return mailer.Send(ctx, msg)
}

func (m *Manager) RegisterMailer(mailer Mailer, driver MailerDriver) error {
	_, ok := m.mailers[mailer]
	if ok {
		return errors.New("mailer already exists")
	}

	m.mailers[mailer] = driver
	return nil
}

func (m *Manager) SetDefaultMailer(mailer Mailer) error {
	if _, ok := m.mailers[mailer]; !ok {
		return errors.New("mailer not found")
	}

	m.defaultMailer = mailer
	return nil
}
