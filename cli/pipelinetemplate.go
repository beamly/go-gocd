package cli

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	ListPipelineTemplatesCommandName   = "list-pipeline-templates"
	ListPipelineTemplatesCommandUsage  = "List Pipeline Templates"
	GetPipelineTemplateCommandName     = "get-pipeline-template"
	GetPipelineTemplateCommandUsage    = "Get Pipeline Templates"
	CreatePipelineTemplateCommandName  = "create-pipeline-template"
	CreatePipelineTemplateCommandUsage = "Create Pipeline Templates"
	UpdatePipelineTemplateCommandName  = "update-pipeline-template"
	UpdatePipelineTemplateCommandUsage = "Update Pipeline template"
	DeletePipelineTemplateCommandName  = "delete-pipeline-template"
	DeletePipelineTemplateCommandUsage = "Delete Pipeline template"
)

// ListPipelineTemplatesAction lists all pipeline templates.
func listPipelineTemplatesAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	ts, resp, err := client.PipelineTemplates.List(context.Background())
	if err != nil {
		return ts, resp, err
	}

	type p struct {
		Name string `json:"name"`
	}

	type ptr struct {
		Name      string `json:"name"`
		Pipelines []*p   `json:"pipelines"`
	}
	responses := []ptr{}
	for _, pt := range ts {
		pt.RemoveLinks()
		ps := []*p{}
		for _, pipe := range pt.Pipelines() {
			ps = append(ps, &p{pipe.Name})
		}

		responses = append(responses, ptr{
			Name:      pt.Name,
			Pipelines: ps,
		})
	}
	return responses, resp, err
}

// GetPipelineTemplateAction checks template-name is provided, and that the response is 2xx.
func getPipelineTemplateAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name string
	if name = c.String("template-name"); name == "" {
		return nil, nil, NewFlagError("template-name")
	}

	pt, resp, err := client.PipelineTemplates.Get(context.Background(), name)
	if resp.HTTP.StatusCode != 404 {
		pt.RemoveLinks()
	}
	return pt, resp, err
}

// CreatePipelineTemplateAction checks stages and template-name is provided, and that the response is 2xx.
func createPipelineTemplateAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	if c.String("template-name") == "" {
		return nil, nil, NewFlagError("template-name")
	}
	if len(c.StringSlice("stage")) < 1 {
		return nil, nil, errors.New("At least 1 '--stage' must be set")
	}

	stages := []*gocd.Stage{}
	for _, stage := range c.StringSlice("stage") {
		st := gocd.Stage{}
		json.Unmarshal([]byte(stage), &st)

		if err := st.Validate(); err != nil {
			return nil, nil, err
		}
		stages = append(stages, &st)
	}

	return client.PipelineTemplates.Create(context.Background(), c.String("template-name"), stages)
}

// UpdatePipelineTemplateAction checks stages, template-name and template-version is provided, and that the response is
// 2xx.
func updatePipelineTemplateAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	if c.String("template-name") == "" {
		return nil, nil, NewFlagError("template-name")
	}
	if c.String("template-version") == "" {
		return nil, nil, NewFlagError("version")
	}
	if len(c.StringSlice("stage")) < 1 {
		return nil, nil, errors.New("At least 1 '--stage' must be set")
	}

	stages := []*gocd.Stage{}
	for _, stage := range c.StringSlice("stage") {
		st := gocd.Stage{}
		json.Unmarshal([]byte(stage), &st)

		if err := st.Validate(); err != nil {
			return nil, nil, err
		}
		stages = append(stages, &st)
	}

	ptr := gocd.PipelineTemplate{
		Version: c.String("template-version"),
		Stages:  stages,
	}

	return client.PipelineTemplates.Update(context.Background(), c.String("template-name"), &ptr)
}

// DeletePipelineTemplateCommand handles the interaction between the cli flags and the action handler for
// delete-pipeline-template and checks a template-name is provided and that the response is a 2xx response.
func deletePipelineTemplateCommand() *cli.Command {
	return &cli.Command{
		Name:     DeletePipelineTemplateCommandName,
		Usage:    DeletePipelineTemplateCommandUsage,
		Category: "Pipeline Templates",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "template-name", Usage: "Pipeline Template name."}},
		Action: func(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
			if c.String("template-name") == "" {
				return nil, nil, NewFlagError("template-name")
			}

			deleteResponse, resp, err := client.PipelineTemplates.Delete(context.Background(), c.String("template-name"))
			if resp.HTTP.StatusCode == 406 {
				err = errors.New(deleteResponse)
			}
			return deleteResponse, resp, err
		},
	}
}

// ListPipelineTemplatesCommand handles the interaction between the cli flags and the action handler for
// list-pipeline-templates
func listPipelineTemplatesCommand() *cli.Command {
	return &cli.Command{
		Name:     ListPipelineTemplatesCommandName,
		Usage:    ListPipelineTemplatesCommandUsage,
		Action:   ActionWrapper(listPipelineTemplatesAction),
		Category: "Pipeline Templates",
	}
}

// GetPipelineTemplateCommand handles the interaction between the cli flags and the action handler for
// get-pipeline-template
func getPipelineTemplateCommand() *cli.Command {
	return &cli.Command{
		Name:     GetPipelineTemplateCommandName,
		Usage:    GetPipelineTemplateCommandUsage,
		Action:   ActionWrapper(getPipelineTemplateAction),
		Category: "Pipeline Templates",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "template-name", Usage: "Name of the Pipeline Template configuration."},
		},
	}
}

// CreatePipelineTemplateCommand handles the interaction between the cli flags and the action handler for
// create-pipeline-template
func createPipelineTemplateCommand() *cli.Command {
	return &cli.Command{
		Name:     CreatePipelineTemplateCommandName,
		Usage:    CreatePipelineTemplateCommandUsage,
		Action:   ActionWrapper(createPipelineTemplateAction),
		Category: "Pipeline Templates",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "template-name", Usage: "Pipeline Template name."},
			cli.StringSliceFlag{Name: "stage", Usage: "JSON encoded stage object."},
		},
	}
}

// UpdatePipelineTemplateCommand handles the interaction between the cli flags and the action handler for
// update-pipeline-template
func updatePipelineTemplateCommand() *cli.Command {
	return &cli.Command{
		Name:     UpdatePipelineTemplateCommandName,
		Usage:    UpdatePipelineTemplateCommandUsage,
		Action:   ActionWrapper(updatePipelineTemplateAction),
		Category: "Pipeline Templates",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "template-version", Usage: "Pipeline template version."},
			cli.StringFlag{Name: "template-name", Usage: "Pipeline Template name."},
			cli.StringSliceFlag{Name: "stage", Usage: "JSON encoded stage object."},
		},
	}
}
