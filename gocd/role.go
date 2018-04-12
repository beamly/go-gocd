package gocd

import (
	"context"
)

// PropertiesService describes Actions which can be performed on agents
type RoleService service

type Role struct {
	Name       string              `json:"name"`
	Type       string              `json:"type"`
	Attributes *RoleAttributesGoCD `json:"attributes"`
}

type RoleAttributesGoCD struct {
	Users []string `json:"users"`
}

type RoleAttributesPlugin struct{}

// Create a role
func (rs *RoleService) Create(ctx context.Context, role *Role) (r *Role, resp *APIResponse, err error) {
	r = &Role{}
	_, resp, err = rs.client.postAction(ctx, &APIClientRequest{
		APIVersion:   apiV1,
		Path:         "admin/security/roles",
		RequestBody:  role,
		ResponseBody: r,
	})

	return
}