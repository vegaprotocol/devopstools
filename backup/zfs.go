package backup

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/vegaprotocol/devopstools/tools"
)

type Pool struct {
	Name       string
	MountPoint string
}

func CheckZfsCommand() error {
	if _, err := exec.LookPath("zfs"); err != nil {
		return fmt.Errorf("zfs command not found: %w", err)
	}

	return nil
}

func ListZfsPools(parent string) ([]Pool, error) {
	args := []string{
		"list",
		"-o", "name,mountpoint",
	}
	out, err := tools.ExecuteBinary("zfs", args, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list zfs pools: %w", err)
	}
	result := []Pool{}

	tokens := strings.Fields(string(out))
	if len(tokens)%2 != 0 {
		return nil, fmt.Errorf("failed to parse output: got invalid number of tokens")
	}

	// Skip first 2 tokens because it is header from the example output:
	// NAME                            MOUNTPOINT
	// vega_pool                       /vega_pool
	// vega_pool/home                  /home/vega
	for i := 2; i < len(tokens); i += 2 {
		if !strings.HasPrefix(tokens[i], parent) {
			continue
		}

		result = append(result, Pool{
			Name:       tokens[i],
			MountPoint: tokens[i+1],
		})
	}

	return result, nil
}

func CreateRecursiveZfsSnapshot(pool string, ID string) error {
	args := []string{
		"snapshot",
		"-r", fmt.Sprintf("%s@%s", pool, ID),
	}
	if out, err := tools.ExecuteBinary("zfs", args, nil); err != nil {
		return fmt.Errorf("failed to create snapshot, out: %s: %w", out, err)
	}

	return nil
}
