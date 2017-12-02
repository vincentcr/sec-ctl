package main

type config struct {
	BindHost      string
	TPIBindPort   uint16
	RESTBindPort  uint16
	Password      string
	StateFilename string
}

var defaultConfig = config{
	BindHost:      "0.0.0.0",
	TPIBindPort:   4025,
	RESTBindPort:  9751,
	Password:      "mock123",
	StateFilename: "mock-tpi-state.json",
}
