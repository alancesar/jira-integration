package gateway

import (
	"fmt"
	"jira-integration/pkg/issue"
	"jira-integration/pkg/jira"
	"strings"
)

type (
	Gateway struct {
		client *jira.Client
	}
)

func New(client *jira.Client) *Gateway {
	return &Gateway{
		client: client,
	}
}

func (g Gateway) StreamAllIssues(args ...string) <-chan issue.Issue {
	issues := make(chan issue.Issue)
	parsedArgs := strings.Join(args, " and ")

	go func() {
		jql := fmt.Sprintf("%s and issuetype = Epic", parsedArgs)
		params := jira.DefaultSearchRequest(jql, 0)
		g.streamAllIssues(params, issues)

		jql = fmt.Sprintf("%s and issuetype in (standardIssueTypes())", parsedArgs)
		params = jira.DefaultSearchRequest(jql, 0)
		g.streamAllIssues(params, issues)

		jql = fmt.Sprintf("%s and issuetype in (subTaskIssueTypes())", parsedArgs)
		params = jira.DefaultSearchRequest(jql, 0)
		g.streamAllIssues(params, issues)

		close(issues)
	}()

	return issues
}

func (g Gateway) streamAllIssues(params jira.SearchRequest, issues chan issue.Issue) {
	response, err := g.client.Search(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i := range response.Issues {
		issues <- i
	}

	if response.StartAt+response.MaxResults < response.Total {
		g.streamAllIssues(jira.SearchRequest{
			Fields:       params.Fields,
			FieldsByKeys: params.FieldsByKeys,
			JQL:          params.JQL,
			MaxResults:   params.MaxResults,
			StartAt:      params.StartAt + params.MaxResults,
		}, issues)
	}
}
