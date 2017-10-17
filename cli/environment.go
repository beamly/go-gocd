package cli

import (
	"context"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
	"strings"
)

// List of command name and descriptions
const (
	ListEnvironmentsCommandName                = "list-environments"
	ListEnvironmentsCommandUsage               = "List all environments"
	GetEnvironmentCommandName                  = "get-environment"
	GetEnvironmentCommandUsage                 = "Get an environment by name"
	AddPipelinesToEnvironmentCommandName       = "add-pipelines-to-environment"
	AddPipelinesToEnvironmentCommandUsage      = "Add one or more pipelines to an environment"
	RemovePipelinesFromEnvironmentCommandName  = "remove-pipelines-from-environment"
	RemovePipelinesFromEnvironmentCommandUsage = "Remove one or more pipelines from an environment"
)

// ListEnvironmentsAction handles the listing of environments
func listEnvironmentsAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	es, resp, err := client.Environments.List(context.Background())
	if err == nil {
		es.RemoveLinks()
	}
	return es, resp, err
}

// GetEnvironmentAction handles the retrieval of environments
func getEnvironmentAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name string
	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}
	e, resp, err := client.Environments.Get(context.Background(), name)
	if err == nil {
		e.RemoveLinks()
	}
	return e, resp, err
}

// AddPipelinesToEnvironmentAction handles the adding of a pipeline to an environment
func addPipelinesToEnvironmentAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var environment, pipelines string

	if environment = c.String("environment-name"); environment == "" {
		return nil, nil, NewFlagError("environment-name")
	}
	if pipelines = c.String("pipeline-names"); pipelines == "" {
		return nil, nil, NewFlagError("pipeline-names")
	}

	e, resp, err := client.Environments.Patch(context.Background(), environment, &gocd.EnvironmentPatchRequest{
		Pipelines: &gocd.PatchStringAction{
			Add: strings.Split(pipelines, ","),
		},
	})
	if err == nil {
		e.RemoveLinks()
	}
	return e, resp, err

}

// RemovePipelinesFromEnvironmentAction handles the removing of a pipeline from an environment
func removePipelinesFromEnvironmentAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var environment, pipelines string

	if environment = c.String("environment-name"); environment == "" {
		return nil, nil, NewFlagError("environment-name")
	}
	if pipelines = c.String("pipeline-names"); pipelines == "" {
		return nil, nil, NewFlagError("pipeline-names")
	}

	e, resp, err := client.Environments.Patch(context.Background(), environment, &gocd.EnvironmentPatchRequest{
		Pipelines: &gocd.PatchStringAction{
			Remove: strings.Split(pipelines, ","),
		},
	})
	if err == nil {
		e.RemoveLinks()
	}
	return e, resp, err
}

// ListEnvironmentsCommand handles definition of cli command
func listEnvironmentsCommand() *cli.Command {
	return &cli.Command{
		Name:     ListEnvironmentsCommandName,
		Usage:    ListEnvironmentsCommandUsage,
		Action:   ActionWrapper(listEnvironmentsAction),
		Category: "Environments",
	}
}

// GetEnvironmentCommand handles definition of cli command
func getEnvironmentCommand() *cli.Command {
	return &cli.Command{
		Name:     GetEnvironmentCommandName,
		Usage:    GetEnvironmentCommandUsage,
		Action:   ActionWrapper(getEnvironmentAction),
		Category: "Environments",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
	}
}

// AddPipelinesToEnvironmentCommand handles definition of cli command
func addPipelinesToEnvironmentCommand() *cli.Command {
	return &cli.Command{
		Name:     AddPipelinesToEnvironmentCommandName,
		Usage:    AddPipelinesToEnvironmentCommandUsage,
		Action:   ActionWrapper(addPipelinesToEnvironmentAction),
		Category: "Environments",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "environment-name"},
			cli.StringFlag{
				Name:  "pipeline-names",
				Usage: "Comma seperated list of pipeline names to add.",
			},
		},
	}
}

// RemovePipelinesFromEnvironmentCommand handles definition of cli command
func removePipelinesFromEnvironmentCommand() *cli.Command {
	return &cli.Command{
		Name:     RemovePipelinesFromEnvironmentCommandName,
		Usage:    RemovePipelinesFromEnvironmentCommandUsage,
		Action:   ActionWrapper(removePipelinesFromEnvironmentAction),
		Category: "Environments",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "environment-name",
			},
			cli.StringFlag{
				Name:  "pipeline-names",
				Usage: "Comma seperated list of pipeline names to remove.",
			},
		},
	}
}
