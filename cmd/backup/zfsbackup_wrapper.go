package backup

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/vegaprotocol/devopstools/backup"
	"github.com/vegaprotocol/devopstools/tools"
)

type SendBackupOutput struct {
	TotalZFSBytes    int `json:"TotalZFSBytes"`
	TotalBackupBytes int `json:"TotalBackupBytes"`
	ElapsedTime      int `json:"ElapsedTime"`
	FilesUploaded    int `json:"FilesUploaded"`
}

func SendBackup(bin string, pool backup.Pool, s3Destination string, full bool) (*SendBackupOutput, error) {
	if !strings.HasPrefix(s3Destination, "s3://") {
		return nil, fmt.Errorf("destination must start with s3://")
	}

	cores := runtime.NumCPU()

	// The command is: zfsbackup-go send --full vega_pool  s3://vega-internal-tm-postgres-backups/n00.devnet1.vega.xyz
	args := []string{"send", "--jsonOutput", "--numCores", fmt.Sprintf("%d", cores)}
	if full {
		args = append(args, "--full")
	} else {
		args = append(args, "--increment")
	}
	args = append(args, pool.Name, s3Destination)

	result := &SendBackupOutput{}
	if _, err := tools.ExecuteBinary(bin, args, result); err != nil {
		return nil, fmt.Errorf("failed to send backup: %w", err)
	}
	fmt.Printf("result: %#v\n\n", result)

	return result, nil
}
