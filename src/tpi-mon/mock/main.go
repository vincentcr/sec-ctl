package main

import (
	"log"
	"tpi-mon/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Panicln(err)
	}

	if err = Run(cfg.Mock.BindHost, cfg.Mock.TPIBindPort, cfg.Mock.RESTBindPort, cfg.Mock.Password, cfg.Mock.StateFilePath); err != nil {
		log.Panicln(err)
	}
}
