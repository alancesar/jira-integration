package api

import (
	"jira-integration/pkg/issue"
	"jira-integration/pkg/jira/internal"
	"jira-integration/pkg/search"
	"jira-integration/pkg/sprint"
	"strconv"
	"time"
)

type (
	JiraDate time.Time
)

type (
	Field struct {
		ID    string `json:"id"`
		Self  string `json:"self,omitempty"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
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
		Description string   `json:"description"`
		ReleaseDate JiraDate `json:"releaseDate"`
		Archived    bool     `json:"archived"`
		Released    bool     `json:"released"`
	}

	NewProjects struct {
		Field
		Child Field `json:"child"`
	}

	Status struct {
		Field
		IconURL  string         `json:"iconUrl"`
		Category StatusCategory `json:"statusCategory"`
	}

	StatusCategory struct {
		ID        uint   `json:"id"`
		Self      string `json:"self,omitempty"`
		Name      string `json:"name,omitempty"`
		Key       string `json:"key"`
		ColorName string `json:"colorName"`
	}

	Account struct {
		Self         string            `json:"self"`
		AccountID    string            `json:"accountId"`
		EmailAddress string            `json:"emailAddress"`
		AvatarURLs   map[string]string `json:"avatarUrls"`
		DisplayName  string            `json:"displayName"`
		Active       bool              `json:"active"`
		TimeZone     string            `json:"timeZone"`
		AccountType  string            `json:"accountType"`
	}

	Priority struct {
		ID      string `json:"id"`
		Self    string `json:"self"`
		IconURL string `json:"iconUrl"`
		Name    string `json:"name"`
	}

	IssueType struct {
		ID             string `json:"id"`
		Self           string `json:"self"`
		Description    string `json:"description"`
		IconURL        string `json:"iconUrl"`
		Name           string `json:"name"`
		Subtask        bool   `json:"subtask"`
		AvatarID       int    `json:"avatarId"`
		HierarchyLevel int    `json:"hierarchyLevel"`
	}

	Fields struct {
		Summary     string       `json:"summary"`
		Status      Status       `json:"status"`
		Priority    Priority     `json:"priority"`
		IssueType   IssueType    `json:"issuetype"`
		Parent      *Issue       `json:"parent,omitempty"`
		Sprints     []Sprint     `json:"customfield_10020,omitempty"`
		FixVersions []FixVersion `json:"fixVersions,omitempty"`
		Labels      []string     `json:"labels,omitempty"`
		Assignee    Account      `json:"assignee"`
		Reporter    Account      `json:"reporter"`
		StoryPoints float32      `json:"customfield_10025,omitempty"`
		NewProjects NewProjects  `json:"customfield_10444,omitempty"`
		Allocation  Field        `json:"customfield_10427,omitempty"`
		Created     string       `json:"created"`
		Updated     string       `json:"updated"`
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
		ID:      internal.ParseStringToUint(i.ID),
		Key:     i.Key,
		Summary: i.Fields.Summary,
		Status: issue.Status{
			ID:   internal.ParseStringToUint(i.Fields.Status.ID),
			Name: i.Fields.Status.Name,
			Category: issue.StatusCategory{
				ID:   i.Fields.Status.Category.ID,
				Name: i.Fields.Status.Category.Name,
			},
		},
		Priority: issue.Priority{
			ID:   internal.ParseStringToUint(i.Fields.Priority.ID),
			Name: i.Fields.Priority.Name,
		},
		Type: issue.Type{
			ID:          internal.ParseStringToUint(i.Fields.IssueType.ID),
			Description: i.Fields.IssueType.Description,
			Name:        i.Fields.IssueType.Name,
			Subtask:     i.Fields.IssueType.Subtask,
		},
		FixVersions: nil,
		Labels:      i.Fields.Labels,
		Assignee:    i.Fields.Assignee.EmailAddress,
		Reporter:    i.Fields.Reporter.EmailAddress,
		StoryPoints: uint(i.Fields.StoryPoints),
		NewProjects: i.Fields.NewProjects.JoinValues(),
		Allocation:  i.Fields.Allocation.Value,
		CreatedAt:   internal.MustParseTimeRFC3339WithTimezone(i.Fields.Created),
		UpdatedAt:   internal.MustParseTimeRFC3339WithTimezone(i.Fields.Updated),
	}

	if i.Fields.Parent != nil {
		parent := i.Fields.Parent.ToDomain()
		output.Parent = &parent
	}

	if i.Fields.Sprints != nil {
		output.Sprints = make([]sprint.Sprint, len(i.Fields.Sprints), len(i.Fields.Sprints))
		for index := range i.Fields.Sprints {
			output.Sprints[index] = i.Fields.Sprints[index].ToDomain()
		}
	}

	if i.Fields.FixVersions != nil {
		output.FixVersions = make([]issue.FixVersion, len(i.Fields.FixVersions), len(i.Fields.FixVersions))
		for index := range i.Fields.FixVersions {
			output.FixVersions[index] = i.Fields.FixVersions[index].ToDomain()
		}
	}

	return output
}

func (s Status) ToDomain() issue.Status {
	return issue.Status{
		ID:   internal.ParseStringToUint(s.ID),
		Name: s.Name,
		Category: issue.StatusCategory{
			ID:   s.Category.ID,
			Name: s.Category.Name,
		},
	}
}

func (i IssueType) ToDomain() issue.Type {
	return issue.Type{
		ID:          internal.ParseStringToUint(i.ID),
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
		ID:          internal.ParseStringToUint(f.ID),
		Name:        f.Name,
		Description: f.Description,
		Archived:    f.Archived,
		Released:    f.Released,
		ReleaseDate: time.Time(f.ReleaseDate),
	}
}

func (p NewProjects) JoinValues() string {
	if p.Value != "" && p.Child.Value != "" {
		return p.Value + " - " + p.Child.Value
	}
	return p.Value
}

func (j *JiraDate) UnmarshalJSON(bytes []byte) error {
	unquoted, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}

	parsed, err := time.Parse("2006-01-02", unquoted)
	if err != nil {
		return err
	}

	*j = JiraDate(parsed)
	return nil
}

func (j *JiraDate) MarshalJSON() ([]byte, error) {
	t := time.Time(*j)
	formatted := t.Format("2006-01-02")
	return []byte(formatted), nil
}
