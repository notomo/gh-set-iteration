package setiteration

import (
	"fmt"
	"io"

	"github.com/cli/go-gh/pkg/api"
)

func Run(
	gql api.GQLClient,
	projectUrl string,
	issueOrPullRequestUrl string,
	iterationFieldName string,
	dryRun bool,
	writer io.Writer,
) error {
	projectDescriptor, err := GetProjectDescriptor(projectUrl)
	if err != nil {
		return err
	}
	project, err := GetProject(gql, *projectDescriptor, iterationFieldName)
	if err != nil {
		return err
	}

	descriptor, err := GetIssueOrPullRequestDescriptor(issueOrPullRequestUrl)
	if err != nil {
		return err
	}
	content, err := GetIssueOrPullRequest(gql, *descriptor)
	if err != nil {
		return err
	}

	startDate, err := ExtractDate(content.Title)
	if err != nil {
		return err
	}
	iteration := project.Field.SelectIteration(startDate)
	if iteration == nil {
		return fmt.Errorf("no matched iteration: startDate=%s", startDate)
	}

	item := project.SelectItem(content.ID)
	if item == nil {
		return fmt.Errorf("no matched item")
	}

	if err := UpdateIteration(
		gql,
		project.ID,
		item.ID,
		project.Field.ID,
		iteration.ID,
		dryRun,
	); err != nil {
		return err
	}

	message := fmt.Sprintf(`
Item is updated:
- iteration title: %s
- iteration start date: %s
`, iteration.Title, startDate)
	if _, err := writer.Write([]byte(message)); err != nil {
		return err
	}

	return nil
}
