package main

import "sec-ctl/pkg/util"

type config struct {
	SiteID string

	TPIHost     string
	TPIPort     uint16
	TPIPassword string

	RESTBindHost string
	RESTBindPort uint16

	CloudWSURL   string
	CloudToken   string
	CloudBaseURL string
}

var defaultConfig = config{
	TPIPort:      4025,
	TPIPassword:  "mock123",
	RESTBindHost: "0.0.0.0",
	RESTBindPort: 9752,
	CloudWSURL:   "ws://localhost:9754",
	CloudBaseURL: "http://localhost:9753",
}

func loadConfig() (config, error) {
	cfg := config{}

	if err := util.LoadConfig(appName, cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}
