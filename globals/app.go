package globals

import (
	"errors"
	"netbeams/config"
	"netbeams/logs"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var App *Application

type Application struct {
	Config   config.BaseConfig
	Logger   logs.Logger
	Services map[string]ServiceCompatible
}

// NewApplication creates a new application instance and returns it as a global reference
func NewApplication(config config.BaseConfig, l *logs.Logger) {
	if App != nil {
		return
	}

	App = &Application{
		Config:   config,
		Logger:   l.Fork("Resource Manager"),
		Services: make(map[string]ServiceCompatible),
	}
}

func (app *Application) AddService(s ServiceCompatible) {
	name := s.GetName()
	app.Logger.Infof("Adding service %s", name)
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
	app.Logger.Info("Starting NetBeam...")
	app.Logger.Info("Starting Server")
	app.Logger.Infof("Node ID: %s", app.Logger.ShortId)

	for name, service := range app.Services {
		app.Logger.Infof("Starting service %s", name)
		err := service.StartService()
		if err != nil {
			app.Logger.Fatal(err)
			app.Shutdown()
		}
	}
	app.Logger.Info("Server started")
}

func (app *Application) Shutdown() {
	app.Logger.Info("Shutting down...")

	for name, service := range app.Services {
		app.Logger.Infof("Stopping service %s", name)
		err := service.StopService()
		if err != nil {
			app.Logger.Fatal(err)
			app.Shutdown()
		}
	}

	app.Logger.Info("Shutdown complete")
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
				if *status != Shutdown && *status != Errored && *status != Restarting && *status != Stopped {
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
	app.Logger.Info("Registering signal handler")
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
