package main

import (
	"netbeams/logs"
)

func main() {
	// Spawn a new logger instance
	logger := logs.NetLogger("test")
	logger.Info("Loading congiguration file")
	logger.Info("Loading data")
	logger.Info("Starting server")
	l2 := logs.NetLogger("test 2")
	l2.Info("All Done")

	logger.Terminate()
	l2.Terminate()
}
