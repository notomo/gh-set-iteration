package setiteration

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
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
	IterationMatchTypeNearest          = IterationMatchType("nearest")
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
	if matchType == IterationMatchTypeNearest {
		return f.selectNearestIteration(targetDate)
	}

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

func (f *ProjectV2IterationField) selectNearestIteration(targetDate string) (*Iteration, error) {
	if len(f.Configuration.Iterations) == 0 {
		return nil, nil
	}

	target, err := time.Parse(time.DateOnly, targetDate)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(f.Configuration.Iterations, func(a, b Iteration) int {
		if n := cmp.Compare(a.StartDate, b.StartDate); n != 0 {
			return -n
		}
		return cmp.Compare(a.ID, b.ID)
	})

	for _, iteration := range f.Configuration.Iterations {
		iteration := iteration

		s, err := time.Parse(time.DateOnly, iteration.StartDate)
		if err != nil {
			return nil, err
		}

		diff := target.Sub(s)
		if diff >= 0 {
			return &iteration, nil
		}
	}

	last := f.Configuration.Iterations[len(f.Configuration.Iterations)-1]
	return &last, nil
}

type ProjectItem struct {
	ID      string
	Content IssueOrPullRequest
}

type Project struct {
	ID             string
	IterationField ProjectV2IterationField
}

type ProjectV2Query struct {
	ID    string
	Field struct {
		ProjectV2IterationField `graphql:"... on ProjectV2IterationField"`
	} `graphql:"field(name: $fieldName)"`
	Items struct {
		Nodes    []ProjectItem
		PageInfo PageInfo
	} `graphql:"items(first: $limit, after: $after, orderBy: {field:POSITION, direction:DESC})"`
}

type GetProjectQuery struct {
	ProjectV2Query `graphql:"projectV2(number: $projectNumber)"`
}

type GetUserProjectQuery struct {
	User GetProjectQuery `graphql:"user(login: $owner)"`
}

type GetOrganizationProjectQuery struct {
	Organization GetProjectQuery `graphql:"organization(login: $owner)"`
}

func GetProject(
	gql *api.GraphQLClient,
	descriptor ProjectDescriptor,
	iterationFieldName string,
	itemLimit int,
	contentID string,
) (*Project, *ProjectItem, error) {
	vars := map[string]interface{}{
		"owner":         graphql.String(descriptor.Owner),
		"projectNumber": graphql.Int(descriptor.Number),
		"fieldName":     graphql.String(iterationFieldName),
	}

	var projectItem *ProjectItem
	itemCount := 0
	each := func(nodes []ProjectItem, pageInfo PageInfo) (PageInfo, int) {
		for _, item := range nodes {
			if item.Content.Content.ID == contentID {
				projectItem = &item
				return PageInfo{HasNextPage: false}, 0
			}
		}
		itemCount += len(nodes)
		return pageInfo, itemCount
	}

	if descriptor.OwnerIsOrganization {
		var query GetOrganizationProjectQuery
		if err := Paginate(gql, "GetProject", &query, vars, func() (PageInfo, int) {
			return each(query.Organization.Items.Nodes, query.Organization.Items.PageInfo)
		}, itemLimit); err != nil {
			return nil, nil, err
		}

		q := query.Organization.ProjectV2Query
		return &Project{
			ID:             q.ID,
			IterationField: q.Field.ProjectV2IterationField,
		}, projectItem, nil
	}

	var query GetUserProjectQuery
	if err := Paginate(gql, "GetProject", &query, vars, func() (PageInfo, int) {
		return each(query.User.Items.Nodes, query.User.Items.PageInfo)
	}, itemLimit); err != nil {
		return nil, nil, err
	}

	q := query.User.ProjectV2Query
	return &Project{
		ID:             q.ID,
		IterationField: q.Field.ProjectV2IterationField,
	}, projectItem, nil
}
