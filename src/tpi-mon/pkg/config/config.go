package config

var defaultConfig = Config{

	Cloud: cloudConfig{
		RESTBindHost: "0.0.0.0",
		RESTBindPort: 9753,
		WSBindHost:   "0.0.0.0",
		WSBindPort:   9754,
	},

	Local: localConfig{
		TPIPort:      4025,
		TPIPassword:  "mock123",
		RESTBindHost: "0.0.0.0",
		RESTBindPort: 9752,
		CloudWSURL:   "ws://localhost:9754",
	},

	Mock: mockConfig{
		BindHost:      "0.0.0.0",
		TPIBindPort:   4025,
		RESTBindPort:  9751,
		Password:      "mock123",
		StateFilename: "mock-tpi-state.json",
	},
}

// Config represents the configuration data
type Config struct {

	// config for local
	Local localConfig

	// config for cloud
	Cloud cloudConfig

	// config for mock
	Mock mockConfig
}

type localConfig struct {
	TPIHost     string
	TPIPort     uint16
	TPIPassword string

	RESTBindHost string
	RESTBindPort uint16

	CloudWSURL string
	CloudToken string
}

type cloudConfig struct {
	RESTBindHost string
	RESTBindPort uint16

	WSBindHost string
	WSBindPort uint16
}

type mockConfig struct {
	BindHost      string
	TPIBindPort   uint16
	RESTBindPort  uint16
	Password      string
	StateFilename string
}
