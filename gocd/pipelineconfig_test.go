package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestPipelineConfig(t *testing.T) {
	if runIntegrationTest(t) {
		input := &Pipeline{
			Name: "new_pipeline",
			Materials: []Material{{
				Type: "git",
				Attributes: MaterialAttributesGit{
					URL:         "git@github.com:sample_repo/example.git",
					Destination: "dest",
					Branch:      "master",
				},
			}},
			Stages: []*Stage{{
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
				}},
			}},
		}

		ctx := context.Background()

		p, _, err := intClient.PipelineConfigs.Create(ctx, "test-group", input)
		assert.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]{32}--gzip$"), p.Version)

		p.RemoveLinks()
		assert.Equal(t, &Pipeline{
			Group:                "test-group",
			Name:                 "new_pipeline",
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
		}, p)

		getP, _, err := intClient.PipelineConfigs.Get(ctx, input.Name)

		getP.RemoveLinks()
		assert.Equal(t, &Pipeline{
			Name:                 "new_pipeline",
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
		}, getP)

		p.LabelTemplate = "Updated_${COUNT}"
		updatedP, _, err := intClient.PipelineConfigs.Update(context.Background(), p.Name, p)
		assert.NoError(t, err)
		assert.NotEqual(t, p.Version, updatedP.Version)
		updatedP.Version = p.Version

		updatedP.RemoveLinks()
		assert.Equal(t, &Pipeline{
			Group:                "test-group",
			Name:                 "new_pipeline",
			LabelTemplate:        p.LabelTemplate,
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
		}, updatedP)

		message, _, err := intClient.PipelineConfigs.Delete(ctx, input.Name)
		assert.Equal(t, "The pipeline 'new_pipeline' was deleted successfully.", message)

	}
}
