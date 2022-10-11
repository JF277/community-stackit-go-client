// package buckets handles creation and management of STACKIT Object Storage buckets

package buckets

import (
	"context"
	"fmt"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/wait"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/internal/common"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
)

// constants
const (
	apiPath     = consts.API_PATH_OBJECT_STORAGE_BUCKET
	apiPathList = consts.API_PATH_OBJECT_STORAGE_BUCKETS
)

// New returns a new handler for the service
func New(c common.Client) *ObjectStorageBucketsService {
	return &ObjectStorageBucketsService{
		Client: c,
	}
}

// ObjectStorageBucketsService is the service that handles
// CRUD functionality for SKE buckets
type ObjectStorageBucketsService common.Service

// BucketResponse is a struct representation of stackit's object storage api response for a bucket
type BucketResponse struct {
	Project string `json:"project"`
	Bucket  Bucket `json:"bucket"`
}

// BucketListResponse is a struct representation of stackit's object storage api response for a bucket list
type BucketListResponse struct {
	Project string   `json:"project"`
	Buckets []Bucket `json:"buckets"`
}

// Bucket holds all the bucket information
type Bucket struct {
	Name                  string `json:"name"`
	Region                string `json:"region"`
	URLPathStyle          string `json:"urlPathStyle"`
	URLVirtualHostedStyle string `json:"urlVirtualHostedStyle"`
}

// List returns the a list of buckets assigned to a project ID
// See also https://api.stackit.schwarz/object-storage-service/openapi.v1.html#operation/get_buckets_v1_project__projectId__buckets_get
func (svc *ObjectStorageBucketsService) List(ctx context.Context, projectID string) (res BucketListResponse, err error) {
	req, err := svc.Client.Request(ctx, http.MethodGet, fmt.Sprintf(apiPathList, projectID), nil)
	if err != nil {
		return
	}

	_, err = svc.Client.Do(req, &res)
	return
}

// Get returns the a bucket by project ID and bucket name
// See also https://api.stackit.schwarz/object-storage-service/openapi.v1.html#operation/get_bucket_v1_project__projectId__bucket__bucketName__get
func (svc *ObjectStorageBucketsService) Get(ctx context.Context, projectID, bucketName string) (res BucketResponse, err error) {
	req, err := svc.Client.Request(ctx, http.MethodGet, fmt.Sprintf(apiPath, projectID, bucketName), nil)
	if err != nil {
		return
	}

	_, err = svc.Client.Do(req, &res)
	return
}

// Create creates a new bucket
// See also https://api.stackit.schwarz/object-storage-service/openapi.v1.html#operation/create_bucket_v1_project__projectId__bucket__bucketName__post
func (svc *ObjectStorageBucketsService) Create(ctx context.Context, projectID, bucketName string) (w *wait.Handler, err error) {
	if err = ValidateBucketName(bucketName); err != nil {
		err = validate.WrapError(err)
		return
	}

	req, err := svc.Client.Request(ctx, http.MethodPost, fmt.Sprintf(apiPath, projectID, bucketName), nil)
	if err != nil {
		return
	}
	_, err = svc.Client.Do(req, nil)
	if err != nil {
		return
	}

	w = wait.New(waitForCreation(ctx, projectID, bucketName, svc))
	return
}

func waitForCreation(ctx context.Context, projectID string, bucketName string, svc *ObjectStorageBucketsService) wait.WaitFn {
	return func() (interface{}, bool, error) {
		res, err := svc.Get(ctx, projectID, bucketName)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				return nil, false, nil
			}
			return nil, false, err
		}
		return res, true, nil
	}
}

// Delete deletes a bucket
// See also https://api.stackit.schwarz/ske-service/openapi.v1.html#operation/SkeService_DeleteBucket
func (svc *ObjectStorageBucketsService) Delete(ctx context.Context, projectID, bucketName string) (w *wait.Handler, err error) {
	req, err := svc.Client.Request(ctx, http.MethodDelete, fmt.Sprintf(apiPath, projectID, bucketName), nil)
	if err != nil {
		return
	}
	_, err = svc.Client.Do(req, nil)
	if err != nil {
		return
	}
	w = wait.New(waitForDeletion(ctx, projectID, bucketName, svc))
	return w, err
}

func waitForDeletion(ctx context.Context, projectID string, bucketName string, svc *ObjectStorageBucketsService) wait.WaitFn {
	return func() (interface{}, bool, error) {
		res, err := svc.Get(ctx, projectID, bucketName)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				return res, true, nil
			}
			return res, false, err
		}
		return res, false, nil
	}
}
