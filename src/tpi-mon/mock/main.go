package main

import (
	"log"
	"os"
	"tpi-mon/pkg/config"
)

var logger = log.New(os.Stderr, "[mock] ", log.LstdFlags|log.Lshortfile)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Panicln(err)
	}

	if err = Run(cfg.Mock.BindHost, cfg.Mock.TPIBindPort, cfg.Mock.RESTBindPort, cfg.Mock.Password, cfg.Mock.StateFilename); err != nil {
		log.Panicln(err)
	}
}
