package config

import (
	"netbeams/logs"
	"os"

	"github.com/pelletier/go-toml"
)

var Configuration BaseConfig

func Load(l *logs.Logger) BaseConfig {
	content, err := os.ReadFile("./ServerConfig.toml")

	if err != nil {
		l.Error("Failed to read config file")
		l.Fatal(err)
	}

	var config BaseConfig
	if err = toml.Unmarshal(content, &config); err != nil {
		panic(err)
	}

	errors := config.Validate()
	if len(errors) > 0 {
		l.Error("Configuration file is invalid")
		for _, e := range errors {
			if e.fatal {
				l.Fatal(e)
			} else if e.warning {
				l.Warn(e.Error())
			} else {
				l.Error(e.Error())
			}
		}
		os.Exit(1)

	}

	return config
}

func GetConfig() BaseConfig {
	return Configuration
}
