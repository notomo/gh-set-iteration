package setiteration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {

	//     November 2023
	// Su Mo Tu We Th Fr Sa
	//
	// 12 13 14 15 16 17 18
	// 19 20 21 22 23 24 25
	// 26 27 28 29 30

	t.Run("valid", func(t *testing.T) {
		cases := []struct {
			name      string
			iteration *Iteration
			date      string
			want      bool
		}{
			{
				name: "equals with start date",
				iteration: &Iteration{
					StartDate: "2023-11-20",
					Duration:  7,
				},
				date: "2023-11-20",
				want: true,
			},
			{
				name: "contains",
				iteration: &Iteration{
					StartDate: "2023-11-20",
					Duration:  7,
				},
				date: "2023-11-23",
				want: true,
			},
			{
				name: "equals with end date",
				iteration: &Iteration{
					StartDate: "2023-11-20",
					Duration:  7,
				},
				date: "2023-11-26",
				want: true,
			},
			{
				name: "after end date",
				iteration: &Iteration{
					StartDate: "2023-11-20",
					Duration:  7,
				},
				date: "2023-11-27",
				want: false,
			},
			{
				name: "before start date",
				iteration: &Iteration{
					StartDate: "2023-11-20",
					Duration:  7,
				},
				date: "2023-11-19",
				want: false,
			},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				got, err := c.iteration.Contains(c.date)
				require.NoError(t, err)
				assert.Equal(t, c.want, got)
			})
		}
	})
}

func TestSelectIteration(t *testing.T) {
	t.Run("nearest match type", func(t *testing.T) {
		field := &ProjectV2IterationField{
			Configuration: struct {
				Iterations          []Iteration
				CompletedIterations []Iteration
			}{
				Iterations: []Iteration{
					{ID: "1", StartDate: "2023-11-01", Duration: 7},
					{ID: "2", StartDate: "2023-11-15", Duration: 7},
					{ID: "3", StartDate: "2023-12-01", Duration: 7},
				},
			},
		}

		cases := []struct {
			name       string
			targetDate string
			expectedID string
		}{
			{
				name:       "before all iterations",
				targetDate: "2023-10-20",
				expectedID: "1",
			},
			{
				name:       "closest to first iteration",
				targetDate: "2023-11-03",
				expectedID: "1",
			},
			{
				name:       "closest to second iteration",
				targetDate: "2023-11-18",
				expectedID: "2",
			},
			{
				name:       "after all iterations",
				targetDate: "2023-12-15",
				expectedID: "3",
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				result, err := field.SelectIteration(c.targetDate, IterationMatchTypeNearest)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, c.expectedID, result.ID)
			})
		}
	})
}
