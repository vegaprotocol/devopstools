package errors

import "errors"

// ErrConnectionNotReady indicated that the network connection to the gRPC server is not ready.
var ErrConnectionNotReady = errors.New("gRPC connection not ready")
