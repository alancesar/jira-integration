package model

import (
	"jira-integration/pkg/issue"
	"time"
)

type (
	Label struct {
		ID   string `gorm:"primarykey"`
		Name string `gorm:"unique"`
	}

	Product struct {
		ID   uint   `gorm:"primarykey"`
		Name string `gorm:"unique"`
	}

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

	Sprints []Sprint

	Account struct {
		ID           string `gorm:"primarykey"`
		EmailAddress string
		AvatarURL    string
		DisplayName  string
		Active       bool
		AccountType  string
	}

	Issue struct {
		ID          uint   `gorm:"primarykey"`
		Key         string `gorm:"index,unique"`
		Summary     string
		Status      string
		IssueType   string
		Project     string
		ParentID    *uint
		Parent      *Issue
		SprintID    *uint
		Sprint      *Sprint
		Labels      []Label `gorm:"many2many:issue_labels;"`
		AssigneeID  *string
		Assignee    *Account
		ReporterID  string
		Reporter    Account
		StoryPoints *uint
		Products    []Product `gorm:"many2many:issue_products;"`
		FixVersion  *string
		Locality    *string
		Changelog   []Changelog
		CreatedAt   time.Time `gorm:"autoCreateTime:false"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime:false"`
	}
)

func (s Sprint) ToDomain() issue.Sprint {
	return issue.Sprint{
		ID:          s.ID,
		Name:        s.Name,
		State:       s.State,
		Goal:        s.Goal,
		StartedAt:   s.StartedAt,
		EndedAt:     s.EndedAt,
		CompletedAt: s.CompletedAt,
	}
}

func (s Sprints) ToDomain() []issue.Sprint {
	output := make([]issue.Sprint, len(s), len(s))
	for i, sprint := range s {
		output[i] = sprint.ToDomain()
	}

	return output
}

func NewIssue(i issue.Issue) *Issue {
	var parent *Issue
	var parentID *uint
	if i.Parent != nil {
		parent = NewIssue(*i.Parent)
		parentID = &parent.ID
	}

	labels := make([]Label, len(i.Labels), len(i.Labels))
	for index, label := range i.Labels {
		labels[index] = NewLabel(label)
	}

	products := make([]Product, len(i.Products), len(i.Products))
	for index, product := range i.Products {
		products[index] = NewProduct(product)
	}

	changelog := make([]Changelog, len(i.Changelog), len(i.Changelog))
	for index, c := range i.Changelog {
		changelog[index] = NewChangelog(c, i.ID)
	}

	var sprintID *uint
	if i.Sprint != nil {
		sprintID = &i.Sprint.ID
	}

	var assigneeID *string
	assignee := NewAccount(i.Assignee)
	if assignee != nil {
		assigneeID = &assignee.ID
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
		AssigneeID:  assigneeID,
		Assignee:    assignee,
		ReporterID:  i.Reporter.ID,
		Reporter:    *NewAccount(&i.Reporter),
		StoryPoints: i.StoryPoints,
		Products:    products,
		FixVersion:  stringToPointer(i.FixVersion),
		Locality:    stringToPointer(i.Locality),
		Changelog:   changelog,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func NewLabel(l issue.Label) Label {
	return Label{
		ID:   l.Hash(),
		Name: string(l),
	}
}

func NewProduct(p issue.Product) Product {
	return Product{
		ID:   p.ID,
		Name: p.Name,
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

func NewAccount(a *issue.Account) *Account {
	if a == nil {
		return nil
	}

	return &Account{
		ID:           a.ID,
		EmailAddress: a.EmailAddress,
		AvatarURL:    a.AvatarURL,
		DisplayName:  a.DisplayName,
		Active:       a.Active,
		AccountType:  a.AccountType,
	}
}

func stringToPointer(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
