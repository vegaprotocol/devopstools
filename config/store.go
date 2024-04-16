package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	vgfs "code.vegaprotocol.io/vega/libs/fs"
	"code.vegaprotocol.io/vega/paths"
)

var ErrFileNotFound = errors.New("file not found")

func Load(path string) (Config, error) {
	found, err := vgfs.FileExists(path)
	if err != nil {
		return Config{}, fmt.Errorf("could not verify file presence: %w", err)
	}
	if !found {
		return Config{}, ErrFileNotFound
	}

	cfg := Config{}

	if err := paths.ReadStructuredFile(path, &cfg); err != nil {
		return Config{}, err
	}

	_, fileName := filepath.Split(path)

	rawName := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))

	cfg.Name = NetworkName(rawName)

	return cfg, nil
}

func SaveConfig(path string, cfg Config) error {
	if err := paths.WriteStructuredFile(path, cfg); err != nil {
		return err
	}
	return nil
}