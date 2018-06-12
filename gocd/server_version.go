package gocd

import (
	"context"
)

var cachedServerVersion *ServerVersion

type ServerVersionService service

type ServerVersionParts struct {
	Major int
	Minor int
	Patch int
}

type ServerVersion struct {
	Version      string `json:"version"`
	VersionParts *ServerVersionParts
	BuildNumber  string `json:"build_number"`
	GitSha       string `json:"git_sha"`
	FullVersion  string `json:"full_version"`
	CommitURL    string `json:"commit_url"`
}

// Get retrieves information about a specific plugin.
func (svs *ServerVersionService) Get(ctx context.Context) (v *ServerVersion, resp *APIResponse, err error) {
	if cachedServerVersion != nil {
		return cachedServerVersion, nil, nil
	}

	v = &ServerVersion{}
	_, resp, err = svs.client.getAction(ctx, &APIClientRequest{
		Path:         "version",
		ResponseBody: v,
		APIVersion:   apiV1,
	})

	err = v.parseVersion()

	cachedServerVersion = v

	return
}
