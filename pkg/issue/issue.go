package issue

import (
	"jira-integration/pkg/sprint"
	"time"
)

type (
	FixVersion struct {
		ID          uint
		Name        string
		Description string
		Archived    bool
		Released    bool
		ReleaseDate time.Time
	}

	Project struct {
		ID   uint   `json:"id"`
		Key  string `json:"key"`
		Name string `json:"name,omitempty"`
	}

	Priority struct {
		ID   uint
		Name string
	}

	Type struct {
		ID          uint
		Description string
		Name        string
		Subtask     bool
	}

	Status struct {
		ID       uint
		Name     string
		Category StatusCategory
	}

	StatusCategory struct {
		ID   uint
		Name string
	}

	Resolution struct {
		ID          uint
		Description string
		Name        string
	}

	Progress struct {
		Progress int
		Total    int
		Percent  int
	}

	Account struct {
		ID           string
		EmailAddress string
		AvatarURL    string
		DisplayName  string
		Active       bool
		TimeZone     string
		AccountType  string
	}

	Issue struct {
		ID                    uint
		Key                   string
		Description           string
		Summary               string
		Status                Status
		Priority              Priority
		Type                  Type
		Project               Project
		Progress              Progress
		AggregateProgress     Progress
		AggregateTimeSpent    int
		AggregateTimeEstimate int
		TimeSpent             int
		Parent                *Issue
		Sprints               []sprint.Sprint
		FixVersions           []FixVersion
		Labels                []string
		Assignee              *Account
		Reporter              Account
		StoryPoints           uint
		NewProjects           string
		Allocation            string
		Resolution            *Resolution
		System                *string
		Squad                 *string
		ResolvedAt            *time.Time
		CreatedAt             time.Time
		UpdatedAt             time.Time
	}
)
