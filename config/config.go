package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var file = "config"

type (
	Config struct {
		v    *viper.Viper
		path string
	}
)

func New() (*Config, error) {
	// Resolve $HOME/.config/app
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %w", err)
	}

	dir := filepath.Join(home, ".config", "daedalus")
	path := filepath.Join(dir, file)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	v := viper.New()

	v.SetConfigName(file)
	v.AddConfigPath(path)

	v.SetConfigType("dotenv")

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("read config: %w", err)
		}
	} else if err := os.Chmod(path, 0o600); err != nil {
		return nil, fmt.Errorf("secure config file: %w", err)
	}

	return &Config{v, path}, nil
}

func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *Config) Set(key string, value any) {
	c.v.Set(key, value)
}

func (c *Config) Save() error {
	dir := filepath.Dir(c.path)
	tmp, err := os.CreateTemp(dir, ".config-*.dotenv")
	if err != nil {
		return fmt.Errorf("create temporary config: %w", err)
	}

	tmpPath := tmp.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("secure temporary config: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("write temporary config: %w", err)
	}

	if err := c.v.WriteConfigAs(tmpPath); err != nil {
		return fmt.Errorf("write temporary config: %w", err)
	}

	if err := os.Rename(tmpPath, c.path); err != nil {
		return fmt.Errorf("replace config: %w", err)
	}

	return nil
}
