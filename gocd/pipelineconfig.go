package gocd

import (
	"context"
	//"fmt"
)

// PipelineConfigsService describes the HAL _link resource for the api response object for a pipelineconfig
type PipelineConfigsService service

// PipelineConfigRequest describes a request object for creating or updating pipelines
type PipelineConfigRequest struct {
	Group    string    `json:"group,omitempty"`
	Pipeline *Pipeline `json:"pipeline"`
}

// Get a single Pipeline object in the GoCD API.
func (pcs *PipelineConfigsService) Get(ctx context.Context, name string) (p *Pipeline, resp *APIResponse, err error) {
	p = &Pipeline{}
	_, resp, err = pcs.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/pipelines/" + name,
		APIVersion:   apiV4,
		ResponseBody: p,
	})

	return
}

// Update a pipeline configuration
func (pcs *PipelineConfigsService) Update(ctx context.Context, name string, p *Pipeline) (pr *Pipeline, resp *APIResponse, err error) {
	pr = &Pipeline{}

	_, resp, err = pcs.client.putAction(ctx, &APIClientRequest{
		Path:       "admin/pipelines/" + name,
		APIVersion: apiV4,
		RequestBody: &PipelineConfigRequest{
			Pipeline: p,
		},
		ResponseBody: pr,
	})

	pr.Group = p.Group

	return
}

// Create a pipeline configuration
func (pcs *PipelineConfigsService) Create(ctx context.Context, group string, p *Pipeline) (pr *Pipeline, resp *APIResponse, err error) {

	pr = &Pipeline{}
	_, resp, err = pcs.client.postAction(ctx, &APIClientRequest{
		Path:       "admin/pipelines",
		APIVersion: apiV4,
		RequestBody: &PipelineConfigRequest{
			Group:    group,
			Pipeline: p,
		},
		ResponseBody: pr,
	})

	pr.Group = group

	return
}

// Delete a pipeline configuration
func (pcs *PipelineConfigsService) Delete(ctx context.Context, name string) (string, *APIResponse, error) {
	return pcs.client.deleteAction(ctx, "admin/pipelines/"+name, apiV4)
}
