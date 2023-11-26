package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/notomo/gh-set-iteration/setiteration"

	"github.com/urfave/cli/v2"
)

const (
	paramProjectUrl         = "project-url"
	paramContentUrl         = "content-url"
	paramIterationField     = "field"
	paramLog                = "log"
	paramDryRun             = "dry-run"
	paramOffsetDays         = "offset-days"
	paramIterationMatchType = "match"
)

func main() {
	app := &cli.App{
		Name: "gh-set-iteration",
		Action: func(c *cli.Context) error {
			opts := &api.ClientOptions{}
			logFilePath := c.String(paramLog)
			if logFilePath != "" {
				f, err := os.Create(logFilePath)
				if err != nil {
					return fmt.Errorf("create log file: %w", err)
				}
				defer f.Close()
				opts.Log = f
				opts.LogVerboseHTTP = true
			}
			gql, err := gh.GQLClient(opts)
			if err != nil {
				return fmt.Errorf("create gql client: %w", err)
			}
			return setiteration.Run(
				gql,
				c.String(paramProjectUrl),
				c.String(paramContentUrl),
				c.String(paramIterationField),
				c.Int(paramOffsetDays),
				setiteration.IterationMatchType(c.String(paramIterationMatchType)),
				c.Bool(paramDryRun),
				os.Stdout,
			)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     paramProjectUrl,
				Value:    "",
				Required: true,
				Usage:    "project url",
			},
			&cli.StringFlag{
				Name:     paramContentUrl,
				Value:    "",
				Required: true,
				Usage:    "issue or pull request url",
			},
			&cli.StringFlag{
				Name:     paramIterationField,
				Value:    "",
				Required: true,
				Usage:    "iteration field name",
			},
			&cli.StringFlag{
				Name:  paramLog,
				Value: "",
				Usage: "log file path",
			},
			&cli.IntFlag{
				Name:  paramOffsetDays,
				Value: 0,
				Usage: "offset days to adjust iteration's start date",
			},
			&cli.StringFlag{
				Name:  paramIterationMatchType,
				Value: string(setiteration.IterationMatchTypeStartDateExactly),
				Usage: `
This changes iteration select behavior.
Iteration match type is the following:
- startDateExactly: match with iteration start_date (default)
- contains: match if date is contains iteration date range
				`,
			},
			&cli.BoolFlag{
				Name:  paramDryRun,
				Value: false,
				Usage: "nothing is updated",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
