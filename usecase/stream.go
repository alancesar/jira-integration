package usecase

import (
	"context"
	"jira-integration/pkg/issue"
)

type (
	StreamIssueClient interface {
		SearchIssueIDsByJQL(ctx context.Context, jql, nextPageToken string) ([]string, string, error)
	}

	StreamUseCase struct {
		client StreamIssueClient
	}
)

func (c StreamUseCase) Execute(ctx context.Context, jql string) {

}

func (c StreamUseCase) search(ctx context.Context, jql, nextPageToken string, ids chan []uint) error {
	response, token, err := c.client.SearchIssueIDsByJQL(ctx, jql, nextPageToken)
	if err != nil {
		return err
	}

	ids <- response
}
