package gdrive

import "net/http"

type GoogleDriveStorage struct {
	Provider string

	DefaultBucket string

	CredentialsFile string

	Project string

	DefaultTimeout string

	Client *http.Client

	EnableDebug bool

	PublicRead bool
}
