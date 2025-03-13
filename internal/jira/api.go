package jira

import (
	"jira-integration/pkg/issue"
	"strconv"
	"time"
)

type (
	Date     time.Time
	DateTime time.Time
)

type (
	Paginated struct {
		NextPageToken string `json:"nextPageToken"`
	}

	Field struct {
		ID    string `json:"id"`
		Self  string `json:"self,omitempty"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
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

	FixVersion struct {
		Field
		Description string `json:"description"`
		Archived    bool   `json:"archived"`
		Released    bool   `json:"released"`
		Date        Date   `json:"releaseDate"`
	}

	Fields struct {
		Summary     string       `json:"summary"`
		Status      Status       `json:"status"`
		Priority    Priority     `json:"priority"`
		IssueType   IssueType    `json:"issuetype"`
		Parent      *Issue       `json:"parent,omitempty"`
		Sprints     []Sprint     `json:"customfield_10020,omitempty"`
		Labels      []string     `json:"labels,omitempty"`
		Assignee    *Account     `json:"assignee,omitempty"`
		Reporter    Account      `json:"reporter"`
		StoryPoints float32      `json:"customfield_10025,omitempty"`
		Product     []Field      `json:"customfield_10693,omitempty"`
		Project     Project      `json:"project"`
		FixVersions []FixVersion `json:"fixVersions"`
		Locality    Field        `json:"customfield_10696"`
		Created     DateTime     `json:"created"`
		Updated     DateTime     `json:"updated"`
	}

	Issue struct {
		ID     string `json:"id"`
		Self   string `json:"self"`
		Key    string `json:"key"`
		Fields Fields `json:"fields"`
	}

	SearchRequest struct {
		Paginated
		Fields     []string `json:"fields"`
		JQL        string   `json:"jql"`
		MaxResults int      `json:"maxResults"`
	}

	ChangelogItem struct {
		Field      string `json:"field"`
		FieldType  string `json:"fieldtype"`
		From       string `json:"from"`
		FromString string `json:"fromString"`
		To         string `json:"to"`
		ToString   string `json:"toString"`
		FieldID    string `json:"fieldId"`
	}

	Changelog struct {
		ID      string          `json:"id"`
		Author  Account         `json:"author"`
		Created int64           `json:"created"`
		Items   []ChangelogItem `json:"items"`
	}

	ChangelogRequest struct {
		Paginated
		FieldIDs       []string `json:"fieldIds"`
		IssueIDsOrKeys []string `json:"issueIdsOrKeys"`
		MaxResults     int      `json:"maxResults"`
	}

	IssueChangelog struct {
		IssueID   string      `json:"issueId"`
		Changelog []Changelog `json:"changeHistories"`
	}

	ChangelogResponse struct {
		Paginated
		IssueChangeLogs []IssueChangelog `json:"issueChangeLogs"`
	}

	SimplifiedIssue struct {
		ID string `json:"id"`
	}

	SearchResponse struct {
		Paginated
		Issues []SimplifiedIssue `json:"issues"`
	}

	GetIssueResponse struct {
		Expand string `json:"expand"`
		Issue
	}
)

func NewChangelogRequest(issueKey, nextPageToken string) ChangelogRequest {
	return ChangelogRequest{
		FieldIDs:       []string{"status"},
		IssueIDsOrKeys: []string{issueKey},
		MaxResults:     defaultMaxResults,
		Paginated: Paginated{
			NextPageToken: nextPageToken,
		},
	}
}

func NewJQLSearchRequest(jql, nextPageToken string) SearchRequest {
	return SearchRequest{
		Fields:     []string{"id"},
		JQL:        jql,
		MaxResults: defaultMaxResults,
		Paginated: Paginated{
			NextPageToken: nextPageToken,
		},
	}
}

func (i Issue) ToDomain() issue.Issue {
	output := issue.Issue{
		ID:          stringToUint(i.ID),
		Key:         i.Key,
		Summary:     i.Fields.Summary,
		Status:      i.Fields.Status.Name,
		IssueType:   i.Fields.IssueType.Name,
		Project:     i.Fields.Project.Name,
		Labels:      i.Fields.Labels,
		Reporter:    i.Fields.Reporter.EmailAddress,
		StoryPoints: uint(i.Fields.StoryPoints),
		Locality:    i.Fields.Locality.Value,
		Changelog:   nil,
		CreatedAt:   time.Time(i.Fields.Created),
		UpdatedAt:   time.Time(i.Fields.Updated),
	}

	if i.Fields.Parent != nil {
		output.Parent = i.Fields.Parent.Key
	}

	if i.Fields.Assignee != nil {
		output.Assignee = i.Fields.Assignee.EmailAddress
	}

	if i.Fields.Sprints != nil {
		output.Sprint = i.Fields.Sprints[len(i.Fields.Sprints)-1].Name
	}

	if i.Fields.Product != nil {
		output.Products = make([]string, len(i.Fields.Product), len(i.Fields.Product))
		for index, product := range i.Fields.Product {
			output.Products[index] = product.Value
		}
	}

	if i.Fields.FixVersions != nil {
		// TODO fixme
		output.FixVersion = i.Fields.FixVersions[len(i.Fields.FixVersions)-1].Name
	}

	return output
}

func (c ChangelogResponse) ToDomain() []issue.Changelog {
	var output []issue.Changelog
	for _, issues := range c.IssueChangeLogs {
		for _, changelog := range issues.Changelog {
			output = append(output, changelog.ToDomain()...)
		}
	}
	return output
}

func (c Changelog) ToDomain() []issue.Changelog {
	output := make([]issue.Changelog, len(c.Items), len(c.Items))
	for i, changelogItem := range c.Items {
		output[i] = issue.Changelog{
			ID:        stringToUint(c.ID),
			Author:    c.Author.EmailAddress,
			From:      changelogItem.FromString,
			To:        changelogItem.ToString,
			CreatedAt: time.UnixMilli(c.Created),
		}
	}

	return output
}

func (s SearchResponse) ToDomain() []string {
	output := make([]string, len(s.Issues), len(s.Issues))
	for i, simplifiedIssue := range s.Issues {
		output[i] = simplifiedIssue.ID
	}

	return output
}

func (d *Date) UnmarshalJSON(bytes []byte) error {
	unquoted, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}

	parsed, err := time.Parse(time.DateOnly, unquoted)
	if err != nil {
		return err
	}

	*d = Date(parsed)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	t := time.Time(*d)
	formatted := t.Format(time.DateOnly)
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
