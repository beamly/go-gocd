package gocd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testPipelineServiceUnPause(t *testing.T) {
	for n, test := range []struct {
		name          string
		v             *ServerVersionRange
		confirmHeader string
		acceptHeader  string
	}{
		{
			name:          "server-version-14.3.0",
			v:             newServerVersionRangeFromString("1.0.0", "14.3.0"),
			confirmHeader: "Confirm",
			acceptHeader:  apiV0,
		},
		{
			name:          "server-version-18.3.0",
			v:             newServerVersionRangeFromString("14.3.0", "18.3.0"),
			confirmHeader: "X-GoCD-Confirm",
			acceptHeader:  apiV1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if runOnlyForServerVersionRange(t, test.v) {

				ctx := context.Background()
				pipelineName := fmt.Sprintf("test-pipeline-un-pause%d", n)

				pausePipeline, _, err := intClient.PipelineConfigs.Create(ctx, mockTestingGroup, &Pipeline{
					Name: pipelineName,
					Materials: []Material{{
						Type: "git",
					}},
					Stages: buildMockPipelineStages(),
				})
				assert.NoError(t, err)
				assert.Equal(t, nil, pausePipeline)

				pp, _, err := intClient.Pipelines.Pause(ctx, pipelineName)
				assert.NoError(t, err)
				assert.True(t, pp)

				deleteResponse, _, err := intClient.PipelineConfigs.Delete(ctx, pipelineName)
				assert.Equal(t, "", deleteResponse)

			}
		})
	}
}
