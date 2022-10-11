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

	buf.WriteString(fmt.Sprintf("from: %s, to: %s\n", s.from, s.to))

	for p, e = rr.NextPart(); e == nil; p, e = rr.NextPart() {
		d, err := io.ReadAll(p.Body)
		if err != nil {
			buf.WriteString(fmt.Sprintf("read part body err: %s\n", err))
			continue
		}
		buf.WriteString(fmt.Sprintf("%s\n", d))
	}
	if e != nil && e != io.EOF {
		fmt.Printf("next part err: %s", err)
	}
	if buf.Len() != 0 {
		sendMsg(buf.String())
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
