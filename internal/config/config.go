/*
 * File: config.go
 * Description: Manages the application's configuration settings, loading them from files or environment variables.
 *              This module centralizes configuration handling to facilitate access to settings like database
 *              credentials, API keys, and service URLs.
 * Usage:
 *   - Load and access configuration settings throughout the application.
 * Dependencies:
 *   - TODO: Viper or a similar package for configuration management.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level configuration structure
type Config struct {
	App               AppConfig               `yaml:"app"`
	Database          DatabaseConfig          `yaml:"database"`
	CloudStorage      CloudStorageConfig      `yaml:"cloudStorage"`
	DataSources       DataSourcesConfig       `yaml:"dataSources"`
	IngestionSettings IngestionSettingsConfig `yaml:"ingestionSettings"`
	Notifications     NotificationsConfig     `yaml:"notifications"`
}

// AppConfig structure for general application settings
type AppConfig struct {
	Port     string `yaml:"port"`
	LogLevel string `yaml:"logLevel"`
}

// DatabaseConfig structure for DB settings
type DatabaseConfig struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"connectionString"`
}

// CloudStorageConfig structure for cloud storage settings
type CloudStorageConfig struct {
	BucketName      string `yaml:"bucketName"`
	CredentialsPath string `yaml:"credentialsPath"`
}

// DataSourcesConfig structure for external data sources
type DataSourcesConfig struct {
	Satellite DataSourceConfig `yaml:"satellite"`
	Weather   DataSourceConfig `yaml:"weather"`
	Soil      DataSourceConfig `yaml:"soil"`
}

// DataSourceConfig generalized for other data sources
type DataSourceConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Schedule    string `yaml:"schedule"`
	TimeZone    string `yaml:"timeZone"`
	HttpMethod  string `yaml:"httpMethod"`
	Endpoint    string `yaml:"endpoint"`
	APIKey      string `yaml:"apiKey"`
	Description string `yaml:"description"`
}

// IngestionSettingsConfig structure for ingestion settings
type IngestionSettingsConfig struct {
	RetryPolicy        RetryPolicyConfig `yaml:"retryPolicy"`
	ParallelIngestions int               `yaml:"parallelIngestions"`
}

// RetryPolicyConfig structure for configuring retry behavior
type RetryPolicyConfig struct {
	MaxRetries      int    `yaml:"maxRetries"`
	BackoffInterval string `yaml:"backoffInterval"`
}

// NotificationsConfig structure for notifications settings
type NotificationsConfig struct {
	EmailService EmailServiceConfig `yaml:"emailService"`
}

// EmailServiceConfig structure for email service settings
type EmailServiceConfig struct {
	Enabled   bool   `yaml:"enabled"`
	SMTPHost  string `yaml:"SMTPHost"`
	SMTPPort  int    `yaml:"SMTPPort"`
	Username  string `yaml:"Username"`
	Password  string `yaml:"Password"`
	FromEmail string `yaml:"FromEmail"`
}

// LoadConfig reads and unmarshals the YAML configuration file using os.ReadFile
func LoadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path) // using os.ReadFile instead of ioutil.ReadFile
	if err != nil {
		log.Fatalf("error reading config file at %s: %v", path, err)
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("error parsing config file: %v", err)
		return nil, err
	}

	return &config, nil
}
