package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	TGBot = &tgBot{}
	SMTP  = &smtp{}

	conf = config{
		TGBot: TGBot,
		SMTP:  SMTP,
	}
)

type tgBot struct {
	Token  string `yaml:"token"`
	ChatID int64  `yaml:"chat_id"`
	Cert   string `yaml:"cert"`
	Domain string `yaml:"domain"`
	Port   int    `yaml:"port"`
}

type smtp struct {
	Domain string `yaml:"domain"`
	Port   int    `yaml:"port"`
}

type config struct {
	TGBot *tgBot `yaml:"tg_bot"`
	SMTP  *smtp  `yaml:"smtp"`
}

func Init(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Warnf("close config file err: %s", err)
		}
	}()

	decoder := yaml.NewDecoder(f)
	return decoder.Decode(&conf)
}
