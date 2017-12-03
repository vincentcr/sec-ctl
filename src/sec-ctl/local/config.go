package main

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

// AppName returns the name of the app to configured
func (cfg *config) AppName() string {
	return "Local"
}

var defaultConfig = config{
	TPIPort:      4025,
	TPIPassword:  "mock123",
	RESTBindHost: "0.0.0.0",
	RESTBindPort: 9752,
	CloudWSURL:   "ws://localhost:9754",
	CloudBaseURL: "http://localhost:9753",
}
