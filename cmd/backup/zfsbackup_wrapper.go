package backup

import (
	"fmt"
	"os"
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

func zfsBackupPrepareEnvVariables(config backup.S3Config) error {
	// Variable names are defined by 3-rd party libary in the zfsbackup-go program

	// We except some env variables are missing except KEY ID and KEY SECRET
	if err := os.Setenv("AWS_REGION", config.Region); err != nil {
		return fmt.Errorf("failed set the AWS_REGION env variable")
	}

	if err := os.Setenv("AWS_S3_CUSTOM_ENDPOINT", config.Endpoint); err != nil {
		return fmt.Errorf("failed set AWS_S3_CUSTOM_ENDPOINT env variable")
	}

	// zfsbackup-go expects `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` as env name for the S3 API credentials, We have to remap it
	if config.AccessKeyIdEnvName != "AWS_ACCESS_KEY_ID" {
		accessKeyId := os.Getenv(config.AccessKeyIdEnvName)
		if err := os.Setenv("AWS_ACCESS_KEY_ID", accessKeyId); err != nil {
			return fmt.Errorf("failed set AWS_ACCESS_KEY_ID env variable")
		}
	}

	if config.AccessKeySecretEnvName != "AWS_SECRET_ACCESS_KEY" {
		accessKeySecret := os.Getenv(config.AccessKeySecretEnvName)
		if err := os.Setenv("AWS_SECRET_ACCESS_KEY", accessKeySecret); err != nil {
			return fmt.Errorf("failed set AWS_SECRET_ACCESS_KEY env variable")
		}
	}

	return nil
}

func zfsBackupSendBackup(bin string, pool backup.Pool, s3Destination string, full bool) (*SendBackupOutput, error) {
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
	return result, nil
}
