package db

import (
	"os"
	"testing"
	"sec-ctl/cloud/config"
)

func TestFoo(t *testing.T) {
}

var db *DB

func TestMain(m *testing.M) {
	var err error

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err = OpenDB(cfg)
	if err != nil {
		panic(err)
	}

	retCode := m.Run()
	os.Exit(retCode)
}
