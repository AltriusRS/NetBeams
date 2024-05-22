package config

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/altriusrs/netbeams/src/types"
	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml"
)

type ConfigService struct {
	types.Service
	configFile     string // The path to the config file
	resourceFolder string // The path to the resource folder
	watcher        *fsnotify.Watcher
	hashes         map[string]string
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
		hashes:         make(map[string]string),
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

			defer f.Close()

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

	resources = append(resources, s.configFile)

	for _, resource := range resources {
		s.Debugf("Hashing %s", resource)
		file, err := os.Open(resource)
		if err != nil {
			s.Error(err.Error())
			continue
		}

		defer file.Close()

		hasher := md5.New()

		if _, err := io.Copy(hasher, file); err != nil {
			s.Error(err.Error())
			continue
		}

		checksum := hasher.Sum(nil)
		hr := hex.EncodeToString(checksum)

		s.hashes[resource] = hr
		s.Debugf("Hash for %s: %s", resource, s.hashes[resource])
	}

	s.Debug("Adding watcher for config file")
	s.watcher.Add(s.configFile)

	s.Debug("Adding watcher for resource folder")
	s.watcher.Add(s.resourceFolder)

	s.SetStatus(types.StatusStarting)

	go s.Watch()

	return types.StatusHealthy, nil
}

func (s *ConfigService) Stop() (types.Status, error) {
	s.watcher.Close()

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
	s.Debugf("File %s changed - %s", event.Name, event.Op.String())

	if event.Name == s.configFile {
		s.Info("Changes detected - Reloading configuration file")
		Load()

		// } else if strings.HasPrefix(event.Name, s.resourceFolder) {
		// TODO: Trigger resource manager reload
	}
}
