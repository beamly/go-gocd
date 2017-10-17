package cli

import (
	"context"
	"encoding/json"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io/ioutil"
)

// List of command name and descriptions
const (
	CreatePipelineConfigCommandName  = "create-pipeline-config"
	CreatePipelineConfigCommandUsage = "Create Pipeline config"
	UpdatePipelineConfigCommandName  = "update-pipeline-config"
	UpdatePipelineConfigCommandUsage = "Update Pipeline config"
	DeletePipelineConfigCommandName  = "delete-pipeline-config"
	DeletePipelineConfigCommandUsage = "Remove Pipeline. This will not delete the pipeline history, which will be stored in the database"
	GetPipelineConfigCommandName     = "get-pipeline-config"
	GetPipelineConfigCommandUsage    = "Get a Pipeline Configuration"
)

// CreatePipelineConfigAction handles the interaction between the cli flags and the action handler for
// create-pipeline-config-action
func createPipelineConfigAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	group := c.String("group")
	if group == "" {
		return nil, nil, NewFlagError("group")
	}

	pipeline := c.String("pipeline-json")
	pipelineFile := c.String("pipeline-file")
	if pipeline == "" && pipelineFile == "" {
		return nil, nil, errors.New("One of '--pipeline-file' or '--pipeline-json' must be specified")
	}

	if pipeline != "" && pipelineFile != "" {
		return nil, nil, errors.New("Only one of '--pipeline-file' or '--pipeline-json' can be specified")
	}

	var pf []byte
	if pipelineFile != "" {
		pf, err = ioutil.ReadFile(pipelineFile)
		if err != nil {
			return nil, nil, err
		}
	} else {
		pf = []byte(pipeline)
	}
	p := &gocd.Pipeline{}
	err = json.Unmarshal(pf, &p)
	if err != nil {
		return nil, nil, err
	}

	return client.PipelineConfigs.Create(context.Background(), group, p)
}

// UpdatePipelineConfigAction handles the interaction between the cli flags and the action handler for
// update-pipeline-config-action
func updatePipelineConfigAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name, version string

	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}

	if version = c.String("pipeline-version"); version == "" {
		return nil, nil, NewFlagError("pipeline-version")
	}

	pipeline := c.String("pipeline")
	pipelineFile := c.String("pipeline-file")
	if pipeline == "" && pipelineFile == "" {
		return nil, nil, errors.New("One of '--pipeline-file' or '--pipeline' must be specified")
	}

	if pipeline != "" && pipelineFile != "" {
		return nil, nil, errors.New("Only one of '--pipeline-file' or '--pipeline' can be specified")
	}

	var pf []byte
	if pipelineFile != "" {
		pf, err = ioutil.ReadFile(pipelineFile)
		if err != nil {
			return nil, nil, err
		}
	} else {
		pf = []byte(pipeline)
	}
	p := &gocd.Pipeline{
		Version: version,
	}
	err = json.Unmarshal(pf, &p)
	if err != nil {
		return nil, nil, err
	}

	return client.PipelineConfigs.Update(context.Background(), name, p)
}

// DeletePipelineConfigAction handles the interaction between the cli flags and the action handler for
// delete-pipeline-config-action
func deletePipelineConfigAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	name := c.String("name")
	if name == "" {
		return nil, nil, NewFlagError("name")
	}

	deleteResponse, resp, err := client.PipelineConfigs.Delete(context.Background(), name)
	if resp.HTTP.StatusCode == 406 {
		err = errors.New(deleteResponse)
	}
	return deleteResponse, resp, err
}

// GetPipelineConfigAction handles the interaction between the cli flags and the action handler for get-pipeline-config
func getPipelineConfigAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	name := c.String("name")
	if name == "" {
		return nil, nil, NewFlagError("name")
	}

	getResponse, resp, err := client.PipelineConfigs.Get(context.Background(), name)
	if resp.HTTP.StatusCode != 404 {
		getResponse.RemoveLinks()
	}
	return getResponse, resp, err
}

// CreatePipelineConfigCommand handles the interaction between the cli flags and the action handler for create-pipeline-config
func createPipelineConfigCommand() *cli.Command {
	return &cli.Command{
		Name:     CreatePipelineConfigCommandName,
		Usage:    CreatePipelineConfigCommandUsage,
		Action:   ActionWrapper(createPipelineConfigAction),
		Category: "Pipeline Configs",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "group"},
			cli.StringFlag{Name: "pipeline-json", Usage: "A JSON string describing the pipeline configuration"},
			cli.StringFlag{Name: "pipeline-file", Usage: "Path to a JSON file describing the pipeline configuration"},
		},
	}
}

// UpdatePipelineConfigCommand handles the interaction between the cli flags and the action handler for update-pipeline-config
func updatePipelineConfigCommand() *cli.Command {
	return &cli.Command{
		Name:     UpdatePipelineConfigCommandName,
		Usage:    UpdatePipelineConfigCommandUsage,
		Action:   ActionWrapper(updatePipelineConfigAction),
		Category: "Pipeline Configs",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
			cli.StringFlag{Name: "pipeline-version"},
			cli.StringFlag{Name: "pipeline"},
			cli.StringFlag{Name: "pipeline-file"},
		},
	}
}

// DeletePipelineConfigCommand handles the interaction between the cli flags and the action handler for delete-pipeline-config
func deletePipelineConfigCommand() *cli.Command {
	return &cli.Command{
		Name:     DeletePipelineConfigCommandName,
		Usage:    DeletePipelineConfigCommandUsage,
		Category: "Pipeline Configs",
		Action:   ActionWrapper(deletePipelineConfigAction),
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
	}
}

// GetPipelineConfigCommand handles the interaction between the cli flags and the action handler for get-pipeline-config
func getPipelineConfigCommand() *cli.Command {
	return &cli.Command{
		Name:     GetPipelineConfigCommandName,
		Usage:    GetPipelineConfigCommandUsage,
		Action:   ActionWrapper(getPipelineConfigAction),
		Category: "Pipeline Configs",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
	}
}
