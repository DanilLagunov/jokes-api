package config

import (
	"flag"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config struct.
type Config struct {
	Server struct {
		Port    string `yaml:"port"`
		Timeout struct {
			ReadHeader time.Duration `yaml:"read_header"`
			Read       time.Duration `yaml:"read"`
			Write      time.Duration `yaml:"write"`
		} `yaml:"timeout"`
	} `yaml:"server"`
	Database struct {
		URI                 string `yaml:"uri"`
		DBName              string `yaml:"db_name"`
		JokesCollectionName string `yaml:"jokes_collection"`
	} `yaml:"db"`
}

// NewConfig creating a new Config object.
func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ParseFlags parse -config flag from CLI.
func ParseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	flag.Parse()

	return configPath, nil
}
