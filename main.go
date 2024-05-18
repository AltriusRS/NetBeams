package main

import (
	"netbeams/config"
	"netbeams/environment"
	"netbeams/logs"
	"netbeams/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	environment.GetBuildContext()

	// Spawn a new logger instance
	logger := logs.NetLogger("Main")
	defer logger.Terminate()

	logger.Info("Welcome to NetBeams v" + environment.Version)

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
	app := server.NewApplication(serverConfig, &logger)

	// Pass application to signal handler to allow graceful shutdown
	signalHandler(&app)

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

// goroutine to handle signals and gracefully shutdown the application
func signalHandler(app *server.Application) {
	channel := make(chan os.Signal, 1)

	exit := make(chan bool, 1)

	signal.Notify(channel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-channel
		app.Logger.Info("Received signal: " + sig.String())
		exit <- true
	}()

	go func() {
		<-exit
		app.Logger.Info("Shutdown signal handler started - Closing in 5 seconds")
		app.Logger.Info("Closing application")
		app.Shutdown()
	}()
}
