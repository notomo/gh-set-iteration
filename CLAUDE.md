# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a GitHub CLI extension called `gh-set-iteration` that automatically sets iteration fields in GitHub Projects based on dates extracted from issue/PR titles. The tool parses dates from titles and assigns the appropriate project iteration.

## Key Commands

### Build and Install
```bash
make build          # Builds the binary as gh-set-iteration
make install        # Builds and installs the extension to gh CLI
```

### Testing
```bash
make test           # Runs all tests with verbose output
go test -v ./...    # Direct test command
```

### Development
```bash
make help           # Shows CLI help
make start          # Installs and runs with example parameters
```

## Architecture

### Main Components

- **main.go**: CLI entry point using urfave/cli/v2, defines all command-line flags and parameters
- **setiteration/run.go**: Core orchestration logic that coordinates the entire workflow
- **setiteration/project.go**: GitHub Projects API interaction, handles project queries and iteration matching
- **setiteration/update_iteration.go**: GraphQL mutations for updating project items
- **setiteration/date.go**: Date extraction and parsing from issue/PR titles
- **setiteration/issue_or_pull_request.go**: GitHub API queries for issues and pull requests

### Key Workflow

1. Parse issue/PR URL and extract content
2. Extract date from issue/PR title using regex patterns
3. Apply date offset if specified
4. Query GitHub project to find matching iteration field
5. Match iteration based on date and match type (exact start date or date range)
6. Update project item with the selected iteration

### GraphQL Integration

The tool heavily uses GitHub's GraphQL API v4 through `github.com/cli/go-gh/v2` for:
- Querying project iterations and configurations
- Finding project items linked to issues/PRs
- Updating iteration field values via mutations

### Date Matching Logic

Two iteration match types are supported:
- `startDateExactly`: Matches iteration with exact start date
- `contains`: Matches if target date falls within iteration date range

## Testing

Tests are located in the `setiteration/` directory with `*_test.go` files. The project uses `github.com/stretchr/testify` for assertions and includes a `gqltest/` subdirectory for GraphQL client testing utilities.