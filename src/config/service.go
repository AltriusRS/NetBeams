package config

import (
	"os"
	"path/filepath"

	"github.com/altriusrs/netbeams/src/crypto"
	"github.com/altriusrs/netbeams/src/types"
	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml"
)

type ConfigService struct {
	types.Service
	configFile     string            // The path to the config file
	resourceFolder string            // The path to the resource folder
	watcher        *fsnotify.Watcher // The file watcher instance
	hash           *string           // The hash of the configuration file (calculated at startup, not on init)
}

func Service() *ConfigService {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	service := &ConfigService{
		Service:        types.SpinUp("Configuration"),
		configFile:     "ServerConfig.toml",
		resourceFolder: Configuration.General.ResourceFolder,
		watcher:        watcher,
		hash:           nil,
	}

	service.RegisterServiceHooks(service.Start, service.Stop, nil)

	return service
}

func (s *ConfigService) Start() (types.Status, error) {
	s.Debug("Calculating hash tables")

	// Check if the config file exists

	if _, err := os.Stat(s.configFile); err != nil {
		if os.IsNotExist(err) {
			s.Info("Configuration file not found - Creating default")
			f, err := os.Create(s.configFile)

			if err != nil {
				s.Error(err.Error())
				return types.StatusErrored, err
			}

			defer func() {
				_ = f.Close()
			}()

			// general := Configuration.General
			// misc := Configuration.Misc
			// netbeams := Configuration.NetBeams
			// authorization := Configuration.Auth

			encoder := toml.NewEncoder(f)
			encoder.Indentation("  ")
			encoder.ArraysWithOneElementPerLine(true)
			encoder.QuoteMapKeys(true)
			encoder.Order(toml.OrderPreserve)

			if err := encoder.Encode(LoadDefault()); err != nil {
				s.Error(err.Error())
				return types.StatusErrored, err
			}
		} else {
			s.Error(err.Error())
			return types.StatusErrored, err
		}
	}

	Load()

	s.resourceFolder = Configuration.General.ResourceFolder

	s.Debugf("folder %s", s.resourceFolder)
	s.Debugf("folder %s", s.resourceFolder+"/Client")
	s.Debugf("folder %s", s.resourceFolder+"/Server")

	if _, err := os.Stat(s.resourceFolder); err != nil {

		if err = os.MkdirAll(s.resourceFolder, 0755); err != nil {
			s.Error(err.Error())
			return types.StatusErrored, err
		}

		if err = os.MkdirAll(s.resourceFolder+"/Client", 0755); err != nil {
			s.Error(err.Error())
			return types.StatusErrored, err
		}

		if err = os.MkdirAll(s.resourceFolder+"/Server", 0755); err != nil {
			s.Error(err.Error())
			return types.StatusErrored, err
		}
	}

	if _, err := os.Stat(s.resourceFolder + "/Client"); err != nil {
		if err = os.MkdirAll(s.resourceFolder+"/Client", 0755); err != nil {
			s.Error(err.Error())
			return types.StatusErrored, err
		}
	}

	if _, err := os.Stat(s.resourceFolder + "/Server"); err != nil {
		if err = os.MkdirAll(s.resourceFolder+"/Server", 0755); err != nil {
			s.Error(err.Error())
			return types.StatusErrored, err
		}
	}

	resources := []string{}

	err := filepath.Walk(s.resourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		resources = append(resources, path)

		return nil
	})

	if err != nil {
		return types.StatusErrored, err
	}

	s.Debug("Hashing configuration file")
	hash, err := crypto.HashFile(s.configFile)

	if err != nil {
		s.Error(err.Error())
		return types.StatusErrored, err
	}

	s.hash = hash
	s.Debugf("Hash for configuration file: %s", *s.hash)

	s.Debug("Adding watcher for config file")
	_ = s.watcher.Add(s.configFile)

	s.SetStatus(types.StatusStarting)

	go s.Watch()

	return types.StatusHealthy, nil
}

func (s *ConfigService) Stop() (types.Status, error) {
	_ = s.watcher.Close()

	return types.StatusShutdown, nil
}

func (s *ConfigService) Watch() {
	for *s.Status == types.StatusHealthy || *s.Status == types.StatusStarting {
		select {
		case event := <-s.watcher.Events:
			s.OnFileChange(event)
		case err := <-s.watcher.Errors:
			s.Error(err.Error())
		}
	}
}

func (s *ConfigService) OnFileChange(event fsnotify.Event) {

	if event.Name == s.configFile {
		check, _, err := crypto.CompareFileHashes(s.configFile, *s.hash)

		if err != nil {
			s.Error(err.Error())
			return
		}

		if !check {
			s.Info("Changes detected - Reloading configuration file")
			Load()
		} else {
			s.Debug("File modification detected but no changes present")
		}
	}
}
