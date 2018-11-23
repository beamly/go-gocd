package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPluginApi(t *testing.T) {
	setup()
	defer teardown()

	t.Run("List", testPluginAPIList)
	t.Run("Get", testPluginAPIGet)
}

func testPluginAPIList(t *testing.T) {
	if runIntegrationTest(t) {
		ctx := context.Background()
		plugins, _, err := intClient.Plugins.List(ctx)
		assert.NoError(t, err)

		assert.NotNil(t, plugins)
		assert.NotNil(t, plugins.Links.Get("Doc"))
		assert.Equal(t, "https://api.gocd.org/#plugin-info", plugins.Links.Get("Doc").URL.String())
		assert.NotNil(t, plugins.Links.Get("Self"))
		assert.Equal(t, "http://127.0.0.1:8153/go/api/admin/plugin_info", plugins.Links.Get("Self").URL.String())

		assert.NotNil(t, plugins.Embedded)
		assert.NotNil(t, plugins.Embedded.PluginInfo)
		assert.Len(t, plugins.Embedded.PluginInfo, 6)
		pi := plugins.Embedded.PluginInfo[0]
		assert.Equal(t, "yum", pi.ID)

		apiVersion, err := client.getAPIVersion(ctx, "admin/plugin_info")
		assert.NoError(t, err)

		switch apiVersion {
		case apiV3:
			assert.Equal(t, "package-repository", pi.Type)
			assert.Equal(t, "active", pi.Status.State)
			assert.Equal(t, "Yum Plugin", pi.About.Name)
			assert.Equal(t, "PACKAGE_SPEC", pi.ExtensionInfo.PackageSettings.Configurations[0].Key)
			assert.Equal(t, "Package Spec", pi.ExtensionInfo.PackageSettings.Configurations[0].Metadata.DisplayName)
			assert.Equal(t, true, pi.ExtensionInfo.PackageSettings.Configurations[0].Metadata.Required)
		case apiV4:
			assert.Equal(t, "active", pi.Status.State)
			assert.Equal(t, "Yum Plugin", pi.About.Name)
			assert.Equal(t, "package-repository", pi.Extensions[0].Type)
			assert.Equal(t, "PACKAGE_SPEC", pi.Extensions[0].PackageSettings.Configurations[0].Key)
			assert.Equal(t, "Package Spec", pi.Extensions[0].PackageSettings.Configurations[0].Metadata.DisplayName)
			assert.Equal(t, true, pi.Extensions[0].PackageSettings.Configurations[0].Metadata.Required)
		default:
			t.Error("Unsupported api version in acceptance tests fo testPluginAPIList")
		}
	}
}

func testPluginAPIGet(t *testing.T) {
	if runIntegrationTest(t) {
		ctx := context.Background()
		plugin, _, err := intClient.Plugins.Get(ctx, "yum")
		assert.NoError(t, err)

		assert.NotNil(t, plugin)
		assert.NotNil(t, plugin.Links.Get("Doc"))
		assert.Equal(t, "https://api.gocd.org/#plugin-info", plugin.Links.Get("Doc").URL.String())
		assert.NotNil(t, plugin.Links.Get("Self"))
		assert.Equal(t, "http://127.0.0.1:8153/go/api/admin/plugin_info/yum", plugin.Links.Get("Self").URL.String())

		assert.Equal(t, "yum", plugin.ID)

		apiVersion, err := client.getAPIVersion(ctx, "admin/plugin_info")
		assert.NoError(t, err)

		switch apiVersion {
		case apiV3:
			assert.Equal(t, "package-repository", plugin.Type)
			assert.Equal(t, "active", plugin.Status.State)
			assert.Equal(t, "Yum Plugin", plugin.About.Name)
			assert.Equal(t, "PACKAGE_SPEC", plugin.ExtensionInfo.PackageSettings.Configurations[0].Key)
			assert.Equal(t, "Package Spec", plugin.ExtensionInfo.PackageSettings.Configurations[0].Metadata.DisplayName)
			assert.Equal(t, true, plugin.ExtensionInfo.PackageSettings.Configurations[0].Metadata.Required)
		case apiV4:
			assert.Equal(t, "active", plugin.Status.State)
			assert.Equal(t, "Yum Plugin", plugin.About.Name)
			assert.Equal(t, "package-repository", plugin.Extensions[0].Type)
			assert.Equal(t, "PACKAGE_SPEC", plugin.Extensions[0].PackageSettings.Configurations[0].Key)
			assert.Equal(t, "Package Spec", plugin.Extensions[0].PackageSettings.Configurations[0].Metadata.DisplayName)
			assert.Equal(t, true, plugin.Extensions[0].PackageSettings.Configurations[0].Metadata.Required)
		default:
			t.Error("Unsupported api version in acceptance tests fo testPluginAPIGet")
		}
	}
}
