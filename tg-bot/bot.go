package tg_bot

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/emersion/go-message/mail"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jdxj/thoughtsense/config"
)

var (
	bot *tgbotapi.BotAPI
)

func Init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(config.TGBot.Token)
	if err != nil {
		logger.Fatalf("init tg bot err: %s", err)
	}
	bot.Debug = true

	link := fmt.Sprintf("https://%s:443/%s", config.TGBot.Domain, config.TGBot.Token)
	// wc, err := tgbotapi.NewWebhookWithCert(link, tgbotapi.FilePath(config.TGBot.Cert))
	wc, err := tgbotapi.NewWebhook(link)
	if err != nil {
		logger.Fatalf("new web hook err: %s", err)
	}
	_, err = bot.Request(wc)
	if err != nil {
		logger.Fatalf("send web hook req err: %s", err)
	}

	wi, err := bot.GetWebhookInfo()
	if err != nil {
		logger.Fatalf("get web hook info err: %s", err)
	}
	if wi.LastErrorDate != 0 {
		logger.Warnf("telegram callback failed: %s", wi.LastErrorMessage)
	}

	updates := bot.ListenForWebhook(fmt.Sprintf("/%s", bot.Token))
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", config.TGBot.Port), nil)
		if err != nil {
			logger.Errorf("listen and serve err: %s", err)
		}
	}()

	go func() {
		logger.Debugf("get updates")
		for update := range updates {
			if msg := update.Message; msg != nil {
				logger.Debugf("%+v\n", msg)
			}
		}
	}()
}

func SendMsg(msg tgbotapi.Chattable) {
	if msg == nil {
		return
	}

	_, err := bot.Send(msg)
	if err != nil {
		logger.Errorf("sendMsg err: %s", err)
	}
}

func SendTxtMsg(txt string) {
	msg := tgbotapi.NewMessage(config.TGBot.ChatID, txt)
	SendMsg(msg)
}

func NewMsg(part *mail.Part) (c tgbotapi.Chattable) {
	logger.Debugf("part header: %v", part.Header)

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
		c = tgbotapi.NewMessage(config.TGBot.ChatID, string(d))

	case "text/html":
		d, err := io.ReadAll(part.Body)
		if err != nil {
			logger.Errorf("read text/html err: %s", err)
			return
		}
		msg := tgbotapi.NewMessage(config.TGBot.ChatID, string(d))
		msg.ParseMode = "HTML"
		c = msg

	case "application/octet-stream":
		switch h := part.Header.(type) {
		case *mail.AttachmentHeader:
			filename, err := h.Filename()
			if err != nil {
				logger.Warnf("get filename err: %s", err)
				return
			}

			c = tgbotapi.NewDocument(config.TGBot.ChatID, tgbotapi.FileReader{
				Name:   filename,
				Reader: part.Body,
			})

		default:
			logger.Warnf("part header %T founded", h)
		}

	default:
		logger.Warnf("%s not handle", ct)
	}
	return
}
