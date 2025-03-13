package issue

import (
	"time"
)

const (
	StatusChangelogField = "status"
)

type (
	Changelog struct {
		ID        uint      `json:"id"`
		Author    string    `json:"author"`
		From      string    `json:"from"`
		To        string    `json:"to"`
		CreatedAt time.Time `json:"created_at"`
	}

	Issue struct {
		ID          uint        `json:"id"`
		Key         string      `json:"key"`
		Summary     string      `json:"summary"`
		Status      string      `json:"status"`
		IssueType   string      `json:"issue_type"`
		Project     string      `json:"project"`
		Parent      string      `json:"parent,omitempty"`
		Sprint      string      `json:"sprint,omitempty"`
		Labels      []string    `json:"labels,omitempty"`
		Assignee    string      `json:"assignee,omitempty"`
		Reporter    string      `json:"reporter,omitempty"`
		StoryPoints uint        `json:"story_points,omitempty"`
		Products    []string    `json:"products,omitempty"`
		FixVersion  string      `json:"fix_version,omitempty"`
		Locality    string      `json:"locality"`
		Changelog   []Changelog `json:"changelog,omitempty"`
		CreatedAt   time.Time   `json:"created_at"`
		UpdatedAt   time.Time   `json:"updated_at"`
	}
)
