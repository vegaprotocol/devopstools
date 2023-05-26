package tools

import (
	"fmt"
	"io"
	"os"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, fmt.Errorf("failed to stat source file: %w", err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("failed to read source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)

	if err != nil {
		return 0, fmt.Errorf("failed to copy content of the source file to destination: %w", err)
	}
	return nBytes, nil
}
