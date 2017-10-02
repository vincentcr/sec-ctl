package main

import (
	"log"
	"tpi-mon/pkg/config"
	"tpi-mon/pkg/rest"
	"tpi-mon/pkg/site"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Panicln(err)
	}

	client := newLocalClient(cfg.Local.TPIHost, cfg.Local.TPIPort, cfg.Local.TPIPassword)
	registry := func(id string) site.Client { return client }

	startCloudConnector(cfg.Local.CloudWSURL, cfg.Local.CloudToken, client)

	rest.Run(registry, cfg.Local.RESTBindHost, cfg.Local.RESTBindPort)

}
