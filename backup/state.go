package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/tools"
)

type BackupEntry struct {
	ID       string
	Date     string
	Block    string
	Origin   string
	Location string
}

type State struct {
	Config Config

	Backups []BackupEntry
}

func (state *State) AddEntry(id, location string, blockHeight int) error {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get server hostname: %w", err)
	}

	if len(id) < 1 {
		return fmt.Errorf("backup id cannot be empty")
	}

	entry := BackupEntry{
		ID:       id,
		Location: location,
		Origin:   hostname,
		Block:    fmt.Sprintf("%d", blockHeight),
		Date:     time.Now().Format(time.RFC3339),
	}

	state.Backups = append(state.Backups, entry)
	return nil
}

func (state State) Save(filePath string) error {
	content, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal backups state")
	}

	if err := os.WriteFile(filePath, content, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write state into file: %w", err)
	}

	return nil
}

func OpenOrCreateNewState(filePath string, config *Config) *State {
	state, err := LoadStateFromFile(filePath)
	if err != nil || state == nil {
		state = &State{
			Backups: []BackupEntry{},
		}
	}

	if config != nil {
		state.Config = *config
	}

	return state
}

func LoadStateFromFile(filePath string) (*State, error) {
	if !tools.FileExists(filePath) {
		return nil, fmt.Errorf("state file does not exist")
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	state := &State{}
	if err := json.Unmarshal(fileBytes, state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal given file into state struct: %w", err)
	}

	return state, nil
}
