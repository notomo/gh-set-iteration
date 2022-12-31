package setiteration

import (
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type Issue struct {
	ID    string
	Title string
}

type PullRequest struct {
	ID    string
	Title string
}

type Content struct {
	ID    string
	Title string
}

type IssueOrPullRequest struct {
	Issue       `graphql:"... on Issue"`
	PullRequest `graphql:"... on PullRequest"`
}

type GetIssueOrPullRequestQuery struct {
	Repository struct {
		IssueOrPullRequest IssueOrPullRequest `graphql:"issueOrPullRequest(number: $issueOrPullRequestNumber)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func GetIssueOrPullRequest(
	gql api.GQLClient,
	descriptor IssueOrPullRequestDescriptor,
) (*Content, error) {
	var query GetIssueOrPullRequestQuery
	vars := map[string]interface{}{
		"owner":                    graphql.String(descriptor.Repository.Owner),
		"name":                     graphql.String(descriptor.Repository.Name),
		"issueOrPullRequestNumber": graphql.Int(descriptor.Number),
	}
	if err := gql.Query("GetIssueOrPullRequest", &query, vars); err != nil {
		return nil, err
	}

	if query.Repository.IssueOrPullRequest.Issue.ID != "" {
		return &Content{
			ID:    query.Repository.IssueOrPullRequest.Issue.ID,
			Title: query.Repository.IssueOrPullRequest.Issue.Title,
		}, nil
	}
	return &Content{
		ID:    query.Repository.IssueOrPullRequest.PullRequest.ID,
		Title: query.Repository.IssueOrPullRequest.PullRequest.Title,
	}, nil
}
