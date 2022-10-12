package main

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	readConfig()
	fmt.Printf("%+v\n", conf)
}

func TestSendMsg(t *testing.T) {
	readConfig()
	newTGBot()

	sendMsg("abc")
}

func TestLogger(t *testing.T) {
	logger.Debugf("dd: %s", "abc")
}
