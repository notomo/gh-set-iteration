package setiteration

import (
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type Iteration struct {
	ID        string
	Title     string
	StartDate string
	Duration  int
}

type Item struct {
	ID      string
	Content IssueOrPullRequest
}

type ProjectV2IterationField struct {
	ID            string
	Name          string
	Configuration struct {
		Iterations []Iteration
	}
}

func (f *ProjectV2IterationField) SelectIteration(startDate string) *Iteration {
	for _, iteration := range f.Configuration.Iterations {
		iteration := iteration
		if iteration.StartDate == startDate {
			return &iteration
		}
	}
	return nil
}

type ProjectV2 struct {
	ID    string
	Field struct {
		ProjectV2IterationField `graphql:"... on ProjectV2IterationField"`
	} `graphql:"field(name: $fieldName)"`
	Items struct {
		Nodes []Item
	} `graphql:"items(last: 100)"`
}

func (f *ProjectV2) SelectItem(contentID string) *Item {
	for _, node := range f.Items.Nodes {
		node := node
		if node.Content.Issue.ID == contentID || node.Content.PullRequest.ID == contentID {
			return &node
		}
	}
	return nil
}

type GetProjectQuery struct {
	ProjectV2 `graphql:"projectV2(number: $projectNumber)"`
}

type GetUserProjectQuery struct {
	User GetProjectQuery `graphql:"user(login: $owner)"`
}

// TODO: org support
func GetProject(
	gql api.GQLClient,
	descriptor ProjectDescriptor,
	iterationFieldName string,
) (*ProjectV2, error) {
	var query GetUserProjectQuery
	vars := map[string]interface{}{
		"owner":         graphql.String(descriptor.Owner),
		"projectNumber": graphql.Int(descriptor.Number),
		"fieldName":     graphql.String(iterationFieldName),
	}
	if err := gql.Query("GetProject", &query, vars); err != nil {
		return nil, err
	}
	return &query.User.ProjectV2, nil
}
