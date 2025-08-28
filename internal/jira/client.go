package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jira-integration/pkg/issue"
	"net/http"
	"net/url"
)

const (
	jiraCloudAPIBasePath = "https://bexs.atlassian.net/rest/api/3"
	jiraAgileAPIBasePath = "https://bexs.atlassian.net/rest/agile/1.0"
	defaultMaxResults    = 500
)

var (
	defaultFieldValues = []string{
		"summary",
		"status",
		"issuetype",
		"parent",
		"labels",
		"assignee",
		"reporter",
		"project",
		"fixVersions",
		"created",
		"updated",
		"customfield_10014",
		"customfield_10020",
		"customfield_10025",
		"customfield_10693",
		"customfield_10696",
	}

	BadStatusErr = errors.New("bad status")
)

type (
	Credentials struct {
		Username string
		Password string
	}

	Client struct {
		credentials Credentials
		httpClient  *http.Client
	}
)

func NewClient(credentials Credentials, client *http.Client) *Client {
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	basicAuthRoundTripper := NewBasicAuthRoundTripper(credentials.Username, credentials.Password, transport)

	return &Client{
		credentials: credentials,
		httpClient: &http.Client{
			Transport: basicAuthRoundTripper,
		},
	}
}

func (c Client) SearchIssuesByJQL(_ context.Context, jql, nextPageToken string) ([]issue.Stamp, string, error) {
	requestURL := fmt.Sprintf("%s/search/jql", jiraCloudAPIBasePath)
	params := NewJQLSearchRequest(jql, nextPageToken)
	rawRequest, err := json.Marshal(&params)
	if err != nil {
		return nil, "", err
	}

	response, err := c.httpClient.Post(requestURL, "application/json", bytes.NewReader(rawRequest))
	if err != nil {
		return nil, "", err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("%w: %s: %d", BadStatusErr, requestURL, response.StatusCode)
	}

	var output SearchResponse
	if err = json.NewDecoder(response.Body).Decode(&output); err != nil {
		return nil, "", err
	}

	return output.ToDomain(), output.NextPageToken, nil
}

func (c Client) GetIssueByID(_ context.Context, issueID uint) (issue.Issue, error) {
	parsedURL, err := url.Parse(fmt.Sprintf("%s/issue/%d", jiraCloudAPIBasePath, issueID))
	if err != nil {
		return issue.Issue{}, err
	}

	query := parsedURL.Query()
	for _, value := range defaultFieldValues {
		query.Add("fields", value)
	}
	parsedURL.RawQuery = query.Encode()

	resp, err := c.httpClient.Get(parsedURL.String())
	if err != nil {
		return issue.Issue{}, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return issue.Issue{}, fmt.Errorf("%w: %s: %d", BadStatusErr, parsedURL, resp.StatusCode)
	}

	var output GetIssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return issue.Issue{}, err
	}

	return output.ToDomain(), nil
}

func (c Client) GetIssueChangelog(_ context.Context, issueKey, nextPageToken string) ([]issue.Changelog, string, error) {
	baseURL := fmt.Sprintf("%s/changelog/bulkfetch", jiraCloudAPIBasePath)
	params := NewChangelogRequest(issueKey, nextPageToken)
	rawRequest, err := json.Marshal(&params)
	if err != nil {
		return nil, "", err
	}

	response, err := c.httpClient.Post(baseURL, "application/json", bytes.NewReader(rawRequest))
	if err != nil {
		return nil, "", err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("%w: %s: %d", BadStatusErr, baseURL, response.StatusCode)
	}

	var output ChangelogResponse
	if err = json.NewDecoder(response.Body).Decode(&output); err != nil {
		return nil, "", err
	}

	return output.ToDomain(), output.NextPageToken, nil
}

func (c Client) GetSprint(_ context.Context, sprintID uint) (*issue.Sprint, error) {
	baseURL := fmt.Sprintf("%s/sprint/%d", jiraAgileAPIBasePath, sprintID)
	response, err := c.httpClient.Get(baseURL)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s: %d", BadStatusErr, baseURL, response.StatusCode)
	}

	var output Sprint
	if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
		return nil, err
	}

	return output.ToDomain(), nil
}
