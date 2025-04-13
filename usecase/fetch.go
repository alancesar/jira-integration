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
		CreateIssue(ctx context.Context, i issue.Issue) error
		UpdateIssue(ctx context.Context, i issue.Issue) error
		IssueExistsByKey(ctx context.Context, issueKey string) (bool, error)
	}

	FetchUseCase struct {
		client IssueClient
		db     IssueDatabase
	}
)

func NewFetchUseCase(client IssueClient, db IssueDatabase) *FetchUseCase {
	return &FetchUseCase{
		client: client,
		db:     db,
	}
}

func (uc FetchUseCase) Execute(ctx context.Context, issueKeyOrID string) error {
	fmt.Println("fetching", issueKeyOrID)
	issueFromClient, err := uc.client.GetIssueByKeyOrID(ctx, issueKeyOrID)
	if err != nil {
		return fmt.Errorf("while fetching issue %s from streamer: %w", issueKeyOrID, err)
	}

	changelog, _, err := uc.client.GetIssueChangelog(ctx, issueFromClient.Key, "")
	if err != nil {
		return fmt.Errorf("while fetching issue %s changelog: %w", issueKeyOrID, err)
	}

	issueFromClient.Changelog = changelog

	if exist, err := uc.db.IssueExistsByKey(ctx, issueFromClient.Key); err != nil {
		return fmt.Errorf("while checking if issue %s exists: %w", issueKeyOrID, err)
	} else if exist {
		return uc.updateIssue(ctx, issueFromClient)
	}

	return uc.createIssue(ctx, issueFromClient)
}

func (uc FetchUseCase) createIssue(ctx context.Context, issueFromClient issue.Issue) error {
	if err := uc.db.CreateIssue(ctx, issueFromClient); err != nil {
		return fmt.Errorf("while creaing issue %s to db: %w", issueFromClient.Key, err)
	}

	return nil
}

func (uc FetchUseCase) updateIssue(ctx context.Context, issueFromClient issue.Issue) error {
	if err := uc.db.UpdateIssue(ctx, issueFromClient); err != nil {
		return fmt.Errorf("while updating issue %s to db: %w", issueFromClient.Key, err)
	}

	return nil
}
