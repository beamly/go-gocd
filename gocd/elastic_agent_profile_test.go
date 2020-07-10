package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElasticAgentProfile(t *testing.T) {
	for _, test := range []struct {
		name          string
		eapCreate     *ElasticAgentProfile
		eapCreateWant *ElasticAgentProfile
		eapGet        string
		eapGetWant    *ElasticAgentProfile
		eapUpdate     *ElasticAgentProfile
		eapUpdateWant *ElasticAgentProfile
		eapListWant   []*ElasticAgentProfile
		eapDelete     string
		eapDeleteWant string
	}{
		{
			name: "basic",
			eapCreate: &ElasticAgentProfile{
				ID:       "mock-id",
				PluginID: "com.example.elasticagent.foocloud",
				Properties: []*ConfigProperty{
					{Key: "Image", Value: "alpine:latest"},
				},
			},
			eapCreateWant: &ElasticAgentProfile{
				ID:       "mock-id",
				PluginID: "com.example.elasticagent.foocloud",
				Properties: []*ConfigProperty{
					{Key: "Image", Value: "alpine:latest"},
					{Key: "Command", Value: ""},
					{Key: "Environment", Value: ""},
					{Key: "MaxMemory", Value: "200M"},
					{Key: "ReservedMemory", Value: "150M"},
				},
			},
			eapGet: "mock-name",
			eapGetWant: &ElasticAgentProfile{
				ID:       "mock-id",
				PluginID: "com.example.elasticagent.foocloud",
				Properties: []*ConfigProperty{
					{Key: "Image", Value: "alpine:latest"},
					{Key: "Command", Value: ""},
					{Key: "Environment", Value: ""},
					{Key: "MaxMemory", Value: "200M"},
					{Key: "ReservedMemory", Value: "150M"},
				},
			},
			eapUpdate: &ElasticAgentProfile{
				ID:       "mock-id",
				PluginID: "com.example.elasticagent.foocloud",
				Properties: []*ConfigProperty{
					{Key: "Image", Value: "alpine:latest"},
					{Key: "Command", Value: ""},
					{Key: "Environment", Value: ""},
					{Key: "MaxMemory", Value: "2G"},
					{Key: "ReservedMemory", Value: "1G"},
				},
			},
			eapUpdateWant: &ElasticAgentProfile{
				ID:       "mock-id",
				PluginID: "com.example.elasticagent.foocloud",
				Properties: []*ConfigProperty{
					{Key: "Image", Value: "alpine:latest"},
					{Key: "Command", Value: ""},
					{Key: "Environment", Value: ""},
					{Key: "MaxMemory", Value: "2G"},
					{Key: "ReservedMemory", Value: "1G"},
				},
			},
			eapListWant: []*ElasticAgentProfile{
				{
					ID:       "mock-id",
					PluginID: "com.example.elasticagent.foocloud",
					Properties: []*ConfigProperty{
						{Key: "Image", Value: "alpine:latest"},
						{Key: "Command", Value: ""},
						{Key: "Environment", Value: ""},
						{Key: "MaxMemory", Value: "2G"},
						{Key: "ReservedMemory", Value: "1G"},
					},
				},
			},
			eapDelete: `mock-id`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if runIntegrationTest(t) {
				ctx := context.Background()

				t.Run("create", func(t *testing.T) {
					eapCreateGot, _, err := intClient.ElasticAgentProfiles.Create(ctx, test.eapCreate)
					if assert.NoError(t, err) {
						assert.Equal(t, test.eapCreateWant, eapCreateGot)
					}
				})

				t.Run("get", func(t *testing.T) {
					eapGetGot, _, err := intClient.ElasticAgentProfiles.Get(ctx, test.eapGet)
					if assert.NoError(t, err) {
						assert.Equal(t, test.eapGetWant, eapGetGot)
					}
				})

				t.Run("update", func(t *testing.T) {
					eapUpdateGot, _, err := intClient.ElasticAgentProfiles.Update(ctx, test.eapUpdate.ID, test.eapUpdate)
					if assert.NoError(t, err) {
						assert.Equal(t, test.eapUpdateWant, eapUpdateGot)
					}
				})

				t.Run("list", func(t *testing.T) {
					eapListGot, _, err := intClient.ElasticAgentProfiles.List(ctx)
					if assert.NoError(t, err) {
						assert.Equal(t, test.eapListWant, eapListGot)
					}
				})

				t.Run("delete", func(t *testing.T) {
					eapDeleteGot, _, err := intClient.ElasticAgentProfiles.Delete(ctx, test.eapDelete)
					if assert.NoError(t, err) {
						assert.Equal(t, test.eapDeleteWant, eapDeleteGot)
					}
				})
			}
		})
	}
}
