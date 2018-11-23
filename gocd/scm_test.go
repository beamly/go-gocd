package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSCM(t *testing.T) {
	if runIntegrationTest(t) {
		ctx := context.Background()

		encPassword, _, err := intClient.Encryption.Encrypt(ctx, "password")
		assert.NoError(t, err)
		for _, test := range []struct {
			name          string
			scmCreate     *SCM
			scmCreateWant *SCM
			scmGet        string
			scmGetWant    *SCM
			scmUpdate     *SCM
			scmUpdateWant *SCM
			scmListWant   []*SCM
		}{
			{
				name: "basic",
				scmCreate: &SCM{
					ID:         "mock-id",
					Name:       "mock-name",
					AutoUpdate: true,
					PluginMetadata: &SCMMetadata{
						ID:      "github.pr",
						Version: "1",
					},
					Configuration: []*SCMConfiguration{
						{Key: "username", Value: "admin"},
						{Key: "password", EncryptedValue: encPassword.EncryptedValue},
						{Key: "url", Value: "https://github.com/gocd/gocd.git"},
					},
				},
				scmCreateWant: &SCM{
					ID:         "mock-id",
					Name:       "mock-name",
					AutoUpdate: true,
					PluginMetadata: &SCMMetadata{
						ID:      "github.pr",
						Version: "1",
					},
					Configuration: []*SCMConfiguration{
						{Key: "username", Value: "admin"},
						{Key: "password", EncryptedValue: encPassword.EncryptedValue},
						{Key: "url", Value: "https://github.com/gocd/gocd.git"},
					},
				},
				scmGet: "mock-name",
				scmGetWant: &SCM{
					ID:         "mock-id",
					Name:       "mock-name",
					AutoUpdate: true,
					PluginMetadata: &SCMMetadata{
						ID:      "github.pr",
						Version: "1",
					},
					Configuration: []*SCMConfiguration{
						{Key: "username", Value: "admin"},
						{Key: "password", EncryptedValue: encPassword.EncryptedValue},
						{Key: "url", Value: "https://github.com/gocd/gocd.git"},
					},
				},
				scmUpdate: &SCM{
					ID:         "mock-id",
					Name:       "mock-name",
					AutoUpdate: true,
					PluginMetadata: &SCMMetadata{
						ID:      "github.pr",
						Version: "1",
					},
					Configuration: []*SCMConfiguration{
						{Key: "username", Value: "admin"},
						{Key: "password", EncryptedValue: encPassword.EncryptedValue},
						{Key: "url", Value: "https://github.com/gocd/other-gocd.git"},
					},
				},
				scmUpdateWant: &SCM{
					ID:         "mock-id",
					Name:       "mock-name",
					AutoUpdate: true,
					PluginMetadata: &SCMMetadata{
						ID:      "github.pr",
						Version: "1",
					},
					Configuration: []*SCMConfiguration{
						{Key: "username", Value: "admin"},
						{Key: "password", EncryptedValue: encPassword.EncryptedValue},
						{Key: "url", Value: "https://github.com/gocd/other-gocd.git"},
					},
				},
				scmListWant: []*SCM{
					{
						ID:   "mock-id",
						Name: "mock-name",
						PluginMetadata: &SCMMetadata{
							ID:      "github.pr",
							Version: "1",
						},
					},
				},
			},
		} {
			t.Run(test.name, func(t *testing.T) {

				var version string
				t.Run("create", func(t *testing.T) {
					scmCreateGot, _, err := intClient.SCM.Create(ctx, test.scmCreate)
					if assert.NoError(t, err) {
						scmCreateGot.RemoveLinks()
						scmCreateGot.Version = ""
						assert.Equal(t, test.scmCreateWant, scmCreateGot)
					}
				})

				t.Run("get", func(t *testing.T) {
					scmGetGot, _, err := intClient.SCM.Get(ctx, test.scmGet)
					if assert.NoError(t, err) {
						scmGetGot.RemoveLinks()
						version = scmGetGot.Version
						scmGetGot.Version = ""
						assert.Equal(t, test.scmGetWant, scmGetGot)
					}
				})

				t.Run("update", func(t *testing.T) {
					test.scmUpdate.Version = version
					scmUpdateGot, _, err := intClient.SCM.Update(ctx, test.scmUpdate.Name, test.scmUpdate)
					if assert.NoError(t, err) {
						scmUpdateGot.RemoveLinks()
						scmUpdateGot.Version = ""
						assert.Equal(t, test.scmUpdateWant, scmUpdateGot)
					}
				})

				t.Run("list", func(t *testing.T) {
					scmListGot, _, err := intClient.SCM.List(ctx)
					if assert.NoError(t, err) {
						for _, scm := range scmListGot {
							scm.Links = nil
						}
						assert.Equal(t, test.scmListWant, scmListGot)
					}
				})
			})
		}
	}
}
