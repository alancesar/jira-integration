package model

import (
	"encoding/json"
	"gorm.io/datatypes"
	"jira-integration/pkg/issue"
	"time"
)

type (
	Changelog struct {
		ID        uint `gorm:"primarykey"`
		IssueID   uint
		From      string
		To        string
		CreatedAt time.Time `gorm:"autoCreateTime:false"`
	}

	Sprint struct {
		ID          uint `gorm:"primarykey"`
		Name        string
		State       string
		Goal        string
		StartedAt   time.Time
		EndedAt     time.Time
		CompletedAt time.Time
	}

	Issue struct {
		ID          uint   `gorm:"primarykey"`
		Key         string `gorm:"unique"`
		Summary     string
		Status      string
		IssueType   string
		Project     string
		ParentID    *uint
		Parent      *Issue
		SprintID    *uint
		Sprint      *Sprint
		Labels      datatypes.JSON
		Assignee    *string
		Reporter    string
		StoryPoints *uint
		Products    datatypes.JSON
		FixVersion  *string
		Locality    *string
		Changelog   []Changelog
		CreatedAt   time.Time `gorm:"autoCreateTime:false"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime:false"`
	}
)

func NewIssue(i issue.Issue) *Issue {
	labels, _ := json.Marshal(i.Labels)
	products, _ := json.Marshal(i.Products)

	var parent *Issue
	var parentID *uint
	if i.Parent != nil {
		parent = NewIssue(*i.Parent)
		parentID = &parent.ID
	}

	changelog := make([]Changelog, len(i.Changelog), len(i.Changelog))
	for index, c := range i.Changelog {
		changelog[index] = NewChangelog(c, i.ID)
	}

	var sprintID *uint
	if i.Sprint != nil {
		sprintID = &i.Sprint.ID
	}

	return &Issue{
		ID:          i.ID,
		Key:         i.Key,
		Summary:     i.Summary,
		Status:      i.Status,
		IssueType:   i.IssueType,
		Project:     i.Project,
		ParentID:    parentID,
		Parent:      parent,
		SprintID:    sprintID,
		Sprint:      NewSprint(i.Sprint),
		Labels:      labels,
		Assignee:    stringToPointer(i.Assignee),
		Reporter:    i.Reporter,
		StoryPoints: uintToPointer(i.StoryPoints),
		Products:    products,
		FixVersion:  stringToPointer(i.FixVersion),
		Locality:    stringToPointer(i.Locality),
		Changelog:   changelog,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func NewChangelog(c issue.Changelog, issueID uint) Changelog {
	return Changelog{
		ID:        c.ID,
		IssueID:   issueID,
		From:      c.From,
		To:        c.To,
		CreatedAt: c.CreatedAt,
	}
}

func NewSprint(sprint *issue.Sprint) *Sprint {
	if sprint == nil {
		return nil
	}

	return &Sprint{
		ID:          sprint.ID,
		Name:        sprint.Name,
		State:       sprint.State,
		Goal:        sprint.Goal,
		StartedAt:   sprint.StartedAt,
		EndedAt:     sprint.EndedAt,
		CompletedAt: sprint.CompletedAt,
	}
}

func stringToPointer(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func uintToPointer(value uint) *uint {
	if value == 0 {
		return nil
	}

	return &value
}
