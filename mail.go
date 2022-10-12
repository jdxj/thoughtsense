package main

import (
	"bytes"
	"fmt"
	"io"

	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-smtp"
)

type Backend struct{}

func (bkd *Backend) Login(_ *smtp.ConnectionState, _, _ string) (smtp.Session, error) {
	return &Session{}, nil
}

func (bkd *Backend) AnonymousLogin(_ *smtp.ConnectionState) (smtp.Session, error) {
	return &Session{}, nil
}

type Session struct {
	from, to string
}

func (s *Session) AuthPlain(username, password string) error {
	return nil
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	rr, err := mail.CreateReader(r)
	if err != nil {
		return err
	}
	defer rr.Close()

	var (
		p   *mail.Part
		e   error
		buf = bytes.NewBuffer(nil)
	)

	subject := "Subject"
	rr.Header.Get(subject)
	buf.WriteString(fmt.Sprintf("from: %s, to: %s\n", s.from, s.to))
	buf.WriteString(fmt.Sprintf("%s: %s\n", subject, rr.Header.Get(subject)))
	sendTxtMsg(buf.String())

	for p, e = rr.NextPart(); e == nil; p, e = rr.NextPart() {
		sendMsg(newMsg(p))
	}
	if e != nil && e != io.EOF {
		logger.Infof("next part err: %s", err)
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
