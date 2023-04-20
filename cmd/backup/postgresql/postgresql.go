package postgresql

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/ini.v1"
)

type Config struct {
	DataDirectory string `ini:"data_directory"`
}

func ReadConfig(location string) (*Config, error) {
	resultConfig := &Config{}

	if _, err := os.Stat(location); err != nil {
		return resultConfig, fmt.Errorf("failed to check if postgresql config(%s) exists: %w", location, err)
	}

	configData, err := ini.Load(location)
	if err != nil {
		return resultConfig, fmt.Errorf("failed to read postgresql file: %w", err)
	}

	if err := configData.MapTo(&resultConfig); err != nil {
		return resultConfig, fmt.Errorf("failed to unmarshal postgresql config: %w", err)
	}

	return resultConfig, nil
}
