package main

import (
	"log"
	"os"
	"tpi-mon/pkg/config"
	"tpi-mon/pkg/rest"
)

var logger = log.New(os.Stderr, "[cloud] ", log.LstdFlags|log.Lshortfile)

func main() {

	cfg, err := config.Load()
	if err != nil {
		logger.Panicln(err)
	}

	s := startServer(cfg.Cloud.WSBindHost, cfg.Cloud.WSBindPort)
	rest.Run(s.GetClient, cfg.Cloud.RESTBindHost, cfg.Cloud.RESTBindPort)
}
