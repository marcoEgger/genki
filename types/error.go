package types

// PanicRecoveredError can be used as default error type when panics have been recovered.
// This can come in handy if using a recovery middleware or simply while encoding errors to transport layer.
type PanicRecoveredError struct{ error }
