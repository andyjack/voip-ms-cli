package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

// https://benaiah.me/posts/configuring-go-apps-with-toml

var configDirName = "voip-ms-cli"

// GetDefaultConfigDir returns a string containing the path to the config dir.
// On most platforms return the expanded ~/.config, but on linux maybe respect
// XDG_CONFIG_HOME env var.
func getDefaultConfigDir() (string, error) {
	var configDirLocation string

	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	xdgconfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgconfig != "" && runtime.GOOS == "linux" {
		configDirLocation = xdgconfig
	} else {
		configDirLocation = filepath.Join(homeDir, ".config", configDirName)
	}

	return configDirLocation, nil
}

// Config specifies our needed config for talking to the voip.ms API
type Config struct {
	Credentials credentials
}

// Credentials section of the config
type credentials struct {
	Email    string
	Password string
}

// LoadConfig returns a pointer to a Config
func loadConfig(filename string) (*Config, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, errors.New("Config file " + filename + " does not exist!")
	} else if err != nil {
		return nil, err
	}

	var c Config
	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return nil, err
	}
	if c.Credentials.Email == "" || c.Credentials.Password == "" {
		return nil, errors.New("config is missing credentials.email or password")
	}
	return &c, nil
}
