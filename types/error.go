package types

import "fmt"

// PanicRecoveredError can be used as default error type when panics have been recovered.
// This can come in handy if using a recovery middleware or simply while encoding errors to transport layer.
type PanicRecoveredError struct{ error }

// NewPanicRecoveredError wraps the error as PanicRecoveredError and adds 'recovered panic:' to the message.
func NewPanicRecoveredError(err error) PanicRecoveredError {
	return PanicRecoveredError{fmt.Errorf("recovered panic: %s", err)}
}
