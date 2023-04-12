package tools

import (
	"fmt"
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
