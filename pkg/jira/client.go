package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jira-integration/pkg/issue"
	"jira-integration/pkg/jira/internal"
	"jira-integration/pkg/jira/internal/agile"
	"jira-integration/pkg/jira/internal/api"
	"jira-integration/pkg/search"
	"jira-integration/pkg/sprint"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	jiraCloudAPIBasePath = "https://bexs.atlassian.net/rest/api/3"
	jiraAgileAPIBasePath = "https://bexs.atlassian.net/rest/agile/1.0"
	defaultMaxResults    = 50
)

type (
	SearchRequest struct {
		Fields       []string `json:"fields"`
		FieldsByKeys bool     `json:"fieldsByKeys"`
		JQL          string   `json:"jql"`
		MaxResults   int      `json:"maxResults"`
		StartAt      int      `json:"startAt"`
	}

	Credentials struct {
		Username string
		Password string
	}

	Client struct {
		credentials Credentials
		httpClient  *http.Client
	}

	Streamer[T any, I Interactor[T]] struct {
		httpClient *http.Client
	}

	Response[T any] struct {
		MaxResults uint `json:"maxResults"`
		StartAt    uint `json:"startAt"`
		IsLast     bool `json:"isLast"`
		Values     []T  `json:"values"`
	}

	Interactor[T any] interface {
		ToDomain() T
	}
)

func NewClient(credentials Credentials, client *http.Client) *Client {
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	basicAuthRoundTripper := internal.NewBasicAuthRoundTripper(credentials.Username, credentials.Password, transport)

	return &Client{
		credentials: credentials,
		httpClient: &http.Client{
			Transport: basicAuthRoundTripper,
		},
	}
}

func NewStreamer[K any, I Interactor[K]](client *http.Client) *Streamer[K, I] {
	return &Streamer[K, I]{
		httpClient: client,
	}
}

func DefaultSearchRequest(jql string, startAt int) SearchRequest {
	return SearchRequest{
		Fields: []string{
			"summary",
			"status",
			"assignee",
			"priority",
			"issuetype",
			"fixVersions",
			"labels",
			"sprint",
			"parent",
			"created",
			"updated",
			"customfield_10025",
			"customfield_10020",
			"customfield_10427",
			"customfield_10444",
		},
		JQL:        jql,
		MaxResults: defaultMaxResults,
		StartAt:    startAt,
	}
}

func (c Client) Search(params SearchRequest) (search.Response, error) {
	requestURL := fmt.Sprintf("%s/search", jiraCloudAPIBasePath)

	rawRequest, err := json.Marshal(&params)
	if err != nil {
		return search.Response{}, err
	}
	body := bytes.NewReader(rawRequest)

	response, err := c.httpClient.Post(requestURL, "application/json", body)
	if err != nil {
		return search.Response{}, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return search.Response{}, fmt.Errorf("bad status from %s: %d", requestURL, response.StatusCode)
	}

	var output api.SearchResponse
	if err = json.NewDecoder(response.Body).Decode(&output); err != nil {
		return search.Response{}, err
	}

	return output.ToDomain(), nil
}

func (c Client) GetStatuses() ([]issue.Status, error) {
	baseURL := fmt.Sprintf("%s/status", jiraCloudAPIBasePath)
	var statuses []api.Status
	if err := doGetRequest(c.httpClient, baseURL, &statuses); err != nil {
		return nil, err
	}

	output := make([]issue.Status, len(statuses), len(statuses))
	for i, s := range statuses {
		output[i] = s.ToDomain()
	}

	return output, nil
}

func (c Client) GetIssueTypes() ([]issue.Type, error) {
	baseURL := fmt.Sprintf("%s/issuetype", jiraCloudAPIBasePath)
	var issueTypes []api.IssueType
	if err := doGetRequest(c.httpClient, baseURL, &issueTypes); err != nil {
		return nil, err
	}

	output := make([]issue.Type, len(issueTypes), len(issueTypes))
	for i, issueType := range issueTypes {
		output[i] = issueType.ToDomain()
	}

	return output, nil
}

func (c Client) StreamSprints(boardID int) <-chan sprint.Sprint {
	baseURL := fmt.Sprintf("%s/board/%d/sprint", jiraAgileAPIBasePath, boardID)
	streamer := NewStreamer[sprint.Sprint, agile.Sprint](c.httpClient)
	return streamer.Stream(baseURL)
}

func (c Client) StreamFixVersions(boardID int) <-chan issue.FixVersion {
	baseURL := fmt.Sprintf("%s/board/%d/version", jiraAgileAPIBasePath, boardID)
	streamer := NewStreamer[issue.FixVersion, agile.FixVersion](c.httpClient)
	return streamer.Stream(baseURL)
}

func (s Streamer[K, I]) Stream(url string) <-chan K {
	items := make(chan K)

	go func() {
		if err := s.stream(url, 0, items); err != nil {
			log.Println(err)
		}
		close(items)
	}()

	return items
}

func (s Streamer[K, I]) stream(baseURL string, startAt uint, items chan<- K) error {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	query := parsedURL.Query()
	query.Add("startAt", strconv.Itoa(int(startAt)))
	parsedURL.RawQuery = query.Encode()

	var output Response[I]
	if err := doGetRequest(s.httpClient, parsedURL.String(), &output); err != nil {
		return err
	}

	for _, value := range output.Values {
		items <- value.ToDomain()
	}

	if !output.IsLast {
		return s.stream(baseURL, startAt+50, items)
	}

	return nil
}

func doGetRequest(client *http.Client, url string, output any) error {
	response, err := client.Get(url)
	if err != nil {
		return err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status from %s: %d", url, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(output); err != nil {
		return err
	}

	return nil
}
