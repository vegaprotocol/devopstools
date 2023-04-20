package tools

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func ChownR(path, userName, groupName string) error {
	userInfo, err := user.Lookup(userName)
	if err != nil {
		return fmt.Errorf("failed to lookup for user %s: %w", userName, err)
	}

	uID, err := strconv.Atoi(userInfo.Uid)
	if err != nil {
		return fmt.Errorf("failed to convert uid from string to int: %w", err)
	}

	groupInfo, err := user.LookupGroup(groupName)
	if err != nil {
		return fmt.Errorf("failed to lookup for group %s: %w", groupName, err)
	}

	gID, err := strconv.Atoi(groupInfo.Gid)
	if err != nil {
		return fmt.Errorf("failed to convert gid from string to int: %w", err)
	}

	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uID, gID)
		}
		return err
	})
}

func RemoveDirectoryContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return fmt.Errorf("failed to open directory: %w", err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("failed to read directory content: %w", err)
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("failed to remove the %s directory: %w", name, err)
		}
	}

	return nil
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

	return nBytes, fmt.Errorf("failed to copy content of the source file to destination: %w", err)
}
