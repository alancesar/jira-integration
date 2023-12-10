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

	Issue struct {
		ID          uint
		Key         string
		Summary     string
		Status      Status
		Priority    Priority
		Type        Type
		Parent      *Issue
		Sprints     []sprint.Sprint
		FixVersions []FixVersion
		Labels      []string
		Assignee    string
		StoryPoints uint
		NewProjects string
		Allocation  string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)
