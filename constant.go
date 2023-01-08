package bifrost

var (
	// providers is a map of the supported providers
	providers = map[string]string{
		"s3":  "Simple Storage Service",
		"gcs": "Google Cloud Storage",
	}
)

// Misc constants
const (
	// SimpleStorageService is the identifier of the S3 provider
	SimpleStorageService = "s3"
	// GoogleCloudStorage is the identifier of the Google Cloud Storage provider
	GoogleCloudStorage = "gcs"

	//GoogleDriveStorage is the identifier of the Google drive provider
	GoogleDriveStorage = "gdrive"

	// BridgeConfigType is the type of the bridge configuration
	bridgeConfigType = "BridgeConfig"
)

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
)

// Options constants.
const (
	// ACL is the option to set the ACL of the file.
	OptACL = "acl"
	// PublicRead is the option to set the ACL of the file to public read.
	ACLPublicRead = "public-read"
	// Private is the option to set the ACL of the file to private.
	ACLPrivate = "private"
	// ContentType is the option to set the content type of the file.
	OptContentType = "content-type"
	// Metadata is the option to set the metadata of the file.
	OptMetadata = "metadata"
)
