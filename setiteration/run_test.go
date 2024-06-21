package setiteration

import (
	"bytes"
	"testing"

	"github.com/notomo/gh-set-iteration/setiteration/gqltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	gql, err := gqltest.New(
		t,

		gqltest.WithQueryOK("GetProject", `
{
  "data": {
    "user": {
      "projectV2": {
        "id": "PVT_11111111111111",
        "field": {
          "id": "PVTIF_22222222222222222222",
          "name": "Iteration",
          "configuration": {
            "completedIterations": [
              {
                "id": "00000a0b",
                "title": "Iteration 0",
                "startDate": "2022-12-25",
                "duration": 7
              }
            ],
            "iterations": [
              {
                "id": "11111a1b",
                "title": "Iteration 1",
                "startDate": "2023-01-01",
                "duration": 7
              },
              {
                "id": "22222a2b",
                "title": "Iteration 2",
                "startDate": "2023-01-08",
                "duration": 7
              }
            ]
          }
        },
        "items": {
          "nodes": [
            {
              "id": "PVTI_11111111111111111111",
              "content": {
                "id": "I_1111111111111111",
                "title": "issue1: 2022-01-01"
              }
            },
            {
              "id": "PVTI_22222222222222222222",
              "content": {
                "id": "I_2222222222222222",
                "title": "issue2: 2023-01-08"
              }
            }
          ]
        }
      }
    }
  }
}
`),

		gqltest.WithQueryOK("GetIssueOrPullRequest", `
{
  "data": {
    "repository": {
      "issueOrPullRequest": {
        "id": "I_2222222222222222",
        "title": "issue2: 2023-01-08"
      }
    }
  }
}
`),

		gqltest.WithMutateOK("UpdateIteration", `
{
  "data": {
    "updateProjectV2ItemFieldValue": {
      "clientMutationId": null
    }
  }
}
`),
	)
	require.NoError(t, err)

	output := &bytes.Buffer{}
	assert.NoError(t, Run(
		gql,
		"https://github.com/users/notomo/projects/1",
		"https://github.com/notomo/example/issues/2",
		"Iteration",
		ContentStateAll,
		-14,
		IterationMatchTypeStartDateExactly,
		false,
		300,
		output,
	))

	want := `
Item is updated:
- iteration title: Iteration 0
- iteration start date: 2022-12-25
`
	got := output.String()
	assert.Equal(t, want, got)
}
