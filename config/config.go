package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config defines the application configuration
type Config struct {
	HTTP    HTTPConfig
	Storage StorageConfig
}

// HTTPConfig defines the configration for the http api
type HTTPConfig struct {
	Port int
}

// StorageConfig defines the database storage information
type StorageConfig struct {
	Driver string
	Config string
}

// ReadConfig reads the json configuration values into the Config Struct
func ReadConfig(file string) (config *Config, err error) {
	fileContent, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileContent, &config)

	if err != nil {
		return nil, err
	}

	return config, err
}
