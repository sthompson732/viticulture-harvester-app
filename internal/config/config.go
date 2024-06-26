/*
 * config.go: Manages configuration settings loaded from YAML files.
 * Centralizes access to all configurable parameters of the application.
 * Usage: Used at startup to configure services and parameters.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */
package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App               AppConfig                   `yaml:"app"`
	Database          DatabaseConfig              `yaml:"database"`
	CloudStorage      CloudStorageConfig          `yaml:"cloudStorage"`
	DataSources       map[string]DataSourceConfig `yaml:"dataSources"` // Changed to a map
	IngestionSettings IngestionSettingsConfig     `yaml:"ingestionSettings"`
	Notifications     NotificationsConfig         `yaml:"notifications"`
	ProjectID         string                      `yaml:"projectID"`
	LocationID        string                      `yaml:"locationID"`
	ValidAPIKeys      []string                    `yaml:"validApiKeys"`
}

type AppConfig struct {
	Port     string `yaml:"port"`
	LogLevel string `yaml:"logLevel"`
}

type DatabaseConfig struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"connectionString"`
}

type CloudStorageConfig struct {
	BucketName      string `yaml:"bucketName"`
	CredentialsPath string `yaml:"credentialsPath"`
}

// DataSourceConfig generalized for all data sources
type DataSourceConfig struct {
	Enabled     bool              `yaml:"enabled"`
	Schedule    string            `yaml:"schedule"`
	TimeZone    string            `yaml:"timeZone"`
	HttpMethod  string            `yaml:"httpMethod"`
	Endpoint    string            `yaml:"endpoint"`
	APIKey      string            `yaml:"apiKey"`
	Params      map[string]string `yaml:"params"`  // Additional parameters for API calls
	Headers     map[string]string `yaml:"headers"` // Custom headers for API calls
	Body        string            `yaml:"body"`    // Added for PUT/PATCH/POST
	Description string            `yaml:"description"`
}

type IngestionSettingsConfig struct {
	RetryPolicy        RetryPolicyConfig `yaml:"retryPolicy"`
	ParallelIngestions int               `yaml:"parallelIngestions"`
}

type RetryPolicyConfig struct {
	MaxRetries      int    `yaml:"maxRetries"`
	BackoffInterval string `yaml:"backoffInterval"`
}

type NotificationsConfig struct {
	EmailService EmailServiceConfig `yaml:"emailService"`
}

type EmailServiceConfig struct {
	Enabled   bool   `yaml:"enabled"`
	SMTPHost  string `yaml:"SMTPHost"`
	SMTPPort  int    `yaml:"SMTPPort"`
	Username  string `yaml:"Username"`
	Password  string `yaml:"Password"`
	FromEmail string `yaml:"FromEmail"`
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file at %s: %v", path, err)
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
		return nil, err
	}

	return &config, nil
}
