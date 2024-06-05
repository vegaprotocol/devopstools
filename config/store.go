package config

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	vgfs "code.vegaprotocol.io/vega/libs/fs"
	"code.vegaprotocol.io/vega/paths"
)

var ErrFileNotFound = errors.New("file not found")

func Load(ctx context.Context, path string) (Config, error) {
	// Get file from the remote
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		var err error
		path, err = downloadConfig(ctx, path)
		if err != nil {
			return Config{}, fmt.Errorf("cannot download file from the url: %s: %w", path, err)
		}
		defer func() {
			// File was downloaded but We do not want to keep local copy
			// in case we are running it in the CI and server is compromised.
			if err := os.RemoveAll(path); err != nil {
				panic(err)
			}
		}()
	}

	return loadLocalFile(path)
}

func downloadConfig(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request with context: %w", err)
	}

	// For github We want to add authorization token if available
	if strings.HasPrefix(url, "https://github.com") || strings.HasPrefix(url, "https://raw.githubusercontent.com") {
		req.Header.Add("Accept", "application/vnd.github.v3.raw")

		token := os.Getenv("GITHUB_TOKEN")
		if len(token) > 0 {
			req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
		}
	}

	cli := http.DefaultClient
	resp, err := cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get config from remote server: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid response code: expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

	tempFile, err := os.CreateTemp(os.TempDir(), "devopstools-config")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file for downloaded config: %w", err)
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save downloaded content into temp file: %w", err)
	}

	stat, err := tempFile.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat created file with config: %w", err)
	}

	return filepath.Join(os.TempDir(), stat.Name()), nil
}

func loadLocalFile(path string) (Config, error) {
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
