package config

import (
	"fmt"
	"os"
	"time"

	"github.com/altriusrs/netbeams/src/logs"
	"github.com/pelletier/go-toml"
)

var Configuration BaseConfig

func Load() {
	l := logs.NetLogger("Config")
	content, err := os.ReadFile("./ServerConfig.toml")

	if err != nil {

		if os.IsNotExist(err) {
			l.Warn("Configuration file not found - Creating default")
			f, err := os.Create("./ServerConfig.toml")

			if err != nil {
				l.Error(err.Error())
				l.Fatal(err)
			}

			defer func() {
				_ = f.Close()
			}()

			encoder := toml.NewEncoder(f)
			encoder.Indentation("  ")
			encoder.ArraysWithOneElementPerLine(true)
			encoder.QuoteMapKeys(true)
			encoder.Order(toml.OrderPreserve)

			if err := encoder.Encode(LoadDefault()); err != nil {
				l.Error(err.Error())
				l.Fatal(err)
			}
			Configuration = LoadDefault()

			content = []byte{}

			_, err = f.Read(content)

			if err != nil {
				l.Error(err.Error())
				l.Fatal(err)
			}

			content = content[:len(content)-1]

			l.Warn("Configuration file created from defaults")
		} else {

			l.Error("Failed to read config file")
			l.Fatal(err)
			return
		}
	}

	var config BaseConfig
	if err = toml.Unmarshal(content, &config); err != nil {
		panic(err)
	}

	errors := config.Validate()
	if len(errors) > 0 {
		l.Error("Configuration file is invalid")
		hasFatal := false
		for _, e := range errors {
			if e.fatal {
				l.Fatal(e)
				hasFatal = true
			} else if e.warning {
				l.Warn(e.Error())
			} else {
				l.Error(e.Error())
			}
		}

		if hasFatal {
			l.Fatal(fmt.Errorf("exiting due to invalid config"))
			os.Exit(1)
		}

	}

	config.Auth.Idle.MaxTimeTime, _ = time.ParseDuration(config.Auth.Idle.MaxTime)
	config.Auth.Kick.IdleDurationTime, _ = time.ParseDuration(config.Auth.Kick.IdleDuration)
	config.Auth.Kick.OnlineDurationTime, _ = time.ParseDuration(config.Auth.Kick.OnlineDuration)
	config.Auth.Kick.AdminDurationTime, _ = time.ParseDuration(config.Auth.Kick.AdminDuration)
	config.Auth.Online.QuotaTime, _ = time.ParseDuration(config.Auth.Online.Quota)

	Configuration = config
}
