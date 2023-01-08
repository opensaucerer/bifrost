package gdrive

import (
	"context"
	"fmt"
	"time"
)

func (g *GoogleDriveStorage) UploadFile(path, filename string, options map[string]interface{}) (types.UploadFile, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx = context.Background()

	if s.DefaultTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(s.DefaultTimeout)*time.Second)
		defer cancel()
	}
	// verify that file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("file does not exist: %s", path),
			ErrorCode: errors.ErrBadRequest,
		}
	}
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// close file
	defer file.Close()

	srv, err := drive.NewService(ctx, option.WithHTTPClient(g.Client))
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}

	f, err := srv.Files.Create(file).Do()
	if err != nil {
		return nil, &error.BifrostError{
			Err: err, 
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}

	return &types.UploadedFile{
		Name:           f.Name,
		Path:           path,
		Size:           f.Size,
		URL:            f.IconLink,
	}, nil
}

func (g *GoogleDriveStorage) Disconnect() error {
	if g.Client != nil {
		return g.Client.Close()
	}
	return nil
}

// Config returns the Google Cloud Storage configuration.
func (g *GoogleDriveStorage) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		Provider:        g.Provider,
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		DefaultTimeout:  g.DefaultTimeout,
		EnableDebug:     g.EnableDebug,
	}
}
