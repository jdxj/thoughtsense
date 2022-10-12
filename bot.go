package main

import (
	"io"
	"strings"

	"github.com/emersion/go-message/mail"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
)

func newTGBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		panic(err)
	}
}

func sendMsg(msg tgbotapi.Chattable) {
	if msg == nil {
		return
	}

	_, err := bot.Send(msg)
	if err != nil {
		logger.Errorf("sendMsg err: %s", err)
	}
}

func sendTxtMsg(txt string) {
	msg := tgbotapi.NewMessage(conf.ChatID, txt)
	sendMsg(msg)
}

func newMsg(part *mail.Part) (c tgbotapi.Chattable) {
	ct := part.Header.Get("Content-Type")
	i := strings.Index(ct, ";")
	if i < 0 {
		return
	}
	ct = ct[:i]
	switch ct {
	case "text/plain":
		d, err := io.ReadAll(part.Body)
		if err != nil {
			logger.Errorf("read text/plain err: %s", err)
			return
		}
		c = tgbotapi.NewMessage(conf.ChatID, string(d))

	case "text/html":
		d, err := io.ReadAll(part.Body)
		if err != nil {
			logger.Errorf("read text/html err: %s", err)
			return
		}
		msg := tgbotapi.NewMessage(conf.ChatID, string(d))
		msg.ParseMode = "HTML"
		return msg

	case "application/octet-stream":
		switch h := part.Header.(type) {
		case *mail.AttachmentHeader:
			filename, err := h.Filename()
			if err != nil {
				logger.Warnf("get filename err: %s", err)
			}

			c = tgbotapi.NewDocument(conf.ChatID, tgbotapi.FileReader{
				Name:   filename,
				Reader: part.Body,
			})
		}
	}
	return
}
