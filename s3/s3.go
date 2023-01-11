package s3

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

func (s *SimpleStorageService) UploadFile(path, filename string, options map[string]interface{}) (*types.UploadedFile, error) {
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

	var params *s3.PutObjectInput = &s3.PutObjectInput{
		Bucket: aws.String(s.DefaultBucket),
		Key:    aws.String(filename),
		Body:   file,
	}
	// check the bridge config for default acl settings
	if s.PublicRead {
		// set public read permissions
		params.ACL = awsTypes.ObjectCannedACLPublicRead
	}
	// set file permissions
	for k, v := range options {
		switch k {
		// check the options map for acl settings
		case config.OptACL:
			if v, ok := v.(string); ok {
				switch v {
				case config.ACLPublicRead:
					params.ACL = awsTypes.ObjectCannedACLPublicRead
				case config.ACLPrivate:
					params.ACL = awsTypes.ObjectCannedACLPrivate
				}

			}
		// set the content type
		case config.OptContentType:
			if v, ok := v.(string); ok {
				params.ContentType = aws.String(v)
			}
		// set object metadata
		case config.OptMetadata:
			if v, ok := v.(map[string]string); ok {
				params.Metadata = v
			}
		}
	}
	// Upload the file to S3
	if _, err := s.Client.PutObject(ctx, params); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// head object details
	obj, err := s.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.DefaultBucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	fmt.Println(obj.Metadata)
	return &types.UploadedFile{
		Name:           filename,
		Bucket:         s.DefaultBucket,
		Path:           path,
		Preview:        fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.DefaultBucket, s.Region, filename),
		Size:           obj.ContentLength,
		ProviderObject: obj,
	}, nil
}

// UploadMultiFile
func (s *SimpleStorageService) UploadMultiFile(requests []*types.UploadFileRequest) ([]*types.UploadedFile, error) {
	return nil, nil
}

// Config returns the s3 configuration.
func (s *SimpleStorageService) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		DefaultBucket:  s.DefaultBucket,
		Region:         s.Region,
		AccessKey:      s.AccessKey,
		SecretKey:      s.SecretKey,
		DefaultTimeout: s.DefaultTimeout,
		EnableDebug:    s.EnableDebug,
		Provider:       s.Provider,
		UseAsync:       s.UseAsync,
	}
}

func (s *SimpleStorageService) Disconnect() error {
	return nil
}
