package mailer

import (
	"flag"
	"fmt"
)

const fmtRfc822Msg = "Mime-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n"

var (
	flags = flag.NewFlagSet("mail", flag.ExitOnError)
	dir   = flags.String("dir", "internal/mailer/tmpl", "directory with mail templates")

	tmplActivationEmail = fmt.Sprintf("%s/%s", *dir, "activation-email.html")
)

type Conf struct {
	Host    string
	Port    int
	User    string
	Pass    string
	Senders *Senders
	Links   *Links
}
