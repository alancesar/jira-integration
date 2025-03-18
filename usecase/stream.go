//go:generate mockgen -source=stream.go -destination=mocks/stream_mock.go -package=mocks

package usecase

import (
	"context"
)

type (
	IssueStreamer interface {
		SearchIssueIDsByJQL(ctx context.Context, jql, nextPageToken string) ([]string, string, error)
	}

	IssuePublisher func(ctx context.Context, issueID string) error

	StreamUseCase struct {
		streamer  IssueStreamer
		publisher IssuePublisher
	}
)

func NewStreamUseCase(streamer IssueStreamer, publisher IssuePublisher) *StreamUseCase {
	return &StreamUseCase{
		streamer:  streamer,
		publisher: publisher,
	}
}

func (c StreamUseCase) Execute(ctx context.Context, jql string) error {
	ids := make(chan string)
	errs := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer close(ids)
		if err := c.search(ctx, jql, "", ids); err != nil {
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
		case id, ok := <-ids:
			if !ok {
				return nil
			}

			if err := c.publisher(ctx, id); err != nil {
				cancel()
				return err
			}
		}
	}
}

func (c StreamUseCase) search(ctx context.Context, jql, nextPageToken string, ids chan string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	response, token, err := c.streamer.SearchIssueIDsByJQL(ctx, jql, nextPageToken)
	if err != nil {
		return err
	}

	for _, id := range response {
		ids <- id
	}

	if token != "" {
		return c.search(ctx, jql, token, ids)
	}

	return nil
}
