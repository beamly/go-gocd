package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestPipelineConfig(t *testing.T) {
	if !runIntegrationTest(t) {
		t.Skip("Skipping acceptance tests as GOCD_ACC not set to 1")
	}
	input := &Pipeline{
		Name: "test_pipeline_config",
		Materials: []Material{{
			Type: "git",
			Attributes: MaterialAttributesGit{
				URL:         "git@github.com:sample_repo/example.git",
				Destination: "dest",
				Branch:      "master",
			},
		}},
		Stages: buildMockPipelineStages(),
	}

	ctx := context.Background()

	p, _, err := intClient.PipelineConfigs.Create(ctx, "test-group", input)
	assert.NoError(t, err)
	assert.Regexp(t, regexp.MustCompile("^[a-f0-9]{32}--gzip$"), p.Version)

	p.RemoveLinks()
	expected := &Pipeline{
		Group:                "test-group",
		Name:                 "test_pipeline_config",
		LabelTemplate:        "${COUNT}",
		Parameters:           make([]*Parameter, 0),
		EnvironmentVariables: make([]*EnvironmentVariable, 0),
		Materials: []Material{{
			Type: "git",
			Attributes: &MaterialAttributesGit{
				URL:         "git@github.com:sample_repo/example.git",
				Destination: "dest",
				Branch:      "master",
				AutoUpdate:  true,
			},
		}},
		Stages: []*Stage{{
			Name: "defaultStage",
			Approval: &Approval{
				Type: "success",
				Authorization: &Authorization{
					Users: []string{},
					Roles: []string{},
				},
			},
			Jobs: []*Job{{
				Name:                 "defaultJob",
				EnvironmentVariables: []*EnvironmentVariable{},
				Resources:            []string{},
				Tasks: []*Task{{
					Type: "exec",
					Attributes: TaskAttributes{
						RunIf:   []string{"passed"},
						Command: "ls",
					},
				}},
				Tabs:      []*Tab{},
				Artifacts: []*Artifact{},
			}},
			EnvironmentVariables: []*EnvironmentVariable{},
		}},
		Version: p.Version,
	}

	// Make sure version-specific defaults are properly set
	apiVersion, err := intClient.getAPIVersion(ctx, "admin/pipelines/:pipeline_name")
	assert.NoError(t, err)
	switch apiVersion {
	case apiV6:
		expected.Origin = &PipelineConfigOrigin{Type: "gocd"}
		fallthrough
	case apiV5:
		expected.LockBehavior = "none"
	}

	assert.Equal(t, expected, p)

	getP, _, err := intClient.PipelineConfigs.Get(ctx, input.Name)

	getP.RemoveLinks()
	expected.Group = ""
	assert.Equal(t, expected, getP)

	// The tests on the update have been commented as it seems there's a problem on 18.7.0 about it
	p.LabelTemplate = "Updated_${COUNT}"
	p.EnvironmentVariables = []*EnvironmentVariable{{Name: "FOO", Value: "bar"}}
	updatedP, _, err := intClient.PipelineConfigs.Update(ctx, p.Name, p)
	assert.NoError(t, err)
	assert.NotEqual(t, p.Version, updatedP.Version)
	updatedP.Version = p.Version

	updatedP.RemoveLinks()
	expected.LabelTemplate = "Updated_${COUNT}"
	expected.EnvironmentVariables = []*EnvironmentVariable{{Name: "FOO", Value: "bar"}}
	expected.Group = "test-group"
	assert.Equal(t, expected, updatedP)

	message, _, err := intClient.PipelineConfigs.Delete(ctx, input.Name)
	assert.Equal(t, "The pipeline 'test_pipeline_config' was deleted successfully.", message)

}

func buildMockPipelineStages() []*Stage {
	return []*Stage{{
		Name: "defaultStage",
		Jobs: []*Job{{
			Name: "defaultJob",
			Tasks: []*Task{{
				Type: "exec",
				Attributes: TaskAttributes{
					RunIf:   []string{"passed"},
					Command: "ls",
				},
			}},
			Tabs:                 make([]*Tab, 0),
			Artifacts:            make([]*Artifact, 0),
			EnvironmentVariables: make([]*EnvironmentVariable, 0),
			Resources:            []string{},
		}},
		Approval: &Approval{
			Type: "success",
			Authorization: &Authorization{
				Users: make([]string, 0),
				Roles: make([]string, 0),
			},
		},
		EnvironmentVariables: make([]*EnvironmentVariable, 0),
	}}
}
