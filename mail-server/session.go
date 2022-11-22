package mail_server

import (
	"bytes"
	"fmt"
	"io"

	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-smtp"

	tgBot "github.com/jdxj/thoughtsense/tg-bot"
)

type Session struct {
	from, to string
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
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
	tgBot.SendTxtMsg(buf.String())

	for p, e = rr.NextPart(); e == nil; p, e = rr.NextPart() {
		msg := tgBot.NewMsg(p)
		if msg != nil {
			tgBot.SendMsg(msg)
		}
	}
	if e != nil && e != io.EOF {
		logger.Infof("next part err: %s", err)
	}
	return nil
}
