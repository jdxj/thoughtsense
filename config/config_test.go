package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("wd: %s\n", wd)

	path := filepath.Join(wd, "conf.yaml")
	err = ReadConfig(path)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	fmt.Printf("%+v\n", TGBot)
	fmt.Printf("%+v\n", SMTP)
}
