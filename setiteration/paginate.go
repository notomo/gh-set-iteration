package setiteration

import (
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type PageInfo struct {
	EndCursor   string `json:"end_cursor"`
	HasNextPage bool   `json:"has_next_page"`
}

const limitPerRequest = 100

func Paginate(
	gql api.GQLClient,
	name string,
	query interface{},
	variables map[string]interface{},
	each func() (PageInfo, int),
	limit int,
) error {
	var cursor *graphql.String
	currentCount := 0
	for limit-currentCount > 0 {
		vars := map[string]interface{}{
			"limit": graphql.Int(limitPerRequest),
			"after": cursor,
		}
		for k, v := range variables {
			vars[k] = v
		}
		if err := gql.Query(name, query, vars); err != nil {
			return err
		}

		pageInfo, count := each()
		if !pageInfo.HasNextPage {
			break
		}
		endCursor := graphql.NewString(graphql.String(pageInfo.EndCursor))
		cursor = endCursor
		currentCount = count
	}
	return nil
}
