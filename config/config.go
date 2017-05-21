package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config defines the application configuration
type Config struct {
	HTTP    HTTPConfig
	Storage StorageConfig
	Queue   QueueConfig
	SMTP    SMTPConfig
	Attachments AttachmentConfig
	Logging LoggingConfig
}

// HTTPConfig defines the configuration for the http api
type HTTPConfig struct {
	Port int
}

// StorageConfig defines the database storage information
type StorageConfig struct {
	Driver string
	Config string
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

// Config for mail attachments
type AttachmentConfig struct {
	Path string
}

// Config for everything about logging
type LoggingConfig struct {
	Sentry SentryConfig
	File FileLoggingConfig
}

// Specific Sentry logging config
type SentryConfig struct {
	DSN string
	Level string
}

// Specifig File logging config
type FileLoggingConfig struct {
	Path string
	Level string
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
