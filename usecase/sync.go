package usecase

import (
	"context"
	"fmt"
	"jira-integration/pkg/issue"
)

type (
	SprintClient interface {
		GetSprint(ctx context.Context, id uint) (*issue.Sprint, error)
	}

	SprintDatabase interface {
		GetSprintsByState(ctx context.Context, states []string) ([]issue.Sprint, error)
		SaveSprint(ctx context.Context, s issue.Sprint) error
	}

	SyncSprintsUseCase struct {
		db     SprintDatabase
		client SprintClient
	}
)

func NewSyncSprintsUseCase(client SprintClient, db SprintDatabase) *SyncSprintsUseCase {
	return &SyncSprintsUseCase{
		client: client,
		db:     db,
	}
}

func (uc SyncSprintsUseCase) Execute(ctx context.Context, states []string) error {
	sprints, err := uc.db.GetSprintsByState(ctx, states)
	if err != nil {
		return err
	}

	for _, s := range sprints {
		retrievedSprint, err := uc.client.GetSprint(ctx, s.ID)
		if err != nil {
			return err
		}

		fmt.Println("syncing sprint", retrievedSprint.Name)

		if err := uc.db.SaveSprint(ctx, *retrievedSprint); err != nil {
			return err
		}
	}

	return nil
}
