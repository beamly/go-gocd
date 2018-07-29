package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRole(t *testing.T) {
	t.Run("GoCD", testRoleGoCD)
}

func testRoleGoCD(t *testing.T) {

	if runIntegrationTest(t) {

		ctx := context.Background()

		roles := []*Role{
			{
				Name: "spacetiger",
				Type: "gocd",
				Attributes: &RoleAttributesGoCD{
					Users: []string{"alice", "bob", "robin"},
				},
			},
			{
				Name: "my-mock-gocd-role",
				Type: "gocd",
				Attributes: &RoleAttributesGoCD{
					Users: []string{"user-one", "user-two"},
				},
			},
			// Currently there's no fixtures to test the plugin roles,
			// so until there is a way, we can not test plugin role types.
			//{
			//	Name: "blackbird",
			//	Type: "plugin",
			//	Attributes: &RoleAttributesGoCD{
			//		AuthConfigId: String("ldap"),
			//		Properties: []*RoleAttributeProperties{
			//			{
			//				Key:   "UserGroupMembershipAttribute",
			//				Value: "memberOf",
			//			},
			//			{
			//				Key:   "GroupIdentifiers",
			//				Value: "ou=admins,ou=groups,ou=system,dc=example,dc=com",
			//			},
			//		},
			//	},
			//},
		}

		// Test role creation
		for _, role := range roles {
			roleResponse, _, err := intClient.Roles.Create(ctx, role)
			assert.NoError(t, err)

			assert.Regexp(t, regexp.MustCompile("^[a-f0-9]{32}--gzip$"), roleResponse.Version)
			role.Version = roleResponse.Version

			assert.Equal(t, role, roleResponse)
		}

		// Test role listing
		rolesResponses, _, err := intClient.Roles.List(ctx)
		assert.NoError(t, err)

		for i, roleResponse := range rolesResponses {
			assert.Regexp(t, regexp.MustCompile("^[a-f0-9]{32}--gzip$"), roles[i].Version)
			roleResponse.Version = roles[i].Version

			assert.Equal(t, roles[i], roleResponse)
		}

		// Test role update
		roles[0].Attributes.Users = []string{"new-admin"}
		roleUpdateResponse, _, err := intClient.Roles.Update(ctx, roles[0].Name, roles[0])
		assert.NoError(t, err)
		updatedRole, _, err := intClient.Roles.Get(ctx, roleUpdateResponse.Name)
		assert.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile("^[a-f0-9]{32}--gzip$"), updatedRole.Version)
		roles[0].Version = updatedRole.Version
		assert.Equal(t, updatedRole, roles[0])

		// Test role delete
		for _, role := range roles {
			result, _, err := intClient.Roles.Delete(ctx, role.Name)
			assert.True(t, result)
			assert.NoError(t, err)
		}
		roleResponse, _, err := intClient.Roles.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, roleResponse)

	} else {
		skipIntegrationtest(t)
	}
}
