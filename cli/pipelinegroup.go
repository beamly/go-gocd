package cli

import (
	"context"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	ListPipelineGroupsCommandName  = "list-pipeline-groups"
	ListPipelineGroupsCommandUsage = "List Pipeline Groups"
)

// ListPipelineGroupsAction handles the interaction between the cli flags and the action handler for
// list-pipeline-groups
func listPipelineGroupsAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	return client.PipelineGroups.List(context.Background(), c.String("group-name"))
}

// ListPipelineGroupsCommand handles the interaction between the cli flags and the action handler for
// list-pipeline-groups
func listPipelineGroupsCommand() *cli.Command {
	return &cli.Command{
		Name:     ListPipelineGroupsCommandName,
		Usage:    ListPipelineGroupsCommandUsage,
		Action:   ActionWrapper(listPipelineGroupsAction),
		Category: "Pipeline Groups",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "group-name"},
		},
	}
}
