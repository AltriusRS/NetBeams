package main

import (
	"fmt"
	"netbeams/config"
	"netbeams/environment"
	"netbeams/globals"
	"netbeams/http"
	"netbeams/logs"
	"netbeams/tcp"
	"netbeams/udp"
)

func main() {
	environment.GetBuildContext()

	// Spawn a new logger instance
	logger := logs.NetLogger("Main")
	defer logger.Terminate()

	logger.Info("Welcome to NetBeams v" + environment.Version)

	if environment.Context.IsDev {
		fmt.Println("\x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m")
		fmt.Println(" DEVELOPER ENVIRONMENT")
		fmt.Println("\x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m \x1b[40;33m⬕⬔\x1b[0m")

	}

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

	// Spawn a new application instance
	globals.NewApplication(serverConfig, &logger)

	// Pass application to signal handler to allow graceful shutdown
	globals.App.RegisterSignalHandler()

	// Spawn the required services
	globals.App.AddService(http.Service())
	globals.App.AddService(tcp.Service())
	globals.App.AddService(udp.Service())

	switch mode {
	case "main":
		logger.Info("Starting main node")

	case "node":
		logger.Error("Node mode not implemented yet")
		logger.Info("Please set the master node to 'localhost' in the configuration file")
		return
		// logger.Info("Starting node")
		// app.StartNode()
	}

	globals.App.Start()
	globals.App.Wait() // Wait for the application to terminate
	logger.Info("Exiting")
}
