package main

import (
	"log"
	"os"
	"sec-ctl/pkg/util"
)

const appName = "Local"

var logger = log.New(os.Stderr, "[local] ", log.LstdFlags|log.Lshortfile)

func main() {

	cfg := config{}
	if err := util.LoadConfig(&cfg, &defaultConfig); err != nil {
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
