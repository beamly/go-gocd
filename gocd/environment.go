package gocd

import (
	"context"
)

// EnvironmentsService exposes calls for interacting with Environment objects in the GoCD API.
type EnvironmentsService service

// EnvironmentsResponse describes the response obejct for a plugin API call.
type EnvironmentsResponse struct {
	Links    *HALLinks             `json:"_links"`
	Embedded *EmbeddedEnvironments `json:"_embedded"`
}

// EmbeddedEnvironments encapsulates the environment struct
type EmbeddedEnvironments struct {
	Environments []*Environment `json:"environments"`
}

// Environment describes a group of pipelines and agents
type Environment struct {
	Links                *HALLinks              `json:"_links,omitempty"`
	Name                 string                 `json:"name"`
	Pipelines            []*Pipeline            `json:"pipelines,omitempty"`
	Agents               []*Agent               `json:"agents,omitempty"`
	EnvironmentVariables []*EnvironmentVariable `json:"environment_variables,omitempty"`
	Version              string                 `json:"version"`
}

// EnvironmentPatchRequest describes the actions to perform on an environment
type EnvironmentPatchRequest struct {
	Pipelines            *PatchStringAction          `json:"pipelines"`
	Agents               *PatchStringAction          `json:"agents"`
	EnvironmentVariables *EnvironmentVariablesAction `json:"environment_variables"`
}

// EnvironmentVariablesAction describes a collection of Environment Variables to add or remove.
type EnvironmentVariablesAction struct {
	Add    []*EnvironmentVariable `json:"add"`
	Remove []string               `json:"remove"`
}

// PatchStringAction describes a collection of resources to add or remove.
type PatchStringAction struct {
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

// List all environments
func (es *EnvironmentsService) List(ctx context.Context) (e *EnvironmentsResponse, resp *APIResponse, err error) {
	e = &EnvironmentsResponse{}
	_, resp, err = es.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/environments",
		ResponseBody: e,
		APIVersion:   apiV2,
	})

	return
}

// Delete an environment
func (es *EnvironmentsService) Delete(ctx context.Context, name string) (string, *APIResponse, error) {
	return es.client.deleteAction(ctx, "admin/environments/"+name, apiV2)
}

// Create an environment
func (es *EnvironmentsService) Create(ctx context.Context, name string) (e *Environment, resp *APIResponse, err error) {
	_, resp, err = es.client.postAction(ctx, &APIClientRequest{
		Path: "admin/environments/",
		RequestBody: Environment{
			Name: name,
		},
		ResponseBody: &e,
		APIVersion:   apiV2,
	})

	return
}

// Get a single environment by name
func (es *EnvironmentsService) Get(ctx context.Context, name string) (e *Environment, resp *APIResponse, err error) {
	e = &Environment{}
	_, resp, err = es.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/environments/" + name,
		ResponseBody: e,
		APIVersion:   apiV2,
	})

	return
}

// Patch an environments configuration by adding or removing pipelines, agents, environment variables
func (es *EnvironmentsService) Patch(ctx context.Context, name string, patch *EnvironmentPatchRequest) (e *Environment, resp *APIResponse, err error) {
	e = &Environment{}
	_, resp, err = es.client.patchAction(ctx, &APIClientRequest{
		Path:         "admin/environments/" + name,
		RequestBody:  patch,
		ResponseBody: e,
		APIVersion:   apiV2,
	})

	return
}
