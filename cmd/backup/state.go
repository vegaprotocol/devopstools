package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/tendermint/tendermint/libs/json"
	"github.com/vegaprotocol/devopstools/cmd/backup/pgbackrest"
	"github.com/vegaprotocol/devopstools/tools"
)

type BackupStatus string

const (
	BackupStatusFailed     BackupStatus = "failed"
	BackupStatusSuccess    BackupStatus = "success"
	BackupStatusInProgress BackupStatus = "in-progress"
	BackupStatusUnknown    BackupStatus = "unknown"
)

type PgBackrestEntry struct {
	Status BackupStatus
	Type   pgbackrest.BackupType
	Label  string

	Started  time.Time
	Finished time.Time
}

type VegaChainEntry struct {
	Location struct {
		Bucket string
		Path   string
	}

	Started    time.Time
	Finished   time.Time
	Status     BackupStatus
	Components struct {
		VegaHome       bool
		TendermintHome bool
		VisorHome      bool
	}
}

type BackupEntry struct {
	ID         uuid.UUID
	Started    time.Time
	Finished   time.Time
	ServerHost string
	Status     BackupStatus
	Postgresql PgBackrestEntry
	VegaChain  VegaChainEntry
}

type State struct {
	encryptionKey    []byte
	LastUpdated      time.Time
	Backups          map[string]BackupEntry
	PgBackrestConfig string

	Locked bool

	localFilePath string
	_remoteS3Path string
}

func (state State) DeepClone() State {
	result := state

	result.Backups = make(map[string]BackupEntry)

	for k, v := range state.Backups {
		result.Backups[k] = v
	}

	return result
}

func (state *State) AddOrModifyEntry(entry BackupEntry, writeLocal bool) error {
	if entry.ID == uuid.Nil {
		return fmt.Errorf("ID for backup entry cannot be empty")
	}

	state.Backups[entry.ID.String()] = entry

	if err := state.WriteLocal(state.localFilePath); err != nil {
		return fmt.Errorf("failed to write backup state to local file(2): %w", err)
	}
	return nil
}

func (state State) AsJSON() (string, error) {
	jsonStr, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal state into JSON: %w", err)
	}

	return string(jsonStr), nil
}

func (state *State) WriteLocal(filePath string) error {
	state.localFilePath = filePath
	state.LastUpdated = time.Now()

	stateCopy := state.DeepClone()

	// Enrypt pgbackrest config before saving it as a plain text
	if stateCopy.PgBackrestConfig != "" && !tools.IsEncrypted(stateCopy.PgBackrestConfig) {
		var err error

		stateCopy.PgBackrestConfig, err = tools.EncryptMessage(stateCopy.encryptionKey, stateCopy.PgBackrestConfig)
		if err != nil {
			return fmt.Errorf("failed to encrypt pgbackrest config: %w", err)
		}
	}

	jsonState, err := stateCopy.AsJSON()
	if err != nil {
		return fmt.Errorf("failed to write state into local file: failed to convert state into json: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(jsonState), os.ModePerm); err != nil {
		return fmt.Errorf("failed to save state into file: %w", err)
	}

	return nil
}

func (state *State) WriteRemote() error {
	// TODO: Implement it
	return nil
}

func (state *State) LoadFromRemote() error {
	// TODO: Implement it
	return nil
}

func LoadFromRemote() (State, error) {
	return State{
		Backups: map[string]BackupEntry{},
	}, nil
}

func LoadFromLocal(encryptionKey, location string) (State, error) {
	result := State{}
	content, err := os.ReadFile(location)
	if err == nil {
		if err := json.Unmarshal(content, &result); err != nil {
			return result, fmt.Errorf("failed to unmarshal file: %w", err)
		}

		result.localFilePath = location
		result.encryptionKey = []byte(encryptionKey)

		// Decrypt pgbackrest config
		if result.PgBackrestConfig != "" && tools.IsEncrypted(result.PgBackrestConfig) && encryptionKey != "" {
			result.PgBackrestConfig, err = tools.DecryptMessage([]byte(encryptionKey), result.PgBackrestConfig)
			if err != nil {
				return result, fmt.Errorf("failed to decrypt pgbackrest config from the state: %w", err)
			}
		}

		return result, nil
	}

	return result, fmt.Errorf("failed to read state file from local: %w", err)
}

func NewEmptyState(encryptionKey string) State {
	return State{
		Backups:       map[string]BackupEntry{},
		encryptionKey: []byte(encryptionKey),
	}
}

// LoadOrCreateNew tries to load state from file otherwise it returns empty state
// Fails only when it should, and program cannot proceed with an error
func LoadOrCreateNew(encryptionKey string, locaLocation string) (State, error) {
	if err := tools.ValidateEncryptionKey(encryptionKey); err != nil {
		return State{}, fmt.Errorf("failed to check the encryption key: %w", err)
	}

	// todo: add support for s3
	localState, err := LoadFromLocal(encryptionKey, locaLocation)
	if err != nil {
		return NewEmptyState(encryptionKey), nil
	}

	localState.localFilePath = locaLocation
	return localState, nil
}

func NewBackupEntry() (BackupEntry, error) {
	var err error

	result := BackupEntry{}
	result.ID = uuid.New()
	result.Started = time.Now()
	result.Status = BackupStatusInProgress
	result.Postgresql.Status = BackupStatusUnknown
	result.VegaChain.Status = BackupStatusUnknown

	result.ServerHost, err = os.Hostname()
	if err != nil {
		return result, fmt.Errorf("failed to get server hostname: %w", err)
	}

	return result, nil
}
