package jira

import (
	"jira-integration/pkg/issue"
	"jira-integration/pkg/jira/internal"
	"jira-integration/pkg/jira/internal/mock"
	"jira-integration/pkg/search"
	"jira-integration/pkg/sprint"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_Search(t *testing.T) {
	type fields struct {
		roundTripper http.RoundTripper
	}
	type args struct {
		params SearchRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    search.Response
		wantErr bool
	}{
		{
			name: "Should return search properly",
			fields: fields{
				roundTripper: mock.NewMockedRoundTripper(mock.SearchResponse, http.StatusOK),
			},
			args: args{
				params: SearchRequest{
					Fields:       []string{},
					FieldsByKeys: false,
					JQL:          "project = MAQ AND key = MAQ-1234",
					MaxResults:   15,
					StartAt:      0,
				},
			},
			want: search.Response{
				StartAt:    0,
				MaxResults: 15,
				Total:      1,
				Issues: []issue.Issue{
					{
						ID:      1234,
						Key:     "MAQ-1234",
						Summary: "Some story",
						Status: issue.Status{
							ID:   10052,
							Name: "Done",
							Category: issue.StatusCategory{
								ID:   3,
								Name: "Done",
							},
						},

						Priority: issue.Priority{
							ID:   3,
							Name: "Medium",
						},
						Type: issue.Type{
							ID:          1234,
							Description: "Some issue type description",
							Name:        "Story",
							Subtask:     false,
						},
						Parent: &issue.Issue{
							ID:      123,
							Key:     "MAQ-123",
							Summary: "Some epic",
							Status: issue.Status{
								ID:   123,
								Name: "In Progress",
								Category: issue.StatusCategory{
									ID:   4,
									Name: "In Progress",
								},
							},
							Priority: issue.Priority{
								ID:   3,
								Name: "Medium",
							},
							Type: issue.Type{
								ID:          10000,
								Description: "A big user story that needs to be broken down. Created by Jira Software - do not edit or delete.",
								Name:        "Epic",
								Subtask:     false,
							},
						},
						Sprints: []sprint.Sprint{
							{
								ID:          1,
								Name:        "Sprint 1",
								State:       "closed",
								Goal:        "Finish some big item",
								StartedAt:   internal.MustParseTimeRFC3339("2023-03-20T15:00:00.0Z"),
								EndedAt:     internal.MustParseTimeRFC3339("2023-03-30T21:00:00.0Z"),
								CompletedAt: internal.MustParseTimeRFC3339("2023-03-30T21:00:00.0Z"),
							},
							{
								ID:        2,
								Name:      "Sprint 2",
								State:     "active",
								Goal:      "Finish another big item",
								StartedAt: internal.MustParseTimeRFC3339("2023-04-03T15:00:00.0Z"),
								EndedAt:   internal.MustParseTimeRFC3339("2023-04-13T21:00:00.0Z"),
							},
						},
						FixVersions: []issue.FixVersion{
							{
								ID:          123,
								Name:        "2023 - Q1",
								Description: "Release 2023/Q1",
								ReleaseDate: internal.MustParseTime("2006-01-02", "2023-03-31"),
								Archived:    false,
								Released:    false,
							},
						},
						Labels:      []string{"Some Label", "Another Label"},
						StoryPoints: 3,
						Assignee:    "some.user@bexsbanco.com.br",
						NewProjects: "Maquininha - FXaaS",
						Allocation:  "Operação",
						CreatedAt:   internal.MustParseTimeRFC3339WithTimezone("2023-04-12T14:00:00.0-0300"),
						UpdatedAt:   internal.MustParseTimeRFC3339WithTimezone("2023-04-13T16:00:00.0-0300"),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(Credentials{}, &http.Client{
				Transport: tt.fields.roundTripper,
			})
			got, err := c.Search(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
		})
	}
}
