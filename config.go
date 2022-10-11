package main

import (
	"flag"
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

var (
	path = flag.String("conf", "conf.yaml", "config path")
)

func readConfig() {
	flag.Parse()

	f, err := os.Open(*path)
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
