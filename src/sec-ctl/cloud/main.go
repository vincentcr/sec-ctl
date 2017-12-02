package main

import (
	"log"
	"os"
	"sec-ctl/cloud/config"
	"sec-ctl/cloud/db"
)

var logger = log.New(os.Stderr, "[cloud] ", log.LstdFlags|log.Lshortfile)

func main() {

	cfg, err := config.Load()
	if err != nil {
		logger.Panicln(err)
	}

	db, err := db.OpenDB(cfg)
	if err != nil {
		logger.Panicln(err)
	}

	queue, err := newQueue(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		logger.Panicln(err)
	}

	registry := newRegistry(db, queue)

	runRESTAPI(registry, db, cfg.RESTBindHost, cfg.RESTBindPort)
}
