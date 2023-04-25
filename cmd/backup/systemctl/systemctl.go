package systemctl

import (
	"fmt"
	"time"

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

func Restart(logger *zap.Logger, service string) error {
	if err := Stop(logger, service); err != nil {
		return fmt.Errorf("failed to stop the %s service %w", service, err)
	}

	if err := tools.Retry(3, 3*time.Second, func() error {
		if err := Start(logger, service); err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("all attempts to start the %s service failed: %w", service, err)
	}

	return nil
}

// ref: https://www.freedesktop.org/software/systemd/man/systemctl.html#Exit%20status
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
