package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

var conf config

type config struct {
	Token  string `yaml:"token"`
	ChatID int64  `yaml:"chat_id"`
	Domain string `yaml:"domain"`
	Port   int    `yaml:"port"`
}

func readConfig() {
	f, err := os.Open("./conf.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
}
