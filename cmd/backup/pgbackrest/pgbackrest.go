package pgbackrest

import (
	"fmt"
	"os"
	"strings"

	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

type BackupType string

const (
	BackupFull        BackupType = "full"
	BackupIncremental BackupType = "incr"
)

const StanzaName = "main_archive"

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

type PgBackrestBackupInfo struct {
	Type      string `json:"type"`
	Label     string `json:"label"`
	Error     bool   `json:"error"`
	Timestamp struct {
		Start uint64 `json:"start"`
		Stop  uint64 `json:"stop"`
	} `json:"timestamp"`
}

type PgBackRestInfo []struct {
	Backup []PgBackrestBackupInfo `json:"backup"`
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

func ReadRawConfig(location string) (string, error) {
	if !tools.FileExists(location) {
		return "", fmt.Errorf("pgbackrest config(%s) does not exists", location)
	}

	content, err := os.ReadFile(location)
	if err != nil {
		return "", fmt.Errorf("failed to read pgbackrest config file: %w", err)
	}

	return string(content), nil
}

func Check(logger zap.Logger, postgresqlUser, pgBackrestBinary string) error {
	args := []string{
		"check",
		"--log-level-console", "info",
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUser(postgresqlUser, pgBackrestBinary, args, nil)
	logger.Debug(string(out))
	if err != nil {
		return fmt.Errorf("failed to check pgbackrest: %w", err)
	}

	return nil
}

func CreateStanza(logger zap.Logger, postgresqlUser, pgBackrestBinary string) error {
	args := []string{
		"stanza-create",
		"--log-level-console", "info",
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUser(postgresqlUser, pgBackrestBinary, args, nil)
	logger.Debug(string(out))

	if err != nil {
		// stanza already exists but it is stopped
		if strings.Contains(err.Error(), "stop file exists for stanza main_archive") {
			return nil
		}

		return fmt.Errorf("failed to create pgbackrest stanza: %w, command output: %s", err, out)
	}

	return nil
}

func Backup(logger zap.Logger, postgresqlUser, pgBackrestBinary string, backupType BackupType) error {
	args := []string{
		"backup",
		"--log-level-console", "info",
		"--stanza", StanzaName,
		"--type", string(backupType),
	}

	out, err := tools.ExecuteBinaryAsUser(postgresqlUser, pgBackrestBinary, args, nil)
	logger.Debug(string(out))
	if err != nil {
		return fmt.Errorf("failed to perform a backup operation pgbackrest stanza: %w, command output: %s", err, out)
	}

	return nil
}

func Start(logger zap.Logger, postgresqlUser, pgBackrestBinary string) error {
	args := []string{
		"start",
		"--log-level-console", "info",
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUser(postgresqlUser, pgBackrestBinary, args, nil)
	logger.Debug(string(out))
	if err != nil {
		return fmt.Errorf("failed to start pgbackrest stanza: %w, command output: %s", err, out)
	}

	return nil
}

func Stop(logger zap.Logger, postgresqlUser, pgBackrestBinary string) error {
	args := []string{
		"stop",
		"--log-level-console", "info",
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUser(postgresqlUser, pgBackrestBinary, args, nil)
	logger.Debug(string(out))
	if err != nil {
		return fmt.Errorf("failed to stop pgbackrest stanza: %w, command output: %s", err, out)
	}

	return nil
}

func Info(logger zap.Logger, postgresqlUser, pgBackrestBinary string) (PgBackRestInfo, error) {
	result := PgBackRestInfo{}
	args := []string{
		"info",
		"--log-level-console", "info",
		"--output", "json",
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUser(postgresqlUser, pgBackrestBinary, args, &result)
	logger.Debug(string(out))
	if err != nil {
		return result, fmt.Errorf("failed to get info about pgbackrest stanza: %w, command output: %s", err, out)
	}

	return result, nil
}

func Restore(logger zap.Logger, postgresqlUser, pgBackrestBinary, label string, delta bool) (PgBackRestInfo, error) {
	result := PgBackRestInfo{}
	args := []string{
		"restore",
		"--type", "none",
		"--log-level-console", "info",
		"--stanza", StanzaName,
		"--set", label,
		"--archive-mode", "off",
	}

	if delta {
		args = append(args, "--delta")
	}

	out, err := tools.ExecuteBinaryAsUser(postgresqlUser, pgBackrestBinary, args, &result)
	logger.Debug(string(out))
	if err != nil {
		return result, fmt.Errorf("failed to restore pgbackrest backup: %w, command output: %s", err, out)
	}

	return result, nil
}

func LastPgBackRestBackupInfo(info PgBackRestInfo, onlySuccessfull bool) *PgBackrestBackupInfo {
	if len(info) < 1 || len(info[0].Backup) < 1 {
		return nil
	}

	var result *PgBackrestBackupInfo
	for idx, backupInfo := range info[0].Backup {
		if result != nil && result.Timestamp.Stop >= backupInfo.Timestamp.Stop {
			continue
		}

		if backupInfo.Error && onlySuccessfull {
			continue
		}

		result = &info[0].Backup[idx]
	}

	return result
}
