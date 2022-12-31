package setiteration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIssueOrPullRequestDescriptor(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cases := []struct {
			name                  string
			issueOrPullRequestUrl string
			want                  IssueOrPullRequestDescriptor
		}{
			{
				name:                  "issue url",
				issueOrPullRequestUrl: "https://github.com/notomo/example/issues/1",
				want: IssueOrPullRequestDescriptor{
					Repository: Repository{
						Owner: "notomo",
						Name:  "example",
					},
					Number: 1,
				},
			},
			{
				name:                  "pull request url",
				issueOrPullRequestUrl: "https://github.com/notomo/example/pull/2",
				want: IssueOrPullRequestDescriptor{
					Repository: Repository{
						Owner: "notomo",
						Name:  "example",
					},
					Number: 2,
				},
			},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				got, err := GetIssueOrPullRequestDescriptor(c.issueOrPullRequestUrl)
				require.NoError(t, err)
				assert.Equal(t, c.want, *got)
			})
		}
	})
}
