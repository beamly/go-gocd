package cli

import (
	"context"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

const (
	listPropertiesCommandName  = "list-properties"
	listPropertiesCommandUsage = "List the properties for a given job."
	createPropertyCommandName  = "create-property"
	createPropertyCommandUsage = "Create a property for a given job."
	propertiesGroup            = "Properties"
)

func createPropertyAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name, value, pipeline, stage, job string
	var pipelineCounter, stageCounter int

	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}

	if value = c.String("value"); value == "" {
		return nil, nil, NewFlagError("value")
	}

	if pipeline = c.String("pipeline"); pipeline == "" {
		return nil, nil, NewFlagError("pipeline")
	}

	if stage = c.String("stage"); stage == "" {
		return nil, nil, NewFlagError("stage")
	}

	if job = c.String("job"); job == "" {
		return nil, nil, NewFlagError("job")
	}

	if pipelineCounter = c.Int("pipeline-counter"); pipelineCounter < 1 {
		return nil, nil, NewFlagError("pipeline-counter")
	}

	if stageCounter = c.Int("stage-counter"); stageCounter < 1 {
		return nil, nil, NewFlagError("stage-counter")
	}

	return client.Properties.Create(context.Background(), name, value, &gocd.PropertyRequest{
		Pipeline:        pipeline,
		PipelineCounter: pipelineCounter,
		Stage:           stage,
		StageCounter:    stageCounter,
		Job:             job,
	})
}

func listPropertiesAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var pipeline, stage, job string
	var pipelineCounter, stageCounter int

	if pipeline = c.String("pipeline"); pipeline == "" {
		return nil, nil, NewFlagError("pipeline")
	}

	if stage = c.String("stage"); stage == "" {
		return nil, nil, NewFlagError("stage")
	}

	if job = c.String("job"); job == "" {
		return nil, nil, NewFlagError("job")
	}

	if pipelineCounter = c.Int("pipeline-counter"); pipelineCounter < 1 {
		return nil, nil, NewFlagError("pipeline-counter")
	}

	if stageCounter = c.Int("stage-counter"); stageCounter < 1 {
		return nil, nil, NewFlagError("stage-counter")
	}

	return client.Properties.List(context.Background(), &gocd.PropertyRequest{
		Pipeline:        pipeline,
		PipelineCounter: pipelineCounter,
		Stage:           stage,
		StageCounter:    stageCounter,
		Job:             job,
		Single:          true,
	})
}

func listPropertiesCommand() *cli.Command {
	return &cli.Command{
		Name:     listPropertiesCommandName,
		Usage:    listPropertiesCommandUsage,
		Action:   ActionWrapper(listPropertiesAction),
		Category: propertiesGroup,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "pipeline", EnvVar: "GO_PIPELINE_NAME"},
			cli.IntFlag{Name: "pipeline-counter", EnvVar: "GO_PIPELINE_COUNTER"},
			cli.StringFlag{Name: "stage", EnvVar: "GO_STAGE_NAME"},
			cli.IntFlag{Name: "stage-counter", EnvVar: "GO_STAGE_COUNTER"},
			cli.StringFlag{Name: "job", EnvVar: "GO_JOB_NAME"},
		},
	}
}

func createPropertyCommand() *cli.Command {
	return &cli.Command{
		Name:     createPropertyCommandName,
		Usage:    createPropertyCommandUsage,
		Action:   ActionWrapper(createPropertyAction),
		Category: propertiesGroup,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
			cli.StringFlag{Name: "value"},
			cli.StringFlag{Name: "pipeline", EnvVar: "GO_PIPELINE_NAME"},
			cli.IntFlag{Name: "pipeline-counter", EnvVar: "GO_PIPELINE_COUNTER"},
			cli.StringFlag{Name: "stage", EnvVar: "GO_STAGE_NAME"},
			cli.IntFlag{Name: "stage-counter", EnvVar: "GO_STAGE_COUNTER"},
			cli.StringFlag{Name: "job", EnvVar: "GO_JOB_NAME"},
		},
	}
}
