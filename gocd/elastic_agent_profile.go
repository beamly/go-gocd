package gocd

import "context"

// RoleService describes Actions which can be performed on roles
type ElasticAgentService service

// Role represents a type of agent/actor who can access resources perform operations
type ElasticAgentProfile struct {
	ID         string            `json:"id"`
	PluginID   string            `json:"plugin_id"`
	Properties []*ConfigProperty `json:"properties"`
	Version    string            `json:"version,omitempty"`
	Links      *HALLinks         `json:"_links,omitempty"`
}

type ElasticAgentProfileListResponse struct {
	Links    *HALLinks `json:"_links,omitempty"`
	Embedded struct {
		Profiles []*ElasticAgentProfile `json:"profiles"`
	} `json:"_embedded"`
}

func (eaps *ElasticAgentService) Create(ctx context.Context, newProfile *ElasticAgentProfile) (profile *ElasticAgentProfile, resp *APIResponse, err error) {
	var ver string
	profile = &ElasticAgentProfile{}
	if ver, err = eaps.client.getAPIVersion(ctx, "elastic/profiles"); err == nil {
		_, resp, err = eaps.client.postAction(ctx, &APIClientRequest{
			Path:         "elastic/profiles",
			APIVersion:   ver,
			RequestBody:  newProfile,
			ResponseBody: profile,
		})
	}
	return
}

func (eaps *ElasticAgentService) Get(ctx context.Context, profileID string) (eap *ElasticAgentProfile, resp *APIResponse, err error) {
	return

}

func (eaps *ElasticAgentService) Update(ctx context.Context, profileID string, profile *ElasticAgentProfile) (eap *ElasticAgentProfile, resp *APIResponse, err error) {
	return

}

func (eaps *ElasticAgentService) List(ctx context.Context) (eapSlice []*ElasticAgentProfile, resp *APIResponse, err error) {
	var ver string
	response := &ElasticAgentProfileListResponse{}
	if ver, err = eaps.client.getAPIVersion(ctx, "elastic/profiles"); err == nil {
		_, resp, err = eaps.client.getAction(ctx, &APIClientRequest{
			Path:         "elastic/profiles",
			APIVersion:   ver,
			ResponseBody: response,
		})
	}
	return response.Embedded.Profiles, resp, err
}

func (eaps *ElasticAgentService) Delete(ctx context.Context, profileName string) (result string, resp *APIResponse, err error) {
	return

}
