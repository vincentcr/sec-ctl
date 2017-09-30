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

	errCh := make(chan error)
	s := startServer(cfg.Cloud.WSBindHost, cfg.Cloud.WSBindPort, errCh)
	rest.Start(s.GetClient, cfg.Cloud.RESTBindHost, cfg.Cloud.RESTBindPort, errCh)

	for {
		select {
		case err := <-errCh:
			log.Panicln(err)
		}
	}

}
