package gocd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResource(t *testing.T) {
	//t.Run("Pipeline", testResourcePipeline)
	//t.Run("PipelineTemplate", testResourcePipelineTemplate)
	t.Run("PipelineGroups", testResourcePipelineGroups)
	t.Run("StageContainer", testResourceStageContainers)
	t.Run("HALContainer", testResourceHALContainers)
	t.Run("Versioned", testResourceVersioned)
	t.Run("Materials", testResourceMaterial)
}

func testResourceVersioned(t *testing.T) {
	vers := map[string]Versioned{
		"Environment":             &Environment{Version: "mock-version1"},
		"PipelineTemplate":        &PipelineTemplate{Version: "mock-version1"},
		"PipelineConfigRequest":   &PipelineConfigRequest{Pipeline: &Pipeline{Version: "mock-version1"}},
		"PipelineTemplateRequest": &PipelineTemplateRequest{Version: "mock-version1"},
		"Role":                    &Role{Version: "mock-version1"},
	}
	for key, ver := range vers {
		t.Run(key, func(t *testing.T) {
			testResourceVersion(t, ver)
		})
	}
}

func testResourceVersion(t *testing.T, ver Versioned) {
	v := ver.GetVersion()
	assert.Equal(t, "mock-version1", v)

	ver.SetVersion("mock-version")
	v = ver.GetVersion()
	assert.Equal(t, "mock-version", v)
}

func testResourceHALContainers(t *testing.T) {
	l := &HALLinks{links: []*HALLink{}}

	hals := map[string]HALContainer{
		"Agent":    &Agent{Links: l},
		"Pipeline": &Pipeline{Links: l},
		"Role":     &Role{Links: l},
		"EnvironmentResponse": &EnvironmentsResponse{
			Links: l,
			Embedded: &EmbeddedEnvironments{
				Environments: []*Environment{
					{Links: l},
				},
			},
		},
		"Environment": &Environment{
			Links: l,
			Pipelines: []*Pipeline{{
				Links: l,
			}},
			Agents: []*Agent{{
				Links: l,
			}},
		},
	}
	for key, hal := range hals {
		t.Run(key, func(t *testing.T) {
			testResourceHALContainer(t, hal)
		})
	}
}

func testResourceHALContainer(t *testing.T, hal HALContainer) {
	assert.NotNil(t, hal.GetLinks())
	hal.RemoveLinks()
	assert.Nil(t, hal.GetLinks())
}

func testResourceStageContainers(t *testing.T) {
	scs := map[string]StageContainer{
		"PipelineTemplate": new(PipelineTemplate),
		"Pipeline":         new(Pipeline),
	}
	for key, sc := range scs {
		t.Run(key, func(tr *testing.T) {
			testResourceStageContainerI(tr, sc)
		})
	}
}

func testResourceStageContainerI(t *testing.T, sc StageContainer) {
	s1 := Stage{Name: "s"}
	s1.CleanWorkingDirectory = false

	s2 := Stage{Name: "s"}
	s2.CleanWorkingDirectory = true

	sc.AddStage(&s1)

	s := sc.GetStage("s")
	assert.False(t, s.CleanWorkingDirectory)

	sc.SetStage(&s2)

	s = sc.GetStage("s")
	assert.True(t, s.CleanWorkingDirectory)

	s3 := Stage{Name: "s3"}
	sc.SetStage(&s3)

	sc.SetStage(&s3)

	s = sc.GetStage("s3")
	assert.Equal(t, s.Name, "s3")

}
