package main

import (
	"log"
	"os"
	"sec-ctl/pkg/util"
)

var logger = log.New(os.Stderr, "[mock] ", log.LstdFlags|log.Lshortfile)

func main() {
	cfg := config{}
	err := util.LoadConfig(&cfg, &defaultConfig)
	if err != nil {
		log.Panicln(err)
	}

	if err = Run(cfg.BindHost, cfg.TPIBindPort, cfg.RESTBindPort, cfg.Password, cfg.StateFilename); err != nil {
		log.Panicln(err)
	}
}
