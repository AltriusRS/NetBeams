package main

import (
	"fmt"

	"github.com/altriusrs/netbeams/src/config"
	"github.com/altriusrs/netbeams/src/environment"
	"github.com/altriusrs/netbeams/src/heartbeat"
	"github.com/altriusrs/netbeams/src/http"
	"github.com/altriusrs/netbeams/src/logs"
	"github.com/altriusrs/netbeams/src/netcheck"
	"github.com/altriusrs/netbeams/src/player_manager"
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
	configuration := config.Service()
	configuration.StartService()

	// Spawn a new application instance
	types.NewApplication()

	// Spawn the netcheck service if required
	if config.Configuration.Auth.Proxy.Enable || config.Configuration.Auth.VPN.Enable {
		logger.Info("Proxy or VPN authentication checking enabled - Loading databases (this will use a lot of memory)")
		logger.Info("Loading IP2Proxy databases")
		netchecker := netcheck.Service()
		failed := netchecker.StartService()
		if failed != nil {
			logger.Error("Failed to start IP2Proxy databases")
			logger.Error(failed.Error())
			return
		}
		logger.Info("IP2Proxy databases loaded")
		types.App.AddService(netchecker)
	}

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

	// Pass application to signal handler to allow graceful shutdown
	types.App.RegisterSignalHandler()

	// Spawn the required services
	types.App.AddService(configuration)
	types.App.AddService(http.Service())
	types.App.AddService(player_manager.Service())
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
