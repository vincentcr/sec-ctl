package config

import "sec-ctl/pkg/util"

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

var defaultConfig = Config{
	RESTBindHost: "0.0.0.0",
	RESTBindPort: 9753,
	WSBindHost:   "0.0.0.0",
	WSBindPort:   9754,

	DBPort:     5432,
	DBPassword: "tpimon_dev",
	DBUsername: "tpimon_dev",
	DBName:     "tpimon_dev",

	RedisPort: 6739,
}

// Load loads the configuration
func Load() (Config, error) {
	cfg := Config{}
	err := util.LoadConfig("Cloud", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
