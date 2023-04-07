package pgbackrest

import (
	"fmt"
	"os"

	"github.com/vegaprotocol/devopstools/tools"
	"gopkg.in/ini.v1"
)

// Example config
// [global]
// repo1-retention-full-type=count
// repo1-retention-full=3
// repo1-type=s3
// repo1-path=/stagnet1/be02.stagnet1.vega.xyz-2023-04-07_19-48
// repo1-s3-region=fra1
// repo1-s3-endpoint=fra1.digitaloceanspaces.com
// repo1-s3-bucket=XXXXXX
// repo1-s3-key=XXXXX
// repo1-s3-key-secret=XXXXXXXX
// log-level-file=debug
// log-level-console=debug

// [global:archive-push]
// compress-level=5

// [main_archive]
// pg1-path=/var/lib/postgresql/14/main

type PgBackrestConfig struct {
	Global struct {
		R1Type        string `ini:"repo1-type"`
		R1Path        string `ini:"repo1-path"`
		R1S3Region    string `ini:"repo1-s3-region"`
		R1S3Endpoint  string `ini:"repo1-s3-endpoint"`
		R1S3Bucket    string `ini:"repo1-s3-bucket"`
		R1S3Key       string `ini:"repo1-s3-key"`
		R1S3KeySecret string `ini:"repo1-s3-key-secret"`
		MainArchive   struct {
			Pg1Path string `ini:"pg1-path"`
		} `ini:"main_archive"`
	} `ini:"global"`
}

func CheckPgBackRestSetup(pgBackrestBinary string, config PgBackrestConfig) error {
	if config.Global.R1Type != "s3" {
		return fmt.Errorf("invalid repository type(repo1-type): got %s, only s3 supported", config.Global.R1Type)
	}

	if _, err := tools.ExecuteBinary("pgbackrest", []string{"version"}, nil); err != nil {
		return fmt.Errorf("failed to execute pgbackrest binary: %w", err)
	}

	return nil
}

func ReadConfig(location string) (PgBackrestConfig, error) {
	resultConfig := PgBackrestConfig{}

	if _, err := os.Stat(location); err != nil {
		return resultConfig, fmt.Errorf("failed to check if pgbackrest config(%s) exists: %w", location, err)
	}

	configData, err := ini.Load(location)
	if err != nil {
		return resultConfig, fmt.Errorf("failed to read pgbackrest file: %w", err)
	}

	if err := configData.MapTo(&resultConfig); err != nil {
		return resultConfig, fmt.Errorf("failed to unmarshal pgbackrest config: %w", err)
	}

	return resultConfig, nil
}
