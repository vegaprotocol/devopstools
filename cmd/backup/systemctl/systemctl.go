package systemctl

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
)

func Start(logger *zap.Logger, service string) error {
	out, err := tools.ExecuteBinary("systemctl", []string{"start", service}, nil)
	logger.Debug(
		"Starting systemctl service",
		zap.String("stdout", string(out)),
		zap.Error(err),
		zap.String("service", service),
	)

	if err != nil {
		return fmt.Errorf("failed to start systemctl service: %w", err)
	}

	return nil
}

func Stop(logger *zap.Logger, service string) error {
	out, err := tools.ExecuteBinary("systemctl", []string{"stop", service}, nil)
	logger.Debug(
		"Stopping systemctl service",
		zap.String("stdout", string(out)),
		zap.Error(err),
		zap.String("service", service),
	)

	if err != nil {
		return fmt.Errorf("failed to stop systemctl service: %w", err)
	}

	return nil
}

func IsRunning(logger *zap.Logger, service string) bool {
	exitCode, out, err := tools.ExecuteBinaryWithExitCode("systemctl", []string{"status", service}, nil)

	logger.Debug(
		"Checking if systemctl service is running",
		zap.String("stdout", string(out)),
		zap.Error(err),
		zap.String("service", service),
	)

	if err != nil || exitCode > 0 {
		return false
	}

	return true
}
