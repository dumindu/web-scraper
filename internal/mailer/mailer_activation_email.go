package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

const (
	titleActivationEmail = "Activation Email"
	fmtActivationLink    = "%s/users/activate?email=%s&token=%s"
)

type ActivationEmail struct {
	ActivationCode string
	ActivationLink string
}

func (ml *Mailer) ActivationMail(userEmail string, token string) error {
	from := ml.Senders.NoReply
	to := userEmail
	subject := titleActivationEmail
	data := &ActivationEmail{
		ActivationCode: token,
		ActivationLink: fmt.Sprintf(fmtActivationLink, ml.Links.WebsiteHost, to, token),
	}

	wr := new(bytes.Buffer)
	t, err := template.ParseFiles(tmplActivationEmail)
	if err != nil {
		return ErrNoTmpl
	}

	if err := t.Execute(wr, data); err != nil {
		return err
	}

	mail := newMail(from, to, subject, wr)
	rfc822Msg := mail.ToRfc822Msg()

	auth := ml.Auth
	if strings.Contains(ml.Addr, "mailhog") {
		auth = nil
	}

	return smtp.SendMail(ml.Addr, auth, from, []string{to}, []byte(rfc822Msg))
}
