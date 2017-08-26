package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config defines the application configuration
type Config struct {
	HTTP    HTTPConfig
	MongoDB MongoDBConfig
	Queue   QueueConfig
	SMTP    SMTPConfig
	Files   FilesConfig
}

// HTTPConfig defines the configuration for the http api
type HTTPConfig struct {
	Port int
}

// QueueConfig defines the rabbit mq connection information
type QueueConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	QueueName string
}

// SMTPConfig defines the SMTP server settings
type SMTPConfig struct {
	Host     string
	Port     int
	Identity string
	User     string
	Password string
}

type FilesConfig struct {
	AttachmentPath string
	BodyPath       string
}

type MongoDBConfig struct {
	Address  []string
	Database string
	User     string
	Password string
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
