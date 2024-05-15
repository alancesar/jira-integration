package issue

import (
	"jira-integration/pkg/sprint"
	"time"
)

type (
	Project struct {
		ID   uint   `json:"id"`
		Key  string `json:"key"`
		Name string `json:"name,omitempty"`
	}

	Product struct {
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
		ID          uint
		Key         string
		Description string
		Summary     string
		Status      Status
		Type        Type
		Project     Project
		Parent      *Issue
		Sprints     []sprint.Sprint
		Labels      []string
		Assignee    *Account
		Reporter    Account
		StoryPoints uint
		Product     []Product
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)
