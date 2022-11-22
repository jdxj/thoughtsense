package tg_bot

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jdxj/thoughtsense/config"
)

func TestMain(t *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(wd, "..", "deploy/conf.yaml")
	fmt.Printf("path: %s\n", path)
	err = config.Init(path)
	if err != nil {
		panic(err)
	}

	Init()

	os.Exit(t.Run())
}

func TestSendTxtMsg(t *testing.T) {
	SendTxtMsg("hello")
}
