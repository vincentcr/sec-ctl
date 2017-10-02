package main

import (
	"log"
	"tpi-mon/pkg/config"
	"tpi-mon/pkg/rest"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Panicln(err)
	}

	s := startServer(cfg.Cloud.WSBindHost, cfg.Cloud.WSBindPort)
	rest.Run(s.GetClient, cfg.Cloud.RESTBindHost, cfg.Cloud.RESTBindPort)
}
