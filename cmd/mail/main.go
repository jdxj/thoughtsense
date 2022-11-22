package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/emersion/go-smtp"

	"github.com/jdxj/thoughtsense/config"
	ms "github.com/jdxj/thoughtsense/mail-server"
	tgBot "github.com/jdxj/thoughtsense/tg-bot"
)

var (
	conf = flag.String("conf", "conf.yaml", "config path")
)

func main() {
	flag.Parse()

	err := config.Init(*conf)
	if err != nil {
		logger.Fatalf("read config err: %s", err)
	}
	tgBot.Init()

	s := smtp.NewServer(&ms.Backend{})
	s.Addr = fmt.Sprintf(":%d", config.SMTP.Port)
	s.Domain = config.SMTP.Domain
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024 * 50 // 50MB
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	logger.Infof("Starting server at: %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatalf("listen and serve err: %s", err)
	}
}
