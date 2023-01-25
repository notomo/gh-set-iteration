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

type ProjectV2IterationField struct {
	ID            string
	Name          string
	Configuration struct {
		Iterations          []Iteration
		CompletedIterations []Iteration
	}
}

func (f *ProjectV2IterationField) SelectIteration(startDate string) *Iteration {
	iterations := []Iteration{}
	iterations = append(iterations, f.Configuration.Iterations...)
	iterations = append(iterations, f.Configuration.CompletedIterations...)
	for _, iteration := range iterations {
		iteration := iteration
		if iteration.StartDate == startDate {
			return &iteration
		}
	}
	return nil
}

type ProjectItem struct {
	ID      string
	Content IssueOrPullRequest
}

type ProjectV2 struct {
	ID    string
	Field struct {
		ProjectV2IterationField `graphql:"... on ProjectV2IterationField"`
	} `graphql:"field(name: $fieldName)"`
	Items struct {
		// TODO: pagenation ?
		Nodes []ProjectItem
	} `graphql:"items(last: 100)"`
}

func (f *ProjectV2) SelectItem(contentID string) *ProjectItem {
	for _, node := range f.Items.Nodes {
		node := node
		if node.Content.Content.ID == contentID {
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

type GetOrganizationProjectQuery struct {
	Organization GetProjectQuery `graphql:"organization(login: $owner)"`
}

func GetProject(
	gql api.GQLClient,
	descriptor ProjectDescriptor,
	iterationFieldName string,
) (*ProjectV2, error) {
	vars := map[string]interface{}{
		"owner":         graphql.String(descriptor.Owner),
		"projectNumber": graphql.Int(descriptor.Number),
		"fieldName":     graphql.String(iterationFieldName),
	}

	if descriptor.OwnerIsOrganization {
		var query GetOrganizationProjectQuery
		if err := gql.Query("GetProject", &query, vars); err != nil {
			return nil, err
		}
		return &query.Organization.ProjectV2, nil
	}

	var query GetUserProjectQuery
	if err := gql.Query("GetProject", &query, vars); err != nil {
		return nil, err
	}
	return &query.User.ProjectV2, nil
}
