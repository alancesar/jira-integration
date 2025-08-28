package jira

import (
	"context"
	"jira-integration/internal/jira/mocks"
	"jira-integration/pkg/issue"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_SearchIssueIDsByJQL(t *testing.T) {
	type fields struct {
		credentials Credentials
		boardID     int
		httpClient  *http.Client
	}
	type args struct {
		ctx           context.Context
		jql           string
		nextPageToken string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		want              []issue.Stamp
		wantNextPageToken string
		wantErr           bool
	}{
		{
			name: "return the issues from the jql search response",
			fields: fields{
				credentials: Credentials{},
				boardID:     0,
				httpClient: &http.Client{
					Transport: mocks.NewMockedRoundTripper(mocks.SearchResponse, http.StatusOK),
				},
			},
			args: args{
				ctx:           context.Background(),
				jql:           "project in (whatever)",
				nextPageToken: "return the ids of the issues return by the search",
			},
			want: []issue.Stamp{
				{
					ID:  1,
					Key: "key-1",
				},
				{
					ID:  2,
					Key: "key-2",
				},
			},
			wantNextPageToken: "next-page-token",
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				credentials: tt.fields.credentials,
				httpClient:  tt.fields.httpClient,
			}
			got, gotNextPageToken, err := c.SearchIssuesByJQL(tt.args.ctx, tt.args.jql, tt.args.nextPageToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchIssuesByJQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchIssuesByJQL() got = %v, want %v", got, tt.want)
			}
			if gotNextPageToken != tt.wantNextPageToken {
				t.Errorf("SearchIssuesByJQL() gotNextPageToken = %v, want %v", gotNextPageToken, tt.wantNextPageToken)
			}
		})
	}
}
