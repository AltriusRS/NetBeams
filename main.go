package main

import (
	"fmt"

	"github.com/altriusrs/netbeams/src/config"
	"github.com/altriusrs/netbeams/src/environment"
	"github.com/altriusrs/netbeams/src/heartbeat"
	"github.com/altriusrs/netbeams/src/http"
	"github.com/altriusrs/netbeams/src/logs"
	"github.com/altriusrs/netbeams/src/tcp"
	"github.com/altriusrs/netbeams/src/types"
	"github.com/altriusrs/netbeams/src/udp"
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
	config.Load()
	logger.Info("Starting server")
	logger.Info("Name: " + config.Configuration.General.Name)
	logger.Infof("Port: %d", config.Configuration.General.Port)
	mode := "main"

	if config.Configuration.NetBeams.MasterNode != "localhost" {
		logger.Info("Main node: " + config.Configuration.NetBeams.MasterNode)
		logger.Info("Switching to node mode")
		mode = "node"
	}
	logger.Info("Mode: " + mode)

	// Spawn a new application instance
	types.NewApplication()

	// Pass application to signal handler to allow graceful shutdown
	types.App.RegisterSignalHandler()

	// Spawn the required services
	// types.App.AddService(config.Service())
	types.App.AddService(http.Service())
	types.App.AddService(tcp.Service())
	types.App.AddService(udp.Service())
	types.App.AddService(heartbeat.Service())

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

	types.App.Start()
	types.App.Wait() // Wait for the application to terminate
	logger.Info("Exiting")
}
