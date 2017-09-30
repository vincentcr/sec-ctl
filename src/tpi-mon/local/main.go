package main

import (
	"log"
	"tpi-mon/pkg/config"
	"tpi-mon/pkg/rest"
	"tpi-mon/pkg/site"
)

func main() {

	errCh := make(chan error)

	cfg, err := config.Load()
	if err != nil {
		log.Panicln(err)
	}

	c := newLocalClient(cfg.Local.TPIHost, cfg.Local.TPIPort, cfg.Local.TPIPassword, errCh)
	r := func(id string) site.Client { return c }

	startCloudConnector(cfg.Local.CloudWSURL, cfg.Local.CloudToken, c)

	rest.Start(r, cfg.Local.RESTBindHost, cfg.Local.RESTBindPort, errCh)

	for {
		select {
		case err := <-errCh:
			log.Panicln(err)
		}
	}
}
