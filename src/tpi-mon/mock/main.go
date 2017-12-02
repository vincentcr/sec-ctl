package main

import (
	"log"
	"os"
	"tpi-mon/pkg/util"
)

var logger = log.New(os.Stderr, "[mock] ", log.LstdFlags|log.Lshortfile)

func main() {
	cfg := config{}
	err := util.LoadConfig("Mock", &cfg)
	if err != nil {
		log.Panicln(err)
	}

	if err = Run(cfg.BindHost, cfg.TPIBindPort, cfg.RESTBindPort, cfg.Password, cfg.StateFilename); err != nil {
		log.Panicln(err)
	}
}
