package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestPipelineConfig(t *testing.T) {
	for _, tt := range []struct {
		name               string
		pipelineCreate     *Pipeline
		pipelineCreateWant *Pipeline
		pipelineGet        string
		pipelineGetWant    *Pipeline
		pipelineUpdate     *Pipeline
		pipelineUpdateWant *Pipeline
		delete             string
		deleteWant         string
		serverVersions     *ServerVersionRange
	}{
		{
			name: "v5",
			serverVersions: newServerVersionRangeFromString(
				SERVER_VERSION_14_3_0,
				SERVER_VERSION_18_6_0,
			),
			pipelineCreate: &Pipeline{
				Name:      "new_pipeline",
				Materials: buildMockMaterials(),
				Stages:    buildMockPipelineStages(),
			},
			pipelineCreateWant: &Pipeline{
				Group:                "test-group",
				Name:                 "new_pipeline",
				LabelTemplate:        "${COUNT}",
				Parameters:           make([]*Parameter, 0),
				EnvironmentVariables: make([]*EnvironmentVariable, 0),
				Materials:            buildMockMaterials(),
				Stages:               buildMockPipelineStages(),
			},
			pipelineGet: "new_pipeline",
			pipelineGetWant: &Pipeline{
				Group:                "test-group",
				Name:                 "new_pipeline",
				LabelTemplate:        "${COUNT}",
				Parameters:           make([]*Parameter, 0),
				EnvironmentVariables: make([]*EnvironmentVariable, 0),
				Materials:            buildMockMaterials(),
				Stages:               buildMockPipelineStages(),
			},
			pipelineUpdateWant: &Pipeline{
				Group:                "test-group",
				Name:                 "new_pipeline",
				LabelTemplate:        "Update ${COUNT}",
				Parameters:           make([]*Parameter, 0),
				EnvironmentVariables: make([]*EnvironmentVariable, 0),
				Materials:            buildMockMaterials(),
				Stages:               buildMockPipelineStages(),
			},
			delete:     "new_pipeline",
			deleteWant: "The pipeline 'new_pipeline' was deleted successfully.",
		},
		{
			name: "v6",
			serverVersions: newServerVersionRangeFromString(
				SERVER_VERSION_18_6_0,
				SERVER_VERSION_MAX,
			),
			pipelineCreate: &Pipeline{
				Name:         "new_pipeline",
				LockBehavior: "lockOnFailure",
				Materials:    buildMockMaterials(),
				Stages:       buildMockPipelineStages(),
			},
			pipelineCreateWant: &Pipeline{
				Group:                "test-group",
				Name:                 "new_pipeline",
				LabelTemplate:        "${COUNT}",
				LockBehavior:         "lockOnFailure",
				Parameters:           make([]*Parameter, 0),
				EnvironmentVariables: make([]*EnvironmentVariable, 0),
				Materials:            buildMockMaterials(),
				Stages:               buildMockPipelineStages(),
			},
			pipelineGet: "new_pipeline",
			pipelineGetWant: &Pipeline{
				Group:                "test-group",
				Name:                 "new_pipeline",
				LabelTemplate:        "${COUNT}",
				LockBehavior:         "lockOnFailure",
				Parameters:           make([]*Parameter, 0),
				EnvironmentVariables: make([]*EnvironmentVariable, 0),
				Materials:            buildMockMaterials(),
				Stages:               buildMockPipelineStages(),
			},
			pipelineUpdateWant: &Pipeline{
				Group:                "test-group",
				Name:                 "new_pipeline",
				LabelTemplate:        "Update ${COUNT}",
				LockBehavior:         "lockOnFailure",
				Parameters:           make([]*Parameter, 0),
				EnvironmentVariables: make([]*EnvironmentVariable, 0),
				Materials:            buildMockMaterials(),
				Stages:               buildMockPipelineStages(),
			},
			delete:     "new_pipeline",
			deleteWant: "The pipeline 'new_pipeline' was deleted successfully.",
		},
	} {
		t.Run("basic", func(t *testing.T) {
			if runIntegrationTest(t) {
				if runOnlyForServerVersionRange(t, tt.serverVersions) {

					var getPipeline *Pipeline
					ctx := context.Background()
					t.Run("create", func(t *testing.T) {
						createPipeline, _, err := intClient.PipelineConfigs.Create(ctx, "test-group", tt.pipelineCreate)
						assert.NoError(t, err)
						assert.Regexp(t, regexp.MustCompile("^[a-f0-9]{32}--gzip$"), createPipeline.Version)

						createPipeline.RemoveLinks()
						assert.Equal(t, tt.pipelineCreateWant, createPipeline)
					})

					t.Run("get", func(t *testing.T) {
						getPipeline, _, err := intClient.PipelineConfigs.Get(ctx, tt.pipelineGet)
						assert.NoError(t, err)

						getPipeline.RemoveLinks()
						assert.Equal(t, tt.pipelineGetWant, getPipeline)
					})

					t.Run("update", func(t *testing.T) {
						updatePipeline, _, err := intClient.PipelineConfigs.Update(context.Background(),
							tt.pipelineUpdate.Name, tt.pipelineUpdate)
						assert.NoError(t, err)
						assert.NotEqual(t, getPipeline.Version, updatePipeline.Version)
						updatePipeline.Version = getPipeline.Version

						updatePipeline.RemoveLinks()
						assert.Equal(t, tt.pipelineUpdateWant, updatePipeline)
					})

					t.Run("delete", func(t *testing.T) {
						message, _, err := intClient.PipelineConfigs.Delete(ctx, tt.delete)
						assert.NoError(t, err)
						assert.Equal(t, tt.deleteWant, message)
					})
				}
			}
		})
	}
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

func buildMockMaterials() []Material {
	return []Material{{
		Type: "git",
		Attributes: MaterialAttributesGit{
			URL:         "git@github.com:sample_repo/example.git",
			Destination: "dest",
			Branch:      "master",
			AutoUpdate:  true,
		},
	}}
}
