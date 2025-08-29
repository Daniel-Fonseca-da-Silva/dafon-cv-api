package errors

// Custom errors for application operations
var (
	// Email related errors
	ErrEmailConfigMissing = &AppError{message: "email configuration is missing"}
	ErrEmailSendFailed    = &AppError{message: "failed to send email"}

	// Configuration related errors
	ErrInvalidPort = &AppError{message: "invalid port configuration"}

	// Database related errors
	ErrDatabaseConnection = &AppError{message: "database connection failed"}

	// User related errors
	ErrUserNotFound      = &AppError{message: "user not found"}
	ErrUserAlreadyExists = &AppError{message: "user already exists"}

	// Authentication related errors
	ErrInvalidCredentials = &AppError{message: "invalid credentials"}
	ErrTokenExpired       = &AppError{message: "token expired"}
	ErrInvalidToken       = &AppError{message: "invalid token"}

	// Worker pool related errors
	ErrQueueFull         = &AppError{message: "queue is full"}
	ErrPoolStopped       = &AppError{message: "worker pool is stopped"}
	ErrWorkerUnavailable = &AppError{message: "no worker available"}
)

// AppError represents an error that can occur during application operations
type AppError struct {
	message string
}

// Error returns the error message
func (e *AppError) Error() string {
	return e.message
}

// NewAppError creates a new application error with a custom message
func NewAppError(message string) *AppError {
	return &AppError{message: message}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, context string) error {
	if appErr, ok := err.(*AppError); ok {
		return &AppError{message: context + ": " + appErr.message}
	}
	return &AppError{message: context + ": " + err.Error()}
}
