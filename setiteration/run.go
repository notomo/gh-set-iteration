package setiteration

import (
	"fmt"
	"io"

	"github.com/cli/go-gh/pkg/api"
)

type ContentState string

var (
	ContentStateAll    = ContentState("all")
	ContentStateOpen   = ContentState("open")
	ContentStateClosed = ContentState("closed")
)

var ErrSkipped = fmt.Errorf("skipped")

func Run(
	gql api.GQLClient,
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
	projectDescriptor, err := GetProjectDescriptor(projectUrl)
	if err != nil {
		return err
	}
	project, err := GetProject(gql, *projectDescriptor, iterationFieldName, itemLimit)
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

	if contentState == ContentStateOpen && content.Closed {
		return fmt.Errorf("content is closed (state filter=open): %w", ErrSkipped)
	}
	if contentState == ContentStateClosed && !content.Closed {
		return fmt.Errorf("content is open (state filter=closed): %w", ErrSkipped)
	}

	extractedDate, err := ExtractDate(content.Title)
	if err != nil {
		return err
	}
	targetDate, err := ShiftDate(extractedDate, offsetDays)
	if err != nil {
		return err
	}

	iteration, err := project.Field.SelectIteration(targetDate, iterationMatchType)
	if err != nil {
		return err
	}
	if iteration == nil {
		return fmt.Errorf("no matched iteration: targetDate=%s", targetDate)
	}

	projectItem := project.SelectItem(content.ID)
	if projectItem == nil {
		return fmt.Errorf("no matched project item")
	}

	if err := UpdateIteration(
		gql,
		project.ID,
		projectItem.ID,
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
`, iteration.Title, iteration.StartDate)
	if _, err := writer.Write([]byte(message)); err != nil {
		return err
	}

	return nil
}
