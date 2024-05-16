package server

import (
	"errors"
	"netbeams/config"
	"netbeams/logs"
	"time"
)

type Status int

const (
	Idle Status = iota
	Starting
	Healthy
	Stopping
	Stopped
	Shutdown
	Errored
)

func (s Status) String() string {
	switch s {
	case Idle:
		return "Idle"
	case Starting:
		return "Starting"
	case Healthy:
		return "Healthy"
	case Stopping:
		return "Stopping"
	case Stopped:
		return "Stopped"
	case Shutdown:
		return "Shutdown"
	case Errored:
		return "Errored"
	default:
		return "Unknown"
	}
}

type Application struct {
	Config    config.BaseConfig
	TCPServer TCPServer
	Logger    logs.Logger

	tasks map[string]*Status
}

func NewApplication(config config.BaseConfig, l *logs.Logger) Application {
	return Application{
		Config: config,
		Logger: l.Fork("Resource Manager"),
		tasks:  make(map[string]*Status),
	}
}

func (app *Application) StartMainNode() {
	app.Logger.Info("Starting NetBeam...")
	app.Logger.Info("Starting Server")
	app.Logger.Infof("Node ID: %s", app.Logger.MachineID)

	// Spawn the TCP server manager instance
	app.TCPServer = NewTCPServer(app.Config.General.Port, &app.Logger, func(s Status) {
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

func (app *Application) GetStatus(name string) *Status {
	return app.tasks[name]
}

func (app *Application) SetStatus(name string, status *Status) {
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
			if *status != Shutdown && *status != Errored {
				active = append(active, name)
			}
			// app.Logger.Debugf("Service '%s' is %s", name, *status)
		}
		// app.Logger.Debugf("Active services: %d", len(active))
		time.Sleep(time.Second) // Sleep for 100ms before checking again
		if len(active) == 0 {
			break
		}
	}
}
