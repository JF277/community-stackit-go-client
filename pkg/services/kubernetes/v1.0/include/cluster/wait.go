package cluster

import (
	"context"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/generated/cluster"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/wait"
)

func (r CreateOrUpdateClusterResponse) WaitHandler(ctx context.Context, c cluster.GetClusterWithResponseInterface, projectID, clusterName string) *wait.Handler {
	return wait.New(func() (res interface{}, done bool, err error) {
		resp, err := c.GetClusterWithResponse(ctx, projectID, clusterName)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusForbidden)) {
				return nil, false, nil
			}
			return nil, false, err
		}
		if resp.StatusCode() == http.StatusForbidden {
			return nil, false, nil
		}
		if resp.StatusCode() == http.StatusInternalServerError {
			return nil, false, nil
		}
		if resp.HasError != nil {
			return nil, false, resp.HasError
		}

		status := *resp.JSON200.Status.Aggregated
		if status == cluster.STATE_HEALTHY || status == cluster.STATE_HIBERNATED {
			return resp, true, nil
		}
		return resp, false, nil
	})
}

func (r DeleteClusterResponse) WaitHandler(ctx context.Context, c cluster.GetClusterWithResponseInterface, projectID, clusterName string) *wait.Handler {
	return wait.New(func() (res interface{}, done bool, err error) {
		resp, err := c.GetClusterWithResponse(ctx, projectID, clusterName)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				return nil, true, nil
			}
			return nil, false, err
		}
		if resp.StatusCode() == http.StatusInternalServerError {
			return nil, false, nil
		}
		if resp.HasError != nil {
			if resp.StatusCode() == http.StatusNotFound {
				return nil, true, nil
			}
			return nil, false, err
		}
		return nil, false, nil
	})
}
