package config

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	configReloadFrequencySeconds = 3
)

var (
	_configfile       = "config.yml"
	_md5sum           = []byte{}
	_config           *CogConfig
	_lastReloadWorked = true // Prevents spam
)

// GetBundleConfigs returns the data wrapper for the "bundles" config section.
func GetBundleConfigs() []BundleConfig {
	return _config.BundleConfigs
}

// GetDockerConfigs returns the data wrapper for the "docker" config section.
func GetDockerConfigs() DockerConfigs {
	return _config.DockerConfigs
}

// GetGlobalConfigs returns the data wrapper for the "global" config section.
func GetGlobalConfigs() GlobalConfigs {
	return _config.GlobalConfigs
}

// GetSlackProviders returns the data wrapper for the "slack" config section.
func GetSlackProviders() []SlackProvider {
	return _config.SlackProviders
}

// Runs at startup but causes problems during testing.
// TODO replace with a call that points to a config doc specidfied by the user.
func init() {
	err := executeFullConfigurationReload()
	if err != nil {
		panic(err.Error())
	}

	ticker := time.NewTicker(configReloadFrequencySeconds * time.Second)

	go func() {
		for range ticker.C {
			err := executeFullConfigurationReload()
			if err != nil {
				if _lastReloadWorked {
					_lastReloadWorked = false
					log.Println(err.Error())
				}
			}
		}
	}()
}

func executeFullConfigurationReload() error {
	md5sum, err := getMd5Sum(_configfile)
	if err != nil {
		return fmt.Errorf("Failed hash file %s: %s", _configfile, err.Error())
	}

	if !slicesAreEqual(_md5sum, md5sum) {
		cp, err := loadConfiguration(_configfile)
		if err != nil {
			return fmt.Errorf("Failed to load config %s: %s", _configfile, err.Error())
		}

		_md5sum = md5sum
		_config = cp
		_lastReloadWorked = true

		log.Printf("Loaded configuration file %s\n", _configfile)
	}

	return nil
}

func getMd5Sum(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return []byte{}, err
	}

	hashBytes := hasher.Sum(nil)

	return hashBytes, nil
}

func loadConfiguration(file string) (*CogConfig, error) {
	// Read file as a byte slice
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cp CogConfig
	err = yaml.Unmarshal(dat, &cp)
	if err != nil {
		return nil, err
	}

	return &cp, nil
}

func slicesAreEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
