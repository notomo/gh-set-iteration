package setiteration

import (
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type UpdateIterationMutation struct {
	UpdateProjectV2ItemFieldValue struct {
		ClientMutationId string
	} `graphql:"updateProjectV2ItemFieldValue(input: { projectId: $projectId itemId: $itemId fieldId: $fieldId value: { iterationId: $iterationId } })"`
}

func UpdateIteration(
	gql api.GQLClient,
	projectId string,
	itemId string,
	fieldId string,
	iterationId string,
	dryRun bool,
) error {
	if dryRun {
		return nil
	}
	var mutation UpdateIterationMutation
	vars := map[string]interface{}{
		"projectId":   graphql.ID(projectId),
		"itemId":      graphql.ID(itemId),
		"fieldId":     graphql.ID(fieldId),
		"iterationId": graphql.String(iterationId),
	}
	if err := gql.Mutate("UpdateIteration", &mutation, vars); err != nil {
		return err
	}
	return nil
}
