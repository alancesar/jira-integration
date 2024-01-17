package gateway

import (
	"context"
	"fmt"
	"jira-integration/pkg/issue"
	"jira-integration/pkg/jira"
	"jira-integration/pkg/sprint"
	"log"
	"strings"
)

type (
	Database interface {
		SaveIssueType(ctx context.Context, it issue.Type) error
		SaveStatus(ctx context.Context, s issue.Status) error
		SaveFixVersion(ctx context.Context, fv issue.FixVersion) error
		SaveIssue(ctx context.Context, i issue.Issue) error
		SaveSprint(ctx context.Context, s sprint.Sprint) error
	}

	Gateway struct {
		client *jira.Client
		db     Database
	}
)

func New(db Database, client *jira.Client) *Gateway {
	return &Gateway{
		db:     db,
		client: client,
	}
}

func (g Gateway) SyncIssues(ctx context.Context, args ...string) {
	issues := make(chan issue.Issue)
	parsedArgs := strings.Join(args, " and ")

	go func() {
		jql := fmt.Sprintf("%s and issuetype = Epic", parsedArgs)
		params := jira.DefaultSearchRequest(jql, 0)
		g.streamIssues(params, issues)

		jql = fmt.Sprintf("%s and issuetype in (standardIssueTypes())", parsedArgs)
		params = jira.DefaultSearchRequest(jql, 0)
		g.streamIssues(params, issues)

		jql = fmt.Sprintf("%s and issuetype in (subTaskIssueTypes())", parsedArgs)
		params = jira.DefaultSearchRequest(jql, 0)
		g.streamIssues(params, issues)

		close(issues)
	}()

	for iss := range issues {
		fmt.Println("fetching issue", iss.Key)
		if err := g.db.SaveIssue(ctx, iss); err != nil {
			log.Println(err)
		}
	}
}

func (g Gateway) SyncDependencies(ctx context.Context) error {
	fixVersions := g.client.StreamFixVersions()
	for fixVersion := range fixVersions {
		fmt.Println("fetching fix version", fixVersion.Name)
		if err := g.db.SaveFixVersion(ctx, fixVersion); err != nil {
			return err
		}
	}

	sprints := g.client.StreamSprints()
	for s := range sprints {
		fmt.Println("fetching sprint", s.Name)
		if err := g.db.SaveSprint(ctx, s); err != nil {
			return err
		}
	}

	return nil
}

func (g Gateway) Setup(ctx context.Context) error {
	issueTypes, err := g.client.GetIssueTypes()
	if err != nil {
		return err
	}

	for _, issueType := range issueTypes {
		fmt.Println("fetching issue type", issueType.Name)
		if err := g.db.SaveIssueType(ctx, issueType); err != nil {
			return err
		}
	}

	statuses, err := g.client.GetStatuses()
	if err != nil {
		return err
	}

	for _, status := range statuses {
		fmt.Println("fetching status", status.Name)
		if err := g.db.SaveStatus(ctx, status); err != nil {
			return err
		}
	}

	return nil
}

func (g Gateway) streamIssues(params jira.SearchRequest, issues chan issue.Issue) {
	response, err := g.client.Search(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i := range response.Issues {
		issues <- i
	}

	if response.StartAt+response.MaxResults < response.Total {
		g.streamIssues(jira.SearchRequest{
			Fields:       params.Fields,
			FieldsByKeys: params.FieldsByKeys,
			JQL:          params.JQL,
			MaxResults:   params.MaxResults,
			StartAt:      params.StartAt + params.MaxResults,
		}, issues)
	}
}
