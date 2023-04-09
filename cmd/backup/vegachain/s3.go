package vegachain

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vegaprotocol/devopstools/tools"
)

type S3Credentials struct {
	Endpoint     string
	Region       string
	AccessKey    string
	AccessSecret string
}

func CheckS3Setup(s3CmdBinary string) error {
	if _, err := tools.ExecuteBinary(s3CmdBinary, []string{"--version"}, nil); err != nil {
		return fmt.Errorf("s3cmd command is missing")
	}

	return nil
}

// We use S3Cmd because golang implementation of s3 sync does not exists.
func S3CmdInit(s3CmdBinary string, creds S3Credentials) error {
	args := []string{
		"--access_key", creds.AccessKey,
		"--secret_key", creds.AccessSecret,
		"--ssl",
		"--no-encrypt",
		"--dump-config",
		"--host", creds.Endpoint,
		"--host-bucket", fmt.Sprintf("%%(bucket)s.%s", creds.Endpoint),
	}

	s3cmdConfig, err := tools.ExecuteBinary(s3CmdBinary, args, nil)
	if err != nil {
		return fmt.Errorf("failed to generate s3cmd config: %w", err)
	}

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user's home directory: %w", err)
	}

	s3ConfigFilePath := filepath.Join(homeDirectory, ".s3cfg")
	if err := os.WriteFile(s3ConfigFilePath, s3cmdConfig, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write s3cmd config: %w", err)
	}

	if _, err := tools.ExecuteBinary(s3CmdBinary, []string{"ls"}, nil); err != nil {
		return fmt.Errorf("given permissions \"%s\" do not gave acces to s3", creds.AccessKey)
	}

	return nil
}

func S3Sync(s3CmdBinary, source, destination string, debug bool) error {
	args := []string{
		"sync",
		source,
		destination,
		"--follow-symlinks",
		"--stop-on-error",
		"--preserve",
		"--delete-removed",
		"--recursive",
		// "--check-md5", "--skip-existing", // TODO: need verification
	}
	if !debug {
		args = append(args, "--quiet")
	}

	if _, err := tools.ExecuteBinary(s3CmdBinary, args, nil); err != nil {
		return fmt.Errorf("failed to sync files from \"%s\" to \"%s\": %w", source, destination, err)
	}

	return nil
}

func S3Copy(s3CmdBinary, source, destination string, debug bool) error {
	args := []string{
		"cp",
		source,
		destination,
		"--follow-symlinks",
		"--stop-on-error",
		"--preserve",
		"--recursive",
		// "--check-md5", "--skip-existing", // TODO: need verification
	}
	if !debug {
		args = append(args, "--quiet")
	}

	if _, err := tools.ExecuteBinary(s3CmdBinary, args, nil); err != nil {
		return fmt.Errorf("failed to copy files from \"%s\" to \"%s\": %w", source, destination, err)
	}

	return nil
}
