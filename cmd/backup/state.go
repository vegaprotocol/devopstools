package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/tendermint/tendermint/libs/json"
)

type BackupStatus string

const (
	BackupStatusFailed     = "failed"
	BackupStatusSuccess    = "success"
	BackupStatusInProgress = "in-progress"
)

type PgBackrestEntry struct {
	Status BackupStatus

	Started  time.Time
	Finished time.Time
}

type BackupEntry struct {
	Started    time.Time
	Finished   time.Time
	ServerHost string
	Status     BackupStatus
}

type State struct {
	LastUpdated time.Time
	Backups     []BackupEntry

	Locked bool

	localFilePath string
	remoteS3Path  string
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
	return State{}, nil
}

func LoadFromLocal(location string) (State, error) {
	result := State{}
	content, err := os.ReadFile(location)
	if err == nil {

		if err := json.Unmarshal(content, &result); err != nil {
			return result, fmt.Errorf("failed to unmarshal file: %w", err)
		}
	}

	return result, fmt.Errorf("failed to read state file from local: %w", err)
}

func NewEmptyState() State {
	return State{}
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
