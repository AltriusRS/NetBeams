package types

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/altriusrs/netbeams/src/logs"
)

var App *Application

type Application struct {
	logs.Logger
	Services map[string]ServiceCompatible
}

// NewApplication creates a new application instance and returns it as a global reference
func NewApplication() {
	if App != nil {
		return
	}

	App = &Application{
		Logger:   logs.NetLogger("Resource Manager"),
		Services: make(map[string]ServiceCompatible),
	}
}

func (app *Application) AddService(s ServiceCompatible) {
	name := s.GetName()
	app.Infof("Adding service %s", name)
	app.Services[name] = s
}

func (app *Application) RemoveService(name string) error {
	service, ok := app.Services[name]
	if !ok {
		return errors.New("service not found")
	}

	err := service.StopService()

	if err != nil {
		return err
	}

	delete(app.Services, name)
	return nil
}

func (app *Application) GetService(name string) any {
	service, ok := app.Services[name]
	if !ok {
		return nil
	}

	return service
}

func (app *Application) Start() {
	app.Info("Starting NetBeam...")

	app.Info("Configuring UPnP")

	app.Info("Starting Server")

	app.Infof("Node ID: %s", app.ShortId)

	for name, service := range app.Services {
		status := service.GetStatus()

		// Skip services that are already running
		if status != nil {
			if *status == StatusStarting || *status == StatusHealthy {
				app.Infof("Skipping service %s as it is already running", name)
				continue
			}
		}

		app.Infof("Starting service %s", name)
		err := service.StartService()
		if err != nil {
			app.Fatal(err)
			app.Shutdown()
		}
	}
	app.Info("Server started")
}

func (app *Application) Shutdown() {
	app.Info("Shutting down...")

	for name, service := range app.Services {
		app.Infof("Stopping service %s", name)
		err := service.StopService()
		if err != nil {
			app.Fatal(err)
			app.Shutdown()
		}
	}

	app.Info("Shutdown complete")
}

func (app *Application) GetStatus(name string) *Status {
	service, ok := app.Services[name]
	if !ok {
		return nil
	}

	return service.GetStatus()
}

func (app *Application) Wait() {
	for {
		active := [][]string{}
		for name, service := range app.Services {

			status := service.GetStatus()

			if status != nil {
				if *status != StatusShutdown && *status != StatusErrored && *status != StatusRestarting && *status != StatusStopped {
					active = append(active, []string{name, status.String()})
				}
			}

		}

		if len(active) == 0 {
			break
		}

		time.Sleep(time.Second) // Sleep for 100ms before checking again
	}
}

func (app *Application) StartService(name string) error {
	service, ok := app.Services[name]
	if !ok {
		return errors.New("service not found")
	}

	service.StartService()

	return nil
}

// goroutine to handle signals and gracefully shutdown the application
func (app *Application) RegisterSignalHandler() {
	app.Info("Registering signal handler")
	channel := make(chan os.Signal, 1)

	exit := make(chan bool, 1)

	signal.Notify(channel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-channel
		app.Info("Received signal: " + sig.String())
		exit <- true
	}()

	go func() {
		<-exit
		app.Info("Shutdown signal handler started - Closing in 5 seconds")
		app.Info("Closing application")
		app.Shutdown()
	}()
}
