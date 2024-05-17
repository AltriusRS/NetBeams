package server

import (
	"errors"
	"netbeams/config"
	"netbeams/globals"
	"netbeams/logs"
	"netbeams/tcp"
	"time"
)

type Application struct {
	Config    config.BaseConfig
	TCPServer tcp.Server
	Logger    logs.Logger

	tasks map[string]*globals.Status
}

func NewApplication(config config.BaseConfig, l *logs.Logger) Application {
	return Application{
		Config: config,
		Logger: l.Fork("Resource Manager"),
		tasks:  make(map[string]*globals.Status),
	}
}

func (app *Application) StartMainNode() {
	app.Logger.Info("Starting NetBeam...")
	app.Logger.Info("Starting Server")
	app.Logger.Infof("Node ID: %s", app.Logger.ShortId)

	// Spawn the TCP server manager instance
	app.TCPServer = tcp.NewServer(app.Config.General.Port, &app.Logger, func(s globals.Status) {
		app.SetStatus("tcp", &s)
	})

	// Start the TCP server
	sucess := app.TCPServer.Start()

	if !sucess {
		app.Logger.Fatal(errors.New("failed to start tcp server"))
	}

	app.Logger.Info("Server started")
}

func (app *Application) Shutdown() {
	app.Logger.Info("Shutting down...")
	app.TCPServer.Stop()
	app.Logger.Info("Shutdown complete")
}

func (app *Application) GetStatus(name string) *globals.Status {
	return app.tasks[name]
}

func (app *Application) SetStatus(name string, status *globals.Status) {
	app.Logger.Debugf("Status of service '%s' changed : %s -> %s", name, app.tasks[name], *status)
	app.tasks[name] = status
}

func (app *Application) Wait() {
	for {
		active := []string{}
		for name, status := range app.tasks {
			if status == nil {
				continue
			}
			if *status != globals.Shutdown && *status != globals.Errored {
				active = append(active, name)
			}
		}

		time.Sleep(time.Second) // Sleep for 100ms before checking again
		if len(active) == 0 {
			break
		}
	}
}
