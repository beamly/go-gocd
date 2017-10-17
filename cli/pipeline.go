package cli

import (
	"context"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	GetPipelineHistoryCommandName   = "get-pipeline-history"
	GetPipelineCommandName          = "get-pipeline"
	GetPipelineHistoryCommandUsage  = "Get Pipeline History"
	GetPipelineCommandUsage         = "Get Pipeline"
	ReleasePipelineLockCommandName  = "release-pipeline-lock"
	ReleasePipelineLockCommandUsage = "Release Pipeline Lock"
	UnpausePipelineCommandName      = "unpause-pipeline"
	UnpausePipelineCommandUsage     = "Unpause Pipeline"
	PausePipelineCommandName        = "pause-pipeline"
	PausePipelineCommandUsage       = "Pause Pipeline"
	GetPipelineStatusCommandName    = "get-pipeline-status"
	GetPipelineStatusCommandUsage   = "Get Pipeline Status"
)

// GetPipelineStatusAction handles the business logic between the command objects and the go-gocd library.
func getPipelineStatusAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name string
	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}
	return client.Pipelines.GetStatus(context.Background(), name, -1)
}

// GetPipelineAction handles the business logic between the command objects and the go-gocd library.
func getPipelineAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name string
	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}
	return client.PipelineConfigs.Get(context.Background(), c.String("name"))
}

// GetPipelineHistoryAction handles the interaction between the cli flags and the action handler for
// get-pipeline-history-action
func getPipelineHistoryAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name string
	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}

	return client.Pipelines.GetHistory(context.Background(), c.String("name"), -1)
}

// PausePipelineAction handles the business logic between the command objects and the go-gocd library.
func pausePipelineAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name string
	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}

	return client.Pipelines.Pause(context.Background(), c.String("name"))
}

// UnpausePipelineAction handles the business logic between the command objects and the go-gocd library.
func unpausePipelineAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name string
	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}

	return client.Pipelines.Unpause(context.Background(), c.String("name"))
}

// ReleasePipelineLockAction handles the business logic between the command objects and the go-gocd library.
func releasePipelineLockAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name string
	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}

	return client.Pipelines.ReleaseLock(context.Background(), c.String("name"))
}

// GetPipelineStatusCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func getPipelineStatusCommand() *cli.Command {
	return &cli.Command{
		Name:     GetPipelineStatusCommandName,
		Usage:    GetPipelineStatusCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: ActionWrapper(getPipelineStatusAction),
	}
}

// PausePipelineCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func pausePipelineCommand() *cli.Command {
	return &cli.Command{
		Name:     PausePipelineCommandName,
		Usage:    PausePipelineCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: ActionWrapper(pausePipelineAction),
	}
}

// UnpausePipelineCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func unpausePipelineCommand() *cli.Command {
	return &cli.Command{
		Name:     UnpausePipelineCommandName,
		Usage:    UnpausePipelineCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: ActionWrapper(unpausePipelineAction),
	}
}

// ReleasePipelineLockCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func releasePipelineLockCommand() *cli.Command {
	return &cli.Command{
		Name:     ReleasePipelineLockCommandName,
		Usage:    ReleasePipelineLockCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: ActionWrapper(releasePipelineLockAction),
	}
}

// GetPipelineCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func getPipelineCommand() *cli.Command {
	return &cli.Command{
		Name:     GetPipelineCommandName,
		Usage:    GetPipelineCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: ActionWrapper(getPipelineAction),
	}
}

// GetPipelineHistoryCommand handles the interaction between the cli flags and the action handler for
// get-pipeline-history-action
func getPipelineHistoryCommand() *cli.Command {
	return &cli.Command{
		Name:     GetPipelineHistoryCommandName,
		Usage:    GetPipelineHistoryCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: ActionWrapper(getPipelineHistoryAction),
	}
}
