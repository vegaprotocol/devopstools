package tools

import (
	"bufio"
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

func ExecuteBinaryWithExitCode(binaryPath string, args []string, v interface{}) (int, []byte, error) {
	command := exec.Command(binaryPath, args...)

	out, err := executeCmd(command, v)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), out, err
		}
	}

	return 0, out, err
}

func ExecuteBinaryAsUserWithRealTimeLogs(userName string, binaryPath string, args []string, logParser func(outputType, logLine string)) ([]byte, error) {
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

	return executeBinaryWithRealTimeLogs(command, logParser)
}

func ExecuteBinaryWithRealTimeLogs(binaryPath string, args []string, logParser func(outputType, logLine string)) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)

	return executeBinaryWithRealTimeLogs(cmd, logParser)
}

func executeBinaryWithRealTimeLogs(cmd *exec.Cmd, logParser func(outputType, logLine string)) ([]byte, error) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s %v: failed to create stderr pipe: %w", cmd.Path, cmd.Args, err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s %v: failed to create stdout pipe: %w", cmd.Path, cmd.Args, err)
	}
	cmd.Start()

	// 256 kB for stdErr buffer
	stdErrBuffer := NewFixedBuffer(make([]byte, 0, 256*1024))

	// 1 MB for stdOut buffer
	stdOutBuffer := NewFixedBuffer(make([]byte, 0, 1024*1024))

	go func() {
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			logLine := scanner.Text()
			stdErrBuffer.Write([]byte(logLine))

			logParser("stderr", logLine)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			logLine := scanner.Text()

			stdOutBuffer.Write([]byte(logLine))
			logParser("stdout", logLine)
		}
	}()

	cmd.Wait()

	if !stdErrBuffer.Empty() {
		return stdOutBuffer.Read(), fmt.Errorf("failed to execute %s %v: %s", cmd.Path, cmd.Args, stdErrBuffer.String())
	}

	return stdOutBuffer.Read(), nil
}
