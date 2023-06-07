package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/backup"
	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
)

type CreateArgs struct {
	*BackupArgs

	fullBackup bool
}

var createArgs CreateArgs

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create backup.",

	Run: func(cmd *cobra.Command, args []string) {
		if err := createBackup(createArgs.Logger, createArgs.configPath, createArgs.fullBackup); err != nil {
			createArgs.Logger.Fatal("failed to create backup", zap.Error(err))
		}
	},
}

func init() {
	createArgs.BackupArgs = &backupArgs

	createCmd.PersistentFlags().BoolVar(
		&createArgs.fullBackup,
		"full",
		false,
		"Create full backup")

	BackupCmd.AddCommand(createCmd)
}

func createBackup(logger *zap.Logger, configFile string, fullBackup bool) error {
	if user, _ := tools.WhoAmI(); user != "root" {
		return fmt.Errorf("you must run this command from the root user, current user: %s", user)
	}

	logger.Info("Loading config", zap.String("file", configFile))
	config, err := backup.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config file path: %w", err)
	}

	if err := config.Check(); err != nil {
		return fmt.Errorf("failed to check config: %w", err)
	}

	logger.Info("Loading node height", zap.String("url", config.CoreRestURL))
	coreHeight, err := getCoreHeight(config.CoreRestURL)
	if err != nil {
		return fmt.Errorf("failed to get core height: %w", err)
	}
	logger.Info("Found node height", zap.Int("height", coreHeight))

	snapshotID := fmt.Sprintf("%s-%d", time.Now().Format("20060102_150405"), coreHeight)
	logger.Info("Generate snapshot ID", zap.String("id", snapshotID))

	logger.Info("Checking ZFS command")
	if err := backup.CheckZfsCommand(); err != nil {
		return fmt.Errorf("failed to check zfs: %w", err)
	}

	logger.Info("Getting ZFS pools", zap.String("file_system", config.FileSystem))
	pools, err := backup.ListZfsPools(config.FileSystem)
	if err != nil {
		return fmt.Errorf("failed to list zfs pools: %w", err)
	}
	logger.Info("Found ZFS pools", zap.Any("pools", pools))

	logger.Info("Creating zfs snapsgot", zap.String("file_system", config.FileSystem), zap.String("id", snapshotID))
	if err := backup.CreateRecursiveZfsSnapshot(config.FileSystem, snapshotID); err != nil {
		return fmt.Errorf("failed to create snapshot with ID %s: %w", snapshotID, err)
	}
	logger.Info("Snapshot created", zap.String("file_system", config.FileSystem), zap.String("id", snapshotID))

	logger.Info("Preparing env variables for sending backup")
	zfsBackupPrepareEnvVariables(config.Destination)

	fullS3Destination := fmt.Sprintf("s3://%s/%s", strings.Trim(config.Destination.Bucket, "/"), strings.Trim(config.Destination.Path, "/"))
	for _, pool := range pools {
		sendStart := time.Now()
		logger.Info("Sending backup, may take some time", zap.String("pool", pool.Name))
		result, err := zfsBackupSendBackup(config.ZfsBackupBinaryPath, pool, fullS3Destination, fullBackup)
		if err != nil {
			return fmt.Errorf("failed to send backup to s3 for pool %s: %w", pool.Name, err)
		}
		sendEnd := time.Now()

		sendDuration := sendEnd.Sub(sendStart)
		logger.Info("Backup sent", zap.String("duration", sendDuration.String()), zap.Int("zfs_size_mb", result.TotalZFSBytes/1024/1024), zap.Int("total_size_mb", result.TotalBackupBytes/1024/1024))
	}

	state := backup.OpenOrCreateNewState(config.StateFilePath, config)
	if err := state.AddEntry(snapshotID, fullS3Destination, coreHeight, pools); err != nil {
		return fmt.Errorf("failed to add backup entry to state: %w", err)
	}
	if err := state.Save(config.StateFilePath); err != nil {
		return fmt.Errorf("failed to save state into disk: %w", err)
	}

	return nil
}

type statisticsResponse struct {
	Statistics struct {
		BlockHeight string `json:"blockHeight"`
	} `json:"statistics"`
}

func getCoreHeight(restUrl string) (int, error) {
	if !strings.HasPrefix(restUrl, "http") {
		restUrl = fmt.Sprintf("http://%s", restUrl)
	}

	resp, err := http.Get(fmt.Sprintf("%s/statistics", restUrl))
	if err != nil {
		return 0, fmt.Errorf("failed to get node height: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read http response: %w", err)
	}

	respStruct := &statisticsResponse{}
	if err := json.Unmarshal(respBytes, respStruct); err != nil {
		return 0, fmt.Errorf("failed to unmarshal /statistics response: %w", err)
	}

	height, err := strconv.Atoi(respStruct.Statistics.BlockHeight)
	if err != nil {
		return 0, fmt.Errorf("failed to convert height into int: %w", err)
	}
	return height, nil
}
