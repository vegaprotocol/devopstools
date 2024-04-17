package errors

import "errors"

var (
	// ErrConnectionNotReady indicated that the network connection to the gRPC server is not ready.
	ErrConnectionNotReady = errors.New("gRPC connection not ready")

	// ErrNil indicates that a nil/null pointer was encountered.
	ErrNil = errors.New("nil pointer")
)
