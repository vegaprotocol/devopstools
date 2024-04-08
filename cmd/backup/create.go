package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vegaprotocol/devopstools/backup"
	"github.com/vegaprotocol/devopstools/tools"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type CreateArgs struct {
	*Args

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
	createArgs.Args = &backupArgs

	createCmd.PersistentFlags().BoolVar(
		&createArgs.fullBackup,
		"full",
		false,
		"Create full backup")

	Cmd.AddCommand(createCmd)
}

func enforceFullBackup(logger *zap.Logger, state *backup.State, conf backup.FullBackupConfig) bool {
	// Empty state and corresponding config
	if conf.WhenEmptyState && state != nil && state.Empty() {
		logger.Info("Enforced full backup, empty state")
		return true
	}

	var (
		lastFullBackup   *backup.Entry
		backupsSinceFull int
	)
	entries := state.SortedBackups()

	for idx, entry := range entries {
		if entry.Full {
			backupsSinceFull = 0
			lastFullBackup = &(entries[idx])
			continue
		}

		backupsSinceFull = backupsSinceFull + 1
	}

	if lastFullBackup == nil {
		logger.Info("Enforced full backup, no full backups")
		return true
	}

	logger.Info(
		fmt.Sprintf(
			"Last full backup created %d backups ago, at %s",
			backupsSinceFull,
			lastFullBackup.Date,
		),
	)

	lastBackupTime, err := time.Parse(backup.EntryTimeFormat, lastFullBackup.Date)
	if err != nil {
		logger.Info("Enforced full backup, error parsing last full backup time", zap.Error(err))
		return true
	}

	lastExpectedFullBackup := time.Now().Add(-conf.EveryTimeDuration)
	if lastBackupTime.Before(lastExpectedFullBackup) {
		logger.Info(
			fmt.Sprintf(
				"Enforced full backup. Last full backup %s, full backup expected every %s",
				lastFullBackup.Date,
				conf.EveryTimeDuration,
			),
		)
		return true
	}

	if backupsSinceFull >= conf.EveryNBackups {
		logger.Info(
			fmt.Sprintf(
				"Enforced full backup. %d incremental backups since last full backup, full backup expected evet %d backups",
				backupsSinceFull,
				conf.EveryNBackups,
			),
		)
		return true
	}

	return false
}

func createBackup(logger *zap.Logger, configFile string, fullBackup bool) error {
	totalStart := time.Now()
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

	state := backup.OpenOrCreateNewState(config.StateFilePath, config)
	mustBeFullBackup := enforceFullBackup(logger, state, config.FullBackup) || fullBackup

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

	logger.Info("Getting ZFS pools", zap.String("file_system", config.PoolName))
	pools, err := backup.ListZfsPools(config.PoolName)
	if err != nil {
		return fmt.Errorf("failed to list zfs pools: %w", err)
	}
	logger.Info("Found ZFS pools", zap.Any("pools", pools))

	logger.Info(
		"Creating zfs snapsgot",
		zap.String("file_system", config.PoolName),
		zap.String("id", snapshotID),
	)
	if err := backup.CreateRecursiveZfsSnapshot(config.PoolName, snapshotID); err != nil {
		return fmt.Errorf("failed to create snapshot with ID %s: %w", snapshotID, err)
	}
	logger.Info(
		"Snapshot created",
		zap.String("file_system", config.PoolName),
		zap.String("id", snapshotID),
	)

	logger.Info("Preparing env variables for sending backup")
	zfsBackupPrepareEnvVariables(config.Destination)

	if mustBeFullBackup {
		logger.Info("Creating full backup")
	} else {
		logger.Info("Craeting incremental backup")
	}

	var (
		backupSizeZFS   uint64 = 0
		backupSizeTotal uint64 = 0
	)
	fullS3Destination := fmt.Sprintf(
		"s3://%s/%s/",
		strings.Trim(config.Destination.Bucket, "/"),
		strings.Trim(config.Destination.Path, "/"),
	)
	for _, pool := range pools {
		sendStart := time.Now()
		logger.Info("Sending backup, may take some time", zap.String("pool", pool.Name))
		result, err := zfsBackupSendBackup(
			config.ZfsBackupBinaryPath,
			pool,
			fullS3Destination,
			mustBeFullBackup,
		)
		if err != nil {
			return fmt.Errorf("failed to send backup to s3 for pool %s: %w", pool.Name, err)
		}
		sendEnd := time.Now()

		sendDuration := sendEnd.Sub(sendStart)
		logger.Info(
			"Backup sent",
			zap.String("duration", sendDuration.String()),
			zap.Int("zfs_size_mb", result.TotalZFSBytes/1024/1024),
			zap.Int("total_size_mb", result.TotalBackupBytes/1024/1024),
		)

		backupSizeZFS = backupSizeZFS + uint64(result.TotalZFSBytes)
		backupSizeTotal = backupSizeTotal + uint64(result.TotalBackupBytes)
	}

	if err := state.AddEntry(snapshotID, fullS3Destination, coreHeight, mustBeFullBackup, pools); err != nil {
		return fmt.Errorf("failed to add backup entry to state: %w", err)
	}
	logger.Info("Saving backup details in the state file", zap.String("path", config.StateFilePath))
	if err := state.Save(config.StateFilePath); err != nil {
		return fmt.Errorf("failed to save state into disk: %w", err)
	}

	totalEnd := time.Now()
	logger.Info(
		"Backup finished with no errors",
		zap.String("duration", totalEnd.Sub(totalStart).String()),
		zap.Uint64("zfs_size_mb", backupSizeZFS/1024/1024),
		zap.Uint64("total_size_mb", backupSizeTotal/1024/1024))

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
