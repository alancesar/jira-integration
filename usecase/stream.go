//go:generate mockgen -source=stream.go -destination=mocks/stream_mock.go -package=mocks

package usecase

import (
	"context"
	"fmt"
	"jira-integration/pkg/issue"
)

type (
	IssueStreamer interface {
		SearchIssuesByJQL(ctx context.Context, jql, nextPageToken string) ([]issue.Stamp, string, error)
	}

	IssuePublisher func(ctx context.Context, issueID uint) error

	StampDatabase interface {
		GetByID(ctx context.Context, issueID uint) (issue.Stamp, bool, error)
	}

	StreamUseCase struct {
		streamer  IssueStreamer
		publisher IssuePublisher
		database  StampDatabase
	}
)

func NewStreamUseCase(streamer IssueStreamer, publisher IssuePublisher, database StampDatabase) *StreamUseCase {
	return &StreamUseCase{
		streamer:  streamer,
		publisher: publisher,
		database:  database,
	}
}

func (c StreamUseCase) Execute(ctx context.Context, jql string) error {
	issues := make(chan issue.Stamp)
	errs := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer close(issues)
		if err := c.search(ctx, jql, "", issues); err != nil {
			errs <- err
			cancel()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errs:
			return err
		case i, ok := <-issues:
			if !ok {
				return nil
			}

			stamp, exists, err := c.database.GetByID(ctx, i.ID)
			if err != nil {
				cancel()
				return err
			}

			if exists && stamp.UpdatedAt.Equal(i.UpdatedAt) {
				fmt.Println("skipping", i.ID)
				continue
			}

			if err := c.publisher(ctx, i.ID); err != nil {
				cancel()
				return err
			}
		}
	}
}

func (c StreamUseCase) search(ctx context.Context, jql, nextPageToken string, issues chan issue.Stamp) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	response, token, err := c.streamer.SearchIssuesByJQL(ctx, jql, nextPageToken)
	if err != nil {
		return err
	}

	for _, r := range response {
		issues <- r
	}

	if token != "" {
		return c.search(ctx, jql, token, issues)
	}

	return nil
}
