package backup

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pelletier/go-toml"
	"github.com/vegaprotocol/devopstools/tools"
)

type S3Config struct {
	Endpoint               string `toml:"endpoint"`
	Region                 string `toml:"region"`
	Bucket                 string `toml:"bucket"`
	Path                   string `toml:"path"`
	AccessKeyIdEnvName     string `toml:"access_key_id_env_name"`
	AccessKeySecretEnvName string `toml:"access_key_secret_env_name"`
}

type Config struct {
	Destination S3Config `toml:"destination"`

	CoreRestURL         string `toml:"core_rest_url"`
	FileSystem          string `toml:"file_system"`
	StateFilePath       string `toml:"state_file"`
	ZfsBackupBinaryPath string `toml:"zfsbackup_binary"`
}

func LoadConfig(filePath string) (*Config, error) {
	if !tools.FileExists(filePath) {
		return nil, fmt.Errorf("config file does not exists")
	}

	configContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	result := &Config{}
	if err := toml.Unmarshal(configContent, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return result, nil
}

func (conf Config) Check() error {
	if conf.CoreRestURL == "" {
		return fmt.Errorf("core_rest_url in the config cannot be empty")
	}

	if conf.FileSystem == "" {
		return fmt.Errorf("file_system in the config cannot be empty")
	}

	if _, err := exec.LookPath(conf.ZfsBackupBinaryPath); err != nil {
		return fmt.Errorf("the given zfsbackup-go binary could not be found: install it from https://github.com/someone1/zfsbackup-go")
	}

	if conf.Destination.Bucket == "" {
		return fmt.Errorf("destination.bucket in the config cannot be empty")
	}

	if conf.Destination.Region == "" {
		return fmt.Errorf("destination.region in the config cannot be empty")
	}

	if conf.Destination.Path == "" {
		return fmt.Errorf("destination.path in the config cannot be empty")
	}

	if conf.Destination.Endpoint == "" {
		return fmt.Errorf("destination.endpoint in the config cannot be empty")
	}

	if conf.Destination.AccessKeyIdEnvName == "" {
		return fmt.Errorf("destination.access_key_id_env_name in the config cannot be empty")
	}

	if conf.Destination.AccessKeySecretEnvName == "" {
		return fmt.Errorf("destination.access_secret_env_name in the config cannot be empty")
	}

	accessKey, _ := os.LookupEnv(conf.Destination.AccessKeyIdEnvName)
	if len(accessKey) < 1 {
		return fmt.Errorf("env variable '%s' not set", conf.Destination.AccessKeyIdEnvName)
	}

	accessSecret, _ := os.LookupEnv(conf.Destination.AccessKeySecretEnvName)
	if len(accessSecret) < 1 {
		return fmt.Errorf("env variable '%s' not set", conf.Destination.AccessKeySecretEnvName)
	}

	return nil
}
