package config

import (
	"fmt"
)

// A ConfigError represents a single error in a config file
type ConfigError struct {
	// The error code
	code uint16

	// The error message
	message string

	// Additional details about the error
	details string

	// Whether the error is a fatal error
	usesDefault bool

	// Whether the error is a fatal error
	fatal bool

	// Whether the error is a warning
	warning bool
}

func (e ConfigError) Error() string {
	code := fmt.Sprintf("%x", e.code)

	for len(code) < 4 {
		code = "0" + code
	}

	return fmt.Sprintf("0x%s : %s - %s", code, e.message, e.details)
}

func (c *BaseConfig) Validate() []ConfigError {
	errors := []ConfigError{}

	baseErrors := c.General.Validate()

	miscErrors := c.Misc.Validate()

	netBeamsErrors := c.NetBeams.Validate()

	errors = append(errors, baseErrors...)
	errors = append(errors, miscErrors...)
	errors = append(errors, netBeamsErrors...)

	if c.General.Port == c.NetBeams.MasterPort {
		errors = append(errors, ConfigError{
			code:        0x1000,
			message:     "Master port and Game port are the same",
			details:     "Master port and Game port should not be the same",
			usesDefault: false,
			fatal:       true,
			warning:     false,
		})
	}

	return errors
}

func (c *GeneralConfig) Validate() []ConfigError {
	errors := []ConfigError{}
	if c.Port < 100 || c.Port > 65535 {
		c.Port = 30814 // default
		errors = append(errors, ConfigError{
			code:        0x0001,
			message:     "Invalid port",
			details:     "Port must be between 100 and 65535 - Will use default value (30814)",
			usesDefault: true,
			fatal:       false,
			warning:     true,
		})
	}

	if c.MaxPlayers < 1 {
		c.MaxPlayers = 10 // default
		errors = append(errors, ConfigError{
			code:        0x0002,
			message:     "Invalid max players",
			usesDefault: false,
			fatal:       false,
			warning:     false,
		})
	} else if c.MaxPlayers > 16 {
		errors = append(errors, ConfigError{
			code:        0x0003,
			message:     "High max player count",
			details:     "It is recommended to keep the max player count at most 16",
			usesDefault: false,
			fatal:       false,
			warning:     true,
		})
	}

	if c.MaxCars < 1 {
		c.MaxCars = 2 // default
		errors = append(errors, ConfigError{
			code:    0x0004,
			message: "Invalid max cars",
			fatal:   false,
			warning: false,
		})
	} else if c.MaxCars > 4 {
		errors = append(errors, ConfigError{
			code:        0x0005,
			message:     "High max player count",
			details:     "It is recommended to keep the vehicle limit at most 4",
			usesDefault: false,
			fatal:       false,
			warning:     true,
		})
	}

	if c.ResourceFolder == "" {
		c.ResourceFolder = "Resources" // default
		errors = append(errors, ConfigError{
			code:        0x0006,
			message:     "Invalid resource folder",
			details:     "Resource folder should be set - default value will be used",
			usesDefault: true,
			fatal:       false,
			warning:     true,
		})
	}

	return errors
}

func (c *MiscConfig) Validate() []ConfigError {
	errors := []ConfigError{}

	// Nothing to validate here right now

	return errors
}

func (c *NetBeamsConfig) Validate() []ConfigError {
	errors := []ConfigError{}

	if c.MasterNode == "" {
		c.MasterNode = "localhost" // default
		errors = append(errors, ConfigError{
			code:        0x0100,
			message:     "Invalid master node",
			details:     "Master node should be set - use localhost as default",
			usesDefault: true,
			fatal:       false,
			warning:     true,
		})
	}

	if c.MasterPort < 100 || c.MasterPort > 65535 {
		c.MasterPort = 30814 // default
		errors = append(errors, ConfigError{
			code:        0x0200,
			message:     "Invalid master port",
			details:     "Master port must be between 100 and 65535 - Will use default value (30814)",
			usesDefault: true,
			fatal:       false,
			warning:     false,
		})
	}

	switch c.LogLevel {
	case "debug", "info", "warn", "error", "fatal":
		// Nothing to validate here
	default:
		c.LogLevel = "info" // default
		errors = append(errors, ConfigError{
			code:        0x0300,
			message:     "Invalid log level",
			details:     "Log level must be one of debug, info, warn, error, fatal",
			usesDefault: true,
			fatal:       false,
			warning:     false,
		})
	}

	return errors
}
