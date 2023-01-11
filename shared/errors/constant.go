package errors

// Error constants.
const (
	// ErrBadRequest is returned when something fails due to client error.
	ErrBadRequest = "bad request"

	// ErrUnauthorized is returned when something fails due to client not being authorized.
	ErrUnauthorized = "unauthorized"

	// ErrInvalidConfig is returned when the config is invalid.
	ErrInvalidConfig = "invalid config"

	// ErrInvalidBucket is returned when the bucket is invalid.
	ErrInvalidBucket = "invalid bucket"

	// ErrInvalidProvider is returned when the provider is invalid.
	ErrInvalidProvider = "invalid provider"

	// ErrInvalidCredentials is returned when the authentication credentials are invalid.
	ErrInvalidCredentials = "invalid credentials"

	// ErrFileOperationFailed is returned when a file operation fails.
	ErrFileOperationFailed = "file operation failed"

	// ErrIncompleteMultiFileUpload is returned when a multifile upload fails on some files.
	ErrIncompleteMultiFileUpload = "incomplete files upload"
)
