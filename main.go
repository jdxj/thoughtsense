package main

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

func main() {
	readConfig()
	newTGBot()

	s := smtp.NewServer(&Backend{})
	s.Addr = fmt.Sprintf(":%d", conf.Port)
	s.Domain = conf.Domain
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
