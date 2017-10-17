package cli

import (
	"context"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	ListScheduledJobsCommandName  = "list-scheduled-jobs"
	ListScheduledJobsCommandUsage = "List Scheduled Jobs"
)

// ListScheduledJobsAction gets a list of agents and return them.
func listScheduledJobsAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	return client.Jobs.ListScheduled(context.Background())
}

// ListScheduledJobsCommand provides interface between handler and action
func listScheduledJobsCommand() *cli.Command {
	return &cli.Command{
		Name:     ListScheduledJobsCommandName,
		Usage:    ListScheduledJobsCommandUsage,
		Action:   ActionWrapper(listScheduledJobsAction),
		Category: "Jobs",
	}
}
