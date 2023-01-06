// package bifrost provides a rainbow bridge for shipping files to cloud storage services.
package bifrost

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/opensaucerer/bifrost/gcs"
	"github.com/opensaucerer/bifrost/pinata"
	"github.com/opensaucerer/bifrost/s3"
	"github.com/opensaucerer/bifrost/shared/errors"
	"google.golang.org/api/option"
)

// NewRainbowBridge returns a new Rainbow Bridge for shipping files to your specified cloud storage service.
func NewRainbowBridge(bc *BridgeConfig) (RainbowBridge, error) {
	// vefify that the config is valid
	if bc == nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("config is nil"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the config is a struct
	t := reflect.TypeOf(bc)
	if t.Elem().Kind() != reflect.Struct {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid config type: %s", t.Elem().Kind()),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the config struct is of valid type
	if t.Elem().Name() != bridgeConfigType {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid config type: %s", t.Elem().Name()),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the config struct has a valid provider
	if bc.Provider == "" {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no provider specified"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the provider is valid
	if _, ok := providers[strings.ToLower(bc.Provider)]; !ok {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid provider: %s", bc.Provider),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the config struct has a valid bucket
	if bc.DefaultBucket == "" {
		// some provider might not require a bucket
		// Just log a warning
		if bc.EnableDebug {
			// TODO: create a logger
			log.Printf("No bucket specified for provider %s. This might cause errors or require you to specify a bucket for each operation.", providers[strings.ToLower(bc.Provider)])
		}
	}

	// Create a new bridge based on the provider
	switch bc.Provider {
	case "s3":
		return newSimpleStorageService(bc)
	case "gcs":
		return newGoogleCloudStorage(bc)
	case "pinata":
		return newPinataIPFS(bc)
	default:
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid provider: %s", bc.Provider),
			ErrorCode: errors.ErrBadRequest,
		}
	}

}
func newPinataIPFS(g *BridgeConfig) (rainbowBridge, error) {
	
	if g.PinataJWT == "" {
			return nil, &errors.BifrostError{
				Err:       fmt.Errorf("jwt is required for authentication"),
				ErrorCode: errors.ErrUnauthorized,
			}
	}

		// return a new Pinata Storage Provider
	return &pinata.PinataIPFSStorage{
		PinataJWT: g.PinataJWT,
		Provider:        providers[strings.ToLower(g.Provider)],
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		DefaultTimeout:  g.DefaultTimeout,
		EnableDebug:     g.EnableDebug,
		PublicRead:      g.PublicRead,
		UseAsync:        g.UseAsync,
	}, nil
}

// newGoogleCloudStorage returns a new client for Google Cloud Storage.
func newGoogleCloudStorage(g *BridgeConfig) (RainbowBridge, error) {
	var client *storage.Client
	var err error
	if g.CredentialsFile != "" {
		// first attempt to authenticate with credentials file
		client, err = storage.NewClient(context.Background(), option.WithCredentialsFile(g.CredentialsFile))
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}
	} else {
		// if no credentials file is specified, attempt to authenticate without credentials file
		client, err = storage.NewClient(context.Background())
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}
	}

	// return a new Google Cloud Storage client
	return &gcs.GoogleCloudStorage{
		Provider:        providers[strings.ToLower(g.Provider)],
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		DefaultTimeout:  g.DefaultTimeout,
		Client:          client,
		EnableDebug:     g.EnableDebug,
		PublicRead:      g.PublicRead,
		UseAsync:        g.UseAsync,
	}, nil
}

// newSimpleStorageService returns a new client for AWS S3
func newSimpleStorageService(g *BridgeConfig) (RainbowBridge, error) {
	var client *s3.Client
	if g.AccessKey != "" && g.SecretKey != "" {
		creds := credentials.NewStaticCredentialsProvider(g.AccessKey, g.SecretKey, "")
		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(creds), config.WithRegion(g.Region))
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}
		client = s3.NewFromConfig(cfg)
	} else {
		// Load AWS Shared Configuration
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}
		client = s3.NewFromConfig(cfg)
	}
	return &bs3.SimpleStorageService{
		Provider:      providers[strings.ToLower(g.Provider)],
		DefaultBucket: g.DefaultBucket,
		Region:        g.Region,
		PublicRead:    g.PublicRead,
		SecretKey:     g.SecretKey,
		AccessKey:     g.AccessKey,
		Client:        client,
	}, nil
}
