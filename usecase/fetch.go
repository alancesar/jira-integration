package usecase

import (
	"context"
	"fmt"
	"jira-integration/pkg/issue"
)

type (
	IssueClient interface {
		GetIssueByID(ctx context.Context, issueKeyOrID uint) (issue.Issue, error)
		GetIssueChangelog(ctx context.Context, issueKeyOrID, nextPageToken string) ([]issue.Changelog, string, error)
	}

	IssueDatabase interface {
		StampDatabase
		CreateIssue(ctx context.Context, i issue.Issue) error
		UpdateIssue(ctx context.Context, i issue.Issue) error
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

func (uc FetchUseCase) Execute(ctx context.Context, issueID uint) error {
	fmt.Println("fetching", issueID)
	issueFromClient, err := uc.client.GetIssueByID(ctx, issueID)
	if err != nil {
		return fmt.Errorf("while fetching issue %d from streamer: %w", issueID, err)
	}

	changelog, _, err := uc.client.GetIssueChangelog(ctx, issueFromClient.Key, "")
	if err != nil {
		return fmt.Errorf("while fetching issue %d changelog: %w", issueID, err)
	}

	issueFromClient.Changelog = changelog

	if _, exist, err := uc.db.GetByID(ctx, issueID); err != nil {
		return fmt.Errorf("while checking if issue %d exists: %w", issueID, err)
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
