package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	// name of the file, where the configuration will reside
	ConfigFileName = ".dotconfig"
)

var (
	PathDotConfig = fmt.Sprintf("%s/%s", HomeDir(), ConfigFileName)
)

type Config struct {
	// folder where the files that are tracked will reside
	DotPath string `json:"dot_path"`

	// map with the individual files that are being tracked
	Files map[string]string `json:"files"`
}

// Constructor for the Config struct
func NewConfig(path string) (*Config, error) {
	c := &Config{}
	err := c.load(path)
	return c, err
}

// Pointer receiver for the Config struct load the config file
func (c *Config) load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(file).Decode(&c); err != nil {
		return err
	}

	return nil
}

// Pointer receiver for the config struct that will save the config file
func (c *Config) Save() error {
	// truncate existing file if it exists
	f, err := os.Create(PathDotConfig)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}
