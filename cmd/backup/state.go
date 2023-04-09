package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/tendermint/tendermint/libs/json"
	"github.com/vegaprotocol/devopstools/cmd/backup/pgbackrest"
)

type BackupStatus string

const (
	BackupStatusFailed     BackupStatus = "failed"
	BackupStatusSuccess    BackupStatus = "success"
	BackupStatusInProgress BackupStatus = "in-progress"
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

	Started  time.Time
	Finished time.Time
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
	LastUpdated time.Time
	Backups     map[string]BackupEntry

	Locked bool

	localFilePath string
	remoteS3Path  string
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
	jsonState, err := state.AsJSON()
	if err != nil {
		return fmt.Errorf("failed to write state into local file: failed to convert state into json: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(jsonState), os.ModePerm); err != nil {
		return fmt.Errorf("failed to save state into file: %w", err)
	}

	return nil
}

func (state *State) WriteRemote() error {
	return nil
}

func (state *State) UpdateFromRemote() error {
	return nil
}

func (state *State) UpdateFromLocal() error {
	return nil
}

func LoadFromRemote() (State, error) {
	return State{
		Backups: map[string]BackupEntry{},
	}, nil
}

func LoadFromLocal(location string) (State, error) {
	result := State{}
	content, err := os.ReadFile(location)
	if err == nil {
		if err := json.Unmarshal(content, &result); err != nil {
			return result, fmt.Errorf("failed to unmarshal file: %w", err)
		}

		return result, nil
	}

	return result, fmt.Errorf("failed to read state file from local: %w", err)
}

func NewEmptyState() State {
	return State{
		Backups: map[string]BackupEntry{},
	}
}

// LoadOrCreateNew tries to load state from file otherwise it returns empty state
func LoadOrCreateNew(locaLocation string) State {
	// todo: add support for s3
	localState, err := LoadFromLocal(locaLocation)
	if err != nil {
		return NewEmptyState()
	}

	return localState
}

func NewBackupEntry() (BackupEntry, error) {
	var err error

	result := BackupEntry{}
	result.ID = uuid.New()
	result.Started = time.Now()
	result.Status = BackupStatusInProgress

	result.ServerHost, err = os.Hostname()
	if err != nil {
		return result, fmt.Errorf("failed to get server hostname: %w", err)
	}

	return result, nil
}
