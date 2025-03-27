package issue

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type (
	Label string

	Product struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	Changelog struct {
		ID        uint      `json:"id"`
		Author    string    `json:"author"`
		From      string    `json:"from"`
		To        string    `json:"to"`
		CreatedAt time.Time `json:"created_at"`
	}

	Sprint struct {
		ID          uint      `json:"id"`
		Name        string    `json:"name"`
		State       string    `json:"state"`
		Goal        string    `json:"goal"`
		StartedAt   time.Time `json:"startDate"`
		EndedAt     time.Time `json:"endDate"`
		CompletedAt time.Time `json:"completeDate,omitempty"`
	}

	Issue struct {
		ID          uint        `json:"id"`
		Key         string      `json:"key"`
		Summary     string      `json:"summary"`
		Status      string      `json:"status"`
		IssueType   string      `json:"issue_type"`
		Project     string      `json:"project"`
		Parent      *Issue      `json:"parent,omitempty"`
		Sprint      *Sprint     `json:"sprint,omitempty"`
		Labels      []Label     `json:"labels,omitempty"`
		Assignee    string      `json:"assignee,omitempty"`
		Reporter    string      `json:"reporter,omitempty"`
		StoryPoints uint        `json:"story_points,omitempty"`
		Products    []Product   `json:"products,omitempty"`
		FixVersion  string      `json:"fix_version,omitempty"`
		Locality    string      `json:"locality"`
		Changelog   []Changelog `json:"changelog,omitempty"`
		CreatedAt   time.Time   `json:"created_at"`
		UpdatedAt   time.Time   `json:"updated_at"`
	}
)

func NewLabels(labels []string) []Label {
	l := make([]Label, len(labels), len(labels))
	for i, label := range labels {
		l[i] = Label(label)
	}
	return l
}

func (l Label) Hash() string {
	hash := md5.Sum([]byte(l))
	return hex.EncodeToString(hash[:])
}
