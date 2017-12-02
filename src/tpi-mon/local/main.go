package main

import (
	"log"
	"os"
)

const appName = "Local"

var logger = log.New(os.Stderr, "[local] ", log.LstdFlags|log.Lshortfile)

func main() {

	cfg, err := loadConfig()
	if err != nil {
		logger.Panicln(err)
	}

	if cfg.SiteID == "" { // unset client id => first time!
		if err := firstTime(&cfg); err != nil {
			logger.Panicln(err)
		}
	}

	site := newLocalSite(cfg.TPIHost, cfg.TPIPort, cfg.TPIPassword, cfg.SiteID)

	startCloudConnector(cfg.CloudWSURL, cfg.CloudToken, site)

	runRESTAPI(site, cfg.RESTBindHost, cfg.RESTBindPort)
}
