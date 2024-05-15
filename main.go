package main

import (
	"netbeams/config"
	"netbeams/logs"
	"netbeams/server"
)

func main() {
	// Spawn a new logger instance
	logger := logs.NetLogger("")
	defer logger.Terminate()
	logger.Info("Loading congiguration file")
	logger.Info("Loading data")
	serverConfig := config.Load(&logger)
	logger.Info("Starting server")
	logger.Info("Name: " + serverConfig.General.Name)
	logger.Infof("Port: %d", serverConfig.General.Port)
	mode := "main"

	if serverConfig.NetBeams.MasterNode != "localhost" {
		logger.Info("Main node: " + serverConfig.NetBeams.MasterNode)
		logger.Info("Switching to node mode")
		mode = "node"
	}
	logger.Info("Mode: " + mode)

	app := server.NewApplication(serverConfig)

	switch mode {
	case "main":
		logger.Info("Starting main node")
		app.StartMainNode()
	case "node":
		logger.Error("Node mode not implemented yet")
		logger.Info("Please set the master node to 'localhost' in the configuration file")
		// logger.Info("Starting node")
		// app.StartNode()
	}

	app.Wait() // Wait for the application to terminate
	logger.Info("Exiting")
}
