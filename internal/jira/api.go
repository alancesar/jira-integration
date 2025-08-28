package jira

import (
	"jira-integration/pkg/issue"
	"strconv"
	"time"
)

const (
	AvatarSize48x48 = "48x48"
	AvatarSize32x32 = "32x32"
	AvatarSize24x24 = "24x24"
	AvatarSize16x16 = "16x16"
)

var (
	avatarSizes = []AvatarSize{
		AvatarSize48x48,
		AvatarSize32x32,
		AvatarSize24x24,
		AvatarSize16x16,
	}
)

type (
	Date     time.Time
	DateTime time.Time

	AvatarSize string
)

type (
	Paginated struct {
		NextPageToken string `json:"nextPageToken,omitempty"`
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

	AvatarURLs map[AvatarSize]string

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

	Sprints     []Sprint
	FixVersions []FixVersion

	Fields struct {
		Summary     string      `json:"summary"`
		Status      Status      `json:"status"`
		Priority    Priority    `json:"priority"`
		IssueType   IssueType   `json:"issuetype"`
		Parent      *Issue      `json:"parent,omitempty"`
		Sprints     Sprints     `json:"customfield_10020,omitempty"`
		Labels      []string    `json:"labels,omitempty"`
		Assignee    *Account    `json:"assignee,omitempty"`
		Reporter    Account     `json:"reporter"`
		StoryPoints float32     `json:"customfield_10025,omitempty"`
		Product     []Field     `json:"customfield_10693,omitempty"`
		Project     Project     `json:"project"`
		FixVersions FixVersions `json:"fixVersions,omitempty"`
		Locality    Field       `json:"customfield_10696"`
		Created     DateTime    `json:"created"`
		Updated     DateTime    `json:"updated"`
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

	SearchResponseIssue struct {
		ID     string `json:"id"`
		Key    string `json:"key"`
		Fields Fields `json:"fields"`
	}

	SearchResponse struct {
		Paginated
		Issues []SearchResponseIssue `json:"issues"`
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
		Fields:     []string{"created", "updated"},
		JQL:        jql,
		MaxResults: defaultMaxResults,
		Paginated: Paginated{
			NextPageToken: nextPageToken,
		},
	}
}

func (i Issue) ToDomain() issue.Issue {
	output := issue.Issue{
		Stamp: issue.Stamp{
			ID:        stringToUint(i.ID),
			Key:       i.Key,
			CreatedAt: time.Time(i.Fields.Created),
			UpdatedAt: time.Time(i.Fields.Updated),
		},
		Summary:     i.Fields.Summary,
		Status:      i.Fields.Status.Name,
		IssueType:   i.Fields.IssueType.Name,
		Project:     i.Fields.Project.Name,
		Labels:      issue.NewLabels(i.Fields.Labels),
		Reporter:    i.Fields.Reporter.ToDomain(),
		StoryPoints: uint(i.Fields.StoryPoints),
		Locality:    i.Fields.Locality.Value,
		Changelog:   nil,
	}

	if i.Fields.Parent != nil {
		parent := i.Fields.Parent.ToDomain()
		output.Parent = &parent
	}

	if i.Fields.Assignee != nil {
		assignee := i.Fields.Assignee.ToDomain()
		output.Assignee = &assignee
	}

	if len(i.Fields.Sprints) != 0 {
		output.Sprint = i.Fields.Sprints.GetLast().ToDomain()
	}

	if len(i.Fields.Product) != 0 {
		output.Products = make([]issue.Product, len(i.Fields.Product), len(i.Fields.Product))
		for index, product := range i.Fields.Product {
			output.Products[index] = issue.Product{
				ID:   stringToUint(product.ID),
				Name: product.Value,
			}
		}
	}

	if len(i.Fields.FixVersions) != 0 {
		lastFixVersion := i.Fields.FixVersions.GetLast()
		output.FixVersion = lastFixVersion.Name
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

func (s SearchResponse) ToDomain() []issue.Stamp {
	output := make([]issue.Stamp, len(s.Issues), len(s.Issues))
	for i, searchIssueResponse := range s.Issues {
		output[i] = issue.Stamp{
			ID:        stringToUint(searchIssueResponse.ID),
			Key:       searchIssueResponse.Key,
			CreatedAt: time.Time(searchIssueResponse.Fields.Created),
			UpdatedAt: time.Time(searchIssueResponse.Fields.Updated),
		}
	}

	return output
}

func (f FixVersions) GetLast() FixVersion {
	if len(f) == 0 {
		return FixVersion{}
	}

	last := f[0]
	for _, fv := range f {
		if time.Time(fv.Date).After(time.Time(last.Date)) {
			last = fv
		}
	}

	return last
}

func (s Sprints) GetLast() *Sprint {
	if len(s) == 0 {
		return nil
	}

	last := s[0]
	for _, sprint := range s {
		if sprint.EndDate.After(last.EndDate) {
			last = sprint
		}
	}

	return &last
}

func (s Sprint) ToDomain() *issue.Sprint {
	return &issue.Sprint{
		ID:          s.ID,
		Name:        s.Name,
		State:       s.State,
		Goal:        s.Goal,
		StartedAt:   s.StartDate,
		EndedAt:     s.EndDate,
		CompletedAt: s.CompleteDate,
	}
}

func (a AvatarURLs) GetLargest() string {
	for _, size := range avatarSizes {
		if url, ok := a[size]; ok {
			return url
		}
	}

	return ""
}

func (a Account) ToDomain() issue.Account {
	return issue.Account{
		ID:           a.AccountID,
		EmailAddress: a.EmailAddress,
		AvatarURL:    a.AvatarURLs.GetLargest(),
		DisplayName:  a.DisplayName,
		Active:       a.Active,
		AccountType:  a.AccountType,
	}
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
