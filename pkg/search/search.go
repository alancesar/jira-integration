package search

import "jira-integration/pkg/issue"

type (
	Response struct {
		StartAt    int
		MaxResults int
		Total      int
		Issues     []issue.Issue
	}
)
