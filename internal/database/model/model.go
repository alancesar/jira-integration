package model

import (
	"encoding/json"
	"gorm.io/datatypes"
	"jira-integration/pkg/issue"
	"jira-integration/pkg/sprint"
	"time"
)

type (
	Account struct {
		ID           string `gorm:"primarykey"`
		EmailAddress string
		AvatarURL    string
		DisplayName  string
		Active       bool
		TimeZone     string
		AccountType  string
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

	Project struct {
		ID   uint `gorm:"primarykey"`
		Key  string
		Name string
	}

	Status struct {
		ID               uint `gorm:"primarykey"`
		Name             string
		StatusCategoryID uint
		StatusCategory   StatusCategory
	}

	StatusCategory struct {
		ID   uint `gorm:"primarykey"`
		Name string
	}

	IssueType struct {
		ID          uint `gorm:"primarykey"`
		Description string
		Name        string
		Subtask     bool
	}

	Product struct {
		ID   uint `gorm:"primarykey"`
		Name string
	}

	Changelog struct {
		ID           uint `gorm:"primarykey"`
		Issue        Issue
		IssueID      uint
		Author       Account
		AuthorID     string
		Field        string
		FromStatus   Status
		FromStatusID uint
		ToStatus     Status
		ToStatusID   uint
		CreatedAt    time.Time `gorm:"autoCreateTime:false"`
	}

	Issue struct {
		ID          uint   `gorm:"primarykey"`
		Key         string `gorm:"unique"`
		Summary     string
		Project     Project
		ProjectID   uint
		StatusID    uint
		Status      Status
		IssueTypeID uint
		IssueType   IssueType
		ParentID    *uint
		Parent      *Issue
		Sprints     []Sprint `gorm:"many2many:issue_sprints"`
		Labels      datatypes.JSON
		AssigneeID  *string
		Assignee    *Account
		ReporterID  string
		Reporter    Account
		StoryPoints *uint
		Products    []Product `gorm:"many2many:issue_products"`
		CreatedAt   time.Time `gorm:"autoCreateTime:false"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime:false"`
	}
)

func NewIssue(i issue.Issue) *Issue {
	labels, _ := json.Marshal(i.Labels)

	output := &Issue{
		ID:          i.ID,
		Key:         i.Key,
		Summary:     i.Summary,
		Project:     NewProject(i.Project),
		ProjectID:   i.Project.ID,
		StatusID:    i.Status.ID,
		Status:      NewStatus(i.Status),
		IssueTypeID: i.Type.ID,
		IssueType:   NewIssueType(i.Type),
		Sprints:     NewSprints(i.Sprints),
		Labels:      labels,
		ReporterID:  i.Reporter.ID,
		Reporter:    NewAccount(i.Reporter),
		StoryPoints: uintToPointer(i.StoryPoints),
		Products:    NewProducts(i.Product),
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}

	if i.Assignee != nil {
		assignee := NewAccount(*i.Assignee)
		output.AssigneeID = &assignee.ID
		output.Assignee = &assignee
	}

	if i.Parent != nil {
		output.Parent = NewIssue(*i.Parent)
		output.ParentID = uintToPointer(i.Parent.ID)
	}

	return output
}

func NewProject(p issue.Project) Project {
	return Project{
		ID:   p.ID,
		Key:  p.Key,
		Name: p.Name,
	}
}

func NewStatus(s issue.Status) Status {
	return Status{
		ID:               s.ID,
		Name:             s.Name,
		StatusCategoryID: s.Category.ID,
		StatusCategory:   NewStatusCategory(s.Category),
	}
}

func NewStatusCategory(c issue.StatusCategory) StatusCategory {
	return StatusCategory{
		ID:   c.ID,
		Name: c.Name,
	}
}

func NewIssueType(i issue.Type) IssueType {
	return IssueType{
		ID:          i.ID,
		Description: i.Description,
		Name:        i.Name,
		Subtask:     i.Subtask,
	}
}

func NewSprint(s sprint.Sprint) Sprint {
	return Sprint{
		ID:          s.ID,
		Name:        s.Name,
		State:       string(s.State),
		Goal:        s.Goal,
		StartedAt:   s.StartedAt,
		EndedAt:     s.EndedAt,
		CompletedAt: s.CompletedAt,
	}
}

func NewSprints(sprints []sprint.Sprint) []Sprint {
	output := make([]Sprint, len(sprints), len(sprints))
	for i := range sprints {
		output[i] = NewSprint(sprints[i])
	}

	return output
}

func NewProduct(p issue.Product) Product {
	return Product{
		ID:   p.ID,
		Name: p.Name,
	}
}

func NewProducts(products []issue.Product) []Product {
	if products == nil {
		return nil
	}

	output := make([]Product, len(products), len(products))
	for i := range products {
		output[i] = NewProduct(products[i])
	}

	return output
}

func NewAccount(a issue.Account) Account {
	return Account{
		ID:           a.ID,
		EmailAddress: a.EmailAddress,
		AvatarURL:    a.AvatarURL,
		DisplayName:  a.DisplayName,
		Active:       a.Active,
		TimeZone:     a.TimeZone,
		AccountType:  a.AccountType,
	}
}

func NewChangelog(i issue.Issue, c issue.Changelog) Changelog {
	return Changelog{
		ID:           c.ID,
		IssueID:      i.ID,
		Issue:        *NewIssue(i),
		Author:       NewAccount(c.Author),
		AuthorID:     c.Author.ID,
		Field:        string(c.Field),
		FromStatusID: c.FromStatusID,
		ToStatusID:   c.ToStatusID,
		CreatedAt:    c.CreatedAt,
	}
}

func uintToPointer(value uint) *uint {
	if value == 0 {
		return nil
	}

	return &value
}
