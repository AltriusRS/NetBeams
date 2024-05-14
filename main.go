package main

import (
	"netbeams/logs"
)

func main() {
	logger := logs.NewLogger("test")
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	logger.Fatal("test")
}
