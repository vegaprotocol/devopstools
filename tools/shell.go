package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

// To execute this function you have to run your program with syscall permissions to switch user
func ExecuteBinaryAsUser(userName, binaryPath string, args []string, v interface{}) ([]byte, error) {
	command := exec.Command(binaryPath, args...)

	u, err := user.Lookup(userName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command %s %v as a %s: %w", binaryPath, args, userName, err)
	}
	uid, err := strconv.ParseInt(u.Uid, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user id from string to int: %w", err)
	}
	gid, err := strconv.ParseInt(u.Gid, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert group id from string to int: %w", err)
	}

	command.SysProcAttr = &syscall.SysProcAttr{}
	command.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}

	return executeCmd(command, v)
}

func ExecuteBinary(binaryPath string, args []string, v interface{}) ([]byte, error) {
	command := exec.Command(binaryPath, args...)

	return executeCmd(command, v)
}

func executeCmd(command *exec.Cmd, v interface{}) ([]byte, error) {
	var stdOut, stErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stErr

	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute binary %s %v with error: %s: %s", command.Path, command.Args, stErr.String(), err.Error())
	}

	if v == nil {
		return stdOut.Bytes(), nil
	}

	if err := json.Unmarshal(stdOut.Bytes(), v); err != nil {
		// TODO Maybe failback to text parsing instead??
		return nil, err
	}

	return nil, nil
}
