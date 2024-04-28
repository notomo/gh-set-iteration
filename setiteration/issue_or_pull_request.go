package setiteration

import (
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type Content struct {
	ID     string
	Title  string
	Closed bool
}

type DummyContent = Content
type IssueOrPullRequest struct {
	Content      `graphql:"... on Issue"`
	DummyContent `graphql:"... on PullRequest"` // dummy to avoid `Content redeclared`
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
	return &query.Repository.IssueOrPullRequest.Content, nil
}
