package types

// UploadedFile is the struct returned by the UploadFile function.
type UploadedFile struct {
	// Name is the name of the file.
	Name string
	// Bucket is the bucket the file was uploaded to.
	Bucket string
	// Local path is the path to the file.
	Path string
	// Size is the size of the file in bytes.
	Size int64
	// URL is the location of the file in the cloud.
	URL string
	// Preview is the URL to a preview of the file.
	Preview string
	// ProviderObject is the object returned by the cloud storage provider.
	// You need to type assert this to the correct type to use it.
	ProviderObject interface{}
	// Done sends a message to signal when an async process is complete.
	Done chan bool
	// Quit receives a message to signal for the exit of an async process.
	Quit chan bool
}

// UploadFileRequest is the request struct for Uploading Multiple Files.
type UploadFileRequest struct {
	// Path is the path to the file in memory.
	Path string
	// Filename is the name to store the file in cloud platform.
	Filename string
	// Options
	Options map[string]interface{}
}
