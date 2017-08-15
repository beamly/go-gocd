package gocd

import "context"

// PipelineGroupsService describes the HAL _link resource for the api response object for a pipeline group response.
type PipelineGroupsService service

// PipelineGroup describes a pipeline group API response.
type PipelineGroup struct {
	Name      string      `json:"name"`
	Pipelines []*Pipeline `json:"pipelines"`
}

// List Pipeline groups
func (pgs *PipelineGroupsService) List(ctx context.Context, name string) ([]*PipelineGroup, *APIResponse, error) {

	req, err := pgs.client.NewRequest("GET", "config/pipeline_groups", nil, "")
	if err != nil {
		return nil, nil, err
	}

	pg := []*PipelineGroup{}
	resp, err := pgs.client.Do(ctx, req, &pg, responseTypeJSON)
	if err != nil {
		return nil, resp, err
	}

	filtered := []*PipelineGroup{}
	if name != "" {
		for _, pipelineGroup := range pg {
			if pipelineGroup.Name == name {
				filtered = append(filtered, pipelineGroup)
			}
		}
	} else {
		filtered = pg
	}

	return filtered, resp, nil
}
