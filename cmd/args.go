package cmd

import (
	"go.uber.org/zap"
)

type RootArgs struct {
	Debug  bool
	Logger *zap.Logger
}
