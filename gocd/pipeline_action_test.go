package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testPipelineServiceUnPause(t *testing.T) {
	if runIntegrationTest(t) {

		ctx := context.Background()
		pipelineName := "test-pipeline-un-pause"

		stages := buildMockPipelineStages()
		mockPipeline := &Pipeline{
			Name:                 pipelineName,
			Group:                mockTestingGroup,
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
			Stages: stages,
		}

		pausePipeline, _, err := intClient.PipelineConfigs.Create(ctx, mockTestingGroup, mockPipeline)
		assert.NoError(t, err)
		pausePipeline.Links = nil
		pausePipeline.Version = ""
		assert.Equal(t, mockPipeline, pausePipeline)

		pp, _, err := intClient.Pipelines.Unpause(ctx, pipelineName)
		assert.NoError(t, err)
		assert.True(t, pp)

		pp, _, err = intClient.Pipelines.Pause(ctx, pipelineName)
		assert.NoError(t, err)
		assert.True(t, pp)

		pp, _, err = intClient.Pipelines.ReleaseLock(context.Background(), pipelineName)
		assert.EqualError(t, err, "Received HTTP Status '406 Not Acceptable'")
		assert.False(t, pp)

		deleteResponse, _, err := intClient.PipelineConfigs.Delete(ctx, pipelineName)
		assert.Equal(t, "The pipeline 'test-pipeline-un-pause' was deleted successfully.", deleteResponse)
	}

}
