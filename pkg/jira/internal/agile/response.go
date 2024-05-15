package agile

import (
	"jira-integration/pkg/sprint"
	"strings"
	"time"
)

type (
	FixVersion struct {
		Self        string
		ID          uint
		ProjectID   uint
		Name        string
		Description string    `json:"description"`
		Archived    bool      `json:"archived"`
		Released    bool      `json:"released"`
		ReleaseDate time.Time `json:"releaseDate"`
	}

	Sprint struct {
		ID            uint      `json:"id"`
		Self          string    `json:"self"`
		State         string    `json:"state"`
		Name          string    `json:"name"`
		StartDate     time.Time `json:"startDate"`
		EndDate       time.Time `json:"endDate"`
		CompleteDate  time.Time `json:"completeDate"`
		OriginBoardID int       `json:"originBoardId"`
		Goal          string    `json:"goal"`
	}
)

func (s Sprint) ToDomain() sprint.Sprint {
	return sprint.Sprint{
		ID:          s.ID,
		Name:        s.Name,
		State:       sprint.State(strings.ToLower(s.State)),
		Goal:        s.Goal,
		StartedAt:   s.StartDate,
		EndedAt:     s.EndDate,
		CompletedAt: s.CompleteDate,
	}
}
