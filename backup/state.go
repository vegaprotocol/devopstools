package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/vegaprotocol/devopstools/tools"
)

const EntryTimeFormat = time.RFC3339

type Entry struct {
	ID       string
	Date     string
	Block    string
	Origin   string
	Location string
	Full     bool
	Pools    []Pool
}

type State struct {
	Config Config

	Backups []Entry
}

func (state State) Empty() bool {
	return len(state.Backups) < 1
}

func (state State) SortedBackups() []Entry {
	entries := state.Backups
	// sort.Slice(people, func(i, j int) bool { return people[i].Name < people[j].Name })
	sort.Slice(entries, func(i, j int) bool {
		iTime, err := time.Parse(EntryTimeFormat, entries[i].Date)
		jTime, err2 := time.Parse(EntryTimeFormat, entries[j].Date)
		if err != nil || err2 != nil {
			return false
		}

		return iTime.Before(jTime)
	})

	return entries
}

func (state *State) AddEntry(id, location string, blockHeight int, full bool, pools []Pool) error {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get server hostname: %w", err)
	}

	if len(id) < 1 {
		return fmt.Errorf("backup id cannot be empty")
	}

	if len(pools) < 1 {
		pools = []Pool{}
	}

	entry := Entry{
		ID:       id,
		Location: location,
		Origin:   hostname,
		Block:    fmt.Sprintf("%d", blockHeight),
		Date:     time.Now().Format(EntryTimeFormat),
		Full:     full,
		Pools:    pools,
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
			Backups: []Entry{},
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
