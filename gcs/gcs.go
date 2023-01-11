package gcs

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

/*
UploadFile uploads a file to Google Cloud Storage and returns an error if one occurs.

Note: UploadFile requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (g *GoogleCloudStorage) UploadFile(path, filename string, options map[string]interface{}) (*types.UploadedFile, error) {
	// create context and add timeout if default timeout is set
	var ctx context.Context
	var cancel context.CancelFunc
	ctx = context.Background()
	if g.DefaultTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(g.DefaultTimeout)*time.Second)
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
	// Upload file to Google Cloud Storage
	obj := g.Client.Bucket(g.DefaultBucket).Object(filename)
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// close writer
	if err := wc.Close(); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// set file permissions
	if options != nil {
		// check the options map for acl settings
		if acl, ok := options[config.OptACL]; ok {
			switch acl.(string) {
			case config.ACLPublicRead:
				// set public read permissions
				if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
					return nil, &errors.BifrostError{
						Err:       err,
						ErrorCode: errors.ErrFileOperationFailed,
					}
				}
			case config.ACLPrivate:
				// set private permissions
				if err := obj.ACL().Set(ctx, storage.AllAuthenticatedUsers, storage.RoleReader); err != nil {
					return nil, &errors.BifrostError{
						Err:       err,
						ErrorCode: errors.ErrFileOperationFailed,
					}
				}
			}
		} else if g.PublicRead { // check the bridge config for default acl settings
			// set public read permissions
			if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
				return nil, &errors.BifrostError{
					Err:       err,
					ErrorCode: errors.ErrFileOperationFailed,
				}
			}
		}
	}
	// set object metadata
	if options != nil {
		// check the options map for metadata settings
		if metadata, ok := options[config.OptMetadata]; ok {
			// set metadata
			if _, err := obj.Update(ctx, storage.ObjectAttrsToUpdate{Metadata: metadata.(map[string]string)}); err != nil {
				return nil, &errors.BifrostError{
					Err:       err,
					ErrorCode: errors.ErrFileOperationFailed,
				}
			}
		}
	}
	// get object attributes
	objAttrs, _ := obj.Attrs(ctx)
	return &types.UploadedFile{
		Name:           objAttrs.Name,
		Bucket:         objAttrs.Bucket,
		Path:           path,
		Size:           objAttrs.Size,
		URL:            objAttrs.MediaLink,
		Preview:        fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name),
		ProviderObject: obj,
	}, nil
}

/*
UploadMultiFile uploads files to Google Cloud Storage and returns an error if all files
are not uploaded successfully. Set EnableDebug to true in bridge.BridgeConfig to see
detailed logs of errors if all files are not uploaded.

Note: UploadMultiFile requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (g *GoogleCloudStorage) UploadMultiFile(requests []*types.UploadFileRequest) ([]*types.UploadedFile, error) {
	var err error
	uploadedFiles := make([]*types.UploadedFile, 0, len(requests))

	for _, request := range requests {
		uploadedFile, err := g.UploadFile(request.Path, request.Filename, request.Options)
		if err != nil {
			if g.EnableDebug {
				// log a failed upload file request
				log.Printf("Upload request for file with path %s failed with err: %s\n", request.Path, err.Error())
				continue
			}
		}
		uploadedFiles = append(uploadedFiles, uploadedFile)
	}

	if len(uploadedFiles) == 0 {
		err = &errors.BifrostError{
			Err:       fmt.Errorf("upload operation failed for all/no files"),
			ErrorCode: errors.ErrBadRequest,
		}
	} else if len(uploadedFiles) < len(requests) {
		err = &errors.BifrostError{
			Err:       fmt.Errorf("upload operation failed for %d files", len(requests)-len(uploadedFiles)),
			ErrorCode: errors.ErrIncompleteMultiFileUpload,
		}
	}

	return uploadedFiles, err
}

/*
Disconnect closes the Google Cloud Storage connection and returns an error if one occurs.

Disconnect should only be called when the connection is no longer needed.
*/
func (g *GoogleCloudStorage) Disconnect() error {
	if g.Client != nil {
		return g.Client.Close()
	}
	return nil
}

// Config returns the Google Cloud Storage configuration.
func (g *GoogleCloudStorage) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		Provider:        g.Provider,
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		DefaultTimeout:  g.DefaultTimeout,
		EnableDebug:     g.EnableDebug,
		UseAsync:        g.UseAsync,
	}
}
