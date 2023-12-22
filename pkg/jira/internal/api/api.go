package api

import (
	"jira-integration/pkg/issue"
	"jira-integration/pkg/search"
	"jira-integration/pkg/sprint"
	"strconv"
	"time"
)

type (
	Date     time.Time
	DateTime time.Time
)

type (
	Field struct {
		ID    string `json:"id"`
		Self  string `json:"self,omitempty"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}

	FieldWithSubfield struct {
		Field
		Child Field `json:"child"`
	}

	Project struct {
		Field
		Key            string     `json:"key"`
		ProjectTypeKey string     `json:"projectTypeKey"`
		Simplified     bool       `json:"simplified"`
		AvatarUrls     AvatarURLs `json:"avatarUrls"`
	}

	Sprint struct {
		ID           uint      `json:"id"`
		Name         string    `json:"name"`
		State        string    `json:"state"`
		BoardID      uint      `json:"boardId"`
		Goal         string    `json:"goal"`
		StartDate    time.Time `json:"startDate"`
		EndDate      time.Time `json:"endDate"`
		CompleteDate time.Time `json:"completeDate,omitempty"`
	}

	FixVersion struct {
		Field
		Description string `json:"description"`
		ReleaseDate Date   `json:"releaseDate"`
		Archived    bool   `json:"archived"`
		Released    bool   `json:"released"`
	}

	Status struct {
		Field
		IconURL  string         `json:"iconUrl"`
		Category StatusCategory `json:"statusCategory"`
	}

	Resolution struct {
		Field
		Description string `json:"description"`
	}

	StatusCategory struct {
		ID        uint   `json:"id"`
		Self      string `json:"self,omitempty"`
		Name      string `json:"name,omitempty"`
		Key       string `json:"key"`
		ColorName string `json:"colorName"`
	}

	Progress struct {
		Progress int `json:"progress"`
		Total    int `json:"total"`
		Percent  int `json:"percent"`
	}

	AvatarURLs map[string]string

	Account struct {
		Self         string     `json:"self"`
		AccountID    string     `json:"accountId"`
		EmailAddress string     `json:"emailAddress"`
		AvatarURLs   AvatarURLs `json:"avatarUrls"`
		DisplayName  string     `json:"displayName"`
		Active       bool       `json:"active"`
		TimeZone     string     `json:"timeZone"`
		AccountType  string     `json:"accountType"`
	}

	Priority struct {
		Field
		IconURL string `json:"iconUrl"`
	}

	IssueType struct {
		Field
		Description    string `json:"description"`
		IconURL        string `json:"iconUrl"`
		Subtask        bool   `json:"subtask"`
		AvatarID       int    `json:"avatarId"`
		HierarchyLevel int    `json:"hierarchyLevel"`
	}

	Fields struct {
		Summary               string             `json:"summary"`
		Status                Status             `json:"status"`
		Priority              Priority           `json:"priority"`
		IssueType             IssueType          `json:"issuetype"`
		Parent                *Issue             `json:"parent,omitempty"`
		Sprints               []Sprint           `json:"customfield_10020,omitempty"`
		FixVersions           []FixVersion       `json:"fixVersions,omitempty"`
		Labels                []string           `json:"labels,omitempty"`
		Assignee              *Account           `json:"assignee,omitempty"`
		Reporter              Account            `json:"reporter"`
		StoryPoints           float32            `json:"customfield_10025,omitempty"`
		NewProjects           FieldWithSubfield  `json:"customfield_10444,omitempty"`
		Allocation            Field              `json:"customfield_10427,omitempty"`
		Resolution            *Resolution        `json:"resolution"`
		Environment           Field              `json:"customfield_10448"`
		Creator               Account            `json:"creator"`
		Squad                 *Field             `json:"customfield_10183"`
		RequestType           string             `json:"customfield_10184"`
		Project               Project            `json:"project"`
		TimeSpent             int                `json:"timespent"`
		AggregateTimeSpent    int                `json:"aggregatetimespent"`
		AggregateTimeEstimate int                `json:"aggregatetimeestimate"`
		Progress              Progress           `json:"progress"`
		AggregateProgress     Progress           `json:"aggregateprogress"`
		System                *FieldWithSubfield `json:"customfield_10231,omitempty"`
		ResolutionDate        *DateTime          `json:"resolutiondate,omitempty"`
		Created               DateTime           `json:"created"`
		Updated               DateTime           `json:"updated"`
	}

	Issue struct {
		ID     string `json:"id"`
		Self   string `json:"self"`
		Key    string `json:"key"`
		Fields Fields `json:"fields"`
	}

	SearchResponse struct {
		Expand     string  `json:"expand"`
		StartAt    int     `json:"startAt"`
		MaxResults int     `json:"maxResults"`
		Total      int     `json:"total"`
		Issues     []Issue `json:"issues"`
	}
)

func (s SearchResponse) ToDomain() search.Response {
	issues := make([]issue.Issue, len(s.Issues), len(s.Issues))
	for i := range s.Issues {
		issues[i] = s.Issues[i].ToDomain()
	}

	return search.Response{
		StartAt:    s.StartAt,
		MaxResults: s.MaxResults,
		Total:      s.Total,
		Issues:     issues,
	}
}

func (i Issue) ToDomain() issue.Issue {
	output := issue.Issue{
		ID:      stringToUint(i.ID),
		Key:     i.Key,
		Summary: i.Fields.Summary,
		Status: issue.Status{
			ID:   stringToUint(i.Fields.Status.ID),
			Name: i.Fields.Status.Name,
			Category: issue.StatusCategory{
				ID:   i.Fields.Status.Category.ID,
				Name: i.Fields.Status.Category.Name,
			},
		},
		Priority: issue.Priority{
			ID:   stringToUint(i.Fields.Priority.ID),
			Name: i.Fields.Priority.Name,
		},
		Type: issue.Type{
			ID:          stringToUint(i.Fields.IssueType.ID),
			Description: i.Fields.IssueType.Description,
			Name:        i.Fields.IssueType.Name,
			Subtask:     i.Fields.IssueType.Subtask,
		},
		Project:               i.Fields.Project.ToDomain(),
		Progress:              i.Fields.Progress.ToDomain(),
		AggregateProgress:     i.Fields.AggregateProgress.ToDomain(),
		AggregateTimeSpent:    i.Fields.AggregateTimeSpent,
		AggregateTimeEstimate: i.Fields.AggregateTimeEstimate,
		TimeSpent:             i.Fields.TimeSpent,
		Labels:                i.Fields.Labels,
		Reporter:              i.Fields.Reporter.ToDomain(),
		StoryPoints:           uint(i.Fields.StoryPoints),
		NewProjects:           i.Fields.NewProjects.JoinValues(),
		Allocation:            i.Fields.Allocation.Value,
		CreatedAt:             time.Time(i.Fields.Created),
		UpdatedAt:             time.Time(i.Fields.Updated),
	}

	if i.Fields.Parent != nil {
		parent := i.Fields.Parent.ToDomain()
		output.Parent = &parent
	}

	if i.Fields.Assignee != nil {
		assignee := i.Fields.Assignee.ToDomain()
		output.Assignee = &assignee
	}

	if i.Fields.Sprints != nil {
		output.Sprints = make([]sprint.Sprint, len(i.Fields.Sprints), len(i.Fields.Sprints))
		for index := range i.Fields.Sprints {
			output.Sprints[index] = i.Fields.Sprints[index].ToDomain()
		}
	}

	if i.Fields.Resolution != nil {
		resolution := i.Fields.Resolution.ToDomain()
		output.Resolution = &resolution
	}

	if i.Fields.System != nil {
		system := i.Fields.System.JoinValues()
		output.System = &system
	}

	if i.Fields.Squad != nil {
		output.Squad = &i.Fields.Squad.Value
	}

	if i.Fields.FixVersions != nil {
		output.FixVersions = make([]issue.FixVersion, len(i.Fields.FixVersions), len(i.Fields.FixVersions))
		for index := range i.Fields.FixVersions {
			output.FixVersions[index] = i.Fields.FixVersions[index].ToDomain()
		}
	}

	if i.Fields.ResolutionDate != nil {
		resolvedAt := time.Time(*i.Fields.ResolutionDate)
		output.ResolvedAt = &resolvedAt
	}

	return output
}

func (s Status) ToDomain() issue.Status {
	return issue.Status{
		ID:   stringToUint(s.ID),
		Name: s.Name,
		Category: issue.StatusCategory{
			ID:   s.Category.ID,
			Name: s.Category.Name,
		},
	}
}

func (r Resolution) ToDomain() issue.Resolution {
	return issue.Resolution{
		ID:          stringToUint(r.ID),
		Description: r.Description,
		Name:        r.Name,
	}
}

func (i IssueType) ToDomain() issue.Type {
	return issue.Type{
		ID:          stringToUint(i.ID),
		Description: i.Description,
		Name:        i.Name,
		Subtask:     i.Subtask,
	}
}

func (s Sprint) ToDomain() sprint.Sprint {
	return sprint.Sprint{
		ID:          s.ID,
		Name:        s.Name,
		State:       sprint.State(s.State),
		Goal:        s.Goal,
		StartedAt:   s.StartDate,
		EndedAt:     s.EndDate,
		CompletedAt: s.CompleteDate,
	}
}

func (f FixVersion) ToDomain() issue.FixVersion {
	return issue.FixVersion{
		ID:          stringToUint(f.ID),
		Name:        f.Name,
		Description: f.Description,
		Archived:    f.Archived,
		Released:    f.Released,
		ReleaseDate: time.Time(f.ReleaseDate),
	}
}

func (au AvatarURLs) Larger() string {
	orderedKeys := []string{"48x48", "32x32", "24x24", "16x16"}
	for _, key := range orderedKeys {
		if avatar, ok := au[key]; ok {
			return avatar
		}
	}

	return ""
}

func (a Account) ToDomain() issue.Account {
	return issue.Account{
		ID:           a.AccountID,
		EmailAddress: a.EmailAddress,
		AvatarURL:    a.AvatarURLs.Larger(),
		DisplayName:  a.DisplayName,
		Active:       a.Active,
		TimeZone:     a.TimeZone,
		AccountType:  a.AccountType,
	}
}

func (p Progress) ToDomain() issue.Progress {
	return issue.Progress{
		Progress: p.Progress,
		Total:    p.Total,
		Percent:  p.Percent,
	}
}

func (p Project) ToDomain() issue.Project {
	return issue.Project{
		ID:   stringToUint(p.ID),
		Key:  p.Key,
		Name: p.Name,
	}
}

func (p FieldWithSubfield) JoinValues() string {
	if p.Value != "" && p.Child.Value != "" {
		return p.Value + " - " + p.Child.Value
	}
	return p.Value
}

func (d *Date) UnmarshalJSON(bytes []byte) error {
	unquoted, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}

	parsed, err := time.Parse("2006-01-02", unquoted)
	if err != nil {
		return err
	}

	*d = Date(parsed)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	t := time.Time(*d)
	formatted := t.Format("2006-01-02")
	return []byte(formatted), nil
}

func (dt *DateTime) UnmarshalJSON(bytes []byte) error {
	unquoted, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}

	parsed, err := time.Parse("2006-01-02T15:04:05.999-0700", unquoted)
	if err != nil {
		return err
	}
	*dt = DateTime(parsed)
	return nil
}

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(*dt)
	formatted := t.Format("2006-01-02T15:04:05.999-0700")
	return []byte(formatted), nil
}

func stringToUint(raw string) uint {
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}

	return uint(parsed)
}
