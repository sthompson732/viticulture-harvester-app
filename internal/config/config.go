package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration structure.
type Config struct {
	App struct {
		Name     string `yaml:"name"`
		LogLevel string `yaml:"logLevel"`
		Port     string `yaml:"port"`
	} `yaml:"app"`
	Database struct {
		Type             string `yaml:"type"`
		ConnectionString string `yaml:"connectionString"`
	} `yaml:"database"`
	CloudStorage struct {
		BucketName      string `yaml:"bucketName"`
		CredentialsPath string `yaml:"credentialsPath"`
	} `yaml:"cloudStorage"`
	DataSources struct {
		Satellite struct {
			Enabled           bool   `yaml:"enabled"`
			APIEndpoint       string `yaml:"apiEndpoint"`
			APIKey            string `yaml:"apiKey"`
			IngestionInterval string `yaml:"ingestionInterval"`
		} `yaml:"satellite"`
		Soil struct {
			Enabled           bool   `yaml:"enabled"`
			APIEndpoint       string `yaml:"apiEndpoint"`
			APIKey            string `yaml:"apiKey"`
			IngestionInterval string `yaml:"ingestionInterval"`
		} `yaml:"soil"`
		// Add other data sources as needed
	} `yaml:"dataSources"`
	// Additional configurations can be added here as needed
}

// LoadConfig reads and unmarshals the YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading config file at %s: %v", path, err)
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("error parsing config file: %v", err)
		return nil, err
	}

	return &config, nil
}
