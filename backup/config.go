package backup

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/vegaprotocol/devopstools/tools"

	"github.com/pelletier/go-toml"
)

type S3Config struct {
	Endpoint               string `toml:"endpoint" comment:"The S3 endpoint"`
	Region                 string `toml:"region" comment:"The S3 region"`
	Bucket                 string `toml:"bucket" comment:"The S3 bucket name where backup is located"`
	Path                   string `toml:"path" comment:"The path on S3 where the backup is located. It is composed with bucket name into s3://<bucket>/<path>"`
	AccessKeyIdEnvName     string `toml:"access_key_id_env_name" comment:"Environment variable name used to source access key id for the S3 API"`
	AccessKeySecretEnvName string `toml:"access_key_secret_env_name" comment:"Environment variable name used to source access key secret for the S3 API"`
}

type FullBackupConfig struct {
	WhenEmptyState    bool          `toml:"when_empty_state" comment:"Full backup is enforced when there is no prevous backups in the state file"`
	EveryNBackups     int           `toml:"every_n_backups" comment:"Defines how ofter(every N backups) full backup must be created"`
	EveryTimeDuration time.Duration `toml:"every_time_duration" comment:"Define how often(in term of times) full backup must be created"`
}

type Config struct {
	Destination S3Config `toml:"destination" comment:"Parent file system to backup. We backup all inherited pools as well.\n\n Let's say We have the following zfs pools:\n\tNAME\n\tvega_pool\n\tvega_pool/home\n\tvega_pool/home/network-history\n\tvega_pool/home/postgresql\n\tvega_pool/home/tendermint_home\n\n When We provide file_system = \"vega_pool\", We backup all of the above zfs pools"`

	FullBackup FullBackupConfig `toml:"full_backup"`

	CoreRestURL         string `toml:"core_rest_url" comment:"Core REST URL where the /statistics endpoint is available"`
	PoolName            string `toml:"pool_name" comment:"Pool name to backup, We backup all inherited pools as well"`
	StateFilePath       string `toml:"state_file" comment:"Path where the state is saved"`
	ZfsBackupBinaryPath string `toml:"zfsbackup_binary" comment:"Path for the binary https://github.com/someone1/zfsbackup-go"`
}

func DefaultConfig() Config {
	return Config{
		CoreRestURL:         "localhost:3003",
		PoolName:            "vega_pool",
		StateFilePath:       "./state.json",
		ZfsBackupBinaryPath: "zfsbackup-go",

		Destination: S3Config{
			AccessKeyIdEnvName:     "AWS_ACCESS_KEY_ID",
			AccessKeySecretEnvName: "AWS_SECRET_ACCESS_KEY",
		},

		FullBackup: FullBackupConfig{
			WhenEmptyState:    true,
			EveryNBackups:     7,
			EveryTimeDuration: time.Hour * 72,
		},
	}
}

func LoadConfig(filePath string) (*Config, error) {
	if !tools.FileExists(filePath) {
		return nil, fmt.Errorf("config file does not exists")
	}

	configContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	result := DefaultConfig()
	if err := toml.Unmarshal(configContent, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &result, nil
}

func (conf Config) Marshal() (string, error) {
	str, err := toml.Marshal(conf)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	return string(str), nil
}

func (conf Config) Check() error {
	if conf.CoreRestURL == "" {
		return fmt.Errorf("core_rest_url in the config cannot be empty")
	}

	if conf.PoolName == "" {
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
