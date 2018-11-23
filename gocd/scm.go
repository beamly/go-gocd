package gocd

import (
	"context"
	"fmt"
)

// SCMsService exposes calls for interacting with SCM objects in the GoCD API.
type SCMsService service

// SCM describes a material for a Plugin Version Control system.
type SCM struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	AutoUpdate     bool                `json:"auto_update"`
	PluginMetadata *SCMMetadata        `json:"plugin_metadata"`
	Configuration  []*SCMConfiguration `json:"configuration"`

	Version string    `json:"version,omitempty"`
	Links   *HALLinks `json:"_links,omitempty"`
}

// SCMMetadata describing the plugin this SCM is describing
type SCMMetadata struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

// SCMConfiguration describing an attribute in an SCM
type SCMConfiguration struct {
	Key            string `json:"key"`
	Value          string `json:"value,omitempty"`
	EncryptedValue string `json:"encrypted_value,omitempty"`
}

// RoleListWrapper describes a container for the result of a role list operation
type SCMListWrapper struct {
	Embedded struct {
		SCMs []*SCM `json:"scms"`
	} `json:"_embedded"`
	Links *HALLinks `json:"_links"`
}

// Create a new SCM
func (scms *SCMsService) Create(ctx context.Context, newSCM *SCM) (scm *SCM, resp *APIResponse, err error) {
	var ver string
	scm = &SCM{}
	if ver, err = scms.client.getAPIVersion(ctx, "admin/scms/"); err == nil {
		_, resp, err = scms.client.postAction(ctx, &APIClientRequest{
			Path:         "admin/scms",
			APIVersion:   ver,
			RequestBody:  newSCM,
			ResponseBody: scm,
		})
	}
	return
}

// Get an SCM by name
func (scms *SCMsService) Get(ctx context.Context, name string) (scm *SCM, resp *APIResponse, err error) {
	var ver string
	scm = &SCM{}
	if ver, err = scms.client.getAPIVersion(ctx, "admin/scms/:material_name"); err == nil {
		_, resp, err = scms.client.getAction(ctx, &APIClientRequest{
			Path:         fmt.Sprintf("admin/scms/%s", name),
			APIVersion:   ver,
			ResponseBody: scm,
		})
	}
	return
}

// Update an SCM by name
func (scms *SCMsService) Update(ctx context.Context, name string, newSCM *SCM) (scm *SCM, resp *APIResponse, err error) {
	var ver string
	scm = &SCM{}
	if ver, err = scms.client.getAPIVersion(ctx, "admin/scms/:material_name"); err == nil {
		_, resp, err = scms.client.putAction(ctx, &APIClientRequest{
			Path:         fmt.Sprintf("admin/scms/%s", name),
			APIVersion:   ver,
			RequestBody:  newSCM,
			ResponseBody: scm,
		})
	}
	return
}

// List all SCMs
func (scms *SCMsService) List(ctx context.Context) (scmSlice []*SCM, resp *APIResponse, err error) {
	var ver string
	scmWrapper := &SCMListWrapper{}
	if ver, err = scms.client.getAPIVersion(ctx, "admin/scms/"); err == nil {
		_, resp, err = scms.client.getAction(ctx, &APIClientRequest{
			Path:         "admin/scms",
			APIVersion:   ver,
			ResponseBody: scmWrapper,
		})
	}
	scmSlice = scmWrapper.Embedded.SCMs
	return
}
