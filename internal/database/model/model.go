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

	FixVersion struct {
		ID          uint `gorm:"primarykey"`
		Name        string
		Description string
		Archived    bool
		Released    bool
		ReleaseDate time.Time
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

	Resolution struct {
		ID          uint `gorm:"primarykey"`
		Name        string
		Description string
	}

	Priority struct {
		ID   uint `gorm:"primarykey"`
		Name string
	}

	IssueType struct {
		ID          uint `gorm:"primarykey"`
		Description string
		Name        string
		Subtask     bool
	}

	Issue struct {
		ID           uint   `gorm:"primarykey"`
		Key          string `gorm:"unique"`
		Summary      string
		Project      Project
		ProjectID    uint
		StatusID     uint
		Resolution   *Resolution
		ResolutionID *uint
		Status       Status
		PriorityID   uint
		Priority     Priority
		IssueTypeID  uint
		IssueType    IssueType
		ParentID     *uint
		Parent       *Issue
		Sprints      []Sprint     `gorm:"many2many:issue_sprints"`
		FixVersions  []FixVersion `gorm:"many2many:issue_fix_versions"`
		Labels       datatypes.JSON
		AssigneeID   *string
		Assignee     *Account
		ReporterID   string
		Reporter     Account
		StoryPoints  *uint
		NewProjects  *string
		Allocation   *string
		TimeSpent    int
		Squad        *string
		System       *string
		CreatedAt    time.Time `gorm:"autoCreateTime:false"`
		UpdatedAt    time.Time `gorm:"autoUpdateTime:false"`
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
		PriorityID:  i.Priority.ID,
		Priority:    NewPriority(i.Priority),
		IssueTypeID: i.Type.ID,
		IssueType:   NewIssueType(i.Type),
		Sprints:     NewSprints(i.Sprints),
		FixVersions: NewFixVersions(i.FixVersions),
		Labels:      labels,
		ReporterID:  i.Reporter.ID,
		Reporter:    NewAccount(i.Reporter),
		StoryPoints: uintToPointer(i.StoryPoints),
		NewProjects: stringToPointer(i.NewProjects),
		Allocation:  stringToPointer(i.Allocation),
		TimeSpent:   i.TimeSpent,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}

	if i.System != nil {
		output.System = i.System
	}

	if i.Squad != nil {
		output.Squad = i.Squad
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

	if i.Resolution != nil {
		resolution := NewResolution(*i.Resolution)
		output.Resolution = &resolution
		output.ResolutionID = &resolution.ID
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

func NewResolution(r issue.Resolution) Resolution {
	return Resolution{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
}

func NewStatusCategory(c issue.StatusCategory) StatusCategory {
	return StatusCategory{
		ID:   c.ID,
		Name: c.Name,
	}
}

func NewPriority(p issue.Priority) Priority {
	return Priority{
		ID:   p.ID,
		Name: p.Name,
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

func NewFixVersion(fv issue.FixVersion) FixVersion {
	return FixVersion{
		ID:          fv.ID,
		Name:        fv.Name,
		Description: fv.Description,
		Archived:    fv.Archived,
		Released:    fv.Released,
		ReleaseDate: fv.ReleaseDate,
	}
}

func NewFixVersions(fvs []issue.FixVersion) []FixVersion {
	output := make([]FixVersion, len(fvs), len(fvs))
	for i := range fvs {
		output[i] = NewFixVersion(fvs[i])
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

func uintToPointer(value uint) *uint {
	if value == 0 {
		return nil
	}

	return &value
}

func stringToPointer(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
