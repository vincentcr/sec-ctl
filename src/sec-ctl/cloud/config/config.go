package config

import "sec-ctl/pkg/util"

// Config represents the cloud configuration options
type Config struct {
	RESTBindHost string
	RESTBindPort uint16

	WSBindHost string
	WSBindPort uint16

	DBHost     string
	DBPort     uint16
	DBUsername string
	DBPassword string
	DBName     string

	RedisHost string
	RedisPort uint16
}

// AppName returns the name of the app being configured
func (cfg *Config) AppName() string {
	return "Cloud"
}

var defaultConfig = Config{
	RESTBindHost: "0.0.0.0",
	RESTBindPort: 9753,
	WSBindHost:   "0.0.0.0",
	WSBindPort:   9754,

	DBPort:     5432,
	DBPassword: "secctl_dev",
	DBUsername: "secctl_dev",
	DBName:     "secctl_dev",

	RedisPort: 6739,
}

// Load loads the configuration
func Load() (Config, error) {
	cfg := defaultConfig
	err := util.LoadConfig(&cfg, &defaultConfig)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
