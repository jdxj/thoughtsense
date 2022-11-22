package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/emersion/go-smtp"

	"github.com/jdxj/thoughtsense/config"
	mail_server "github.com/jdxj/thoughtsense/mail-server"
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

	tg_bot.newTGBot()

	s := smtp.NewServer(&mail_server.Backend{})
	s.Addr = fmt.Sprintf(":%d", config.SMTP.Port)
	s.Domain = config.SMTP.Domain
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	logger.Infof("Starting server at: %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatalf("listen and serve err: %s", err)
	}
}
