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
		StartedAt   time.Time `json:"start_date"`
		EndedAt     time.Time `json:"end_date"`
		CompletedAt time.Time `json:"complete_date,omitempty"`
	}

	Account struct {
		ID           string `json:"id"`
		EmailAddress string `json:"email_address"`
		AvatarURL    string `json:"avatar_url"`
		DisplayName  string `json:"display_name"`
		Active       bool   `json:"active"`
		AccountType  string `json:"account_type"`
	}

	Stamp struct {
		ID        uint      `json:"id"`
		Key       string    `json:"key"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	Issue struct {
		Stamp
		Summary     string      `json:"summary"`
		Status      string      `json:"status"`
		IssueType   string      `json:"issue_type"`
		Project     string      `json:"project"`
		Parent      *Issue      `json:"parent,omitempty"`
		Sprint      *Sprint     `json:"sprint,omitempty"`
		Labels      []Label     `json:"labels,omitempty"`
		Assignee    *Account    `json:"assignee,omitempty"`
		Reporter    Account     `json:"reporter"`
		StoryPoints uint        `json:"story_points,omitempty"`
		Products    []Product   `json:"products,omitempty"`
		FixVersion  string      `json:"fix_version,omitempty"`
		Locality    string      `json:"locality"`
		Changelog   []Changelog `json:"changelog,omitempty"`
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
