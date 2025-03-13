package usecase

import (
	"context"
	"fmt"
	"jira-integration/pkg/issue"
)

type (
	IssueClient interface {
		GetIssueByKeyOrID(ctx context.Context, issueKeyOrID string) (issue.Issue, error)
		GetIssueChangelog(ctx context.Context, issueKeyOrID, nextPageToken string) ([]issue.Changelog, string, error)
	}

	IssueDatabase interface {
		SaveIssue(ctx context.Context, i issue.Issue) error
	}

	FetchUseCase struct {
		db     IssueDatabase
		client IssueClient
	}
)

func (uc FetchUseCase) Execute(ctx context.Context, issueKeyOrID string) error {
	issueFromClient, err := uc.client.GetIssueByKeyOrID(ctx, issueKeyOrID)
	if err != nil {
		return fmt.Errorf("while fetching issue from client: %w", err)
	}

	changelog, _, err := uc.client.GetIssueChangelog(ctx, issueFromClient.Key, "")
	if err != nil {
		return fmt.Errorf("while fetching issue changelog: %w", err)
	}

	issueFromClient.Changelog = changelog
	if err := uc.db.SaveIssue(ctx, issueFromClient); err != nil {
		return fmt.Errorf("while saving issue to db: %w", err)
	}

	return nil
}
