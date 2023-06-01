package tools

import (
	"os"
	"os/user"
	"runtime"
)

func UserHomeDir() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}

	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

func CurrentUsername() string {
	if currentUser, err := user.Current(); err == nil {
		return currentUser.Username
	}

	return ""
}
