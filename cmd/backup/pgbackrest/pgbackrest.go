package pgbackrest

import (
	"fmt"
	"os"
	"regexp"
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

const ConfigBackupFile = "/tmp/pgbackrest.conf.bk"

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

func BackupConfig(location string, force bool) error {
	if force {
		if err := os.RemoveAll(ConfigBackupFile); err != nil {
			return fmt.Errorf("failed to remove backup file when force flag is given: %w", err)
		}
	} else if tools.FileExists(ConfigBackupFile) {
		return nil
	}

	if _, err := tools.CopyFile(location, ConfigBackupFile); err != nil {
		return fmt.Errorf("failed to backup file: %w", err)
	}

	return nil
}

func RestoreConfigFromBackup(location string) error {
	if !tools.FileExists(ConfigBackupFile) {
		return fmt.Errorf("backup file does not exists")
	}

	if err := os.RemoveAll(location); err != nil {
		return fmt.Errorf("failed to config file: %w", err)
	}

	if _, err := tools.CopyFile(ConfigBackupFile, location); err != nil {
		return fmt.Errorf("failed to backup file: %w", err)
	}

	return nil
}

func UpdatePgbackrestConfig(logger *zap.Logger, location string, params map[string]string) error {
	configContent, err := os.ReadFile(location)
	if err != nil {
		return fmt.Errorf("failed to read pgbackrest config: %w", err)
	}

	for k, v := range params {
		logger.Debug("Update pgbackrest config", zap.String("key", k), zap.String("value", v))
		re := regexp.MustCompilePOSIX(fmt.Sprintf(`^[\\s]*%s[\\s]*=.*$`, k))
		configContent = re.ReplaceAll(configContent, []byte(fmt.Sprintf("%s=%s", k, v)))
	}

	logger.Debug("New config content", zap.String("content", string(configContent)))

	if err := os.WriteFile(location, configContent, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write new config content: %w", err)
	}

	return nil
}

func Check(logger zap.Logger, postgresqlUser, pgBackrestBinary string) error {
	consoleLevel := "info"
	if logger.Level() == zap.DebugLevel {
		consoleLevel = "debug"
	}

	args := []string{
		"check",
		"--log-level-console", consoleLevel,
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUserWithRealTimeLogs(postgresqlUser, pgBackrestBinary, args, func(outputType, logLine string) {
		logger.Debug(logLine, zap.String("command", pgBackrestBinary), zap.String("source", outputType))
	})

	if err != nil {
		return fmt.Errorf("failed to check pgbackrest: %w, command output: %s", err, out)
	}

	return nil
}

func CreateStanza(logger zap.Logger, postgresqlUser, pgBackrestBinary string) error {
	consoleLevel := "info"
	if logger.Level() == zap.DebugLevel {
		consoleLevel = "debug"
	}

	args := []string{
		"stanza-create",
		"--log-level-console", consoleLevel,
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUserWithRealTimeLogs(postgresqlUser, pgBackrestBinary, args, func(outputType, logLine string) {
		logger.Debug(logLine, zap.String("command", pgBackrestBinary), zap.String("source", outputType))
	})

	if err != nil {
		// stanza already exists but it is stopped
		if strings.Contains(err.Error(), "stop file exists for stanza main_archive") {
			return nil
		}

		return fmt.Errorf("failed to create pgbackrest stanza: %w, command output: %s", err, out)
	}

	return nil
}

func UpgradeStanza(logger zap.Logger, postgresqlUser, pgBackrestBinary string) error {
	consoleLevel := "info"
	if logger.Level() == zap.DebugLevel {
		consoleLevel = "debug"
	}

	args := []string{
		"stanza-upgrade",
		"--log-level-console", consoleLevel,
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUserWithRealTimeLogs(postgresqlUser, pgBackrestBinary, args, func(outputType, logLine string) {
		logger.Debug(logLine, zap.String("command", pgBackrestBinary), zap.String("source", outputType))
	})

	if err != nil {
		return fmt.Errorf("failed to upgrade pgbackrest stanza: %w, command output: %s", err, out)
	}

	return nil
}

func Backup(logger zap.Logger, postgresqlUser, pgBackrestBinary string, backupType BackupType) error {
	consoleLevel := "info"
	if logger.Level() == zap.DebugLevel {
		consoleLevel = "debug"
	}

	args := []string{
		"backup",
		"--log-level-console", consoleLevel,
		"--stanza", StanzaName,
		"--type", string(backupType),
	}

	out, err := tools.ExecuteBinaryAsUserWithRealTimeLogs(postgresqlUser, pgBackrestBinary, args, func(outputType, logLine string) {
		logger.Debug(logLine, zap.String("command", pgBackrestBinary), zap.String("source", outputType))
	})

	if err != nil {
		return fmt.Errorf("failed to perform a backup operation pgbackrest stanza: %w, command output: %s", err, out)
	}

	return nil
}

func Start(logger zap.Logger, postgresqlUser, pgBackrestBinary string) error {
	consoleLevel := "info"
	if logger.Level() == zap.DebugLevel {
		consoleLevel = "debug"
	}

	args := []string{
		"start",
		"--log-level-console", consoleLevel,
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUser(postgresqlUser, pgBackrestBinary, args, nil)

	if err != nil {
		// Ignore warning
		if strings.Contains(err.Error(), "WARN") && strings.Contains(string(out), "completed successfully") {
			return nil
		}

		return fmt.Errorf("failed to start pgbackrest stanza: %w, command output: %s", err, out)
	}

	return nil
}

func Stop(logger zap.Logger, postgresqlUser, pgBackrestBinary string) error {
	consoleLevel := "info"
	if logger.Level() == zap.DebugLevel {
		consoleLevel = "debug"
	}

	args := []string{
		"stop",
		"--log-level-console", consoleLevel,
		"--stanza", StanzaName,
	}

	out, err := tools.ExecuteBinaryAsUserWithRealTimeLogs(postgresqlUser, pgBackrestBinary, args, func(outputType, logLine string) {
		logger.Debug(logLine, zap.String("command", pgBackrestBinary), zap.String("source", outputType))
	})

	if err != nil {
		return fmt.Errorf("failed to stop pgbackrest stanza: %w, command output: %s", err, out)
	}

	return nil
}

func Info(logger zap.Logger, postgresqlUser, pgBackrestBinary string) (PgBackRestInfo, error) {
	consoleLevel := "info"
	if logger.Level() == zap.DebugLevel {
		consoleLevel = "debug"
	}

	result := PgBackRestInfo{}
	args := []string{
		"info",
		"--log-level-console", consoleLevel,
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

func Restore(logger zap.Logger, postgresqlUser, pgBackrestBinary, label string, delta bool) error {
	consoleLevel := "info"
	if logger.Level() == zap.DebugLevel {
		consoleLevel = "debug"
	}

	args := []string{
		"restore",
		// "--type", "none",
		"--log-level-console", consoleLevel,
		"--stanza", StanzaName,
		"--set", label,
		"--archive-mode", "off",
	}

	if delta {
		args = append(args, "--delta")
	}

	out, err := tools.ExecuteBinaryAsUserWithRealTimeLogs(postgresqlUser, pgBackrestBinary, args, func(outputType, logLine string) {
		logger.Debug(logLine, zap.String("command", pgBackrestBinary), zap.String("source", outputType))
	})

	if err != nil {
		return fmt.Errorf("failed to restore pgbackrest backup: %w, command output: %s", err, out)
	}

	return nil
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
