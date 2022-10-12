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

	for p, e = rr.NextPart(); e == nil; p, e = rr.NextPart() {
		d, err := io.ReadAll(p.Body)
		if err != nil {
			buf.WriteString(fmt.Sprintf("read part body err: %s\n", err))
			continue
		}

		var m map[string][]string
		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			m = h.Map()
		case *mail.AttachmentHeader:
			m = h.Map()
		}
		for k, v := range m {
			buf.WriteString(fmt.Sprintf("%s: %v\n", k, v))
		}
		buf.WriteString(fmt.Sprintf("%s\n", d))
	}
	if e != nil && e != io.EOF {
		logger.Infof("next part err: %s", err)
	}
	if buf.Len() != 0 {
		sendMsg(buf.String())
		logger.Debugf("mail: %s", buf.String())
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
