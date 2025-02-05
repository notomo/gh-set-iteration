package setiteration

import (
	"fmt"
	"io"

	"github.com/cli/go-gh/v2/pkg/api"
)

type ContentState string

var (
	ContentStateAll    = ContentState("all")
	ContentStateOpen   = ContentState("open")
	ContentStateClosed = ContentState("closed")
)

var ErrSkipped = fmt.Errorf("skipped")

func Run(
	gql *api.GraphQLClient,
	projectUrl string,
	issueOrPullRequestUrl string,
	iterationFieldName string,
	contentState ContentState,
	offsetDays int,
	iterationMatchType IterationMatchType,
	dryRun bool,
	itemLimit int,
	writer io.Writer,
) error {
	content, err := getContent(gql, issueOrPullRequestUrl, contentState)
	if err != nil {
		return err
	}

	extractedDate, err := ExtractDate(content.Title)
	if err != nil {
		return err
	}

	targetDate, err := ShiftDate(extractedDate, offsetDays)
	if err != nil {
		return err
	}

	projectDescriptor, err := GetProjectDescriptor(projectUrl)
	if err != nil {
		return err
	}

	project, projectItem, err := GetProject(gql, *projectDescriptor, iterationFieldName, itemLimit, content.ID)
	if err != nil {
		return err
	}
	if projectItem == nil {
		return fmt.Errorf("no matched project item")
	}

	iteration, err := project.IterationField.SelectIteration(targetDate, iterationMatchType)
	if err != nil {
		return err
	}
	if iteration == nil {
		return fmt.Errorf("no matched iteration: targetDate=%s", targetDate)
	}

	if err := UpdateIteration(
		gql,
		project.ID,
		projectItem.ID,
		project.IterationField.ID,
		iteration.ID,
		dryRun,
	); err != nil {
		return err
	}

	message := fmt.Sprintf(`
Item is updated:
- iteration title: %s
- iteration start date: %s
`, iteration.Title, iteration.StartDate)
	if _, err := writer.Write([]byte(message)); err != nil {
		return err
	}

	return nil
}

func getContent(
	gql *api.GraphQLClient,
	issueOrPullRequestUrl string,
	contentState ContentState,
) (*Content, error) {
	descriptor, err := GetIssueOrPullRequestDescriptor(issueOrPullRequestUrl)
	if err != nil {
		return nil, err
	}

	content, err := GetIssueOrPullRequest(gql, *descriptor)
	if err != nil {
		return nil, err
	}

	if contentState == ContentStateOpen && content.Closed {
		return nil, fmt.Errorf("content is closed (state filter=open): %w", ErrSkipped)
	}

	if contentState == ContentStateClosed && !content.Closed {
		return nil, fmt.Errorf("content is open (state filter=closed): %w", ErrSkipped)
	}

	return content, nil
}
