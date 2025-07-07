package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
)

var ErrNoTmpl = errors.New("no tmpl")

type Mailer struct {
	Addr    string
	Auth    smtp.Auth
	Senders *Senders
	Links   *Links
}

type Senders struct {
	NoReply string
}

type Links struct {
	WebsiteHost string
}

func New(conf *Conf) *Mailer {
	return &Mailer{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Auth:    smtp.PlainAuth("", conf.User, conf.Pass, conf.Host),
		Senders: conf.Senders,
		Links:   conf.Links,
	}
}

type mail struct {
	from    string
	to      string
	subject string
	body    *bytes.Buffer
}

func newMail(from string, to string, subject string, body *bytes.Buffer) *mail {
	return &mail{
		from:    from,
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (m *mail) ToRfc822Msg() string {
	return fmt.Sprintf(fmtRfc822Msg, m.from, m.to, m.subject, m.body)
}
