package cluster

import (
	"context"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/generated/cluster"
	"github.com/pkg/errors"
	"testing"
)

var fn func(ctx context.Context, projectID string, clusterName string, reqEditors ...cluster.RequestEditorFn) (*cluster.GetClusterResponse, error)

type mockClientWithResponses struct {
	cluster.GetClusterWithResponseInterface
}

func (mc mockClientWithResponses) GetClusterWithResponse(ctx context.Context, projectID string, clusterName string, reqEditors ...cluster.RequestEditorFn) (*cluster.GetClusterResponse, error) {
	return fn(ctx, projectID, clusterName, reqEditors...)
}

func TestCreateOrUpdateClusterResponse_WaitHandler(t *testing.T) {

	type fields struct {
		ClientWithResponsesInterface cluster.ClientWithResponsesInterface
	}
	type args struct {
		ctx         context.Context
		clientFn    func(ctx context.Context, projectID string, clusterName string, reqEditors ...cluster.RequestEditorFn) (*cluster.GetClusterResponse, error)
		projectID   string
		clusterName string
	}
	type want struct {
		err error
		res interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "fancy error",
			fields: fields{},
			args: args{
				ctx: nil,
				clientFn: func(ctx context.Context, projectID string, clusterName string, reqEditors ...cluster.RequestEditorFn) (*cluster.GetClusterResponse, error) {
					return nil, errors.New("some fancy error")
				},
				projectID:   "",
				clusterName: "",
			},
			want: want{
				err: errors.New("some fancy error"),
				res: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CreateOrUpdateClusterResponse{
				ClientWithResponsesInterface: tt.fields.ClientWithResponsesInterface,
			}

			mc := mockClientWithResponses{}
			fn = tt.args.clientFn

			h := r.WaitHandler(tt.args.ctx, mc, tt.args.projectID, tt.args.clusterName)

			res, err := h.Wait()
			if err.Error() != tt.want.err.Error() {
				t.Errorf("err = %v, want %v", err, tt.want.err)
			}

			// !!! StatusForbidden results in 30min wait time (default timeout)

			_ = err
			_ = res

		})
	}
}
