package setiteration

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

type Repository struct {
	Owner string
	Name  string
}

type IssueOrPullRequestDescriptor struct {
	Repository Repository
	Number     int
}

var issueOrPullRequestUrlRegex = regexp.MustCompile(`^/([^/]+)/([^/]+)/(pull|issues)/(\d+)`)

func GetIssueOrPullRequestDescriptor(
	issueOrPullRequestUrl string,
) (*IssueOrPullRequestDescriptor, error) {
	u, err := url.Parse(issueOrPullRequestUrl)
	if err != nil {
		return nil, err
	}

	matches := issueOrPullRequestUrlRegex.FindStringSubmatch(u.Path)
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid issue or pull request url path: %s", u.Path)
	}

	number, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, err
	}

	return &IssueOrPullRequestDescriptor{
		Repository: Repository{
			Owner: matches[1],
			Name:  matches[2],
		},
		Number: number,
	}, nil
}
