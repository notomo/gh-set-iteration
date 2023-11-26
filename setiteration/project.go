package setiteration

import (
	"fmt"
	"time"

	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type Iteration struct {
	ID        string
	Title     string
	StartDate string
	Duration  int
}

func (iteration *Iteration) Contains(date string) (bool, error) {
	iterationStart, err := time.Parse(time.DateOnly, iteration.StartDate)
	if err != nil {
		return false, err
	}

	iterationEnd := iterationStart.AddDate(0, 0, iteration.Duration-1)

	target, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return false, err
	}

	return iterationStart.Add(-time.Nanosecond).Before(target) && iterationEnd.Add(time.Nanosecond).After(target), nil
}

type IterationMatchType string

var (
	IterationMatchTypeStartDateExactly = IterationMatchType("startDateExactly")
	IterationMatchTypeContains         = IterationMatchType("contains")
)

func matchIteration(
	iteration Iteration,
	targetDate string,
	matchType IterationMatchType,
) (bool, error) {
	switch matchType {
	case IterationMatchTypeStartDateExactly:
		return iteration.StartDate == targetDate, nil
	case IterationMatchTypeContains:
		return iteration.Contains(targetDate)
	}
	return false, fmt.Errorf("unexpected iteration match type: %s", matchType)
}

type ProjectV2IterationField struct {
	ID            string
	Name          string
	Configuration struct {
		Iterations          []Iteration
		CompletedIterations []Iteration
	}
}

func (f *ProjectV2IterationField) SelectIteration(
	targetDate string,
	matchType IterationMatchType,
) (*Iteration, error) {
	iterations := []Iteration{}
	iterations = append(iterations, f.Configuration.Iterations...)
	iterations = append(iterations, f.Configuration.CompletedIterations...)
	for _, iteration := range iterations {
		iteration := iteration

		matched, err := matchIteration(iteration, targetDate, matchType)
		if err != nil {
			return nil, err
		}
		if matched {
			return &iteration, nil
		}
	}
	return nil, nil
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
