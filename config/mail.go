package config

import (
	"errors"
	"net"
	"net/smtp"
)

// Mail represents plantext mail net/smtp
// for this software, `to` will only contain one email
// from is given in config
// TODO: add HTML mail
type Mail struct {
	to      string
	subject string
	body    string
}

// NewMail returns a Mail struct
func NewMail(to string, subj string, body string) Mail {
	return Mail{
		to:      to,
		subject: subj,
		body:    body,
	}
}

// Send prepares msg body for net/SendMail
func (m *Mail) Send() error {
	host, _, err := net.SplitHostPort(App.ESrv)
	if err != nil {
		return errors.New("send: splithostport failed: " + err.Error())
	}
	auth := smtp.PlainAuth(NAME, App.EMail, App.EPass, host)
	msg := []byte("To: " + m.to + "\r\nSubject: " + m.subject + "\r\n" + m.body)
	err = smtp.SendMail(App.ESrv, auth, App.EMail, []string{m.to}, msg)
	if err != nil {
		return errors.New("send: sendmail: " + err.Error())
	}
	return nil
}
