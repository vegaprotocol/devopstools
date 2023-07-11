package vegacmd

import (
	"fmt"
	"sort"

	"github.com/vegaprotocol/devopstools/tools"
)

type CoreSnashotInput struct {
	VegaHome       string
	SnapshotDbPath string
}

func ListCoreSnapshots(vegaBinary string, input CoreSnashotInput) (*CoreToolsSnapshots, error) {
	result := &CoreToolsSnapshots{}

	args := []string{
		"tools",
		"snapshot",
		"--output", "json",
	}
	if input.SnapshotDbPath != "" {
		args = append(args, "--db-path", input.SnapshotDbPath)
	} else {
		args = append(args, "--home", input.VegaHome)
	}

	if _, err := tools.ExecuteBinary(vegaBinary, args, result); err != nil {
		return nil, fmt.Errorf("failed to execute vega %v: %w", args, err)
	}

	return result, nil
}

func LatestCoreSnapshot(vegaBinary string, input CoreSnashotInput) (*CoreToolsSnapshot, error) {
	snapshots, err := ListCoreSnapshots(vegaBinary, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}

	snapshotsList := snapshots.Snapshots
	sort.Slice(snapshotsList, func(i, j int) bool {
		return snapshotsList[i].Height > snapshotsList[j].Height
	})

	if len(snapshotsList) < 1 {
		return nil, fmt.Errorf("no core snapshots found")
	}

	return &snapshotsList[0], nil
}
