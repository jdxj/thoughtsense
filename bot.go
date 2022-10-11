package main

import (
	"fmt"

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

func sendMsg(txt string) {
	msg := tgbotapi.NewMessage(conf.ChatID, txt)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Printf("sendMsg err: %s\n", err)
	}
}
