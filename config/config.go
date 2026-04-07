package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var file = "config"

type (
	Config struct {
		v *viper.Viper
	}
)

func New() (*Config, error) {
	// Resolve $HOME/.config/app
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "daedalus")
	configFile := filepath.Join(configDir, file)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	v := viper.New()

	v.SetConfigName(file)
	v.SetConfigType("dotenv")
	v.AddConfigPath(configDir)

	var fileLookupError viper.ConfigFileNotFoundError
	if err := v.ReadInConfig(); err != nil {
		if errors.As(err, &fileLookupError) {
			// Create the config file if it doesn't exist
			if err := v.WriteConfigAs(configFile); err != nil {
				return nil, fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	return &Config{v}, nil
}

func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *Config) Set(key string, value any) {
	c.v.Set(key, value)
}

func (c *Config) Save() error {
	return c.v.WriteConfig()
}
