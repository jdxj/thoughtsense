package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"

	"github.com/jdxj/thoughtsense/config"
)

var (
	conf = flag.String("conf", "conf.yaml", "config path")
)

func main() {
	flag.Parse()

	err := config.ReadConfig(*conf)
	if err != nil {
		logger.Fatalf("read config err: %s", err)
	}

	newTGBot()

	s := smtp.NewServer(&Backend{})
	s.Addr = fmt.Sprintf(":%d", config.SMTP.Port)
	s.Domain = config.SMTP.Domain
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
